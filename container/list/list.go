package list

type Node struct {
	Value      interface{}
	prev, next *Node
}

func NewNode(v interface{}) *Node {
	return &Node{ //摒弃c++用new的方式
		Value: v,
		prev:  nil,
		next:  nil,
	}
}
func (n *Node) Next() *Node {
	return n.next
}

type List struct {
	len  int //the len should be read-only
	head *Node
	tail *Node
}

func New() *List { return new(List) }

func (l *List) Len() int {
	return l.len
}

/**
 * Get the first node of the list
 */
func (l *List) Front() *Node {
	return l.head
}

/**
 * Get the last node of the list
 */
func (l *List) Back() *Node {
	return l.tail
}

/**
 * Append node to the tail of the list
 */
func (l *List) PushBack(v interface{}) {
	if l.len == 0 {
		l.Init(v)
	} else {
		n := NewNode(v)
		l.tail.next = n
		n.prev = l.tail
		l.tail = n
		l.len++
	}
}

/**
 * Notice:
 * It is mainly used for helping the PushBack() be simpler.
 * Of course you can also used it for other purpose.
 */
func (l *List) Init(v interface{}) {
	l.head = NewNode(v)
	l.tail = l.head
}

/**
 * Insert a new node after n.
 */
func (l *List) Insert(v interface{}, n *Node) {
	newN := NewNode(v)
	if n.next != nil { //if n != l.tail
		n.next.prev = newN
		newN.next = n.next
	} else {
		l.tail = newN
	}
	n.next = newN
	newN.prev = n
	l.len++
}

/**
 * Append node to the head of the list
 */
func (l *List) PushFront(v interface{}) {
	newN := NewNode(v)
	l.head.prev = newN
	newN.next = l.head
	l.head = newN
	l.len++
}

/**
 * Remove the node of the list by using the pointer of list nod
 */
func (l *List) Remove(n *Node) {
	l.len--
	if l.len == 0 {
		l.head = nil
		l.tail = nil
		return
	}
	if n == l.head {
		n.next.prev = nil
		l.head = n.next
		n.next = nil //It is necessary to ensure the new head can be recycled
		return
	} else if n == l.tail {
		n.prev.next = nil
		l.tail = n.prev
		n.prev = nil
		return
	} else {
		n.prev.next = nil
		n.prev = nil
		n.next.prev = nil
		n.next = nil
		return
	}
}
