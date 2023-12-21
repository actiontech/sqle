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
	chmac "crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/tjfoc/gmsm/sm3"
)

type iHasher interface {
	hash(data []byte) ([]byte, error)
	hashHexString(data []byte) (string, error)
	hmac(data []byte, key []byte) ([]byte, error)
}

type sm3Hasher struct {
}

func (h sm3Hasher) hash(data []byte) ([]byte, error) {
	hash := sm3.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (h sm3Hasher) hashHexString(data []byte) (string, error) {
	hash, err := h.hash(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}

func (h sm3Hasher) hmac(data []byte, key []byte) ([]byte, error) {
	hash := chmac.New(sm3.New, key)
	if _, err := hash.Write(data); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

type sha256Hasher struct {
}

func (h sha256Hasher) hash(data []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}

func (h sha256Hasher) hashHexString(data []byte) (string, error) {
	hash, err := h.hash(data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash), nil
}

func (h sha256Hasher) hmac(data []byte, key []byte) ([]byte, error) {
	hash := chmac.New(sha256.New, key)
	if _, err := hash.Write(data); err != nil {
		return nil, err
	}
	return hash.Sum(nil), nil
}
