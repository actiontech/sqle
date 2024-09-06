package util

import (
	"sync"
	"time"
)

type item struct {
	value     string
	createdAt int64
}

type TTLMap struct {
	ttl int64
	m   map[string]*item
	l   sync.Mutex
}

func NewTTLMap(ttl int64) (m *TTLMap) {
	m = &TTLMap{ttl: ttl, m: make(map[string]*item)}
	return
}

func (m *TTLMap) Len() int {
	return len(m.m)
}

func (m *TTLMap) Put(k, v string) {
	m.l.Lock()
	defer m.l.Unlock()
	it, ok := m.m[k]
	if !ok {
		it = &item{value: v}
		m.m[k] = it
	}
	it.createdAt = time.Now().Unix()
}

func (m *TTLMap) Get(k string) (v string) {
	m.l.Lock()
	defer m.l.Unlock()
	if it, ok := m.m[k]; ok {
		if time.Now().Unix()-it.createdAt > m.ttl {
			delete(m.m, k)
		} else {
			v = it.value
		}
	}
	return
}
