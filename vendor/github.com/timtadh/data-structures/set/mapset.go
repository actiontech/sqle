package set

import (
	"fmt"
)

import (
	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/types"
)

type MapSet struct {
	Set types.Set
}

func mapEntry(item types.Hashable) (*types.MapEntry, error) {
	if me, ok := item.(*types.MapEntry); ok {
		return me, nil
	} else {
		return nil, errors.Errorf("Must pass a *types.MapEntry, got %T", item)
	}
}

func asMapEntry(item types.Hashable) *types.MapEntry {
	if me, ok := item.(*types.MapEntry); ok {
		return me
	} else {
		return &types.MapEntry{item, nil}
	}
}

func NewMapSet(set types.Set) *MapSet {
	return &MapSet{set}
}

func (m *MapSet) Equals(b types.Equatable) bool {
	if o, ok := b.(*MapSet); ok {
		return m.Set.Equals(o.Set)
	} else {
		return false
	}
}

func (m *MapSet) Size() int {
	return m.Set.Size()
}

func (m *MapSet) Keys() types.KIterator {
	return types.MakeKeysIterator(m)
}

func (m *MapSet) Values() types.Iterator {
	return types.MakeValuesIterator(m)
}

func (m *MapSet) Iterate() (kvit types.KVIterator) {
	items := m.Items()
	kvit = func() (key types.Hashable, value interface{}, _ types.KVIterator) {
		var item types.Hashable
		item, items = items()
		if items == nil {
			return nil, nil, nil
		}
		me := item.(*types.MapEntry)
		return me.Key, me.Value, kvit
	}
	return kvit
}

func (m *MapSet) Items() types.KIterator {
	return m.Set.Items()
}

func (m *MapSet) Has(key types.Hashable) bool {
	return m.Set.Has(asMapEntry(key))
}

func (m *MapSet) Add(item types.Hashable) (err error) {
	me, err := mapEntry(item)
	if err != nil {
		return err
	}
	return m.Set.Add(me)
}

func (m *MapSet) Put(key types.Hashable, value interface{}) (err error) {
	return m.Add(&types.MapEntry{key, value})
}

func (m *MapSet) Get(key types.Hashable) (value interface{}, err error) {
	item, err := m.Set.Item(asMapEntry(key))
	if err != nil {
		return nil, err
	}
	me, err := mapEntry(item)
	if err != nil {
		return nil, err
	}
	return me.Value, nil
}

func (m *MapSet) Item(key types.Hashable) (me types.Hashable, err error) {
	item, err := m.Set.Item(asMapEntry(key))
	if err != nil {
		return nil, err
	}
	me, err = mapEntry(item)
	if err != nil {
		return nil, err
	}
	return me, nil
}

func (m *MapSet) Remove(key types.Hashable) (value interface{}, err error) {
	item, err := m.Get(asMapEntry(key))
	if err != nil {
		return nil, err
	}
	err = m.Delete(key)
	if err != nil {
		return nil, err
	}
	return item, nil
}

func (m *MapSet) Delete(item types.Hashable) (err error) {
	return m.Set.Delete(asMapEntry(item))
}

func (m *MapSet) Extend(items types.KIterator) (err error) {
	for item, items := items(); items != nil; item, items = items() {
		err := m.Add(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MapSet) Union(b types.Set) (types.Set, error) {
	return Union(m, b)
}

func (m *MapSet) Intersect(b types.Set) (types.Set, error) {
	return Intersect(m, b)
}

func (m *MapSet) Subtract(b types.Set) (types.Set, error) {
	return Subtract(m, b)
}

func (m *MapSet) Subset(b types.Set) bool {
	return Subset(m, b)
}

func (m *MapSet) Superset(b types.Set) bool {
	return Superset(m, b)
}

func (m *MapSet) ProperSubset(b types.Set) bool {
	return ProperSubset(m, b)
}

func (m *MapSet) ProperSuperset(b types.Set) bool {
	return ProperSuperset(m, b)
}

func (m *MapSet) String() string {
	return fmt.Sprintf("%v", m.Set)
}
