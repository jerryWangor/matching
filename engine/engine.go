package engine

import (
	"fmt"
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils"
	"matching/utils/common"
	"matching/utils/enum"
	cache "matching/utils/redis"
	"time"
)

// 开启一个新的引擎
func NewEngine(symbol string, price decimal.Decimal) error {
	if OrderChanMap[symbol] != nil {
		return common.Errors(fmt.Sprintf("%s 引擎重复开启，请先关闭已开启的引擎", symbol))
	}

	OrderChanMap[symbol] = make(chan model.Order, 1000)
	go Run(symbol, price)

	cache.SaveSymbol(symbol)
	cache.SavePrice(symbol, price)

	return nil
}

// 引擎启动入口（协程处理）
func Run(symbol string, price decimal.Decimal) {
	lastTradePrice := price

	book := &model.OrderBook{}
	book.Init()

	utils.LogInfo("engine %s is running", symbol)
	for {
		order, ok := <-OrderChanMap[symbol]
		if !ok {
			// 如果通道关闭就关闭引擎
			utils.LogInfo("engine %s is closed", symbol)
			delete(OrderChanMap, symbol)
			cache.RemoveSymbol(symbol)
			return
		}
		utils.LogInfo("engine %s receive an order: %s", symbol, order.ToJson())
		switch order.Action {
		case enum.ActionCreate:
			dealCreate(&order, book, &lastTradePrice)
		case enum.ActionCancel:
			dealCancel(&order, book)
		}
	}
}

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

func dealLimit(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
	switch order.Side {
	case enum.SideBuy:
		dealBuyLimit(order, book, lastTradePrice)
	case enum.SideSell:
		dealSellLimit(order, book, lastTradePrice)
	}
}

func dealBuyLimit(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {
LOOP:
	headOrder := book.GetHeadSellOrder()
	if headOrder == nil || order.Price.LessThan(headOrder.Price) {
		book.AddBuyOrder(order)
		utils.LogInfo("engine %s, a order has added to the orderbook: %s", order.Symbol, order.ToJson())
	} else {
		matchTrade(headOrder, order, book, lastTradePrice)
		if order.Amount.IsPositive() {
			goto LOOP
		}
	}
}

func matchTrade(headOrder *model.Order, order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {

}

func dealSellLimit(order *model.Order, book *model.OrderBook, lastTradePrice *decimal.Decimal) {

}

func dealCancel(order *model.Order, book *model.OrderBook) {

}

// 分发订单
func Dispatch(order model.Order) error {
	if OrderChanMap[order.Symbol] == nil {
		return common.Errors(fmt.Sprintf("%s 引擎未启动", order.Symbol))
	}


	// 挂单判断存在不
	if order.Action == enum.ActionCreate {
		if cache.OrderExist(order.Symbol, order.OrderId, order.Action.String()) {
			return common.Errors(fmt.Sprintf("%s-%s 订单已存在，不能重复挂单", order.Symbol, order.OrderId))
		}
	} else {
		// 撤单如果订单不存在就没法撤
		if !cache.OrderExist(order.Symbol, order.OrderId, enum.ActionCreate.String()) {
			return common.Errors(fmt.Sprintf("%s-%s 订单不存在，无法撤单", order.Symbol, order.OrderId))
		}
	}

	order.Timestamp = time.Now().UnixMicro()
	ordermap, err := order.ToMap()
	if err != nil {
		return common.Errors(fmt.Sprintf("%s-%s 订单不存在，无法撤单", order.Symbol, order.OrderId))
	}

	cache.SaveOrder(ordermap)
	OrderChanMap[order.Symbol] <- order

	return nil
}

// 关闭引擎
func CloseEngine(symbol string) error {
	if OrderChanMap[symbol] == nil {
		return common.Errors(fmt.Sprintf("%s 引擎未启动", symbol))
	}

	// 关闭通道，引擎里面同步做操作
	close(OrderChanMap[symbol])

	return nil
}