package model

import (
	"container/list"
	"matching/utils/enum"
)

/**
	type Element
	type Element struct {
			Value interface{}   //在元素中存储的值
	}
	func (e *Element) Next() *Element  //返回该元素的下一个元素，如果没有下一个元素则返回nil
	func (e *Element) Prev() *Element//返回该元素的前一个元素，如果没有前一个元素则返回nil。
	type List
	func New() *List //返回一个初始化的list
	func (l *List) Back() *Element //获取list l的最后一个元素
	func (l *List) Front() *Element //获取list l的第一个元素
	func (l *List) Init() *List  //list l初始化或者清除list l
	func (l *List) InsertAfter(v interface{}, mark *Element) *Element  //在list l中元素mark之后插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
	func (l *List) InsertBefore(v interface{}, mark *Element) *Element//在list l中元素mark之前插入一个值为v的元素，并返回该元素，如果mark不是list中元素，则list不改变。
	func (l *List) Len() int //获取list l的长度
	func (l *List) MoveAfter(e, mark *Element)  //将元素e移动到元素mark之后，如果元素e或者mark不属于list l，或者e==mark，则list l不改变。
	func (l *List) MoveBefore(e, mark *Element)//将元素e移动到元素mark之前，如果元素e或者mark不属于list l，或者e==mark，则list l不改变。
	func (l *List) MoveToBack(e *Element)//将元素e移动到list l的末尾，如果e不属于list l，则list不改变。
	func (l *List) MoveToFront(e *Element)//将元素e移动到list l的首部，如果e不属于list l，则list不改变。
	func (l *List) PushBack(v interface{}) *Element//在list l的末尾插入值为v的元素，并返回该元素。
	func (l *List) PushBackList(other *List)//在list l的尾部插入另外一个list，其中l和other可以相等。
	func (l *List) PushFront(v interface{}) *Element//在list l的首部插入值为v的元素，并返回该元素。
	func (l *List) PushFrontList(other *List)//在list l的首部插入另外一个list，其中l和other可以相等。
	func (l *List) Remove(e *Element) interface{}//如果元素e属于list l，将其从list中删除，并返回元素e的值。
 */

/**
	sortBy 指定价格排序的方向，买单队列是降序的，而卖单队列则是升序的。
	parentList 保存整个二维链表的所有订单，第一维以价格排序，第二维以时间排序。
	elementMap 则是 Key 为价格、Value 为第二维订单链表的键值对。
 */
type orderQueue struct {
	sortBy     enum.SortDirection
	parentList *list.List
	elementMap map[string]*list.List // 主要是用来查询top5，用price作为价格更方便
}

// 初始化函数
func (q *orderQueue) init(sortBy enum.SortDirection) {
	q.sortBy = sortBy
	q.parentList = list.New()
	q.elementMap = make(map[string]*list.List)
}

// 把订单插入到链表中
func (q *orderQueue) addOrder(order *Order) {

	// 如果队列长度是0，就直接放到第一个
	if q.parentList.Len() == 0 {
		q.parentList.PushFront(order)
		return
	}

	// 买单队列是按照价格降序的，当前价格的订单<=那条订单，就插入到那条订单的后面
	if q.sortBy == enum.SortDesc {
		for e := q.parentList.Back(); e != nil; e = e.Prev() {
			price := e.Value.(*Order).Price
			if order.Price.LessThanOrEqual(price) {
				q.parentList.InsertAfter(order, e)
			} else {
				q.parentList.InsertBefore(order, e)
			}
		}
	} else {
		// 卖单队列是按照价格升序的，找到<=当前价格的订单，排到前面
		for e := q.parentList.Back(); e != nil; e = e.Prev() {
			price := e.Value.(*Order).Price
			if order.Price.GreaterThanOrEqual(price) {
				q.parentList.InsertAfter(order, e)
			} else {
				q.parentList.InsertBefore(order, e)
			}
		}
	}

	// 插入价格map
	price := order.Price.String()
	if _, ok := q.elementMap[price]; !ok {
		q.elementMap[price] = list.New()
	}
	q.elementMap[price].PushBack(order)
}

// 从委托账本中查询头部订单
func (q *orderQueue) getHeadOrder() *Order {
	front := q.parentList.Front()
	if front != nil {
		// 这里先强制转换，后面改成断言
		return q.parentList.Front().Value.(*Order)
	}
	return nil
}

// 删除头部订单
func (q *orderQueue) popHeadOrder() *Order {
	front := q.parentList.Front()
	if front != nil {
		// 从交易委托账本中删除
		return q.parentList.Remove(front).(*Order)
	}
	return nil
}

// 删除指定订单
func (q *orderQueue) removeOrder(order *Order) bool {
	// 循环链表，找到订单并删除
	for e := q.parentList.Front(); e != nil; e = e.Next() {
		orderId := e.Value.(*Order).OrderId
		if orderId == order.OrderId {
			q.parentList.Remove(e)
			return true
		}
	}
	return false
}

// 读取深度价格是为了方便处理 market-opponent、market-top5、market-top10 等类型的订单时判断上限价格。
func (q *orderQueue) getDepthPrice(depth int) (string, int) {
	if q.parentList.Len() == 0 {
		return "", 0
	}
	p := q.parentList.Front()
	i := 1
	for ; i < depth; i++ {
		t := p.Next()
		if t != nil {
			p = t
		} else {
			break
		}
	}
	o := p.Value.(*list.List).Front().Value.(*Order)
	return o.Price.String(), i
}

