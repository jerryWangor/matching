package engine

import (
	"fmt"
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/common"
	"matching/utils/enum"
	"matching/utils/log"
	"matching/utils/mq"
	"time"
)

// NewEngine
// 开启一个新的引擎
func NewEngine(symbol string, price decimal.Decimal) error {
	if OrderChanMap[symbol] != nil {
		return common.Errors(fmt.Sprintf("%s 引擎重复开启，请先关闭已开启的引擎", symbol))
	}

	OrderChanMap[symbol] = make(chan model.Order, 100)
	go Run(symbol, price)

	cache.SaveSymbol(symbol)
	cache.SavePrice(symbol, price)

	// k线图
	go handleKData(symbol, KDataPriceMap[symbol])

	return nil
}

// Run
// 引擎启动入口（协程处理）
func Run(symbol string, price decimal.Decimal) {

	// 初始化交易委托账本
	book := &model.OrderBook{}
	book.Init()

	// 把该交易标账本放到总账本里面
	AllOrderBookMap[symbol] = book

	// 初始化价格
	KDataPriceMap[symbol] = &model.KDataPrice{
		TopPrice: price,
		BottomPrice: price,
		NowPrice: price,
	}

	log.Info("engine %s is running", symbol)
	for {
		// 监听订单通道进行操作
		order, ok := <-OrderChanMap[symbol]
		if !ok {
			// 如果通道关闭就关闭引擎
			log.Info("engine %s is closed", symbol)
			delete(OrderChanMap, symbol)
			delete(AllOrderBookMap, symbol)
			cache.RemoveSymbol(symbol)
			// 关闭K线图
			StopKDataChan<-true
			return
		}
		log.Info("engine %s receive an order: %s", symbol, common.ToJson(order))
		// 如果价格等于0，设置成第一个买单的价格
		if KDataPriceMap[symbol].NowPrice.IsZero() && order.Side == enum.SideBuy {
			KDataPriceMap[symbol].TopPrice = order.Price
			KDataPriceMap[symbol].BottomPrice = order.Price
			KDataPriceMap[symbol].NowPrice = order.Price
			cache.SavePrice(symbol, KDataPriceMap[symbol].NowPrice)
		}
		switch order.Action {
		case enum.ActionCreate:
			dealCreate(&order, book, KDataPriceMap[symbol])
		case enum.ActionCancel:
			dealCancel(&order, book)
		}

		// top先同步进行处理，后期可以考虑用通道异步处理
		handleTopN(symbol, &KDataPriceMap[symbol].NowPrice, book, 5)
	}
}

// 挂单处理
func dealCreate(order *model.Order, book *model.OrderBook, kDataPrice *model.KDataPrice) {
	switch order.Type {
	case enum.TypeLimit:
		dealLimit(order, book, kDataPrice)
	//case enum.TypeLimitIoc:
	//	dealLimitIoc(order, book, kDataPrice)
	//case enum.TypeMarket:
	//	dealMarket(order, book, kDataPrice)
	//case enum.TypeMarketTop5:
	//	dealMarketTop5(order, book, kDataPrice)
	//case enum.TypeMarketTop10:
	//	dealMarketTop10(order, book, kDataPrice)
	//case enum.TypeMarketOpponent:
	//	dealMarketOpponent(order, book, kDataPrice)
	}
}

// 普通限价处理
func dealLimit(order *model.Order, book *model.OrderBook, kDataPrice *model.KDataPrice) {
	switch order.Side {
	case enum.SideBuy:
		dealBuyLimit(order, book, kDataPrice)
	case enum.SideSell:
		dealSellLimit(order, book, kDataPrice)
	}
}

// 买单处理
func dealBuyLimit(order *model.Order, book *model.OrderBook, kDataPrice *model.KDataPrice) {
LOOP:
	headOrder := book.GetHeadSellOrder()
	// 买单价格小于卖单价格，不能成交
	if headOrder == nil || order.Price.LessThan(headOrder.Price) {
		book.AddBuyOrder(order)
		log.Info("engine %s, add a buy order to the orderbook: %s", order.Symbol, common.ToJson(order))
	} else {
		matchTrade(headOrder, order, book, kDataPrice)
		if order.Amount.IsPositive() {
			goto LOOP
		}
	}
}

