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

func (self *SyncBool) Get() bool {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.value
}

func (self *SyncBool) Set(value bool) {
	self.mutex.Lock()
	self.value = value
	self.mutex.Unlock()
}

func (self *SyncBool) Invert() bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.value = !self.value
	return self.value
}

func (self *SyncBool) Swap(value bool) bool {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	result := self.value
	self.value = value
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

func (self *SyncInt) Get() int {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.value
}

func (self *SyncInt) Set(value int) {
	self.mutex.Lock()
	self.value = value
	self.mutex.Unlock()
}

func (self *SyncInt) Add(value int) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.value += value
	return self.value
}

func (self *SyncInt) Mul(value int) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.value *= value
	return self.value
}

func (self *SyncInt) Swap(value int) int {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	result := self.value
	self.value = value
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

func (self *SyncString) Get() string {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.value
}

func (self *SyncString) Set(value string) {
	self.mutex.Lock()
	self.value = value
	self.mutex.Unlock()
}

func (self *SyncString) Append(value string) string {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.value += value
	return self.value
}

func (self *SyncString) Swap(value string) string {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	result := self.value
	self.value = value
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

func (self *SyncFloat) Get() float64 {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.value
}

func (self *SyncFloat) Set(value float64) {
	self.mutex.Lock()
	self.value = value
	self.mutex.Unlock()
}

func (self *SyncFloat) Add(value float64) float64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.value += value
	return self.value
}

func (self *SyncFloat) Mul(value float64) float64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	self.value *= value
	return self.value
}

func (self *SyncFloat) Swap(value float64) float64 {
	self.mutex.Lock()
	defer self.mutex.Unlock()
	result := self.value
	self.value = value
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

func (self *SyncMap) Has(key string) bool {
	self.mutex.RLock()
	_, ok := self.m[key]
	self.mutex.RUnlock()
	return ok
}

func (self *SyncMap) Get(key string) interface{} {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.m[key]
}

func (self *SyncMap) Add(key string, value interface{}) {
	self.mutex.Lock()
	self.m[key] = value
	self.mutex.Unlock()
}

func (self *SyncMap) Delete(key string) {
	self.mutex.Lock()
	delete(self.m, key)
	self.mutex.Unlock()
}

func (self *SyncMap) Int(key string) *SyncInt {
	return self.Get(key).(*SyncInt)
}

func (self *SyncMap) AddInt(key string, value int) {
	self.Add(key, NewSyncInt(value))
}

func (self *SyncMap) Float(key string) *SyncFloat {
	return self.Get(key).(*SyncFloat)
}

func (self *SyncMap) AddFloat(key string, value float64) {
	self.Add(key, NewSyncFloat(value))
}

func (self *SyncMap) Bool(key string) *SyncBool {
	return self.Get(key).(*SyncBool)
}

func (self *SyncMap) AddBool(key string, value bool) {
	self.Add(key, NewSyncBool(value))
}

func (self *SyncMap) String(key string) *SyncString {
	return self.Get(key).(*SyncString)
}

func (self *SyncMap) AddString(key string, value string) {
	self.Add(key, NewSyncString(value))
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

func (self *SyncStringMap) Has(key string) bool {
	self.mutex.RLock()
	_, ok := self.m[key]
	self.mutex.RUnlock()
	return ok
}

func (self *SyncStringMap) Get(key string) string {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.m[key]
}

func (self *SyncStringMap) Add(key string, value string) {
	self.mutex.Lock()
	self.m[key] = value
	self.mutex.Unlock()
}

func (self *SyncStringMap) Delete(key string) {
	self.mutex.Lock()
	delete(self.m, key)
	self.mutex.Unlock()
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

func (self *SyncPoolMap) Has(key string) bool {
	self.mutex.RLock()
	_, ok := self.m[key]
	self.mutex.RUnlock()
	return ok
}

func (self *SyncPoolMap) Get(key string) *sync.Pool {
	self.mutex.RLock()
	defer self.mutex.RUnlock()
	return self.m[key]
}

func (self *SyncPoolMap) Add(key string, value *sync.Pool) {
	self.mutex.Lock()
	self.m[key] = value
	self.mutex.Unlock()
}

func (self *SyncPoolMap) GetOrAddNew(key string, newFunc func() interface{}) *sync.Pool {
	self.mutex.Lock()
	pool := self.m[key]
	if pool == nil {
		pool = &sync.Pool{New: newFunc}
		self.m[key] = pool
	}
	self.mutex.Unlock()
	return pool
}

func (self *SyncPoolMap) Delete(key string) {
	self.mutex.Lock()
	delete(self.m, key)
	self.mutex.Unlock()
}
