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
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/hex"
	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"math/big"
)

const (
	sdkEcdsaP256Sha256 = "SDK-ECDSA-P256-SHA256"
)

var (
	one           = big.NewInt(1)
	curveP256     = elliptic.P256()
	p256nMinusTwo = new(big.Int).Sub(new(big.Int).Set(curveP256.Params().N), big.NewInt(2))
)

type P256SHA256Signer struct {
}

func (s P256SHA256Signer) Sign(req *request.DefaultHttpRequest, ak, sk string) (map[string]string, error) {
	err := checkAKSK(ak, sk)
	if err != nil {
		return nil, err
	}

	processContentHeader(req, xSdkContentSha256)
	originalHeaders := req.GetHeaderParams()
	t := extractTime(originalHeaders)
	headerDate := t.UTC().Format(BasicDateFormat)
	originalHeaders[HeaderXDate] = headerDate
	additionalHeaders := map[string]string{HeaderXDate: headerDate}

	signedHeaders := extractSignedHeaders(originalHeaders)

	cr, err := canonicalRequest(req, signedHeaders, xSdkContentSha256, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	sts, err := stringToSign(sdkEcdsaP256Sha256, cr, t, sha256HasherInst)
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

	additionalHeaders[HeaderAuthorization] = authHeaderValue(sdkEcdsaP256Sha256, sig, ak, signedHeaders)
	return additionalHeaders, nil
}

// GetSigningKey get the derived key based on ak and sk.
func (s P256SHA256Signer) GetSigningKey(ak, sk string) (ISigningKey, error) {
	privInt, err := derivePrivateInt(sdkEcdsaP256Sha256, ak, sk, p256nMinusTwo, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	return s.deriveSigningKey(privInt), nil
}

func (s P256SHA256Signer) deriveSigningKey(priv *big.Int) ISigningKey {
	privateKey := new(ecdsa.PrivateKey)
	privateKey.PublicKey.Curve = curveP256
	privateKey.D = priv
	privateKey.PublicKey.X, privateKey.PublicKey.Y = curveP256.ScalarBaseMult(priv.Bytes())
	return P256SigningKey{privateKey: privateKey}
}

func derivePrivateInt(alg, ak, sk string, nMinusTwo *big.Int, hasher iHasher) (*big.Int, error) {
	context := bytes.NewBuffer(make([]byte, 0, len(ak)+1))
	data := bytes.NewBuffer(nil)

	for counter := 0; counter <= 0xff; counter++ {
		context.Reset()
		data.Reset()

		context.WriteString(ak)
		context.WriteByte(byte(counter))

		data.Write([]byte{0x00, 0x00, 0x00, 0x01})
		data.WriteString(alg)
		data.WriteByte(0x00)
		data.Write(context.Bytes())
		data.Write([]byte{0x00, 0x00, 0x01, 0x00})

		hmacBytes, err := hasher.hmac(data.Bytes(), []byte(sk))
		if err != nil {
			return nil, err
		}

		candidate := new(big.Int).SetBytes(hmacBytes)
		if candidate.Cmp(nMinusTwo) <= 0 {
			return candidate.Add(candidate, one), nil
		}
	}
	return nil, errors.New("derive candidate failed, counter out of range")
}

// signStringToSign Create the Signature.
func signStringToSign(stringToSign string, signingKey ISigningKey) (string, error) {
	sig, err := signingKey.Sign([]byte(stringToSign))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(sig), nil
}