// 卖单处理
func dealSellLimit(order *model.Order, book *model.OrderBook, kDataPrice *model.KDataPrice) {
LOOP:
	headOrder := book.GetHeadBuyOrder()
	// 卖单价格大于买单价格，不能成交
	if headOrder == nil || order.Price.GreaterThan(headOrder.Price) {
		book.AddSellOrder(order)
		log.Info("engine %s, add a sell order to the orderbook: %s", order.Symbol, common.ToJson(order))
	} else {
		matchTrade(headOrder, order, book, kDataPrice)
		if order.Amount.IsPositive() {
			goto LOOP
		}
	}
}

// 撮合订单
func matchTrade(headOrder *model.Order, order *model.Order, book *model.OrderBook, kDataPrice *model.KDataPrice) {
	// 将头部订单和当前订单进行撮合，然后更新交易委托账本
	var trade *model.Trade
	var useAmount decimal.Decimal

	// 如果当前单数量>=头部订单数量，那么就把头部订单完全吃掉
	if order.Amount.GreaterThanOrEqual(headOrder.Amount) {
		common.Debugs("当前订单数量>=头部订单")
		useAmount = headOrder.Amount
		order.Amount = order.Amount.Sub(headOrder.Amount)
		if order.Side == enum.SideBuy {
			// 把头部订单从交易委托账本中去掉
			book.PopHeadSellOrder()
			// 删除element账本中该头部订单
			book.RemoveSellElementOrder(headOrder)
		} else {
			// 把头部订单从交易委托账本中去掉
			book.PopHeadBuyOrder()
			// 删除element账本中该头部订单
			book.RemoveBuyElementOrder(headOrder)
		}
		// 把头部订单从缓存中去掉
		cache.RemoveOrder(*headOrder)
		// 更新当前订单缓存
		cache.UpdateOrder(*order)
	} else {
		common.Debugs("当前订单数量<头部订单")
		// 如果买单数量<头部订单数量，那么就只消耗了买单数量
		useAmount = order.Amount
		order.Amount = order.Amount.Sub(order.Amount)
		headOrder.Amount = headOrder.Amount.Sub(useAmount)
		//common.Debugs(fmt.Sprintf("order amount：%s， head amount：%s", order.Amount.String(), headOrder.Amount.String()))
		// 删除当前订单缓存
		cache.RemoveOrder(*order)
		// 更新头部订单缓存
		cache.UpdateOrder(*headOrder)
		var res error
		if order.Side == enum.SideBuy {
			// 更新账本头部订单
			res = book.UpdateHeadSellOrder(headOrder)
			// 更新element账本中该头部订单数量
			book.UpdateSellElementOrder(headOrder)
		} else {
			// 更新账本头部订单
			res = book.UpdateHeadBuyOrder(headOrder)
			// 更新element账本中该头部订单数量
			book.UpdateBuyElementOrder(headOrder)
		}
		if res != nil {
			common.Debugs("更新账本头部订单失败")
		}
	}

	// 更新价格，这里是取地址，如果后面有问题可能要修改
	kDataPrice.NowPrice = headOrder.Price
	common.Debugs(fmt.Sprintf("当前价格更新为：%s", kDataPrice.NowPrice.String()))
	cache.SavePrice(order.Symbol, kDataPrice.NowPrice)
	// 判断价格
	if kDataPrice.NowPrice.GreaterThan(kDataPrice.TopPrice) {
		kDataPrice.TopPrice = kDataPrice.NowPrice
	}
	if kDataPrice.NowPrice.LessThan(kDataPrice.BottomPrice) {
		kDataPrice.BottomPrice = kDataPrice.NowPrice
	}

	// 生成交易记录，推给消息队列
	trade = &model.Trade{
		MakerId: headOrder.OrderId,
		TakerId: order.OrderId,
		TakerSide: order.Side,
		Amount: useAmount,
		Price: headOrder.Price,
		Timestamp: time.Now().UnixMicro(),
	}
	tradeMap := trade.ToMap()
	common.Debugs("交易记录：" + common.ToJson(trade))
	mq.SendTradeResult(order.Symbol, tradeMap)
}

