package types

import (
	"fmt"
)

type MapEntry struct {
	Key   Hashable
	Value interface{}
}

func (m *MapEntry) Equals(other Equatable) bool {
	if o, ok := other.(*MapEntry); ok {
		return m.Key.Equals(o.Key)
	} else {
		return m.Key.Equals(other)
	}
}

func (m *MapEntry) Less(other Sortable) bool {
	if o, ok := other.(*MapEntry); ok {
		return m.Key.Less(o.Key)
	} else {
		return m.Key.Less(other)
	}
}

func (m *MapEntry) Hash() int {
	return m.Key.Hash()
}

func (m *MapEntry) String() string {
	return fmt.Sprintf("<MapEntry %v: %v>", m.Key, m.Value)
}
