package util

import (
	"testing"
	"time"
)

func TestDurationCounter(t *testing.T) {
	dc := NewDurationCounter(2 * time.Second)
	dc.Inc()
	if 1 != dc.GetCount() {
		t.Fatal("count should = 1")
	}
	time.Sleep(3 * time.Second)
	if 0 != dc.GetCount() {
		t.Fatal("count should = 0")
	}
}
