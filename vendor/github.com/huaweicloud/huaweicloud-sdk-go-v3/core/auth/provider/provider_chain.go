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
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"strings"
)

type CredentialProviderChain struct {
	providers []ICredentialProvider
}

// NewCredentialProviderChain return a credential provider chain
func NewCredentialProviderChain(providers []ICredentialProvider) *CredentialProviderChain {
	return &CredentialProviderChain{providers: providers}
}

// DefaultCredentialProviderChain return a default credential provider chain
// Supported credential types: basic, global
// Default order: environment variables -> profile -> metadata
func DefaultCredentialProviderChain(credentialType string) *CredentialProviderChain {
	providers := []ICredentialProvider{
		NewEnvCredentialProvider(credentialType),
		NewProfileCredentialProvider(credentialType),
		NewMetadataCredentialProvider(credentialType),
	}
	return NewCredentialProviderChain(providers)
}

// BasicCredentialProviderChain return a provider chain for basic.Credentials
func BasicCredentialProviderChain() *CredentialProviderChain {
	providers := []ICredentialProvider{
		BasicCredentialEnvProvider(),
		BasicCredentialProfileProvider(),
		BasicCredentialMetadataProvider(),
	}
	return NewCredentialProviderChain(providers)
}

// GlobalCredentialProviderChain return a provider chain for global.Credentials
func GlobalCredentialProviderChain() *CredentialProviderChain {
	providers := []ICredentialProvider{
		GlobalCredentialEnvProvider(),
		GlobalCredentialProfileProvider(),
		GlobalCredentialMetadataProvider(),
	}
	return NewCredentialProviderChain(providers)
}

// GetCredentials get basic.Credentials or global.Credentials in providers
// In the order of providers, return the first found credentials
// If credentials not found in every providers, return a error of all providers
func (p *CredentialProviderChain) GetCredentials() (auth.ICredential, error) {
	var errs []string
	for _, provider := range p.providers {
		credential, err := provider.GetCredentials()
		if err == nil {
			return credential, nil
		}
		errs = append(errs, err.Error())
	}
	return nil, sdkerr.NewCredentialsTypeError("unable to get credential in providers:\n" + strings.Join(errs, "\n"))
}
