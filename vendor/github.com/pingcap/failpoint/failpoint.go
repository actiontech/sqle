// Copyright 2019 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package failpoint

import (
	"context"
	"sync"
)

const failpointCtxKey HookKey = "__failpoint_ctx_key__"

type (
	// HookKey represents the type of failpoint hook function key in context
	HookKey string

	// Value represents value that retrieved from failpoint terms.
	// It can be used as following types:
	// 1. val.(int)      // GO_FAILPOINTS="failpoint-name=return(1)"
	// 2. val.(string)   // GO_FAILPOINTS="failpoint-name=return('1')"
	// 3. val.(bool)     // GO_FAILPOINTS="failpoint-name=return(true)"
	Value interface{}

	// Hook is used to filter failpoint, if the hook returns false and the
	// failpoint will not to be evaluated.
	Hook func(ctx context.Context, fpname string) bool

	// Failpoint is a point to inject a failure
	Failpoint struct {
		mu       sync.RWMutex
		t        *terms
		waitChan chan struct{}
	}
)

// Pause will pause until the failpoint is disabled.
func (fp *Failpoint) Pause() {
	<-fp.waitChan
}

// Enable sets a failpoint to a given failpoint description.
func (fp *Failpoint) Enable(inTerms string) error {
	t, err := newTerms(inTerms, fp)
	if err != nil {
		return err
	}
	fp.mu.Lock()
	fp.t = t
	fp.waitChan = make(chan struct{})
	fp.mu.Unlock()
	return nil
}

// EnableWith enables and locks the failpoint, the lock prevents
// the failpoint to be evaluated. It invokes the action while holding
// the lock. It is useful when enables a panic failpoint
// and does some post actions before the failpoint being evaluated.
func (fp *Failpoint) EnableWith(inTerms string, action func() error) error {
	t, err := newTerms(inTerms, fp)
	if err != nil {
		return err
	}
	fp.mu.Lock()
	defer fp.mu.Unlock()
	fp.t = t
	fp.waitChan = make(chan struct{})
	if err := action(); err != nil {
		return err
	}
	return nil
}

// Disable stops a failpoint
func (fp *Failpoint) Disable() error {
	select {
	case <-fp.waitChan:
		return ErrDisabled
	default:
		close(fp.waitChan)
	}

	fp.mu.Lock()
	defer fp.mu.Unlock()
	if fp.t == nil {
		return ErrDisabled
	}
	fp.t = nil
	return nil
}

// Eval evaluates a failpoint's value, It will return the evaluated value or
// an error if the failpoint is disabled or failed to eval
func (fp *Failpoint) Eval() (Value, error) {
	fp.mu.RLock()
	defer fp.mu.RUnlock()
	if fp.t == nil {
		return nil, ErrDisabled
	}
	v, err := fp.t.eval()
	if err != nil {
		return nil, err
	}
	return v, nil
}
