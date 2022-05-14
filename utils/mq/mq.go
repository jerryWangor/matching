package mq

import (
	mq "matching/utils/redis"
)

func SendCancelResult(symbol, orderId string, ok bool) {
	mq.SendCancelResult(symbol, orderId, ok)
}

func SendTrade(symbol string, trade map[string]interface{}) {
	mq.SendTrade(symbol, trade)
}
