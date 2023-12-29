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

package provider

import (
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"reflect"
)

const (
	basicCredentialType  = "basic"
	globalCredentialType = "global"

	credentialsAttr   = "Credentials"
	akAttr            = "AK"
	skAttr            = "SK"
	securityTokenAttr = "SecurityToken"
	idpIdAttr         = "IdpId"
	idTokenFileAttr   = "IdTokenFile"
	iamEndpointAttr   = "IamEndpoint"
)

type ICredentialProvider interface {
	GetCredentials() (auth.ICredential, error)
}

type commonAttrs struct {
	ak            string
	sk            string
	securityToken string
	idpId         string
	idTokenFile   string
	iamEndpoint   string
}

func fillCommonAttrs(builder interface{}, common commonAttrs) error {
	v := reflect.ValueOf(builder)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	v = v.FieldByName(credentialsAttr)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if common.iamEndpoint != "" {
		v.FieldByName(iamEndpointAttr).SetString(common.iamEndpoint)
	}
	if common.idpId != "" && common.idTokenFile != "" {
		v.FieldByName(idpIdAttr).SetString(common.idpId)
		v.FieldByName(idTokenFileAttr).SetString(common.idTokenFile)
		return nil
	} else if common.ak != "" && common.sk != "" {
		v.FieldByName(akAttr).SetString(common.ak)
		v.FieldByName(skAttr).SetString(common.sk)
		v.FieldByName(securityTokenAttr).SetString(common.securityToken)
		return nil
	}
	return sdkerr.NewCredentialsTypeError(fmt.Sprintf("%s&%s or %s&%s does not exist",
		akAttr, skAttr, idpIdAttr, idTokenFileAttr))
}
