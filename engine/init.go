package engine

import (
	"fmt"
	"matching/model"
	"matching/utils"
	cache "matching/utils/redis"
)

var OrderChanMap map[string]chan model.Order

// 初始化，从redis中恢复一些东西
func Init() {

	// 定义订单map通道
	OrderChanMap = make(map[string]chan model.Order)

	// 从redis中查询所有已开启的交易标引擎，并重新开启
	symbols := cache.GetSymbols()
	for _, symbol := range symbols {
		price := cache.GetPrice(symbol)
		if e := NewEngine(symbol, price); e != nil {
			utils.LogError(fmt.Sprintf("交易标：%s，价格：%s 开启失败", symbol, price))
			continue
		}

		orderIds := cache.GetOrderIdsWithAction(symbol)
		for _, orderId := range orderIds {
			mapOrder := cache.GetOrder(symbol, orderId)
			order := model.Order{}
			order.FromMap(mapOrder)
			OrderChanMap[order.Symbol] <- order
		}
	}

	// 从redis中查询所有定序订单，并生成订单委托账本


}