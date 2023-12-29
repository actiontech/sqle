// Copyright 2022 Huawei Technologies Co.,Ltd.
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

package invoker

import (
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/exchange"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/invoker/retry"
)

type RetryChecker func(interface{}, error) bool

type BaseInvoker struct {
	Exchange *exchange.SdkExchange

	client  *core.HcHttpClient
	request interface{}
	meta    *def.HttpRequestDef
	headers map[string]string

	retryTimes      int
	retryChecker    RetryChecker
	backoffStrategy retry.Strategy
}

func NewBaseInvoker(client *core.HcHttpClient, request interface{}, meta *def.HttpRequestDef) *BaseInvoker {
	exch := &exchange.SdkExchange{
		ApiReference: &exchange.ApiReference{
			Method: meta.Method,
			Path:   meta.Path,
		},
		Attributes: make(map[string]interface{}),
	}

	return &BaseInvoker{
		Exchange: exch,
		client:   client,
		request:  request,
		meta:     meta,
		headers:  make(map[string]string),
	}
}

func (b *BaseInvoker) ReplaceCredentialWhen(fun func(auth.ICredential) auth.ICredential) *BaseInvoker {
	b.client.WithCredential(fun(b.client.GetCredential()))
	return b
}

func (b *BaseInvoker) AddHeader(headers map[string]string) *BaseInvoker {
	b.headers = headers
	return b
}

func (b *BaseInvoker) WithRetry(retryTimes int, checker RetryChecker, backoffStrategy retry.Strategy) *BaseInvoker {
	b.retryTimes = retryTimes
	b.retryChecker = checker
	b.backoffStrategy = backoffStrategy
	return b
}

func (b *BaseInvoker) Invoke() (interface{}, error) {
	if b.retryTimes != 0 && b.retryChecker != nil {
		var execTimes int
		var resp interface{}
		var err error
		for {
			if execTimes == b.retryTimes {
				break
			}
			resp, err = b.client.PreInvoke(b.headers).SyncInvoke(b.request, b.meta, b.Exchange)
			execTimes += 1

			if b.retryChecker(resp, err) {
				time.Sleep(time.Duration(b.backoffStrategy.ComputeDelayBeforeNextRetry(int32(execTimes))) * time.Millisecond)
			} else {
				break
			}
		}
		return resp, err
	} else {
		return b.client.PreInvoke(b.headers).SyncInvoke(b.request, b.meta, b.Exchange)
	}
}
