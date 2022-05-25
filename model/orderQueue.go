package model

import (
	"container/list"
	"matching/utils/common"
	"matching/utils/enum"
	"sort"
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
	elementMap 则是 Key 为价格、Value 二维map订单。
 */
type orderQueue struct {
	sortBy     enum.SortDirection
	parentList *list.List
	elementMap map[float64]map[string]*Order // 主要是用来查询top5，用price作为价格更方便
}

// 初始化函数
func (q *orderQueue) init(sortBy enum.SortDirection) {
	q.sortBy = sortBy
	q.parentList = list.New()
	q.elementMap = make(map[float64]map[string]*Order)
}

// 把订单插入到链表中
func (q *orderQueue) addOrder(order *Order) {

	// 先插入价格map
	price, _ := order.Price.Float64()
	if _, ok := q.elementMap[price]; !ok {
		q.elementMap[price] = make(map[string]*Order)
	}
	q.elementMap[price][order.OrderId] = order

	// 如果队列长度是0，就直接放到第一个
	if q.parentList.Len() == 0 {
		q.parentList.PushFront(order)
		return
	}

	var eKey *list.Element
	for e := q.parentList.Front(); e != nil; e = e.Next() {
		price := e.Value.(*Order).Price
		// 取出链表订单，判断当前订单价格是否>=链表订单价格，如果成立，则记录为当前插入eKey，一直循环直到遇到比当前价格小的订单，就插入eKey的后面
		// 卖单如果小于头部卖单就直接放到前面
		if (order.Side == enum.SideBuy && order.Price.GreaterThan(price)) || (order.Side == enum.SideSell && order.Price.LessThan(price)) {
			q.parentList.InsertBefore(order, e)
			// 注意要设置eKey为nil
			eKey = nil
			break
		} else {
			// 卖单如果等于链表订单就设置eKey
			if order.Price.Equals(price) {
				eKey = e
			} else {
				// 如果当前订单价格>链表订单价格，就一直设置eKey，直到遇到比他大的或者最后一个值
				eKey = e
			}
		}
	}

	// 插入到指定order后面
	if eKey != nil {
		q.parentList.InsertAfter(order, eKey)
	}
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

// 更新头部订单
func (q *orderQueue) updateHeadOrder(order *Order) error {
	// 先删除头部订单，在插入
	popOrder := q.popHeadOrder()
	if popOrder == nil {
		return common.Errors("头部订单删除失败")
	}
	nowOrder := q.parentList.PushFront(order)
	if nowOrder == nil {
		return common.Errors("头部订单更新失败")
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

// 展示所有订单
func (q *orderQueue) showAllOrder() {
	for e := q.parentList.Front(); e != nil; e = e.Next() {
		order := e.Value.(*Order)
		common.Debugs(common.ToJson(order))
	}
}

// 更新element账本中的订单
func (q *orderQueue) updateElementOrder(order *Order) bool {
	price, _ := order.Price.Float64()
	orderId := order.OrderId
	// 判断是否在map中
	_, err := q.elementMap[price][orderId]
	if !err {
		return false
	}
	q.elementMap[price][orderId] = order

	return true
}

// 删除element账本中的订单
func (q *orderQueue) removeElementOrder(order *Order) bool {
	price, _ := order.Price.Float64()
	orderId := order.OrderId
	// 判断是否在map中
	if _, ok := q.elementMap[price][orderId]; !ok {
		return false
	}

	delete(q.elementMap[price], orderId)
	// 如果没有了就置空
	if len(q.elementMap[price]) == 0 {
		q.elementMap[price] = nil
	}

	return true
}

// 获取topN的价格和数量
func (q *orderQueue) getTopN(nowPrice float64, num int) *list.List {
	var keys []float64
	// 用有序list保存返回数据
	topMap := list.New()
	// 把elementMap的key进行排序
	for k := range q.elementMap {
		keys = append(keys, k)
	}
	// 卖升买降
	if q.sortBy == enum.SortDesc {
		sort.Sort(sort.Reverse(sort.Float64Slice(keys))) // 降序
	} else {
		sort.Float64s(keys) // 升序
	}
	// 循环keys，找到当前价格的N档
	//fmt.Println("keys", keys, nowPrice)
	for _, v := range keys {
		num--
		if num < 0 {
			break
		}
		if q.sortBy == enum.SortDesc {
			// 买单判断小于
			if v <= nowPrice {
				// 循环所有订单，把数量相加
				amount := 0.0
				for _, val := range q.elementMap[v] {
					a, _ := val.Amount.Float64()
					amount += a
				}
				sTopN := PriceTopN{
					Price: v,
					Amount: amount,
				}
				topMap.PushBack(sTopN)
			}
		} else {
			// 卖单判断大于
			if v >= nowPrice {
				// 循环所有订单，把数量相加
				amount := 0.0
				for _, val := range q.elementMap[v] {
					a, _ := val.Amount.Float64()
					amount += a
				}
				sTopN := PriceTopN{
					Price: v,
					Amount: amount,
				}
				//fmt.Println("卖单数据：", sTopN.Price, sTopN.Amount)
				topMap.PushBack(sTopN)
			}
		}
	}
	return topMap
}

// 获取top价格
func (q *orderQueue) getElementMap() map[float64]map[string]*Order {
	return q.elementMap
}

// 读取深度价格是为了方便处理 market-opponent、market-top5、market-top10 等类型的订单时判断上限价格。
//func (q *orderQueue) getDepthPrice(depth int) (string, int) {
//	if q.parentList.Len() == 0 {
//		return "", 0
//	}
//	p := q.parentList.Front()
//	i := 1
//	for ; i < depth; i++ {
//		t := p.Next()
//		if t != nil {
//			p = t
//		} else {
//			break
//		}
//	}
//	o := p.Value.(*Order)
//	return o.Price.String(), i
//}

