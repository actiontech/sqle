package hashtable

import . "github.com/timtadh/data-structures/types"
import . "github.com/timtadh/data-structures/errors"

type entry struct {
	key   Hashable
	value interface{}
	next  *entry
}

type Hash struct {
	table []*entry
	size  int
}

func (self *entry) Put(key Hashable, value interface{}) (e *entry, appended bool) {
	if self == nil {
		return &entry{key, value, nil}, true
	}
	if self.key.Equals(key) {
		self.value = value
		return self, false
	} else {
		self.next, appended = self.next.Put(key, value)
		return self, appended
	}
}

func (self *entry) Get(key Hashable) (has bool, value interface{}) {
	if self == nil {
		return false, nil
	} else if self.key.Equals(key) {
		return true, self.value
	} else {
		return self.next.Get(key)
	}
}

func (self *entry) Remove(key Hashable) *entry {
	if self == nil {
		panic(Errors["not-found-in-bucket"](key))
	}
	if self.key.Equals(key) {
		return self.next
	} else {
		self.next = self.next.Remove(key)
		return self
	}
}

func NewHashTable(initial_size int) *Hash {
	return &Hash{
		table: make([]*entry, initial_size),
		size:  0,
	}
}

func (self *Hash) bucket(key Hashable) int {
	return key.Hash() % len(self.table)
}

func (self *Hash) Size() int { return self.size }

func (self *Hash) Put(key Hashable, value interface{}) (err error) {
	bucket := self.bucket(key)
	var appended bool
	self.table[bucket], appended = self.table[bucket].Put(key, value)
	if appended {
		self.size += 1
	}
	if self.size*2 > len(self.table) {
		return self.expand()
	}
	return nil
}

func (self *Hash) expand() error {
	table := self.table
	self.table = make([]*entry, len(table)*2)
	self.size = 0
	for _, E := range table {
		for e := E; e != nil; e = e.next {
			if err := self.Put(e.key, e.value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (self *Hash) Get(key Hashable) (value interface{}, err error) {
	bucket := self.bucket(key)
	if has, value := self.table[bucket].Get(key); has {
		return value, nil
	} else {
		return nil, Errors["not-found"](key)
	}
}

func (self *Hash) Has(key Hashable) (has bool) {
	has, _ = self.table[self.bucket(key)].Get(key)
	return
}

func (self *Hash) Remove(key Hashable) (value interface{}, err error) {
	bucket := self.bucket(key)
	has, value := self.table[bucket].Get(key)
	if !has {
		return nil, Errors["not-found"](key)
	}
	self.table[bucket] = self.table[bucket].Remove(key)
	self.size -= 1
	return value, nil
}

func (self *Hash) Iterate() KVIterator {
	table := self.table
	i := -1
	var e *entry
	var kv_iterator KVIterator
	kv_iterator = func() (key Hashable, val interface{}, next KVIterator) {
		for e == nil {
			i++
			if i >= len(table) {
				return nil, nil, nil
			}
			e = table[i]
		}
		key = e.key
		val = e.value
		e = e.next
		return key, val, kv_iterator
	}
	return kv_iterator
}

func (self *Hash) Items() (vi KIterator) {
	return MakeItemsIterator(self)
}

func (self *Hash) Keys() KIterator {
	return MakeKeysIterator(self)
}

func (self *Hash) Values() Iterator {
	return MakeValuesIterator(self)
}
