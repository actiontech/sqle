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
// https://github.com/pingcap/tidb/tree/cc5e161ac06827589c4966674597c137cc9e809c/store/tikv/util/ts_set.go
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

import "sync"

// TSSet is a set of timestamps.
type TSSet struct {
	sync.RWMutex
	m map[uint64]struct{}
}

// Put puts timestamps into the map.
func (s *TSSet) Put(tss ...uint64) {
	s.Lock()
	defer s.Unlock()

	// Lazy initialization..
	// Most of the time, there is no transaction lock conflict.
	// So allocate this in advance is unnecessary and bad for performance.
	if s.m == nil {
		s.m = make(map[uint64]struct{}, 5)
	}

	for _, ts := range tss {
		s.m[ts] = struct{}{}
	}
}

// GetAll returns all timestamps in the set.
func (s *TSSet) GetAll() []uint64 {
	s.RLock()
	defer s.RUnlock()
	if len(s.m) == 0 {
		return nil
	}
	ret := make([]uint64, 0, len(s.m))
	for ts := range s.m {
		ret = append(ret, ts)
	}
	return ret
}
