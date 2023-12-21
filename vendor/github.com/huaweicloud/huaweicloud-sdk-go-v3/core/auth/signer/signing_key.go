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
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/asn1"
	"github.com/tjfoc/gmsm/sm2"
	"math/big"
)

type ISigningKey interface {
	Sign(data []byte) ([]byte, error)
	Verify(signature, data []byte) bool
}

type ecSignature struct {
	R *big.Int
	S *big.Int
}

type P256SigningKey struct {
	privateKey *ecdsa.PrivateKey
}

func (k P256SigningKey) Sign(data []byte) ([]byte, error) {
	hashed, err := sha256HasherInst.hash(data)
	if err != nil {
		return nil, err
	}
	r, s, err := ecdsa.Sign(rand.Reader, k.privateKey, hashed)
	if err != nil {
		return nil, err
	}

	ecSig := ecSignature{
		R: r,
		S: s,
	}
	sig, err := asn1.Marshal(ecSig)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (k P256SigningKey) Verify(signature, data []byte) bool {
	ecSig := ecSignature{}
	_, err := asn1.Unmarshal(signature, &ecSig)
	if err != nil {
		return false
	}

	hashed, err := sha256HasherInst.hash(data)
	if err != nil {
		return false
	}

	publicKey := &k.privateKey.PublicKey
	return ecdsa.Verify(publicKey, hashed, ecSig.R, ecSig.S)
}

type SM2SigningKey struct {
	privateKey *sm2.PrivateKey
}

func (k SM2SigningKey) Sign(data []byte) ([]byte, error) {
	r, s, err := sm2.Sm2Sign(k.privateKey, data, []byte{}, rand.Reader)
	if err != nil {
		return nil, err
	}

	ecSig := ecSignature{
		R: r,
		S: s,
	}
	sig, err := asn1.Marshal(ecSig)
	if err != nil {
		return nil, err
	}
	return sig, nil
}

func (k SM2SigningKey) Verify(signature, data []byte) bool {
	ecSig := ecSignature{}
	_, err := asn1.Unmarshal(signature, &ecSig)
	if err != nil {
		return false
	}

	return sm2.Sm2Verify(&k.privateKey.PublicKey, data, []byte{}, ecSig.R, ecSig.S)
}
