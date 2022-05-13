package model

import "matching/utils/enum"

// 交易委托账本，是一个二维链表
type OrderBook struct {
	buyOrderQueue orderQueue
	sellOrderQueue orderQueue
}

// 账本初始化
func (o *OrderBook) Init() {
	o.buyOrderQueue.init(enum.SortDesc)
	o.sellOrderQueue.init(enum.SortAsc)
}

// 增加买单
func (o *OrderBook) AddBuyOrder(order *Order) {

}

// 增加卖单
func (o *OrderBook) AddSellOrder(order *Order) {

}

// 获取头部买单
func (o *OrderBook) GetHeadBuyOrder() {

}

// 获取头部卖单
func (o *OrderBook) GetHeadSellOrder() *Order {
	return o.sellOrderQueue.getHeadOrder()
}

// 抛出头部买单
func (o *OrderBook) PopHeadBuyOrder() {

}

// 抛出头部买单
func (o *OrderBook) PopHeadSellOrder() {

}

// 删除买单
func (o *OrderBook) RemoveHeadBuyOrder(order *Order) {

}

// 删除买单
func (o *OrderBook) RemoveHeadSellOrder(order *Order) {

}