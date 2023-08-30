package set

import (
	"fmt"
	"strings"

	"github.com/timtadh/data-structures/errors"
	"github.com/timtadh/data-structures/types"
)

type SetMap struct {
	types.Map
}

func NewSetMap(m types.Map) *SetMap {
	return &SetMap{m}
}

func (s *SetMap) String() string {
	if s.Size() <= 0 {
		return "{}"
	}
	items := make([]string, 0, s.Size())
	for item, next := s.Items()(); next != nil; item, next = next() {
		items = append(items, fmt.Sprintf("%v", item))
	}
	return "{" + strings.Join(items, ", ") + "}"
}

func (s *SetMap) Items() types.KIterator {
	return s.Keys()
}

// unimplemented
func (s *SetMap) Item(item types.Hashable) (types.Hashable, error) {
	return nil, errors.Errorf("un-implemented")
}

// unimplemented
func (s *SetMap) Equals(o types.Equatable) bool {
	panic(errors.Errorf("un-implemented"))
}

func (s *SetMap) Add(item types.Hashable) (err error) {
	return s.Put(item, nil)
}

func (s *SetMap) Delete(item types.Hashable) (err error) {
	_, err = s.Remove(item)
	return err
}

func (s *SetMap) Extend(items types.KIterator) (err error) {
	for item, next := items(); next != nil; item, next = next() {
		err := s.Add(item)
		if err != nil {
			return err
		}
	}
	return nil
}

// Unions s with o and returns a new SetMap (with a LinearHash)
func (s *SetMap) Union(other types.Set) (types.Set, error) {
	return Union(s, other)
}

// Unions s with o and returns a new SetMap (with a LinearHash)
func (s *SetMap) Intersect(other types.Set) (types.Set, error) {
	return Intersect(s, other)
}

// Unions s with o and returns a new SetMap (with a LinearHash)
func (s *SetMap) Subtract(other types.Set) (types.Set, error) {
	return Subtract(s, other)
}

// Is s a subset of o?
func (s *SetMap) Subset(o types.Set) bool {
	return Subset(s, o)
}

// Is s a proper subset of o?
func (s *SetMap) ProperSubset(o types.Set) bool {
	return ProperSubset(s, o)
}

// Is s a superset of o?
func (s *SetMap) Superset(o types.Set) bool {
	return Superset(s, o)
}

// Is s a proper superset of o?
func (s *SetMap) ProperSuperset(o types.Set) bool {
	return ProperSuperset(s, o)
}
