package model

import (
	"matching/utils/enum"
)

// OrderBook 交易委托账本，是一个二维链表
type OrderBook struct {
	buyOrderQueue orderQueue
	sellOrderQueue orderQueue
}

// Init 账本初始化
func (o *OrderBook) Init() {
	o.buyOrderQueue.init(enum.SortDesc)
	o.sellOrderQueue.init(enum.SortAsc)
}

// AddBuyOrder 增加买单
func (o *OrderBook) AddBuyOrder(order *Order) {
	o.buyOrderQueue.addOrder(order)
}

// AddSellOrder 增加卖单
func (o *OrderBook) AddSellOrder(order *Order) {
	o.sellOrderQueue.addOrder(order)
}

// GetHeadBuyOrder 获取头部买单
func (o *OrderBook) GetHeadBuyOrder() *Order {
	return o.buyOrderQueue.getHeadOrder()
}

// GetHeadSellOrder 获取头部卖单
func (o *OrderBook) GetHeadSellOrder() *Order {
	return o.sellOrderQueue.getHeadOrder()
}

// PopHeadBuyOrder 抛出头部买单
func (o *OrderBook) PopHeadBuyOrder() *Order {
	return o.buyOrderQueue.popHeadOrder()
}

// PopHeadSellOrder 抛出头部买单
func (o *OrderBook) PopHeadSellOrder() *Order {
	return o.sellOrderQueue.popHeadOrder()
}

// RemoveHeadBuyOrder 删除买单
func (o *OrderBook) RemoveHeadBuyOrder() {
	o.buyOrderQueue.popHeadOrder()
}

// RemoveHeadSellOrder 删除买单
func (o *OrderBook) RemoveHeadSellOrder() {
	o.sellOrderQueue.popHeadOrder()
}

// RemoveBuyOrder 删除指定买单
func (o *OrderBook) RemoveBuyOrder(order *Order) {
	o.buyOrderQueue.removeOrder(order)
}

// RemoveSellOrder 删除指定卖单
func (o *OrderBook) RemoveSellOrder(order *Order) {
	o.sellOrderQueue.removeOrder(order)
}