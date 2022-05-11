package model

// 交易委托账本，是一个二维链表
type OrderBook struct {
	buyOrderQueue orderQueue
	sellOrderQueue orderQueue
}

// 账本初始化
func (o *OrderBook) Init() {

}

// 增加买单
func AddBuyOrder(order *Order) {

}

// 增加卖单
func AddSellOrder(order *Order) {

}

// 获取头部买单
func GetHeadBuyOrder() {

}

// 获取头部卖单
func GetHeadSellOrder() {

}

// 抛出头部买单
func PopHeadBuyOrder() {

}

// 抛出头部买单
func PopHeadSellOrder() {

}

// 删除买单
func RemoveHeadBuyOrder(order *Order) {

}

// 删除买单
func RemoveHeadSellOrder(order *Order) {

}