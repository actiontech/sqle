package dry

import (
	"sync"
)

///////////////////////////////////////////////////////////////////////////////
// SyncBool

type SyncBool struct {
	mutex sync.RWMutex
	value bool
}

func NewSyncBool(value bool) *SyncBool {
	return &SyncBool{value: value}
}

func (s *SyncBool) Get() bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

func (s *SyncBool) Set(value bool) {
	s.mutex.Lock()
	s.value = value
	s.mutex.Unlock()
}

func (s *SyncBool) Invert() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value = !s.value
	return s.value
}

func (s *SyncBool) Swap(value bool) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := s.value
	s.value = value
	return result
}

///////////////////////////////////////////////////////////////////////////////
// SyncInt

type SyncInt struct {
	mutex sync.RWMutex
	value int
}

func NewSyncInt(value int) *SyncInt {
	return &SyncInt{value: value}
}

func (s *SyncInt) Get() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

func (s *SyncInt) Set(value int) {
	s.mutex.Lock()
	s.value = value
	s.mutex.Unlock()
}

func (s *SyncInt) Add(value int) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value += value
	return s.value
}

func (s *SyncInt) Mul(value int) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value *= value
	return s.value
}

func (s *SyncInt) Swap(value int) int {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := s.value
	s.value = value
	return result
}

///////////////////////////////////////////////////////////////////////////////
// SyncString

type SyncString struct {
	mutex sync.RWMutex
	value string
}

func NewSyncString(value string) *SyncString {
	return &SyncString{value: value}
}

func (s *SyncString) Get() string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

func (s *SyncString) Set(value string) {
	s.mutex.Lock()
	s.value = value
	s.mutex.Unlock()
}

func (s *SyncString) Append(value string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value += value
	return s.value
}

func (s *SyncString) Swap(value string) string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := s.value
	s.value = value
	return result
}

///////////////////////////////////////////////////////////////////////////////
// SyncFloat

type SyncFloat struct {
	mutex sync.RWMutex
	value float64
}

func NewSyncFloat(value float64) *SyncFloat {
	return &SyncFloat{value: value}
}

func (s *SyncFloat) Get() float64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.value
}

func (s *SyncFloat) Set(value float64) {
	s.mutex.Lock()
	s.value = value
	s.mutex.Unlock()
}

func (s *SyncFloat) Add(value float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value += value
	return s.value
}

func (s *SyncFloat) Mul(value float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.value *= value
	return s.value
}

func (s *SyncFloat) Swap(value float64) float64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	result := s.value
	s.value = value
	return result
}

///////////////////////////////////////////////////////////////////////////////
// SyncMap

type SyncMap struct {
	mutex sync.RWMutex
	m     map[string]interface{}
}

func NewSyncMap() *SyncMap {
	return &SyncMap{m: make(map[string]interface{})}
}

func (s *SyncMap) Has(key string) bool {
	s.mutex.RLock()
	_, ok := s.m[key]
	s.mutex.RUnlock()
	return ok
}

func (s *SyncMap) Get(key string) interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.m[key]
}

func (s *SyncMap) Add(key string, value interface{}) {
	s.mutex.Lock()
	s.m[key] = value
	s.mutex.Unlock()
}

func (s *SyncMap) Delete(key string) {
	s.mutex.Lock()
	delete(s.m, key)
	s.mutex.Unlock()
}

func (s *SyncMap) Int(key string) *SyncInt {
	return s.Get(key).(*SyncInt)
}

func (s *SyncMap) AddInt(key string, value int) {
	s.Add(key, NewSyncInt(value))
}

func (s *SyncMap) Float(key string) *SyncFloat {
	return s.Get(key).(*SyncFloat)
}

func (s *SyncMap) AddFloat(key string, value float64) {
	s.Add(key, NewSyncFloat(value))
}

func (s *SyncMap) Bool(key string) *SyncBool {
	return s.Get(key).(*SyncBool)
}

func (s *SyncMap) AddBool(key string, value bool) {
	s.Add(key, NewSyncBool(value))
}

func (s *SyncMap) String(key string) *SyncString {
	return s.Get(key).(*SyncString)
}

func (s *SyncMap) AddString(key string, value string) {
	s.Add(key, NewSyncString(value))
}

///////////////////////////////////////////////////////////////////////////////
// SyncStringMap

type SyncStringMap struct {
	mutex sync.RWMutex
	m     map[string]string
}

func NewSyncStringMap() *SyncStringMap {
	return &SyncStringMap{m: make(map[string]string)}
}

func (s *SyncStringMap) Has(key string) bool {
	s.mutex.RLock()
	_, ok := s.m[key]
	s.mutex.RUnlock()
	return ok
}

func (s *SyncStringMap) Get(key string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.m[key]
}

func (s *SyncStringMap) Add(key string, value string) {
	s.mutex.Lock()
	s.m[key] = value
	s.mutex.Unlock()
}

func (s *SyncStringMap) Delete(key string) {
	s.mutex.Lock()
	delete(s.m, key)
	s.mutex.Unlock()
}

///////////////////////////////////////////////////////////////////////////////
// SyncPoolMap

type SyncPoolMap struct {
	mutex sync.RWMutex
	m     map[string]*sync.Pool
}

func NewSyncPoolMap() *SyncPoolMap {
	return &SyncPoolMap{m: make(map[string]*sync.Pool)}
}

func (s *SyncPoolMap) Has(key string) bool {
	s.mutex.RLock()
	_, ok := s.m[key]
	s.mutex.RUnlock()
	return ok
}

func (s *SyncPoolMap) Get(key string) *sync.Pool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.m[key]
}

func (s *SyncPoolMap) Add(key string, value *sync.Pool) {
	s.mutex.Lock()
	s.m[key] = value
	s.mutex.Unlock()
}

func (s *SyncPoolMap) GetOrAddNew(key string, newFunc func() interface{}) *sync.Pool {
	s.mutex.Lock()
	pool := s.m[key]
	if pool == nil {
		pool = &sync.Pool{New: newFunc}
		s.m[key] = pool
	}
	s.mutex.Unlock()
	return pool
}

func (s *SyncPoolMap) Delete(key string) {
	s.mutex.Lock()
	delete(s.m, key)
	s.mutex.Unlock()
}
