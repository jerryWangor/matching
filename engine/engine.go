package engine

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"matching/model"
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

}

// 分发订单
func Dispatch(order model.Order) error {
	if OrderChanMap[order.Symbol] == nil {
		return common.Errors(fmt.Sprintf("%s 引擎未启动", order.Symbol))
	}

	if order.Action == enum.ActionCreate {
		if cache.OrderExist(order.Symbol, order.OrderId, order.Action.String()) {

		}
	} else {
		if !cache.OrderExist(order.Symbol, order.OrderId, enum.ActionCreate.String()) {
			return errcode.OrderNotFound
		}
	}

	order.Timestamp = time.Now().UnixNano() / 1e3
	cache.SaveOrder(order.ToMap())
	engine.ChanMap[order.Symbol] <- order

	return nil
}