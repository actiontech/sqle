package util

import (
	"time"
)

// Difference from time.Tick: the first tick sent immediately
func Tick(d time.Duration) <-chan time.Time {
	if d <= 0 {
		return nil
	}
	c := time.NewTicker(d).C
	ret := make(chan time.Time, 1)
	go func() {
		ret <- time.Now()
		for {
			select {
			case t := <-c:
				ret <- t
			}
		}
	}()
	return ret
}
