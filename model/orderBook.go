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

// PopHeadSellOrder 抛出头部卖单
func (o *OrderBook) PopHeadSellOrder() *Order {
	return o.sellOrderQueue.popHeadOrder()
}

// UpdateHeadBuyOrder 更新头部买单
func (o *OrderBook) UpdateHeadBuyOrder(order *Order) error {
	return o.buyOrderQueue.updateHeadOrder(order)
}

// UpdateHeadSellOrder 更新头部卖单
func (o *OrderBook) UpdateHeadSellOrder(order *Order) error {
	return o.sellOrderQueue.updateHeadOrder(order)
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
func (o *OrderBook) RemoveBuyOrder(order *Order) bool {
	return o.buyOrderQueue.removeOrder(order)
}

// RemoveSellOrder 删除指定卖单
func (o *OrderBook) RemoveSellOrder(order *Order) bool {
	return o.sellOrderQueue.removeOrder(order)
}

func (o *OrderBook) ShowAllBuyOrder() {
	o.buyOrderQueue.showAllOrder()
}

func (o *OrderBook) ShowAllSellOrder() {
	o.sellOrderQueue.showAllOrder()
}