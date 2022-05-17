package mq

import (
	mq "matching/utils/redis"
)

// SendCancelResult 发送撤单结果消息队列
func SendCancelResult(symbol, orderId string, ok bool) {
	mq.SendCancelResult(symbol, orderId, ok)
}

// SendTradeResult 发送撮合结果消息队列
func SendTradeResult(symbol string, trade map[string]interface{}) {
	mq.SendTrade(symbol, trade)
}

// GetCancelResult 读取撤单结果消息队列
func GetCancelResult(symbol, orderId string) map[string]string {
	return mq.GetCancelResult(symbol)
}

// GetTradeResult 读取撮合结果消息队列
func GetTradeResult(symbol, orderId string) map[string]string {
	return mq.GetTradeResult(symbol)
}


