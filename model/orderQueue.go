package model

import (
	"container/list"
	"matching/enum"
)

/**
	sortBy 指定价格排序的方向，买单队列是降序的，而卖单队列则是升序的。
	parentList 保存整个二维链表的所有订单，第一维以价格排序，第二维以时间排序。
	elementMap 则是 Key 为价格、Value 为第二维订单链表的键值对。
 */
type orderQueue struct {
	sortBy     enum.SortDirection
	parentList *list.List
	elementMap map[string]*list.Element
}

func (q *orderQueue) init(sortBy enum.SortDirection) {
	q.sortBy = sortBy
	q.parentList = list.New()
	q.elementMap = make(map[string]*list.Element)
}

func (q *orderQueue) addOrder(order Order) {

}

func (q *orderQueue) getHeadOrder() {

}

func (q *orderQueue) popHeadOrder() {

}

func (q *orderQueue) removeOrder(order Order) {

}


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
			break;
		}
	}
	o := p.Value.(*list.List).Front().Value.(*Order)
	return o.Price.String(), i
}

