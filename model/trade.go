package model

import (
	"github.com/shopspring/decimal"
	"matching/utils/enum"
	"strconv"
)

// Trade 交易记录
type Trade struct {
	MakerId string `json:"makerId" comment:"挂单的订单ID"`
	TakerId string `json:"takerId" comment:"吃单的订单ID"`
	TakerSide enum.OrderSide `json:"takerSide" comment:"买还是卖"`
	Amount decimal.Decimal `json:"amount" comment:"交易成功的数量"`
	Price decimal.Decimal `json:"price" comment:"当前交易价格"`
	Timestamp int64 `json:"timestamp" comment:"交易时间"`
}

// FromMap Map转结构体
func (t *Trade) FromMap(tradeMap map[string]interface{}) {

	// 这里可能有问题
	takerSide, _ := strconv.Atoi(tradeMap["takerSide"].(string))
	amount, _ := decimal.NewFromString(tradeMap["amount"].(string))
	price, _ := decimal.NewFromString(tradeMap["price"].(string))
	timestamp, _ := strconv.ParseInt(tradeMap["timestamp"].(string), 10, 64)

	t.MakerId = tradeMap["makerId"].(string)
	t.TakerId = tradeMap["takerId"].(string)
	t.TakerSide = enum.OrderSide(takerSide)
	t.Amount = amount
	t.Price = price
	t.Timestamp = timestamp
}

// ToMap 结构体转成map
func (t *Trade) ToMap() map[string]interface{} {
	var tradeMap = make(map[string]interface{})
	tradeMap["makerId"] = t.MakerId
	tradeMap["takerId"] = t.TakerId
	tradeMap["takerSide"] = int(t.TakerSide)
	tradeMap["amount"] = t.Amount.String()
	tradeMap["price"] = t.Price.String()
	tradeMap["timestamp"], _ = strconv.ParseFloat(strconv.FormatInt(t.Timestamp,10), 64)

	return tradeMap
}