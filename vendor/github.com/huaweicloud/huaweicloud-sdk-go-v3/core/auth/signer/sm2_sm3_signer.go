// Copyright 2023 Huawei Technologies Co.,Ltd.
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
	_ "crypto/ecdsa"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

const (
	sdkSm2Sm3 = "SDK-SM2-SM3"
)

var (
	curveSm2     = sm2.P256Sm2()
	sm2nMinusTwo = new(big.Int).Sub(new(big.Int).Set(curveSm2.Params().N), big.NewInt(2))
)

type SM2SM3Signer struct {
}

func (s SM2SM3Signer) Sign(req *request.DefaultHttpRequest, ak, sk string) (map[string]string, error) {
	err := checkAKSK(ak, sk)
	if err != nil {
		return nil, err
	}

	processContentHeader(req, xSdkContentSm3)
	originalHeaders := req.GetHeaderParams()
	t := extractTime(originalHeaders)
	headerDate := t.UTC().Format(BasicDateFormat)
	originalHeaders[HeaderXDate] = headerDate
	additionalHeaders := map[string]string{HeaderXDate: headerDate}

	signedHeaders := extractSignedHeaders(originalHeaders)

	cr, err := canonicalRequest(req, signedHeaders, xSdkContentSm3, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	sts, err := stringToSign(sdkSm2Sm3, cr, t, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	signingKey, err := s.GetSigningKey(ak, sk)
	if err != nil {
		return nil, err
	}

	sig, err := signStringToSign(sts, signingKey)
	if err != nil {
		return nil, err
	}

	additionalHeaders[HeaderAuthorization] = authHeaderValue(sdkSm2Sm3, sig, ak, signedHeaders)
	return additionalHeaders, nil
}

func (s SM2SM3Signer) GetSigningKey(ak, sk string) (ISigningKey, error) {
	privInt, err := derivePrivateInt(sdkSm2Sm3, ak, sk, sm2nMinusTwo, sm3HasherInst)
	if err != nil {
		return nil, err
	}

	return s.deriveSigningKey(privInt), nil
}

func (s SM2SM3Signer) deriveSigningKey(priv *big.Int) ISigningKey {
	privateKey := new(sm2.PrivateKey)
	privateKey.PublicKey.Curve = curveSm2
	privateKey.D = priv
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curveSm2.ScalarBaseMult(priv.Bytes())
	return SM2SigningKey{privateKey: privateKey}
}
