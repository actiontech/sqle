package list

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"log"
	"sort"
	"strings"
)

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/types"
)

type MList struct {
	List
	MarshalItem   types.ItemMarshal
	UnmarshalItem types.ItemUnmarshal
}

func NewMList(list *List, marshal types.ItemMarshal, unmarshal types.ItemUnmarshal) *MList {
	return &MList{
		List:          *list,
		MarshalItem:   marshal,
		UnmarshalItem: unmarshal,
	}
}

func (m *MList) MarshalBinary() ([]byte, error) {
	items := make([][]byte, 0, m.Size())
	_cap := make([]byte, 4)
	size := make([]byte, 4)
	binary.LittleEndian.PutUint32(size, uint32(m.Size()))
	if m.List.fixed {
		binary.LittleEndian.PutUint32(_cap, uint32(cap(m.List.list)))
		items = append(items, []byte{1})
	} else {
		binary.LittleEndian.PutUint32(_cap, uint32(m.Size()))
		items = append(items, []byte{0})
	}
	items = append(items, _cap)
	items = append(items, size)
	for item, next := m.Items()(); next != nil; item, next = next() {
		b, err := m.MarshalItem(item)
		if err != nil {
			return nil, err
		}
		size := make([]byte, 4)
		binary.LittleEndian.PutUint32(size, uint32(len(b)))
		items = append(items, size, b)
	}
	return bytes.Join(items, []byte{}), nil
}

func (m *MList) UnmarshalBinary(bytes []byte) error {
	m.List.fixed = bytes[0] == 1
	_cap := int(binary.LittleEndian.Uint32(bytes[1:5]))
	size := int(binary.LittleEndian.Uint32(bytes[5:9]))
	off := 9
	m.list = make([]types.Hashable, 0, _cap)
	for i := 0; i < size; i++ {
		s := off
		e := off + 4
		size := int(binary.LittleEndian.Uint32(bytes[s:e]))
		s = e
		e = s + size
		item, err := m.UnmarshalItem(bytes[s:e])
		if err != nil {
			return err
		}
		m.Append(item)
		off = e
	}
	return nil
}

type Sortable struct {
	List
}

func NewSortable(list *List) sort.Interface {
	return &Sortable{*list}
}

func (s *Sortable) Len() int {
	return s.Size()
}

func (s *Sortable) Less(i, j int) bool {
	a, err := s.Get(i)
	if err != nil {
		log.Panic(err)
	}
	b, err := s.Get(j)
	if err != nil {
		log.Panic(err)
	}
	return a.Less(b)
}

func (s *Sortable) Swap(i, j int) {
	a, err := s.Get(i)
	if err != nil {
		log.Panic(err)
	}
	b, err := s.Get(j)
	if err != nil {
		log.Panic(err)
	}
	err = s.Set(i, b)
	if err != nil {
		log.Panic(err)
	}
	err = s.Set(j, a)
	if err != nil {
		log.Panic(err)
	}
}

type List struct {
	list  []types.Hashable
	fixed bool
}

// Creates a list.
func New(initialSize int) *List {
	return newList(initialSize, false)
}

// Creates a Fixed Size list.
func Fixed(size int) *List {
	return newList(size, true)
}

func newList(initialSize int, fixedSize bool) *List {
	return &List{
		list:  make([]types.Hashable, 0, initialSize),
		fixed: fixedSize,
	}
}

func FromSlice(list []types.Hashable) *List {
	l := &List{
		list: make([]types.Hashable, len(list)),
	}
	copy(l.list, list)
	return l
}

func (l *List) Copy() *List {
	list := make([]types.Hashable, len(l.list), cap(l.list))
	copy(list, l.list)
	return &List{list, l.fixed}
}

func (l *List) Clear() {
	l.list = l.list[:0]
}

func (l *List) Size() int {
	return len(l.list)
}

func (l *List) Full() bool {
	return l.fixed && cap(l.list) == len(l.list)
}

func (l *List) Empty() bool {
	return len(l.list) == 0
}

func (l *List) Has(item types.Hashable) (has bool) {
	for i := range l.list {
		if l.list[i].Equals(item) {
			return true
		}
	}
	return false
}

func (l *List) Equals(b types.Equatable) bool {
	if o, ok := b.(types.IterableContainer); ok {
		return Equals(l, o)
	} else {
		return false
	}
}

func Equals(a, b types.IterableContainer) bool {
	if a.Size() != b.Size() {
		return false
	}
	ca, ai := a.Items()()
	cb, bi := b.Items()()
	for ai != nil || bi != nil {
		if !ca.Equals(cb) {
			return false
		}
		ca, ai = ai()
		cb, bi = bi()
	}
	return true
}

