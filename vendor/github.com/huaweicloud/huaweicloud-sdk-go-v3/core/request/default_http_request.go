// Copyright 2020 Huawei Technologies Co.,Ltd.
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
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/signer/algorithm"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/progress"
	"go.mongodb.org/mongo-driver/bson"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/def"
)

type DefaultHttpRequest struct {
	endpoint string
	path     string
	method   string

	queryParams  map[string]interface{}
	pathParams   map[string]string
	headerParams map[string]string
	formParams   map[string]def.FormData
	body         interface{}

	autoFilledPathParams map[string]string
	progressListener     progress.Listener
	progressInterval     int64
	signingAlgorithm     algorithm.SigningAlgorithm
}

func (httpRequest *DefaultHttpRequest) fillParamsInPath() *DefaultHttpRequest {
	for key, value := range httpRequest.pathParams {
		httpRequest.path = strings.ReplaceAll(httpRequest.path, "{"+key+"}", value)
	}
	for key, value := range httpRequest.autoFilledPathParams {
		httpRequest.path = strings.ReplaceAll(httpRequest.path, "{"+key+"}", value)
	}
	return httpRequest
}

func (httpRequest *DefaultHttpRequest) Builder() *HttpRequestBuilder {
	httpRequestBuilder := HttpRequestBuilder{httpRequest: httpRequest}
	return &httpRequestBuilder
}

func (httpRequest *DefaultHttpRequest) GetEndpoint() string {
	return httpRequest.endpoint
}

func (httpRequest *DefaultHttpRequest) GetPath() string {
	return httpRequest.path
}

func (httpRequest *DefaultHttpRequest) GetSigningAlgorithm() algorithm.SigningAlgorithm {
	return httpRequest.signingAlgorithm
}

func (httpRequest *DefaultHttpRequest) GetMethod() string {
	return httpRequest.method
}

func (httpRequest *DefaultHttpRequest) GetQueryParams() map[string]interface{} {
	return httpRequest.queryParams
}

func (httpRequest *DefaultHttpRequest) GetHeaderParams() map[string]string {
	return httpRequest.headerParams
}

func (httpRequest *DefaultHttpRequest) GetPathPrams() map[string]string {
	return httpRequest.pathParams
}

func (httpRequest *DefaultHttpRequest) GetFormPrams() map[string]def.FormData {
	return httpRequest.formParams
}

func (httpRequest *DefaultHttpRequest) GetBody() interface{} {
	return httpRequest.body
}

func (httpRequest *DefaultHttpRequest) GetBodyToBytes() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}

	if httpRequest.body != nil {
		v := reflect.ValueOf(httpRequest.body)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		if v.Kind() == reflect.String {
			buf.WriteString(v.Interface().(string))
		} else {
			var err error
			if httpRequest.headerParams["Content-Type"] == "application/xml" {
				encoder := xml.NewEncoder(buf)
				err = encoder.Encode(httpRequest.body)
			} else if httpRequest.headerParams["Content-Type"] == "application/bson" {
				buffer, err := bson.Marshal(httpRequest.body)
				if err != nil {
					return nil, err
				}
				buf.Write(buffer)
			} else {
				encoder := json.NewEncoder(buf)
				encoder.SetEscapeHTML(false)
				err = encoder.Encode(httpRequest.body)
			}
			if err != nil {
				return nil, err
			}
		}
	}

	return buf, nil
}

func (httpRequest *DefaultHttpRequest) GetProgressListener() progress.Listener {
	return httpRequest.progressListener
}

func (httpRequest *DefaultHttpRequest) GetProgressInterval() int64 {
	return httpRequest.progressInterval
}

func (httpRequest *DefaultHttpRequest) AddQueryParam(key string, value string) {
	httpRequest.queryParams[key] = value
}

func (httpRequest *DefaultHttpRequest) AddPathParam(key string, value string) {
	httpRequest.pathParams[key] = value
}

func (httpRequest *DefaultHttpRequest) AddHeaderParam(key string, value string) {
	httpRequest.headerParams[key] = value
}

func (httpRequest *DefaultHttpRequest) AddFormParam(key string, value def.FormData) {
	httpRequest.formParams[key] = value
}

func (httpRequest *DefaultHttpRequest) ConvertRequest() (*http.Request, error) {
	t := reflect.TypeOf(httpRequest.body)
	if t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	var req *http.Request
	var err error
	if httpRequest.body != nil && t != nil && t.Name() == "File" {
		req, err = httpRequest.convertStreamBody(err, req)
		if err != nil {
			return nil, err
		}
	} else if len(httpRequest.GetFormPrams()) != 0 {
		req, err = httpRequest.covertFormBody()
		if err != nil {
			return nil, err
		}
	} else {
		var buf *bytes.Buffer

		buf, err = httpRequest.GetBodyToBytes()
		if err != nil {
			return nil, err
		}

		req, err = http.NewRequest(httpRequest.GetMethod(), httpRequest.GetEndpoint(), buf)
		if err != nil {
			return nil, err
		}
	}

	httpRequest.fillPath(req)
	httpRequest.fillQueryParams(req)
	httpRequest.fillHeaderParams(req)

	return req, nil
}

