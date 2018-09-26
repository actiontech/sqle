package util

import (
	"time"
)

func PeriodicalChan(fn func() bool) (sigChan chan bool, exitChan chan bool) {
	return PeriodicalChan2(fn, 500)
}

func PeriodicalChan2(fn func() bool, intervalMs int) (sigChan chan bool, exitChan chan bool) {
	sigChan = make(chan bool, 0)
	exitChan = make(chan bool, 1)
	go func() {
		for {
			if fn() {
				sigChan <- true
				return
			}
			select {
			case <-exitChan:
				return
			default:
				time.Sleep(time.Duration(intervalMs) * time.Millisecond)
			}
		}
	}()
	return sigChan, exitChan
}

func WaitUntil(fn func() bool, timeoutSeconds int) chan bool {
	ch := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			ch <- true
		}
	}()
	return WaitUntilOrChan(fn, ch)
}

func WaitUntilOrChan(fn func() bool, timeoutChan chan bool) chan bool {
	ch := make(chan bool, 1)
	a, aq := PeriodicalChan(fn)
	go func() {
		defer func() {
			aq <- true
		}()
		select {
		case <-a:
			ch <- true
		case <-timeoutChan:
			ch <- false
		}
	}()
	return ch
}

func WaitWhile(fn func() bool, timeoutSeconds int) chan bool {
	ch := make(chan bool, 1)
	a, aq := PeriodicalChan(func() bool {
		return !fn()
	})
	go func() {
		defer func() {
			aq <- true
		}()
		select {
		case <-a:
			ch <- false
		case <-time.After(time.Duration(timeoutSeconds) * time.Second):
			ch <- true
		}
	}()
	return ch
}

func Timeout(fn func(), timeoutSeconds int) (result chan bool) {
	result = make(chan bool, 1)
	fnFinished := make(chan bool, 1)
	go func() {
		fn()
		fnFinished <- true
	}()
	select {
	case <-fnFinished:
		result <- true
	case <-time.After(time.Duration(timeoutSeconds) * time.Second):
		result <- false
	}
	return
}
