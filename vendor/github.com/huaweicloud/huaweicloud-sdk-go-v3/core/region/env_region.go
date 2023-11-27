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

package region

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	regionEnvPrefix = "HUAWEICLOUD_SDK_REGION"
)

var (
	envOnce  sync.Once
	envCache *EnvCache
)

type EnvCache struct {
	value map[string]*Region
}

func getEnvCache() *EnvCache {
	envOnce.Do(func() {
		envCache = &EnvCache{value: make(map[string]*Region)}
	})

	return envCache
}

type EnvProvider struct {
	serviceName string
}

func NewEnvProvider(serviceName string) *EnvProvider {
	return &EnvProvider{serviceName: strings.ToUpper(serviceName)}
}

func (p *EnvProvider) GetRegion(regionId string) *Region {
	if reg, ok := getEnvCache().value[p.serviceName+regionId]; ok {
		return reg
	}

	envName := fmt.Sprintf("%s_%s_%s", regionEnvPrefix, p.serviceName, strings.ToUpper(strings.Replace(regionId, "-", "_", -1)))
	endpoint := os.Getenv(envName)
	if endpoint == "" {
		return nil
	}

	endpoints := strings.Split(endpoint, ",")
	reg := NewRegion(regionId, endpoints...)
	getEnvCache().value[p.serviceName+regionId] = reg
	return reg
}