func (httpRequest *DefaultHttpRequest) covertFormBody() (*http.Request, error) {
	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	sortedKeys := make([]string, 0, len(httpRequest.GetFormPrams()))
	for k, v := range httpRequest.GetFormPrams() {
		if _, ok := v.(*def.FilePart); ok {
			sortedKeys = append(sortedKeys, k)
		} else {
			sortedKeys = append([]string{k}, sortedKeys...)
		}
	}

	for _, k := range sortedKeys {
		if err := httpRequest.GetFormPrams()[k].Write(bodyWriter, k); err != nil {
			return nil, err
		}
	}

	contentType := bodyWriter.FormDataContentType()
	if err := bodyWriter.Close(); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(httpRequest.GetMethod(), httpRequest.GetEndpoint(), bodyBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-type", contentType)
	return req, nil
}

func (httpRequest *DefaultHttpRequest) getContentLength() int64 {
	contentLength := int64(-1)
	if value, ok := httpRequest.GetHeaderParams()["Content-Length"]; ok {
		parseInt, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			contentLength = parseInt
		}
	}

	if contentLength == -1 {
		if value, ok := httpRequest.GetBody().(os.File); ok {
			if stat, err := value.Stat(); err == nil {
				contentLength = stat.Size()
			}
		}
	}

	return contentLength
}

func (httpRequest *DefaultHttpRequest) convertStreamBody(err error, req *http.Request) (*http.Request, error) {
	if f, ok := httpRequest.body.(os.File); !ok {
		return nil, errors.New("failed to get stream request body")
	} else {
		var reader io.Reader
		if httpRequest.progressListener != nil {
			reader = progress.NewTeeReader(&f, nil, httpRequest.getContentLength(), httpRequest.progressListener, httpRequest.progressInterval)
		} else {
			reader = &f
		}

		req, err = http.NewRequest(httpRequest.GetMethod(), httpRequest.GetEndpoint(), reader)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

func (httpRequest *DefaultHttpRequest) fillHeaderParams(req *http.Request) {
	if len(httpRequest.GetHeaderParams()) == 0 {
		return
	}

	for key, value := range httpRequest.GetHeaderParams() {
		if strings.EqualFold(key, "Content-type") && req.Header.Get("Content-type") != "" {
			continue
		}
		req.Header.Add(key, value)
	}
}

func (httpRequest *DefaultHttpRequest) fillQueryParams(req *http.Request) {
	if len(httpRequest.GetQueryParams()) == 0 {
		return
	}

	q := req.URL.Query()
	for key, value := range httpRequest.GetQueryParams() {
		valueWithType, ok := value.(reflect.Value)
		if !ok {
			continue
		}

		if valueWithType.Kind() == reflect.Slice {
			params := httpRequest.CanonicalSliceQueryParamsToMulti(valueWithType)
			for _, param := range params {
				q.Add(key, param)
			}
		} else if valueWithType.Kind() == reflect.Map {
			params := httpRequest.CanonicalMapQueryParams(key, valueWithType)
			for _, param := range params {
				for k, v := range param {
					q.Add(k, v)
				}
			}
		} else {
			q.Add(key, httpRequest.CanonicalStringQueryParams(valueWithType))
		}
	}

	req.URL.RawQuery = strings.ReplaceAll(strings.ReplaceAll(strings.Trim(q.Encode(), "="), "=&", "&"), "+", "%20")
}

func (httpRequest *DefaultHttpRequest) CanonicalStringQueryParams(value reflect.Value) string {
	return fmt.Sprintf("%v", value)
}

func (httpRequest *DefaultHttpRequest) CanonicalSliceQueryParamsToMulti(value reflect.Value) []string {
	params := make([]string, 0)

	for i := 0; i < value.Len(); i++ {
		if value.Index(i).Kind() == reflect.Struct {
			methodByName := value.Index(i).MethodByName("Value")
			if methodByName.IsValid() {
				value := converter.ConvertInterfaceToString(methodByName.Call([]reflect.Value{})[0].Interface())
				params = append(params, value)
			} else {
				v, e := json.Marshal(value.Interface())
				if e == nil {
					if strings.HasPrefix(string(v), "\"") {
						params = append(params, strings.Trim(string(v), "\""))
					} else {
						params = append(params, string(v))
					}
				}
			}
		} else {
			params = append(params, httpRequest.CanonicalStringQueryParams(value.Index(i)))
		}
	}

	return params
}

func (httpRequest *DefaultHttpRequest) CanonicalMapQueryParams(key string, value reflect.Value) []map[string]string {
	queryParams := make([]map[string]string, 0)

	for _, k := range value.MapKeys() {
		if value.MapIndex(k).Kind() == reflect.Struct {
			v, e := json.Marshal(value.Interface())
			if e == nil {
				if strings.HasPrefix(string(v), "\"") {
					queryParams = append(queryParams, map[string]string{
						key: strings.Trim(string(v), "\""),
					})
				} else {
					queryParams = append(queryParams, map[string]string{
						key: string(v),
					})
				}
			}
		} else if value.MapIndex(k).Kind() == reflect.Slice {
			params := httpRequest.CanonicalSliceQueryParamsToMulti(value.MapIndex(k))
			if len(params) == 0 {
				queryParams = append(queryParams, map[string]string{
					fmt.Sprintf("%s[%s]", key, k): "",
				})
				continue
			}
			for _, paramValue := range httpRequest.CanonicalSliceQueryParamsToMulti(value.MapIndex(k)) {
				queryParams = append(queryParams, map[string]string{
					fmt.Sprintf("%s[%s]", key, k): paramValue,
				})
			}
		} else if value.MapIndex(k).Kind() == reflect.Map {
			queryParams = append(queryParams, httpRequest.CanonicalMapQueryParams(fmt.Sprintf("%s[%s]", key, k), value.MapIndex(k))...)
		} else {
			queryParams = append(queryParams, map[string]string{
				fmt.Sprintf("%s[%s]", key, k): httpRequest.CanonicalStringQueryParams(value.MapIndex(k)),
			})
		}
	}

	return queryParams
}

func (httpRequest *DefaultHttpRequest) fillPath(req *http.Request) {
	if "" != httpRequest.GetPath() {
		req.URL.Path = httpRequest.GetPath()
	}
}
