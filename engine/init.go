package engine

import (
	"fmt"
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/log"
	"strings"
)

var OrderChanMap map[string]chan model.Order

var AllOrderBookMap map[string]*model.OrderBook

var StopKDataChan chan bool

var KDataPriceMap map[string]*model.KDataPrice

// Test相关
var AllOrderAmountMap map[string]decimal.Decimal // 用于记录接收到的买单和卖单数量，撮合会更新

// Init 初始化，从redis中恢复一些东西
func Init() {

	AllOrderAmountMap = make(map[string]decimal.Decimal)
	AllOrderAmountMap["buy"] = decimal.NewFromFloat(0)
	AllOrderAmountMap["sell"] = decimal.NewFromFloat(0)

	// 定义订单map通道
	OrderChanMap = make(map[string]chan model.Order)

	// 定义所有的交易委托账本map
	AllOrderBookMap = make(map[string]*model.OrderBook)

	// 定义k线图开关通道
	StopKDataChan = make(chan bool, 1)

	// 定义价格map
	KDataPriceMap = make(map[string]*model.KDataPrice)

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

	// 开启kdata线程
	runKData()

}