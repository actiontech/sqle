// Copyright 2023 Huawei Technologies Co.,Ltd.
//
// Licensed to the Apache Software Foundation (ASF) under one
// or more contributor license agreements.  See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership.  The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package progress

import (
	"io"
)

const defaultProgressInterval = int64(102400)

func NewTeeReader(reader io.Reader, writer io.Writer, totalBytes int64, progressListener Listener, progressInterval int64) *TeeReader {
	if progressInterval <= 0 {
		progressInterval = defaultProgressInterval
	}

	return &TeeReader{
		reader:           reader,
		writer:           writer,
		progressInterval: progressInterval,
		totalBytes:       totalBytes,
		listener:         progressListener,
	}
}

type TeeReader struct {
	reader           io.Reader
	writer           io.Writer
	cacheBytes       int64
	progressInterval int64
	transferredBytes int64
	totalBytes       int64
	listener         Listener
}

func (r *TeeReader) Read(p []byte) (int, error) {
	if r.transferredBytes == 0 {
		event := NewEvent(TransferStartedEvent, r.transferredBytes, r.totalBytes, nil)
		r.listener.ProgressChanged(event)
	}

	n, err := r.reader.Read(p)
	if err != nil && err != io.EOF {
		event := NewEvent(TransferFailedEvent, r.transferredBytes, r.totalBytes, err)
		r.listener.ProgressChanged(event)
	}

	if n > 0 {
		n64 := int64(n)
		r.transferredBytes += n64
		if r.writer != nil {
			if n, err := r.writer.Write(p[:n]); err != nil {
				return n, err
			}
		}

		r.cacheBytes += n64
		if r.cacheBytes >= r.progressInterval || r.transferredBytes == r.totalBytes {
			r.cacheBytes = 0
			event := NewEvent(TransferDataEvent, r.transferredBytes, r.totalBytes, nil)
			r.listener.ProgressChanged(event)
		}
	}

	if err == io.EOF {
		r.cacheBytes = 0
		event := NewEvent(TransferCompletedEvent, r.transferredBytes, r.totalBytes, nil)
		r.listener.ProgressChanged(event)
	}

	return n, err
}

func (r *TeeReader) Close() error {
	if closer, ok := r.reader.(io.ReadCloser); ok {
		return closer.Close()
	}
	return nil
}
