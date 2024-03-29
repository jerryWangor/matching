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

	// topN从elementMap来
	data := make(map[string]interface{})
	fprice, _ := price.Float64()
	//fmt.Println("分割线", fprice)
	buyData := make([]model.PriceTopN, 0)
	buyList := book.GetBuyTopN(fprice, num)
	if buyList.Len() >0 {
		for e := buyList.Front(); e != nil; e = e.Next() {
			topData := e.Value.(model.PriceTopN)
			buyData = append(buyData, topData)
		}
	}
	data["buy"] = buyData

	sellData := make([]model.PriceTopN, 0)
	sellList := book.GetSellTopN(fprice, num)
	if sellList.Len() >0 {
		for e := sellList.Back(); e != nil; e = e.Prev() {
			topData := e.Value.(model.PriceTopN)
			sellData = append(sellData, topData)
		}
	}
	data["sell"] = sellData

	data["now"] = KDataPriceMap[symbol].NowPrice.String()

	//fmt.Println("topdata", data)
	cache.SetTopN(symbol, num, data)
}

// 生成k线图，1分钟生成一次
func handleKData(data *map[string]*model.KDataPrice) {
	defer func() { recover() }()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			currentSecond := time.Now().Second()
			timestamp := time.Now().Unix() - int64(currentSecond)

			for k, v := range *data {
				common.Debugs("开始计算K线图：" + k)

				// 取出当前的实时价格
				kData := model.KData{
					TopPrice: v.TopPrice,
					BottomPrice: v.BottomPrice,
					NowPrice: v.NowPrice,
					Timestamp: timestamp,
				}
				kDataJson := common.ToJson(kData)
				cache.SetKData(k, timestamp, kDataJson)
			}

		case <-StopKDataChan:
			common.Debugs("K线图计算线程关闭")
			return
		}
	}
}