package handler

import (
	"github.com/gin-gonic/gin"
	"matching/model"
	"matching/utils/cache"
	"matching/utils/common"
)

func ShowLogs(c *gin.Context) {
	// 该函数为调试函数，主要是打印引擎的相关情况

	// 缓存相关
	common.Debugs("----------Start 打印缓存相关数据----------")
	symbols := cache.GetSymbols()
	for _, v := range symbols {
		price := cache.GetPrice(v)
		common.Debugs("交易标："+ v + "，价格：" + price.String() + ",包含订单如下：")
		// 获取该交易标缓存的所有订单
		orderIds := cache.GetOrderIdsWithSymbol(v)
		for _, orderId := range orderIds {
			mapOrder := cache.GetOrder(v, orderId)
			common.Debugs(formatOrderString(mapOrder))
		}


	}
	common.Debugs("----------End 打印缓存相关数据----------")

	// 消息队列相关
	common.Debugs("----------Start 打印消息队列相关数据----------")
	common.Debugs("----------End 打印消息队列相关数据----------")

	// 交易委托账本
	common.Debugs("----------Start 打印交易委托账本相关数据----------")
	common.Debugs("----------End 打印交易委托账本相关数据----------")

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