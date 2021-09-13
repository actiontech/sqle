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

import (
	"bufio"
	"io"
	"os"
	"sync"
)

// SimpleFileReader reads lines from the single file from the start until EOF.
type SimpleFileReader struct {
	// file Read/Close calls must be synchronized
	m sync.Mutex
	f *os.File
	r *bufio.Reader
}

// NewSimpleFileReader creates new SimpleFileReader.
func NewSimpleFileReader(filename string) (*SimpleFileReader, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return &SimpleFileReader{
		f: f,
		r: bufio.NewReader(f),
	}, nil
}

// NextLine implements Reader interface.
func (r *SimpleFileReader) NextLine() (string, error) {
	// TODO handle partial line reads as in ContinuousFileReader if needed

	r.m.Lock()
	l, err := r.r.ReadString('\n')
	r.m.Unlock()
	return l, err
}

// Close implements Reader interface.
func (r *SimpleFileReader) Close() error {
	r.m.Lock()
	defer r.m.Unlock()

	return r.f.Close()
}

// Metrics implements Reader interface.
func (r *SimpleFileReader) Metrics() *ReaderMetrics {
	r.m.Lock()
	defer r.m.Unlock()

	var m ReaderMetrics
	fi, err := r.f.Stat()
	if err == nil {
		m.InputSize = fi.Size()
	}
	pos, err := r.f.Seek(0, io.SeekCurrent)
	if err == nil {
		m.InputPos = pos
	}
	return &m
}

// check interfaces
var (
	_ Reader = (*SimpleFileReader)(nil)
)
