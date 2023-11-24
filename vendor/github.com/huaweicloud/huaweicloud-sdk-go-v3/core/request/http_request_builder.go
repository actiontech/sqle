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

package request

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/signer/algorithm"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/progress"
	"reflect"
	"strings"
)

type HttpRequestBuilder struct {
	httpRequest *DefaultHttpRequest
}

func NewHttpRequestBuilder() *HttpRequestBuilder {
	httpRequest := &DefaultHttpRequest{
		queryParams:          make(map[string]interface{}),
		headerParams:         make(map[string]string),
		pathParams:           make(map[string]string),
		autoFilledPathParams: make(map[string]string),
		formParams:           make(map[string]def.FormData),
	}
	httpRequestBuilder := &HttpRequestBuilder{
		httpRequest: httpRequest,
	}
	return httpRequestBuilder
}

func (builder *HttpRequestBuilder) WithEndpoint(endpoint string) *HttpRequestBuilder {
	builder.httpRequest.endpoint = endpoint
	return builder
}

func (builder *HttpRequestBuilder) WithPath(path string) *HttpRequestBuilder {
	builder.httpRequest.path = path
	return builder
}

func (builder *HttpRequestBuilder) WithMethod(method string) *HttpRequestBuilder {
	builder.httpRequest.method = method
	return builder
}

func (builder *HttpRequestBuilder) WithSigningAlgorithm(signingAlgorithm algorithm.SigningAlgorithm) *HttpRequestBuilder {
	builder.httpRequest.signingAlgorithm = signingAlgorithm
	return builder
}

func (builder *HttpRequestBuilder) AddQueryParam(key string, value interface{}) *HttpRequestBuilder {
	builder.httpRequest.queryParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddPathParam(key string, value string) *HttpRequestBuilder {
	builder.httpRequest.pathParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddAutoFilledPathParam(key string, value string) *HttpRequestBuilder {
	builder.httpRequest.autoFilledPathParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddHeaderParam(key string, value string) *HttpRequestBuilder {
	builder.httpRequest.headerParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) AddFormParam(key string, value def.FormData) *HttpRequestBuilder {
	builder.httpRequest.formParams[key] = value
	return builder
}

func (builder *HttpRequestBuilder) WithBody(kind string, body interface{}) *HttpRequestBuilder {
	if kind == "multipart" {
		v := reflect.ValueOf(body)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		t := reflect.TypeOf(body)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}

		fieldNum := t.NumField()
		for i := 0; i < fieldNum; i++ {
			jsonTag := t.Field(i).Tag.Get("json")
			if jsonTag != "" {
				if v.FieldByName(t.Field(i).Name).IsNil() && strings.Contains(jsonTag, "omitempty") {
					continue
				}
				builder.AddFormParam(strings.Split(jsonTag, ",")[0], v.FieldByName(t.Field(i).Name).Interface().(def.FormData))
			} else {
				builder.AddFormParam(t.Field(i).Name, v.FieldByName(t.Field(i).Name).Interface().(def.FormData))
			}
		}
	} else {
		builder.httpRequest.body = body
	}

	return builder
}

func (builder *HttpRequestBuilder) WithProgressListener(progressListener progress.Listener) *HttpRequestBuilder {
	builder.httpRequest.progressListener = progressListener
	return builder
}

func (builder *HttpRequestBuilder) WithProgressInterval(progressInterval int64) *HttpRequestBuilder {
	builder.httpRequest.progressInterval = progressInterval
	return builder
}

func (builder *HttpRequestBuilder) Build() *DefaultHttpRequest {
	return builder.httpRequest.fillParamsInPath()
}
