package cache

import (
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils/common"
	"matching/utils/enum"
	"matching/utils/log"
	cache "matching/utils/redis"
	"strconv"
)

// SaveSymbol 保存交易标
func SaveSymbol(symbol string) {
	cache.SaveSymbol(symbol)
}

// RemoveSymbol 删除交易标
func RemoveSymbol(symbol string) {
	cache.RemoveSymbol(symbol)
}

// HasSymbol 判断缓存中是否有交易标
func HasSymbol(symbol string) bool {
	return cache.HasSymbol(symbol)
}

// GetSymbols 获取所有的交易标
func GetSymbols() []string {
	return cache.GetSymbols()
}

// SavePrice 保存交易标价格
func SavePrice(symbol string, price decimal.Decimal) {
	cache.SavePrice(symbol, price)
}

// GetPrice 获取交易标价格
func GetPrice(symbol string) decimal.Decimal {
	return cache.GetPrice(symbol)
}

// RemovePrice 删除交易标价格
func RemovePrice(symbol string) {
	cache.RemovePrice(symbol)
}

// SaveOrder 保存订单
func SaveOrder(order model.Order) {
	// 这里不用ToMap，手动进行转换
	var orderMap = make(map[string]interface{})
	orderMap["symbol"] = order.Symbol
	orderMap["orderId"] = order.OrderId
	orderMap["action"] = int(order.Action)
	orderMap["type"] = int(order.Type)
	orderMap["side"] = int(order.Side)
	orderMap["amount"] = order.Amount.String()
	orderMap["price"] = order.Price.String()
	orderMap["timestamp"], _ = strconv.ParseFloat(strconv.FormatInt(order.Timestamp,10), 64)
	cache.SaveOrder(orderMap)
}

// GetOrder 获取订单
func GetOrder(symbol string, orderid string) model.Order {
	orderMap := cache.GetOrder(symbol, orderid)

	// 这里可能有问题
	action, _ := strconv.Atoi(orderMap["action"].(string))
	otype, _ := strconv.Atoi(orderMap["type"].(string))
	side, _ := strconv.Atoi(orderMap["side"].(string))
	amount, _ := decimal.NewFromString(orderMap["amount"].(string))
	price, _ := decimal.NewFromString(orderMap["price"].(string))
	timestamp, _ := strconv.ParseInt(orderMap["timestamp"].(string), 10, 64)

	order := model.Order {
		Symbol: orderMap["symbol"].(string),
		OrderId: orderMap["orderId"].(string),
		Action: enum.OrderAction(action),
		Type: enum.OrderType(otype),
		Side: enum.OrderSide(side),
		Amount: amount,
		Price: price,
		Timestamp: timestamp, // 这里可能有精度问题
	}
	return order
}

// UpdateOrder 更新订单
func UpdateOrder(order model.Order) {
	maporder, err := common.ToMap(order)
	if err != nil {
		log.Error("订单更新失败")
	}
	cache.UpdateOrder(maporder)
}

// RemoveOrder 删除订单
func RemoveOrder(order model.Order) {
	maporder, err := common.ToMap(order)
	if err != nil {
		log.Error("订单删除失败")
	}
	cache.RemoveOrder(maporder)
}

// OrderExist 判断订单是否存在
func OrderExist(symbol string, orderId string) bool {
	return cache.OrderExist(symbol, orderId)
}

// GetOrderIdsWithSymbol 获取交易标下的所有订单IDS
func GetOrderIdsWithSymbol(symbol string) []string {
	return cache.GetOrderIdsWithSymbol(symbol)
}

// SendCancelResult 发送撤单结果消息队列
func SendCancelResult(symbol, orderId string, ok bool) {
	cache.SendCancelResult(symbol, orderId, ok)
}

// SendTrade 发送交易记录消息队列
func SendTrade(symbol string, trade map[string]interface{}) {
	cache.SendTrade(symbol, trade)
}
