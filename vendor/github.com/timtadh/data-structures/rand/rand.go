package rand

import (
	"math/rand"
	"sync"
)

// ThreadSafeRand provides a thread safe version of math/rand.Rand using
// the same technique used in the math/rand package to make the top level
// functions thread safe.
func ThreadSafeRand(seed int64) *rand.Rand {
	return rand.New(&lockedSource{src: rand.NewSource(seed).(rand.Source64)})
}

// from: https://golang.org/src/math/rand/rand.go?s=8161:8175#L317
type lockedSource struct {
	lk  sync.Mutex
	src rand.Source64
}

func (r *lockedSource) Int63() (n int64) {
	r.lk.Lock()
	n = r.src.Int63()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Uint64() (n uint64) {
	r.lk.Lock()
	n = r.src.Uint64()
	r.lk.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.lk.Lock()
	r.src.Seed(seed)
	r.lk.Unlock()
}

// seedPos implements Seed for a lockedSource without a race condiiton.
func (r *lockedSource) seedPos(seed int64, readPos *int8) {
	r.lk.Lock()
	r.src.Seed(seed)
	*readPos = 0
	r.lk.Unlock()
}

// read implements Read for a lockedSource without a race condition.
func (r *lockedSource) read(p []byte, readVal *int64, readPos *int8) (n int, err error) {
	r.lk.Lock()
	n, err = read(p, r.src.Int63, readVal, readPos)
	r.lk.Unlock()
	return
}

func read(p []byte, int63 func() int64, readVal *int64, readPos *int8) (n int, err error) {
	pos := *readPos
	val := *readVal
	for n = 0; n < len(p); n++ {
		if pos == 0 {
			val = int63()
			pos = 7
		}
		p[n] = byte(val)
		val >>= 8
		pos--
	}
	*readPos = pos
	*readVal = val
	return
}
