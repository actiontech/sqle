// Copyright 2021 TiKV Authors
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

// NOTE: The code in this file is based on code from the
// TiDB project, licensed under the Apache License v 2.0
//
// https://github.com/pingcap/tidb/tree/cc5e161ac06827589c4966674597c137cc9e809c/store/tikv/util/rate_limit.go
//

// Copyright 2021 PingCAP, Inc.
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

package util

// RateLimit wraps a fix sized channel to control concurrency.
type RateLimit struct {
	capacity int
	token    chan struct{}
}

// NewRateLimit creates a limit controller with capacity n.
func NewRateLimit(n int) *RateLimit {
	return &RateLimit{
		capacity: n,
		token:    make(chan struct{}, n),
	}
}

// GetToken acquires a token.
func (r *RateLimit) GetToken(done <-chan struct{}) (exit bool) {
	select {
	case <-done:
		return true
	case r.token <- struct{}{}:
		return false
	}
}

// PutToken puts a token back.
func (r *RateLimit) PutToken() {
	select {
	case <-r.token:
	default:
		panic("put a redundant token")
	}
}

// GetCapacity returns the token capacity.
func (r *RateLimit) GetCapacity() int {
	return r.capacity
}
