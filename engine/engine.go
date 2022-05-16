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

	return nil
}

// Run
// 引擎启动入口（协程处理）
func Run(symbol string, price decimal.Decimal) {
	lastTradePrice := price

	// 初始化交易委托账本
	book := &model.OrderBook{}
	book.Init()

	// 把该交易标账本放到总账本里面
	AllOrderBookMap[symbol] = book

	log.Info("engine %s is running", symbol)
	for {
		// 监听订单通道进行操作
		order, ok := <-OrderChanMap[symbol]
		fmt.Println(order)
		if !ok {
			// 如果通道关闭就关闭引擎
			log.Info("engine %s is closed", symbol)
			delete(OrderChanMap, symbol)
			cache.RemoveSymbol(symbol)
			return
		}
		log.Info("engine %s receive an order: %s", symbol, common.ToJson(order))
		switch order.Action {
		case enum.ActionCreate:
			dealCreate(&order, book, &lastTradePrice)
		case enum.ActionCancel:
			dealCancel(&order, book)
		}
	}
}

// 挂单处理
func dealCreate(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
	switch order.Type {
	case enum.TypeLimit:
		dealLimit(order, book, lastTradePrice)
	//case enum.TypeLimitIoc:
	//	dealLimitIoc(order, book, lastTradePrice)
	//case enum.TypeMarket:
	//	dealMarket(order, book, lastTradePrice)
	//case enum.TypeMarketTop5:
	//	dealMarketTop5(order, book, lastTradePrice)
	//case enum.TypeMarketTop10:
	//	dealMarketTop10(order, book, lastTradePrice)
	//case enum.TypeMarketOpponent:
	//	dealMarketOpponent(order, book, lastTradePrice)
	}
}

// 普通限价处理
func dealLimit(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
	switch order.Side {
	case enum.SideBuy:
		dealBuyLimit(order, book, lastTradePrice)
	case enum.SideSell:
		dealSellLimit(order, book, lastTradePrice)
	}
}

// 买单处理
func dealBuyLimit(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
LOOP:
	headOrder := book.GetHeadSellOrder()
	// 买单价格小于卖单价格，不能成交
	if headOrder == nil || order.Price.LessThan(headOrder.Price) {
		book.AddBuyOrder(order)
		log.Info("engine %s, a order has added to the orderbook: %s", order.Symbol, common.ToJson(order))
	} else {
		matchTrade(headOrder, order, book, lastTradePrice)
		if order.Amount.IsPositive() {
			goto LOOP
		}
	}
}

// 卖单处理
func dealSellLimit(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
LOOP:
	headOrder := book.GetHeadBuyOrder()
	// 卖单价格大于买单价格，不能成交
	if headOrder == nil || order.Price.GreaterThan(headOrder.Price) {
		book.AddSellOrder(order)
		log.Info("engine %s, a order has added to the orderbook: %s", order.Symbol, common.ToJson(order))
	} else {
		matchTrade(headOrder, order, book, lastTradePrice)
		if order.Amount.IsPositive() {
			goto LOOP
		}
	}
}

// 撮合订单
func matchTrade(headOrder *model.Order, order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
	// 将头部订单和当前订单进行撮合，然后更新交易委托账本

	var trade *model.Trade
	var useAmount decimal.Decimal

	// 判断订单是买还是卖
	if order.Side == enum.SideBuy {
		// 买单数量去吃头部单的数量
		// 如果买单数量>=头部订单数量，那么就把头部订单完全吃掉
		if order.Amount.GreaterThanOrEqual(headOrder.Amount) {
			useAmount = headOrder.Amount
			order.Amount = order.Amount.Sub(headOrder.Amount)
			// 把头部订单从交易委托账本中去掉
			book.PopHeadSellOrder()
			// 把头部订单从缓存中去掉
			cache.RemoveOrder(*headOrder)
		} else {
			// 如果买单数量<头部订单数量，那么就只消耗了买单数量
			useAmount = order.Amount
			headOrder.Amount = headOrder.Amount.Sub(order.Amount)
		}
	} else {
		// 卖单数量去吃头部单的数量
		// 如果卖单数量>=头部订单数量，那么就把头部订单完全吃掉
		if order.Amount.GreaterThanOrEqual(headOrder.Amount) {
			useAmount = headOrder.Amount
			order.Amount = order.Amount.Sub(headOrder.Amount)
			// 把头部订单从交易委托账本中去掉
			book.PopHeadBuyOrder()
			// 把头部订单从缓存中去掉
			cache.RemoveOrder(*headOrder)
		} else {
			// 如果买单数量<头部订单数量，那么就只消耗了买单数量
			useAmount = order.Amount
			headOrder.Amount = headOrder.Amount.Sub(order.Amount)
		}
	}

	// 更新价格
	lastTradePrice = &headOrder.Price

	// 生成交易记录，推给消息队列
	trade = &model.Trade{
		MakerId: headOrder.OrderId,
		TakerId: order.OrderId,
		TakerSide: order.Side,
		Amount: useAmount,
		Price: order.Price,
		Timestamp: time.Now().UnixMicro(),
	}
	mapTrade, _ := common.ToMap(trade)
	mq.SendTrade(order.Symbol, mapTrade)
}

// 撤单处理
func dealCancel(order *model.Order, book *model.OrderBook) {
	// 撤单直接删除redis，从交易委托账本里面移除
	cache.RemoveOrder(*order)
	if order.Side == enum.SideBuy {
		book.RemoveBuyOrder(order)
	} else {
		book.RemoveSellOrder(order)
	}
	// 发送消息队列
	mq.SendCancelResult(order.Symbol, order.OrderId, true)
}

// Dispatch 分发订单
func Dispatch(order model.Order) error {
	if OrderChanMap[order.Symbol] == nil {
		return common.Errors(fmt.Sprintf("%s 引擎未启动", order.Symbol))
	}

	// 挂单判断存在不
	if order.Action == enum.ActionCreate {
		if cache.OrderExist(order.Symbol, order.OrderId) {
			return common.Errors(fmt.Sprintf("%s-%s 订单已存在，不能重复挂单", order.Symbol, order.OrderId))
		}
	} else {
		// 撤单如果订单不存在就没法撤
		if !cache.OrderExist(order.Symbol, order.OrderId) {
			return common.Errors(fmt.Sprintf("%s-%s 订单不存在，无法撤单", order.Symbol, order.OrderId))
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