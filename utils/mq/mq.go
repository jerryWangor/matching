package mq

import (
	mq "matching/utils/redis"
)

// 队列操作
/**
队列操作
其中，matching:cancelresults:{symbol} 就是撤单结果的 MQ 所属的 Key，
matching:trades:{symbol} 则是成交记录的 MQ 所属的 Key。
可以看到，我们还根据不同 symbol 分不同 MQ，这样还方便下游服务可以根据需要实现分布式订阅不同 symbol 的 MQ。
*/
func SendCancelResult(symbol, orderId string, ok bool) {
	mq.SendCancelResult(symbol, orderId, ok)
}

func SendTrade(symbol string, trade map[string]interface{}) {
	mq.SendTrade(symbol, trade)
}
