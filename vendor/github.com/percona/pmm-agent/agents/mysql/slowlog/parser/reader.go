// pmm-agent
// Copyright 2019 Percona LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

// ReaderMetrics contains Reader metrics.
type ReaderMetrics struct {
	InputSize int64
	InputPos  int64
}

// A Reader reads lines from the underlying source.
//
// Implementation should allow concurrent calls to different methods.
type Reader interface {
	// NextLine reads full lines from the underlying source and returns them (including the last '\n').
	// If the full line can't be read because of EOF, reader implementation may decide to return it,
	// or block and wait for new data to arrive. Other errors should be returned without blocking.
	// NextLine also should not block when the source is closed, but it may return buffered data while it has it.
	NextLine() (string, error)

	// Close closes the underlying source. A caller should continue to call NextLine until error is returned.
	Close() error

	// Metrics returns current metrics.
	Metrics() *ReaderMetrics
}
