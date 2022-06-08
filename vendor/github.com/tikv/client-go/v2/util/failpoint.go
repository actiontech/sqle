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
// https://github.com/pingcap/tidb/tree/cc5e161ac06827589c4966674597c137cc9e809c/store/tikv/util/failpoint.go
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

import (
	"errors"

	"github.com/pingcap/failpoint"
)

const failpointPrefix = "tikvclient/"

var (
	failpointsEnabled bool
	errDisabled       = errors.New("failpoints are disabled")
)

// EnableFailpoints enables use of failpoints.
// It should be called before using client to avoid data race.
func EnableFailpoints() {
	failpointsEnabled = true
}

// EvalFailpoint injects code for testing. It is used to replace `failpoint.Inject`
// to make it possible to be used in a library.
func EvalFailpoint(name string) (interface{}, error) {
	if !failpointsEnabled {
		return nil, errDisabled
	}
	return failpoint.Eval(failpointPrefix + name)
}
