package engine

import (
	"github.com/shopspring/decimal"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/common"
	"strconv"
	"time"
)

// 重新生成topN，保存redis
func handleTopN(symbol string, price *decimal.Decimal, book *model.OrderBook, num int) {
	// topN从elementMap来
	data := make(map[string]interface{})
	fprice, _ := price.Float64()
	buyList := book.GetBuyTopN(fprice, num)
	if buyList.Len() >0 {
		for e := buyList.Front(); e != nil; e = e.Next() {
			topData := e.Value.(*model.PriceTopN)
			sprice := strconv.FormatFloat(topData.Price, 'E', -1, 64)
			data[sprice] = topData.Amount
		}
	}
	// 判断当前价格是否在data中
	nowPrice := price.String()
	if _, err := data[nowPrice]; err {
		data[nowPrice] = 0.0
	}
	sellList := book.GetBuyTopN(fprice, num)
	if sellList.Len() >0 {
		for e := sellList.Front(); e != nil; e = e.Next() {
			topData := e.Value.(*model.PriceTopN)
			sprice := strconv.FormatFloat(topData.Price, 'E', -1, 64)
			data[sprice] = topData.Amount
		}
	}

	cache.SetTopN(symbol, num, data)
}

// 生成k线图，1分钟生成一次
func handleKData(symbol string, kDataPrice *model.KDataPrice) {
	defer func() { recover() }()
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			common.Debugs("开始计算K线图")
			topPrice, _ := kDataPrice.TopPrice.Float64()
			bottomPrice, _ := kDataPrice.BottomPrice.Float64()
			nowPrice, _ := kDataPrice.NowPrice.Float64()

			currentSecond := time.Now().Second()
			time := time.Now().UnixMicro() - int64(currentSecond)
			// 转成float64
			timestamp := float64(time)
			// 取出当前的实时价格
			kData := model.KData{
				TopPrice: topPrice,
				BottomPrice: bottomPrice,
				NowPrice: nowPrice,
				Timestamp: timestamp,
			}
			kDataJson := common.ToJson(kData)
			cache.SetKData(symbol, timestamp, kDataJson)

		case <-StopKDataChan:
			common.Debugs("K线图计算线程关闭")
			return
		}
	}
}