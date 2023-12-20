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
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/global"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdkerr"
	"gopkg.in/ini.v1"
	"os"
	"path/filepath"
	"strings"
)

const (
	credentialsFileEnvName = "HUAWEICLOUD_SDK_CREDENTIALS_FILE"
	defaultDir             = ".huaweicloud"
	defaultFile            = "credentials"

	akName            = "ak"
	skName            = "sk"
	projectIdName     = "project_id"
	domainIdName      = "domain_id"
	securityTokenName = "security_token"
	iamEndpointName   = "iam_endpoint"
	idpIdName         = "idp_id"
	idTokenFileName   = "id_token_file"
)

type ProfileCredentialProvider struct {
	credentialType string
}

// NewProfileCredentialProvider return a profile credential provider
// Supported credential types: basic, global
func NewProfileCredentialProvider(credentialType string) *ProfileCredentialProvider {
	return &ProfileCredentialProvider{credentialType: strings.ToLower(credentialType)}
}

// BasicCredentialProfileProvider return a profile provider for basic.Credentials
func BasicCredentialProfileProvider() *ProfileCredentialProvider {
	return NewProfileCredentialProvider(basicCredentialType)
}

// GlobalCredentialProfileProvider return a profile provider for global.Credentials
func GlobalCredentialProfileProvider() *ProfileCredentialProvider {
	return NewProfileCredentialProvider(globalCredentialType)
}

// GetCredentials get basic.Credentials or global.Credentials from profile
func (p *ProfileCredentialProvider) GetCredentials() (auth.ICredential, error) {
	filePath, err := getCredentialsFilePath()
	if err != nil {
		return nil, err
	}
	file, err := ini.Load(filePath)
	if err != nil {
		return nil, err
	}

	section := file.Section(p.credentialType)
	if section == nil {
		return nil, sdkerr.NewCredentialsTypeError(fmt.Sprintf("credential type '%s' does not exist in '%s'", p.credentialType, filePath))
	}

	if strings.HasPrefix(p.credentialType, basicCredentialType) {
		builder := basic.NewCredentialsBuilder().WithProjectId(section.Key(projectIdName).String())
		err := fillCommonAttrs(builder, getCommonAttrsFromProfile(section))
		if err != nil {
			return nil, err
		}
		return builder.Build(), nil
	} else if strings.HasPrefix(p.credentialType, globalCredentialType) {
		builder := global.NewCredentialsBuilder().WithDomainId(section.Key(domainIdName).String())
		err := fillCommonAttrs(builder, getCommonAttrsFromProfile(section))
		if err != nil {
			return nil, err
		}
		return builder.Build(), nil
	}
	return nil, sdkerr.NewCredentialsTypeError("unsupported credential type: " + p.credentialType)
}

func getCredentialsFilePath() (string, error) {
	if path := os.Getenv(credentialsFileEnvName); path != "" {
		return path, nil
	}

	dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, defaultDir, defaultFile), nil
}

func getCommonAttrsFromProfile(section *ini.Section) commonAttrs {
	return commonAttrs{
		ak:            section.Key(akName).String(),
		sk:            section.Key(skName).String(),
		securityToken: section.Key(securityTokenName).String(),
		idpId:         section.Key(idpIdName).String(),
		idTokenFile:   section.Key(idTokenFileName).String(),
		iamEndpoint:   section.Key(iamEndpointName).String(),
	}
}