// 撤单处理
func dealCancel(order *model.Order, book *model.OrderBook) {
	// 撤单直接删除redis，从交易委托账本里面移除
	var result, eResult, cres bool
	// 从缓存取出订单
	newOrder := cache.GetOrder(order.Symbol, order.OrderId, enum.ActionCreate.String())
	// 先删除账本里面的，这里买单如果传action=1，side=1的话，就不会删除账本里的，然后就会出问题
	if order.Side == enum.SideBuy {
		result = book.RemoveBuyOrder(&newOrder)
		// 删除element账本中该头部订单
		eResult = book.RemoveBuyElementOrder(&newOrder)
	} else {
		result = book.RemoveSellOrder(&newOrder)
		// 删除element账本中该头部订单
		eResult = book.RemoveSellElementOrder(&newOrder)
	}
	if result == true {
		common.Debugs(fmt.Sprintf("交易委托账本中的订单删除成功！%s-%s", order.Symbol, order.OrderId))
	} else {
		log.Error(fmt.Sprintf("交易委托账本中的订单删除失败！%s-%s", order.Symbol, order.OrderId))
	}
	if eResult == true {
		common.Debugs(fmt.Sprintf("TopN数据中的订单删除成功！%s-%s", order.Symbol, order.OrderId))
	} else {
		log.Error(fmt.Sprintf("TopN数据中的订单删除失败！%s-%s", order.Symbol, order.OrderId))
	}

	// 缓存删除状态
	cres = cache.RemoveOrder(newOrder)
	if cres == true {
		common.Debugs(fmt.Sprintf("订单缓存删除成功：%s-%s", order.Symbol, order.OrderId))
	} else {
		common.Debugs(fmt.Sprintf("订单缓存删除失败：%s-%s", order.Symbol, order.OrderId))
	}

	// 账本和缓存都删除成功了才算撤单成功
	if result == true && cres == true && eResult == true {
		// 删除缓存中的cancel订单
		cache.RemoveOrder(*order)
		mq.SendCancelResult(order.Symbol, order.OrderId, true)
	} else {
		mq.SendCancelResult(order.Symbol, order.OrderId, false)
	}

}

// Dispatch 分发订单
func Dispatch(order model.Order) error {
	if OrderChanMap[order.Symbol] == nil {
		return common.Errors(fmt.Sprintf("%s 引擎未启动", order.Symbol))
	}

	// 挂单判断存在不
	if order.Action == enum.ActionCreate {
		if cache.OrderExist(order.Symbol, order.OrderId, enum.ActionCreate.String()) {
			return common.Errors(fmt.Sprintf("%s-%s 订单已存在，不能重复挂单", order.Symbol, order.OrderId))
		}
	} else {
		// 撤单如果挂单不存在就没法撤
		if !cache.OrderExist(order.Symbol, order.OrderId, enum.ActionCreate.String()) {
			return common.Errors(fmt.Sprintf("%s-%s 订单不存在，无法撤单", order.Symbol, order.OrderId))
		}
		// 如果撤单已存在，不能重复撤单
		if cache.OrderExist(order.Symbol, order.OrderId, enum.ActionCancel.String()) {
			return common.Errors(fmt.Sprintf("%s-%s 撤单已存在，请勿重复撤单！", order.Symbol, order.OrderId))
		}
	}

	order.Timestamp = time.Now().UnixMicro()
	cache.SaveOrder(order)
	OrderChanMap[order.Symbol] <- order

	return nil
}

// CloseEngine 关闭引擎
func CloseEngine(symbol string) error {
	if OrderChanMap[symbol] == nil {
		return common.Errors(fmt.Sprintf("%s 引擎未启动", symbol))
	}

	// 关闭通道，引擎里面同步做操作
	close(OrderChanMap[symbol])

	return nil
}