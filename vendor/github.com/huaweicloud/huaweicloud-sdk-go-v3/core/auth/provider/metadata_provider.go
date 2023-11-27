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
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"strings"
)

type MetadataCredentialProvider struct {
	credentialType string
}

// NewMetadataCredentialProvider return a metadata credential provider
// Supported credential types: basic, global
func NewMetadataCredentialProvider(credentialType string) *MetadataCredentialProvider {
	return &MetadataCredentialProvider{credentialType: strings.ToLower(credentialType)}
}

// BasicCredentialMetadataProvider return a metadata provider for basic.Credentials
func BasicCredentialMetadataProvider() *MetadataCredentialProvider {
	return NewMetadataCredentialProvider(basicCredentialType)
}

// GlobalCredentialMetadataProvider return a metadata provider for global.Credentials
func GlobalCredentialMetadataProvider() *MetadataCredentialProvider {
	return NewMetadataCredentialProvider(globalCredentialType)
}

// GetCredentials get basic.Credentials or global.Credentials from the instance's metadata
func (p *MetadataCredentialProvider) GetCredentials() (auth.ICredential, error) {
	if p.credentialType == "" {
		return nil, sdkerr.NewCredentialsTypeError("credential type is empty")
	}

	if strings.HasPrefix(p.credentialType, basicCredentialType) {
		credentials := basic.NewCredentialsBuilder().Build()
		err := credentials.UpdateSecurityTokenFromMetadata()
		if err != nil {
			return nil, err
		}
		return credentials, nil
	} else if strings.HasPrefix(p.credentialType, globalCredentialType) {
		credentials := global.NewCredentialsBuilder().Build()
		err := credentials.UpdateSecurityTokenFromMetadata()
		if err != nil {
			return nil, err
		}
		return credentials, nil
	}

	return nil, sdkerr.NewCredentialsTypeError("unsupported credential type: " + p.credentialType)
}
