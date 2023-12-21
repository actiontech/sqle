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

package auth

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"regexp"
	"strings"
)

const DefaultEndpointReg = "^[a-z][a-z0-9-]+(\\.[a-z]{2,}-[a-z]+-\\d{1,2})?\\.(my)?(huaweicloud|myhwclouds).(com|cn)"

type IDerivedCredential interface {
	ProcessDerivedAuthParams(derivedAuthServiceName, regionId string) ICredential
	IsDerivedAuth(httpRequest *request.DefaultHttpRequest) bool
	ICredential
}

func GetDefaultDerivedPredicate() func(*request.DefaultHttpRequest) bool {
	return func(httpRequest *request.DefaultHttpRequest) bool {
		matched, err := regexp.MatchString(DefaultEndpointReg, strings.Replace(httpRequest.GetEndpoint(), "https://", "", 1))
		if err != nil {
			return true
		} else {
			return !matched
		}
	}
}
