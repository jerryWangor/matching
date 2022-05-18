package engine

import (
	"fmt"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/log"
	"strings"
)

var OrderChanMap map[string]chan model.Order

var AllOrderBookMap map[string]*model.OrderBook

// Init 初始化，从redis中恢复一些东西
func Init() {

	// 定义订单map通道
	OrderChanMap = make(map[string]chan model.Order)

	// 定义所有的交易委托账本map
	AllOrderBookMap = make(map[string]*model.OrderBook)

	// 从redis中查询所有已开启的交易标引擎，并重新开启
	symbols := cache.GetSymbols()
	for _, symbol := range symbols {
		price := cache.GetPrice(symbol)
		if e := NewEngine(symbol, price); e != nil {
			log.Error(fmt.Sprintf("交易标：%s，价格：%s 开启失败", symbol, price))
			continue
		}

		// 获取该交易标缓存的所有订单
		orders := cache.GetOrderIdsWithSymbol(symbol)
		for _, val := range orders {
			orderArr := strings.Split(val,":")
			order := cache.GetOrder(symbol, orderArr[0], orderArr[1])
			OrderChanMap[order.Symbol] <- order
		}
	}
}