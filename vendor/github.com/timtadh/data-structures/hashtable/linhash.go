package hashtable

import (
	"github.com/timtadh/data-structures/tree/avl"
	. "github.com/timtadh/data-structures/types"
)

const (
	UTILIZATION       = .75
	RECORDS_PER_BLOCK = 16
)

type bst struct {
	hash  int
	key   Hashable
	value interface{}
	left  *bst
	right *bst
}

type LinearHash struct {
	table []*avl.AvlNode
	n     uint
	r     uint
	i     uint
}

func NewLinearHash() *LinearHash {
	N := uint(32)
	I := uint(5)
	return &LinearHash{
		table: make([]*avl.AvlNode, N),
		n:     N,
		r:     0,
		i:     I,
	}
}

func (self *LinearHash) bucket(key Hashable) uint {
	m := uint(key.Hash() & ((1 << self.i) - 1))
	if m < self.n {
		return m
	} else {
		return m ^ (1 << (self.i - 1))
	}
}

func (self *LinearHash) Size() int {
	return int(self.r)
}

func (self *LinearHash) Put(key Hashable, value interface{}) (err error) {
	var updated bool
	bkt_idx := self.bucket(key)
	self.table[bkt_idx], updated = self.table[bkt_idx].Put(key, value)
	if !updated {
		self.r += 1
	}
	if float64(self.r) > UTILIZATION*float64(self.n)*float64(RECORDS_PER_BLOCK) {
		return self.split()
	}
	return nil
}

func (self *LinearHash) Get(key Hashable) (value interface{}, err error) {
	bkt_idx := self.bucket(key)
	return self.table[bkt_idx].Get(key)
}

func (self *LinearHash) Has(key Hashable) bool {
	bkt_idx := self.bucket(key)
	return self.table[bkt_idx].Has(key)
}

func (self *LinearHash) Remove(key Hashable) (value interface{}, err error) {
	bkt_idx := self.bucket(key)
	self.table[bkt_idx], value, err = self.table[bkt_idx].Remove(key)
	if err == nil {
		self.r -= 1
	}
	return
}

func (self *LinearHash) split() (err error) {
	bkt_idx := self.n % (1 << (self.i - 1))
	old_bkt := self.table[bkt_idx]
	var bkt_a, bkt_b *avl.AvlNode
	self.n += 1
	if self.n > (1 << self.i) {
		self.i += 1
	}
	for key, value, next := old_bkt.Iterate()(); next != nil; key, value, next = next() {
		if self.bucket(key.(Hashable)) == bkt_idx {
			bkt_a, _ = bkt_a.Put(key.(Hashable), value)
		} else {
			bkt_b, _ = bkt_b.Put(key.(Hashable), value)
		}
	}
	self.table[bkt_idx] = bkt_a
	self.table = append(self.table, bkt_b)
	return nil
}

func (self *LinearHash) Iterate() KVIterator {
	table := self.table
	i := 0
	iter := table[i].Iterate()
	var kv_iterator KVIterator
	kv_iterator = func() (key Hashable, val interface{}, next KVIterator) {
		key, val, iter = iter()
		for iter == nil {
			i++
			if i >= len(table) {
				return nil, nil, nil
			}
			key, val, iter = table[i].Iterate()()
		}
		return key, val, kv_iterator
	}
	return kv_iterator
}

func (self *LinearHash) Items() (vi KIterator) {
	return MakeItemsIterator(self)
}

func (self *LinearHash) Keys() KIterator {
	return MakeKeysIterator(self)
}

func (self *LinearHash) Values() Iterator {
	return MakeValuesIterator(self)
}
