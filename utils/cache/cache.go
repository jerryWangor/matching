package cache

import (
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils"
	"matching/utils/common"
	cache "matching/utils/redis"
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
	maporder, err := common.ToMap(order)
	if err != nil {
		utils.LogError("订单保存失败")
	}
	cache.SaveOrder(maporder)
}

// GetOrder 获取订单
func GetOrder(symbol string, orderid string) map[string]interface{} {
	return cache.GetOrder(symbol, orderid)
}

// UpdateOrder 更新订单
func UpdateOrder(order model.Order) {
	maporder, err := common.ToMap(order)
	if err != nil {
		utils.LogError("订单更新失败")
	}
	cache.UpdateOrder(maporder)
}

// RemoveOrder 删除订单
func RemoveOrder(order model.Order) {
	maporder, err := common.ToMap(order)
	if err != nil {
		utils.LogError("订单删除失败")
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
