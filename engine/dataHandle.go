package engine

import (
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/common"
	"time"
)

// 重新生成topN，保存redis
func handleTopN(symbol string, price *decimal.Decimal, book *model.OrderBook, num int) {
	// topN直接从交易委托账本中生成，这里可能需要优化，后期再说
	// 理论上头部买单！=头部卖单价格，基本逻辑就是
	// 如果头部买单=当前价格
	//		头部卖单>当前价格，从买单里面取出num个，从卖单里面取出num个
	// 如果头部买单<当前价格
	//		头部卖单=当前价格，从买单里面取出num个，从卖单里面取出num个
	//		头部卖单>当前价格，从买单里面取出num个，从卖单里面取出num-1个
	// 总结就是都取num个

	var topMap = make(map[string]interface{})
	for i:=0; i<num; i++ {
		order := book.GetHeadBuyOrder()
		if order != nil {
			p := order.Price.String()
			topMap[p] = order.Amount.String()
		}
	}
	for i:=0; i<num; i++ {
		order := book.GetHeadSellOrder()
		if order != nil {
			p := order.Price.String()
			topMap[p] = order.Amount.String()
		}
	}
	fprice := price.String()
	if _, ok := topMap[fprice]; !ok {
		topMap[fprice] = "0"
	}

	cache.SetTopN(symbol, num, topMap)
}

// 重新生成k线图
func handleKData() {
	defer func() { recover() }()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			common.Debugs("开始计算K线图")
			// 从交易记录中读取最近一分钟的数据


		case <-stopKDataChan:
			common.Debugs("K线图计算线程关闭")
			return
		}
	}
}