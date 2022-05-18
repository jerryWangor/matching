package cache

import (
	"github.com/shopspring/decimal"
	"matching/model"
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
	orderMap := order.ToMap()
	cache.SaveOrder(orderMap)
}

// GetOrder 获取订单
func GetOrder(symbol string, orderid string, action string) model.Order {
	orderMap := cache.GetOrder(symbol, orderid, action)
	var order model.Order
	order.FromMap(orderMap)
	return order
}

// UpdateOrder 更新订单
func UpdateOrder(order model.Order) {
	orderMap := order.ToMap()
	cache.UpdateOrder(orderMap)
}

// RemoveOrder 删除订单
func RemoveOrder(order model.Order) bool {
	orderMap := order.ToMap()
	err := cache.RemoveOrder(orderMap)
	if err != nil {
		common.Errors(err.Error())
		return false
	} else {
		return true
	}
}

// OrderExist 判断订单是否存在
func OrderExist(symbol, orderId, action string) bool {
	return cache.OrderExist(symbol, orderId, action)
}

// GetOrderIdsWithSymbol 获取交易标下的所有订单IDS
func GetOrderIdsWithSymbol(symbol string) []string {
	return cache.GetOrderIdsWithSymbol(symbol)
}