func (l *List) Less(b types.Sortable) bool {
	if o, ok := b.(types.IterableContainer); ok {
		return Less(l, o)
	} else {
		return false
	}
}

func Less(a, b types.IterableContainer) bool {
	if a.Size() < b.Size() {
		return true
	} else if a.Size() > b.Size() {
		return false
	}
	ca, ai := a.Items()()
	cb, bi := b.Items()()
	for ai != nil || bi != nil {
		if ca.Less(cb) {
			return true
		} else if !ca.Equals(cb) {
			return false
		}
		ca, ai = ai()
		cb, bi = bi()
	}
	return false
}

func (l *List) Hash() int {
	h := fnv.New32a()
	if len(l.list) == 0 {
		return 0
	}
	bs := make([]byte, 4)
	for _, item := range l.list {
		binary.LittleEndian.PutUint32(bs, uint32(item.Hash()))
		h.Write(bs)
	}
	return int(h.Sum32())
}

func Hash(a types.ListIterable) int {
	h := fnv.New32a()
	bs := make([]byte, 4)
	for item, next := a.Items()(); next != nil; item, next = next() {
		binary.LittleEndian.PutUint32(bs, uint32(item.Hash()))
		h.Write(bs)
	}
	return int(h.Sum32())
}

func (l *List) Items() (it types.KIterator) {
	i := 0
	return func() (item types.Hashable, next types.KIterator) {
		if i < len(l.list) {
			item = l.list[i]
			i++
			return item, it
		}
		return nil, nil
	}
}

func (l *List) ItemsInReverse() (it types.KIterator) {
	i := len(l.list) - 1
	return func() (item types.Hashable, next types.KIterator) {
		if i >= 0 {
			item = l.list[i]
			i--
			return item, it
		}
		return nil, nil
	}
}

func (l *List) Get(i int) (item types.Hashable, err error) {
	if i < 0 || i >= len(l.list) {
		return nil, errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	return l.list[i], nil
}

func (l *List) Set(i int, item types.Hashable) (err error) {
	if i < 0 || i >= len(l.list) {
		return errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	l.list[i] = item
	return nil
}

func (l *List) Push(item types.Hashable) error {
	return l.Append(item)
}

func (l *List) Append(item types.Hashable) error {
	return l.Insert(len(l.list), item)
}

func (l *List) Insert(i int, item types.Hashable) error {
	if i < 0 || i > len(l.list) {
		return errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	if len(l.list) == cap(l.list) {
		if err := l.expand(); err != nil {
			return err
		}
	}
	l.list = l.list[:len(l.list)+1]
	for j := len(l.list) - 1; j > 0; j-- {
		if j == i {
			l.list[i] = item
			break
		}
		l.list[j] = l.list[j-1]
	}
	if i == 0 {
		l.list[i] = item
	}
	return nil
}

func (l *List) Extend(it types.KIterator) (err error) {
	for item, next := it(); next != nil; item, next = next() {
		if err := l.Append(item); err != nil {
			return err
		}
	}
	return nil
}

func (l *List) Pop() (item types.Hashable, err error) {
	item, err = l.Get(len(l.list) - 1)
	if err != nil {
		return nil, err
	}
	err = l.Remove(len(l.list) - 1)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (l *List) Remove(i int) error {
	if i < 0 || i >= len(l.list) {
		return errors.Errorf("Access out of bounds. len(*List) = %v, idx = %v", len(l.list), i)
	}
	dst := l.list[i : len(l.list)-1]
	src := l.list[i+1 : len(l.list)]
	copy(dst, src)
	l.list = l.list[:len(l.list)-1]
	if err := l.shrink(); err != nil {
		return err
	}
	return nil
}

func (l *List) String() string {
	if len(l.list) <= 0 {
		return "{}"
	}
	items := make([]string, 0, len(l.list))
	for _, item := range l.list {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return "{" + strings.Join(items, ", ") + "}"
}

func (l *List) expand() error {
	if l.fixed {
		return errors.Errorf("Fixed size list is full!")
	}
	list := l.list
	if cap(list) < 100 && cap(list) != 0 {
		l.list = make([]types.Hashable, len(list), cap(list)*2)
	} else {
		l.list = make([]types.Hashable, len(list), cap(list)+100)
	}
	copy(l.list, list)
	return nil
}

func (l *List) shrink() error {
	if l.fixed {
		return nil
	}
	if (len(l.list)-1)*2 >= cap(l.list) || cap(l.list)/2 <= 10 {
		return nil
	}
	list := l.list
	l.list = make([]types.Hashable, len(list), cap(list)/2+1)
	copy(l.list, list)
	return nil
}
