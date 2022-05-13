package cache

import (
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils"
	cache "matching/utils/redis"
)

// 缓存操作
// 交易标函数
func SaveSymbol(symbol string) {
	cache.SaveSymbol(symbol)
}

func RemoveSymbol(symbol string) {
	cache.RemoveSymbol(symbol)
}

func HasSymbol(symbol string) bool {
	return cache.HasSymbol(symbol)
}

func GetSymbols() []string {
	return cache.GetSymbols()
}

// 价格函数
func SavePrice(symbol string, price decimal.Decimal) {
	cache.SavePrice(symbol, price)
}

func GetPrice(symbol string) decimal.Decimal {
	return cache.GetPrice(symbol)
}

func RemovePrice(symbol string) {
	cache.RemovePrice(symbol)
}

// 订单函数
func SaveOrder(order model.Order) {
	maporder, err := order.ToMap()
	if err != nil {
		utils.LogError("订单保存失败")
	}
	cache.SaveOrder(maporder)
}

func GetOrder(symbol string, orderid string) map[string]interface{} {
	return cache.GetOrder(symbol, orderid)
}

func UpdateOrder() {

}

func RemoveOrder() {

}

func OrderExist() {

}

func GetOrderIdsWithAction(symbol string) []string {
	return cache.GetOrderIdsWithAction(symbol)
}