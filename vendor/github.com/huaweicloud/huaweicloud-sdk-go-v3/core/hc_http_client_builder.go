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

package core

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/provider"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/config"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/impl"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/region"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"reflect"
	"strings"
)

type HcHttpClientBuilder struct {
	CredentialsType        []string
	derivedAuthServiceName string
	credentials            auth.ICredential
	endpoints              []string
	httpConfig             *config.HttpConfig
	region                 *region.Region
	errorHandler           sdkerr.ErrorHandler
}

func NewHcHttpClientBuilder() *HcHttpClientBuilder {
	hcHttpClientBuilder := &HcHttpClientBuilder{
		CredentialsType: []string{"basic.Credentials"},
		errorHandler:    sdkerr.DefaultErrorHandler{},
	}
	return hcHttpClientBuilder
}

func (builder *HcHttpClientBuilder) WithCredentialsType(credentialsType string) *HcHttpClientBuilder {
	builder.CredentialsType = strings.Split(credentialsType, ",")
	return builder
}

func (builder *HcHttpClientBuilder) WithDerivedAuthServiceName(derivedAuthServiceName string) *HcHttpClientBuilder {
	builder.derivedAuthServiceName = derivedAuthServiceName
	return builder
}

// Deprecated: As of 0.1.27, because of the support of the multi-endpoint feature, use WithEndpoints instead
func (builder *HcHttpClientBuilder) WithEndpoint(endpoint string) *HcHttpClientBuilder {
	return builder.WithEndpoints([]string{endpoint})
}

func (builder *HcHttpClientBuilder) WithEndpoints(endpoints []string) *HcHttpClientBuilder {
	builder.endpoints = endpoints
	return builder
}

func (builder *HcHttpClientBuilder) WithRegion(region *region.Region) *HcHttpClientBuilder {
	builder.region = region
	return builder
}

func (builder *HcHttpClientBuilder) WithHttpConfig(httpConfig *config.HttpConfig) *HcHttpClientBuilder {
	builder.httpConfig = httpConfig
	return builder
}

func (builder *HcHttpClientBuilder) WithCredential(iCredential auth.ICredential) *HcHttpClientBuilder {
	builder.credentials = iCredential
	return builder
}

func (builder *HcHttpClientBuilder) WithErrorHandler(errorHandler sdkerr.ErrorHandler) *HcHttpClientBuilder {
	builder.errorHandler = errorHandler
	return builder
}

func (builder *HcHttpClientBuilder) Build() *HcHttpClient {
	if builder.httpConfig == nil {
		builder.httpConfig = config.DefaultHttpConfig()
	}

	defaultHttpClient := impl.NewDefaultHttpClient(builder.httpConfig)

	if builder.credentials == nil {
		p := provider.DefaultCredentialProviderChain(builder.CredentialsType[0])
		credentials, err := p.GetCredentials()
		if err != nil {
			panic(err)
		}
		builder.credentials = credentials
	}

	t := reflect.TypeOf(builder.credentials)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	givenCredentialsType := t.String()
	match := false
	for _, credentialsType := range builder.CredentialsType {
		if credentialsType == givenCredentialsType {
			match = true
			break
		}
	}
	if !match {
		panic(fmt.Sprintf("Need credential type is %s, actually is %s", builder.CredentialsType, givenCredentialsType))
	}

	if builder.region != nil {
		builder.endpoints = builder.region.Endpoints
		builder.credentials.ProcessAuthParams(defaultHttpClient, builder.region.Id)

		if credential, ok := builder.credentials.(auth.IDerivedCredential); ok {
			credential.ProcessDerivedAuthParams(builder.derivedAuthServiceName, builder.region.Id)
		}
	}

	for index, endpoint := range builder.endpoints {
		if !strings.HasPrefix(endpoint, "http") {
			builder.endpoints[index] = "https://" + endpoint
		}
	}

	hcHttpClient := NewHcHttpClient(defaultHttpClient).
		WithEndpoints(builder.endpoints).
		WithCredential(builder.credentials).
		WithErrorHandler(builder.errorHandler)
	return hcHttpClient
}
