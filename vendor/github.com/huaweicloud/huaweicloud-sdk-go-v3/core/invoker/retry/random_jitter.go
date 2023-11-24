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

package retry

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"
)

// RandomJitter 全抖动指数退避 delay = random(0, Exponential)
type RandomJitter struct {
	*strategyBase
}

func NewRandomJitter() *RandomJitter {
	return &RandomJitter{strategyBase: newStrategyBase()}
}

func (r *RandomJitter) ComputeDelayBeforeNextRetry(retryTimes int32) int32 {
	expoDelay := getExponentialDelay(r.GetBaseDelay(), retryTimes)
	return utils.RandInt32(1, expoDelay)
}
