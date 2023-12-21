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

package signer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"golang.org/x/crypto/hkdf"
	"io"
	"strings"
	"time"
)

const (
	DerivedDateFormat = "20060102"
	AlgorithmV11      = "V11-HMAC-SHA256"
)

type DerivedSigner struct {
}

// Sign SignRequest set Authorization header
func (s DerivedSigner) Sign(r *request.DefaultHttpRequest, ak, sk, derivedAuthServiceName, regionId string) (map[string]string, error) {
	err := checkAKSK(ak, sk)
	if err != nil {
		return nil, err
	}
	if derivedAuthServiceName == "" {
		return nil, errors.New("DerivedAuthServiceName is required in credentials when using derived auth")
	}
	if regionId == "" {
		return nil, errors.New("RegionId is required in credentials when using derived auth")
	}

	originalHeaders := r.GetHeaderParams()
	t := extractTime(originalHeaders)
	headerDate := t.UTC().Format(BasicDateFormat)
	originalHeaders[HeaderXDate] = headerDate
	additionalHeaders := map[string]string{HeaderXDate: headerDate}

	signedHeaders := extractSignedHeaders(originalHeaders)
	cr, err := canonicalRequest(r, signedHeaders, xSdkContentSha256, sha256HasherInst)
	if err != nil {
		return nil, err
	}
	info := t.UTC().Format(DerivedDateFormat) + "/" + regionId + "/" + derivedAuthServiceName
	sts, err := s.stringToSign(cr, info, t)
	if err != nil {
		return nil, err
	}
	derivationKey, err := s.getDerivationKey(ak, sk, info)
	if err != nil {
		return nil, err
	}
	sig, err := s.signStringToSign(sts, []byte(derivationKey))
	if err != nil {
		return nil, err
	}
	additionalHeaders[HeaderAuthorization] = s.authHeaderValue(sig, ak, info, signedHeaders)
	return additionalHeaders, nil
}

// signStringToSign Create the Signature.
func (s DerivedSigner) signStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hm, err := sha256HasherInst.hmac([]byte(stringToSign), signingKey)
	return fmt.Sprintf("%x", hm), err
}

// authHeaderValue Get the finalized value for the "Authorization" header.
// The signature parameter is the output from signStringToSign
func (s DerivedSigner) authHeaderValue(signature, accessKey, info string, signedHeaders []string) string {
	return fmt.Sprintf("%s Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		AlgorithmV11,
		accessKey,
		info,
		strings.Join(signedHeaders, ";"),
		signature)
}

// getDerivationKey Get the derivation key for derived credential.
func (s DerivedSigner) getDerivationKey(accessKey, secretKey, info string) (string, error) {
	hash := sha256.New
	derivationKeyReader := hkdf.New(hash, []byte(secretKey), []byte(accessKey), []byte(info))
	derivationKey := make([]byte, 32)
	_, err := io.ReadFull(derivationKeyReader, derivationKey)
	return hex.EncodeToString(derivationKey), err
}

// stringToSign Create a "String to Sign".
func (s DerivedSigner) stringToSign(canonicalRequest, info string, t time.Time) (string, error) {
	canonicalRequestHash, err := sha256HasherInst.hashHexString([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s", AlgorithmV11, t.UTC().Format(BasicDateFormat), info, canonicalRequestHash), nil
}
