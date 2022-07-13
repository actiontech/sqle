package set

import (
	"github.com/timtadh/data-structures/hashtable"
	"github.com/timtadh/data-structures/types"
)

func newSetBestType(a types.Set, sizeHint int) types.Set {
	switch a.(type) {
	case *MapSet:
		return NewMapSet(NewSortedSet(sizeHint))
	case *SortedSet:
		return NewSortedSet(sizeHint)
	case *SetMap:
		return NewSetMap(hashtable.NewLinearHash())
	default:
		return NewSortedSet(sizeHint)
	}
}

func Union(a, b types.Set) (types.Set, error) {
	c := newSetBestType(a, a.Size()+b.Size())
	err := c.Extend(a.Items())
	if err != nil {
		return nil, err
	}
	err = c.Extend(b.Items())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func Intersect(a, b types.Set) (types.Set, error) {
	c := newSetBestType(a, a.Size()+b.Size())
	for item, next := a.Items()(); next != nil; item, next = next() {
		if b.Has(item) {
			err := c.Add(item)
			if err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

// Unions s with o and returns a new Sorted Set
func Subtract(a, b types.Set) (types.Set, error) {
	c := newSetBestType(a, a.Size()+b.Size())
	for item, next := a.Items()(); next != nil; item, next = next() {
		if !b.Has(item) {
			err := c.Add(item)
			if err != nil {
				return nil, err
			}
		}
	}
	return c, nil
}

func Subset(a, b types.Set) bool {
	if a.Size() > b.Size() {
		return false
	}
	for item, next := a.Items()(); next != nil; item, next = next() {
		if !b.Has(item) {
			return false
		}
	}
	return true
}

func ProperSubset(a, b types.Set) bool {
	if a.Size() >= b.Size() {
		return false
	}
	return Subset(a, b)
}

func Superset(a, b types.Set) bool {
	return Subset(b, a)
}

func ProperSuperset(a, b types.Set) bool {
	return ProperSubset(b, a)
}
