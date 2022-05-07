package utils

import "fmt"

type Object interface {}

type Node struct {
	Data Object // 定义数据域
	Next *Node // 定义地址域(指向下一个表的地址)
}

type SingleLinkList struct {
	headNode *Node // 头节点
}

// 判断是否为空的单链表
func (this *SingleLinkList) IsEmpty() bool {
	if this.headNode == nil {
		return true
	} else {
		return false
	}
}

// 获取链表的长度
func (this *SingleLinkList) Length() int {
	// 获取链表头结点
	cur := this.headNode
	// 定义一个计数器，初始值为0
	count := 0
	for cur != nil {
		// 如果头结点不为空，则count++
		count++
		cur = cur.Next
	}
	return count
}

// 从链表头部添加数据
func (this *SingleLinkList) Add(data Object) *Node {
	node := &Node{Data: data} // 创建一个数据
	node.Next = this.headNode // 把之前的头部数据作为next
	this.headNode = node
	return node
}

// 从链表尾部添加数据
func (this *SingleLinkList) Append(data Object) {
	// 创建一个新元素，通过传入data参数进行数据域的赋值
	node := &Node{Data: data}
	if this.IsEmpty() { // 如果该链表为空，直接将新元素作为头结点
		this.headNode = node
	} else {
		cur := this.headNode // 储存头结点
		for cur.Next != nil { // 判断是否尾节点，如果是nil则是尾节点
			cur = cur.Next // 链表进行位移，直到cur获取到尾节点
		}
		cur.Next = node // 此时cur为尾节点，将其地址指向新创建的节点
	}
}

// 在指定位置添加元素，index是下标，this.headNode的index=0
func (this *SingleLinkList) Insert(index int, data Object) {
	if(index < 0) { // 小于0就是头部插入
		this.Add(data)
	} else if index > this.Length() { // 尾部插入
		this.Append(data)
	} else {
		pre := this.headNode
		count := 0
		for count < (index - 1) { // 用于控制位移的链表数目
			pre = pre.Next
			count++
		}
		// 当循环结束后，pre指向index-1的位置
		node := &Node{Data: data}
		node.Next = pre.Next
		pre.Next = node
	}
}

// 删除链表指定值的元素
func (this *SingleLinkList) Remove(data Object) {
	pre := this.headNode // 定义pre变量存储头结点
	if pre.Data == data {
		this.headNode = pre.Next // 如果该节点数据是要删除的，就指定第二个节点是新的头结点
	} else {
		for pre.Next != nil { // 一直遍历到最后一个节点
			if pre.Next.Data == data { // 如果pre.Next的节点数据等于data，就删除该节点，删除之后记得把该节点之后的节点地址赋值给pre.Next
				pre.Next = pre.Next.Next
			} else { // 如果pre.Next的节点的数据等于data，那么进行节点位移，继续循环
				pre = pre.Next
			}
		}
	}
}

// 删除指定位置的元素
func (this *SingleLinkList) RemoveAtIndex(index int) {
	pre := this.headNode
	if index <= 0 { // index小于等于0就说明删除头结点
		this.headNode = pre.Next // 第二个节点作为头结点
	} else if index > this.Length() {
		fmt.Println("超出链表长度")
		return
	} else {
		count := 0
		for count != (index-1) && pre.Next != nil { // 开始遍历，如果index=1，则直接删除跳出循环，如果index>1，则进行链表位移，然后删除
			count++
			pre = pre.Next
		}
		pre.Next = pre.Next.Next
	}
}

// 查看链表是否包含某个元素
func (this *SingleLinkList) Contain(data Object) bool {
	cur := this.headNode
	for cur != nil {
		if cur.Data == data {
			return true
		}
		cur = cur.Next
	}
	return false
}

// 遍历所有节点
func (this *SingleLinkList) ShowList() {
	if !this.IsEmpty() {
		cur := this.headNode
		for {
			fmt.Printf("\t%v", cur.Data)
			if cur.Next != nil {
				cur = cur.Next
			} else {
				break
			}
		}
	}
}

// 测试
//func main() {
//	list := SingleLinkList{}
//	// 添加数据
//	list.Append(1)
//	list.Append(2)
//	list.Append("name")
//	// 头部加元素
//	list.Add(0)
//	//
//	fmt.Println("是否存在某个元素：", list.Contain(2))
//	list.Insert(3, 3)
//	list.Remove("name")
//	list.RemoveAtIndex(2)
//
//	fmt.Println("链表长度：", list.Length())
//	fmt.Println("链表数据：")
//	list.ShowList()
//}