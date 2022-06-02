package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"matching/engine"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/common"
	"matching/utils/mq"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ShowLogs(c *gin.Context) {
	// 该函数为调试函数，主要是打印引擎的相关情况

	//time.Sleep(3 * time.Second)
	//_, err := http.PostForm("http://127.0.0.1:8080/order/handleOrder",
	//	url.Values{"key": {"Value"}, "id": {"123"}})
	//if err != nil {
	//	fmt.Println(err)
	//}

	//_, err := http.Post("http://127.0.0.1:8080/order/handleOrder",
	//	"application/json",
	//	strings.NewReader("time=test&sign=ab123123"))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//return


	// 缓存相关
	common.Debugs("----------Start 打印缓存相关数据----------")
	symbols := cache.GetSymbols()

	for _, v := range symbols {
		price := cache.GetPrice(v)
		common.Debugs("交易标："+ v + "，价格：" + price.String() + ",包含订单如下：")
		// 获取该交易标缓存的所有订单
		orders := cache.GetOrderIdsWithSymbol(v)
		for _, val := range orders {
			orderArr := strings.Split(val,":")
			mapOrder := cache.GetOrder(v, orderArr[0], orderArr[1])
			common.Debugs(formatOrderString(mapOrder))
		}
	}
	common.Debugs("----------End 打印缓存相关数据----------")

	// 消息队列相关
	common.Debugs("----------Start 打印消息队列相关数据----------")
	common.Debugs("撤单结果消息队列：")
	for _, v := range symbols {
		strMap := mq.GetCancelResult(v)
		for id, str := range strMap {
			common.Debugs("stream id：" + id + "，value：" + str)
		}
	}
	common.Debugs("撮合结果消息队列：")
	for _, v := range symbols {
		strMap := mq.GetTradeResult(v)
		for id, str := range strMap {
			common.Debugs("stream id：" + id + "，value：" + str)
		}
	}
	common.Debugs("----------End 打印消息队列相关数据----------")

	// 交易委托账本
	common.Debugs("----------Start 打印交易委托账本相关数据----------")
	for k, v := range engine.AllOrderBookMap {
		common.Debugs(fmt.Sprintf("交易标：%s，账本如下：", k))

		common.Debugs("卖单：")
		v.ShowAllSellOrder()

		common.Debugs("买单：")
		// 循环一直到nil
		v.ShowAllBuyOrder()

		common.Debugs("卖单TopN数据：")
		sellElementMap := v.GetSellElementMap()
		if len(sellElementMap) >0 {
			var keys []float64
			for time, _ := range sellElementMap {
				keys = append(keys, time)
			}
			sort.Sort(sort.Reverse(sort.Float64Slice(keys))) // 降序
			for _, m := range keys {
				common.Debugs("价格：" + strconv.FormatFloat(m, 'E', -1, 64))
				for oid, order := range sellElementMap[m] {
					common.Debugs("订单号：" + oid + "，数量：" + order.Amount.String())
				}
			}
		}



		common.Debugs("-------------------------")
		common.Debugs("买单TopN数据：")
		buyElementMap := v.GetBuyElementMap()
		if len(buyElementMap) >0 {
			var keys []float64
			for time, _ := range buyElementMap {
				keys = append(keys, time)
			}
			sort.Sort(sort.Reverse(sort.Float64Slice(keys))) // 降序
			for _, m := range keys {
				common.Debugs("价格：" + strconv.FormatFloat(m, 'E', -1, 64))
				for oid, order := range buyElementMap[m] {
					common.Debugs("订单号：" + oid + "，数量：" + order.Amount.String())
				}
			}
		}


	}
	common.Debugs("----------End 打印交易委托账本相关数据----------")

	// K线图
	common.Debugs("----------Start TopN  K线图----------")
	for _, v := range symbols {
		common.Debugs("Top100：")
		topMap := cache.GetTopN(v, 100)
		common.Debugs(common.ToJson(topMap))

		// 取最近10分钟的k线图
		common.Debugs("K线图：")
		currentSecond := time.Now().Second()
		time2 := time.Now().Unix() - int64(currentSecond)
		time1 := time2 - 600
		kData := cache.GetKData(v, time1, time2)
		for _, v := range kData {
			common.Debugs("最高价：" + v.TopPrice.String() + "，最低价：" + v.BottomPrice.String() + "，当前价：" + v.NowPrice.String() + "，时间：" + common.TimeStampToString(v.Timestamp))
		}
	}
	common.Debugs("----------End TopN  K线图----------")

	//common.Debugs("----------Start 订单数量存量----------")
	//common.Debugs("买单存量：" + engine.AllOrderAmountMap[enum.SideBuy.String()].String())
	//common.Debugs("卖单存量：" + engine.AllOrderAmountMap[enum.SideSell.String()].String())
	//common.Debugs("----------End 订单数量存量----------")



}

// 格式化订单输出
func formatOrderString(order model.Order) string {
	return  "交易标：" 	+ order.Symbol + "，" +
			"订单ID：" 	+ order.OrderId + "，" +
			"下单类型：" 	+ order.Action.String() + "，" +
			"竞价类型：" 	+ order.Type.String() + "，" +
			"买/卖：" 	+ order.Side.String() + "，" +
			"数量：" 	+ order.Amount.String() + "，" +
			"价格：" 	+ order.Price.String() + "，" +
			"时间：" 	+ common.TimeStampToString(order.Timestamp)
}