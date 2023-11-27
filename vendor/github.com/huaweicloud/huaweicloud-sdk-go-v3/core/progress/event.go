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

type EventType int

const (
	TransferStartedEvent EventType = iota
	TransferDataEvent
	TransferCompletedEvent
	TransferFailedEvent
)

type Event struct {
	Type             EventType
	TransferredBytes int64
	TotalBytes       int64
	Err              error
}

func NewEvent(eventType EventType, transferredBytes int64, totalBytes int64, err error) *Event {
	return &Event{
		Type:             eventType,
		TransferredBytes: transferredBytes,
		TotalBytes:       totalBytes,
		Err:              err,
	}
}
