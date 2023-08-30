package linked

import (
	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/list"
	"github.com/timtadh/data-structures/set"
	"github.com/timtadh/data-structures/types"
)

type UniqueDeque struct {
	queue *LinkedList
	set   types.Set
}

// A double ended queue that only allows unique items inside. Constructed from a
// doubly linked list and a linear hash table.
func NewUniqueDeque() *UniqueDeque {
	return &UniqueDeque{
		queue: New(),
		set:   set.NewSetMap(hashtable.NewLinearHash()),
	}
}

func (l *UniqueDeque) Size() int {
	return l.queue.Size()
}

func (l *UniqueDeque) Items() (it types.KIterator) {
	return l.queue.Items()
}

func (l *UniqueDeque) Backwards() (it types.KIterator) {
	return l.queue.Backwards()
}

func (l *UniqueDeque) Has(item types.Hashable) bool {
	return l.set.Has(item)
}

func (l *UniqueDeque) Push(item types.Hashable) (err error) {
	return l.EnqueBack(item)
}

func (l *UniqueDeque) Pop() (item types.Hashable, err error) {
	return l.DequeBack()
}

func (l *UniqueDeque) EnqueFront(item types.Hashable) (err error) {
	if l.Has(item) {
		return nil
	}
	err = l.queue.EnqueFront(item)
	if err != nil {
		return err
	}
	return l.set.Add(item)
}

func (l *UniqueDeque) EnqueBack(item types.Hashable) (err error) {
	if l.Has(item) {
		return nil
	}
	err = l.queue.EnqueBack(item)
	if err != nil {
		return err
	}
	return l.set.Add(item)
}

func (l *UniqueDeque) DequeFront() (item types.Hashable, err error) {
	item, err = l.queue.DequeFront()
	if err != nil {
		return nil, err
	}
	err = l.set.Delete(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (l *UniqueDeque) DequeBack() (item types.Hashable, err error) {
	item, err = l.queue.DequeBack()
	if err != nil {
		return nil, err
	}
	err = l.set.Delete(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (l *UniqueDeque) First() (item types.Hashable) {
	return l.queue.First()
}

func (l *UniqueDeque) Last() (item types.Hashable) {
	return l.queue.Last()
}

// Can be compared to any types.IterableContainer
func (l *UniqueDeque) Equals(b types.Equatable) bool {
	return l.queue.Equals(b)
}

// Can be compared to any types.IterableContainer
func (l *UniqueDeque) Less(b types.Sortable) bool {
	return l.queue.Less(b)
}

func (l *UniqueDeque) Hash() int {
	return list.Hash(l.queue)
}

func (l *UniqueDeque) String() string {
	return l.queue.String()
}
