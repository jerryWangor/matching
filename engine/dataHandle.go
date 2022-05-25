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

	//fmt.Println("分割线")
	// topN从elementMap来
	data := make(map[string]interface{})
	fprice, _ := price.Float64()
	buyList := book.GetBuyTopN(fprice, num)
	//fmt.Println("买单数量", buyList.Len())
	if buyList.Len() >0 {
		for e := buyList.Front(); e != nil; e = e.Next() {
			topData := e.Value.(model.PriceTopN)
			sprice := strconv.FormatFloat(topData.Price, 'E', -1, 64)
			data[sprice] = topData.Amount
		}
	}
	// 判断当前价格是否在data中
	nowPrice := price.String()
	if _, ok := data[nowPrice]; !ok {
		data[nowPrice] = 0.0
	}
	sellList := book.GetSellTopN(fprice, num)
	if sellList.Len() >0 {
		for e := sellList.Front(); e != nil; e = e.Next() {
			topData := e.Value.(model.PriceTopN)
			sprice := strconv.FormatFloat(topData.Price, 'E', -1, 64)
			data[sprice] = topData.Amount
		}
	}

	//fmt.Println("topdata", data)
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
			currentSecond := time.Now().Second()
			timestamp := time.Now().Unix() - int64(currentSecond)
			// 取出当前的实时价格
			kData := model.KData{
				TopPrice: kDataPrice.TopPrice,
				BottomPrice: kDataPrice.BottomPrice,
				NowPrice: kDataPrice.NowPrice,
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