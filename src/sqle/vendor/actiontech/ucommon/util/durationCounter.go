package util

import (
	"sync"
	"time"
)

type DurationCounter struct {
	duration time.Duration
	count    int
	mutex    *sync.Mutex
	decTimes []time.Time
}

func NewDurationCounter(duration time.Duration) *DurationCounter {
	ret := DurationCounter{}
	ret.duration = duration
	ret.mutex = &sync.Mutex{}
	ret.count = 0
	ret.decTimes = []time.Time{}
	return &ret
}

func (d *DurationCounter) GetCount() (ret int) {
	d.mutex.Lock()
	now := time.Now()
	for len(d.decTimes) > 0 {
		if !d.decTimes[0].After(now) {
			d.count--
			d.decTimes = d.decTimes[1:]
		} else {
			break
		}
	}
	ret = d.count
	d.mutex.Unlock()
	return ret
}

func (d *DurationCounter) Inc() (ret int) {
	d.mutex.Lock()
	now := time.Now()
	d.count = d.count + 1
	d.decTimes = append(d.decTimes, now.Add(d.duration))
	ret = d.count
	d.mutex.Unlock()
	return ret
}

func (d *DurationCounter) CreateNew() *DurationCounter {
	return NewDurationCounter(d.duration)
}

func (d *DurationCounter) Reset() {
	d.mutex.Lock()
	d.count = 0
	d.decTimes = []time.Time{}
	d.mutex.Unlock()
}
