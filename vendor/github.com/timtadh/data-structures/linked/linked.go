package linked

import (
	"fmt"
	"strings"
)

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/types"
)

// A doubly linked list node.
type Node struct {
	Data       types.Hashable
	Next, Prev *Node
}

// Compares the Data of the node to the passed element.
func (n *Node) Equals(b types.Equatable) bool {
	switch x := b.(type) {
	case *Node:
		return n.Data.Equals(x.Data)
	default:
		return n.Data.Equals(b)
	}
}

// Compares the Data of the node to the passed element.
func (n *Node) Less(b types.Sortable) bool {
	switch x := b.(type) {
	case *Node:
		return n.Data.Less(x.Data)
	default:
		return n.Data.Less(b)
	}
}

// Hashes the Data of the node to the passed element.
func (n *Node) Hash() int {
	return n.Data.Hash()
}

// A doubly linked list. There is no synchronization.
// The fields are publically accessible to allow for easy customization.
type LinkedList struct {
	Length int
	Head   *Node
	Tail   *Node
}

func New() *LinkedList {
	return &LinkedList{
		Length: 0,
		Head:   nil,
		Tail:   nil,
	}
}

func (l *LinkedList) Size() int {
	return l.Length
}

func (l *LinkedList) Items() (it types.KIterator) {
	cur := l.Head
	it = func() (item types.Hashable, _ types.KIterator) {
		if cur == nil {
			return nil, nil
		}
		item = cur.Data
		cur = cur.Next
		return item, it
	}
	return it
}

func (l *LinkedList) Backwards() (it types.KIterator) {
	cur := l.Tail
	it = func() (item types.Hashable, _ types.KIterator) {
		if cur == nil {
			return nil, nil
		}
		item = cur.Data
		cur = cur.Prev
		return item, it
	}
	return it
}

func (l *LinkedList) Has(item types.Hashable) bool {
	for x, next := l.Items()(); next != nil; x, next = next() {
		if x.Equals(item) {
			return true
		}
	}
	return false
}

func (l *LinkedList) Push(item types.Hashable) (err error) {
	return l.EnqueBack(item)
}

func (l *LinkedList) Pop() (item types.Hashable, err error) {
	return l.DequeBack()
}

func (l *LinkedList) EnqueFront(item types.Hashable) (err error) {
	n := &Node{Data: item, Next: l.Head}
	if l.Head != nil {
		l.Head.Prev = n
	} else {
		l.Tail = n
	}
	l.Head = n
	l.Length++
	return nil
}

func (l *LinkedList) EnqueBack(item types.Hashable) (err error) {
	n := &Node{Data: item, Prev: l.Tail}
	if l.Tail != nil {
		l.Tail.Next = n
	} else {
		l.Head = n
	}
	l.Tail = n
	l.Length++
	return nil
}

func (l *LinkedList) DequeFront() (item types.Hashable, err error) {
	if l.Head == nil {
		return nil, errors.Errorf("List is empty")
	}
	item = l.Head.Data
	l.Head = l.Head.Next
	if l.Head != nil {
		l.Head.Prev = nil
	} else {
		l.Tail = nil
	}
	l.Length--
	return item, nil
}

func (l *LinkedList) DequeBack() (item types.Hashable, err error) {
	if l.Tail == nil {
		return nil, errors.Errorf("List is empty")
	}
	item = l.Tail.Data
	l.Tail = l.Tail.Prev
	if l.Tail != nil {
		l.Tail.Next = nil
	} else {
		l.Head = nil
	}
	l.Length--
	return item, nil
}

func (l *LinkedList) First() (item types.Hashable) {
	if l.Head == nil {
		return nil
	}
	return l.Head.Data
}

func (l *LinkedList) Last() (item types.Hashable) {
	if l.Tail == nil {
		return nil
	}
	return l.Tail.Data
}

// Can be compared to any types.IterableContainer
func (l *LinkedList) Equals(b types.Equatable) bool {
	if o, ok := b.(types.IterableContainer); ok {
		return list.Equals(l, o)
	} else {
		return false
	}
}

// Can be compared to any types.IterableContainer
func (l *LinkedList) Less(b types.Sortable) bool {
	if o, ok := b.(types.IterableContainer); ok {
		return list.Less(l, o)
	} else {
		return false
	}
}

func (l *LinkedList) Hash() int {
	return list.Hash(l)
}

func (l *LinkedList) String() string {
	if l.Length <= 0 {
		return "{}"
	}
	items := make([]string, 0, l.Length)
	for item, next := l.Items()(); next != nil; item, next = next() {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return "{" + strings.Join(items, ", ") + "}"
}
