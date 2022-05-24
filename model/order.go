package model

import (
	"github.com/shopspring/decimal"
	"matching/utils/enum"
	"strconv"
)

type Order struct {
	//Accid int `json:"accid" comment:"账号ID"`
	Symbol string `json:"symbol" comment:"交易标"`
	OrderId string `json:"orderId" comment:"订单ID"`
	ChildOrderId string `json:"childOrderId" comment:"子订单ID，如果被拆单可能用到"`
	Action enum.OrderAction `json:"action" comment:"挂单还是撤单"`
	Type enum.OrderType `json:"type" comment:"竞价类型"`
	Side enum.OrderSide `json:"side" comment:"买/卖"`
	Amount decimal.Decimal `json:"amount" comment:"数量"`
	Price decimal.Decimal `json:"price" comment:"价格"`
	Timestamp int64 `json:"timestamp" comment:"时间"`
}

// FromMap Map转结构体
func (o *Order) FromMap(orderMap map[string]interface{}) {

	// 这里可能有问题
	action, _ := strconv.Atoi(orderMap["action"].(string))
	otype, _ := strconv.Atoi(orderMap["type"].(string))
	side, _ := strconv.Atoi(orderMap["side"].(string))
	amount, _ := decimal.NewFromString(orderMap["amount"].(string))
	price, _ := decimal.NewFromString(orderMap["price"].(string))
	timestamp, _ := strconv.ParseInt(orderMap["timestamp"].(string), 10, 64)

	o.Symbol = orderMap["symbol"].(string)
	o.OrderId = orderMap["orderId"].(string)
	o.Action = enum.OrderAction(action)
	o.Type = enum.OrderType(otype)
	o.Side = enum.OrderSide(side)
	o.Amount = amount
	o.Price = price
	o.Timestamp = timestamp
}

// ToMap 结构体转成map
func (o *Order) ToMap() map[string]interface{} {
	var orderMap = make(map[string]interface{})
	orderMap["symbol"] = o.Symbol
	orderMap["orderId"] = o.OrderId
	orderMap["action"] = int(o.Action)
	orderMap["type"] = int(o.Type)
	orderMap["side"] = int(o.Side)
	orderMap["amount"] = o.Amount.String()
	orderMap["price"] = o.Price.String()
	orderMap["timestamp"], _ = strconv.ParseFloat(strconv.FormatInt(o.Timestamp,10), 64)

	return orderMap
}