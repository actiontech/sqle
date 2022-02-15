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

// +build gofuzz

// See https://github.com/dvyukov/go-fuzz

package parser

import (
	"bufio"
	"bytes"
	"io"

	"github.com/percona/go-mysql/log"
)

type bytesReader struct {
	r *bufio.Reader
}

func newBytesReader(b []byte) (*bytesReader, error) {
	return &bytesReader{
		r: bufio.NewReader(bytes.NewReader(b)),
	}, nil
}

func (r *bytesReader) NextLine() (string, error) {
	return r.r.ReadString('\n')
}

func (r *bytesReader) Close() error {
	panic("not reached")
}

func (r *bytesReader) Metrics() *ReaderMetrics {
	panic("not reached")
}

func Fuzz(data []byte) int {
	r, err := newBytesReader(data)
	if err != nil {
		panic(err)
	}
	p := NewSlowLogParser(r, log.Options{})

	go p.Run()

	for p.Parse() != nil {
	}

	if p.Err() == io.EOF {
		return 1
	}
	return 0
}
