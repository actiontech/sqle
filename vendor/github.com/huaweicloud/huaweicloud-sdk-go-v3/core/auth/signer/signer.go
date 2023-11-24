// HWS API Gateway Signature
// based on https://github.com/datastream/aws/blob/master/signv4.go
// Copyright (c) 2014, Xianjie
// License that can be found in the LICENSE file

package signer

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/request"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"time"
)

const (
	sdkHmacSha256     = "SDK-HMAC-SHA256"
	xSdkContentSha256 = "X-Sdk-Content-Sha256"
)

type Signer struct {
}

// Sign SignRequest set Authorization header
func (s Signer) Sign(req *request.DefaultHttpRequest, ak, sk string) (map[string]string, error) {
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

	sts, err := stringToSign(sdkHmacSha256, cr, t, sha256HasherInst)
	if err != nil {
		return nil, err
	}

	sig, err := s.signStringToSign(sts, []byte(sk))
	if err != nil {
		return nil, err
	}

	additionalHeaders[HeaderAuthorization] = authHeaderValue(sdkHmacSha256, sig, ak, signedHeaders)
	return additionalHeaders, nil
}

// signStringToSign Create the Signature.
func (s Signer) signStringToSign(stringToSign string, signingKey []byte) (string, error) {
	hmac, err := sha256HasherInst.hmac([]byte(stringToSign), signingKey)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hmac), nil
}

const (
	BasicDateFormat     = "20060102T150405Z"
	HeaderXDate         = "X-Sdk-Date"
	HeaderHost          = "host"
	HeaderAuthorization = "Authorization"
)

func checkAKSK(ak, sk string) error {
	if ak == "" {
		return errors.New("ak is required in credentials")
	}
	if sk == "" {
		return errors.New("sk is required in credentials")
	}

	return nil
}

// stringToSign Create a "String to Sign".
func stringToSign(alg, canonicalRequest string, t time.Time, hasher iHasher) (string, error) {
	canonicalRequestHash, err := hasher.hashHexString([]byte(canonicalRequest))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s", alg, t.UTC().Format(BasicDateFormat), canonicalRequestHash), nil
}

// authHeaderValue Get the finalized value for the "Authorization" header.
// The signature parameter is the output from stringToSign
func authHeaderValue(alg, sig, ak string, signedHeaders []string) string {
	return fmt.Sprintf("%s Access=%s, SignedHeaders=%s, Signature=%s",
		alg,
		ak,
		strings.Join(signedHeaders, ";"),
		sig)
}

func processContentHeader(req *request.DefaultHttpRequest, contentHeader string) {
	if contentType, ok := req.GetHeaderParams()["Content-Type"]; ok && !strings.Contains(contentType, "application/json") {
		req.AddHeaderParam(contentHeader, "UNSIGNED-PAYLOAD")
	}
}

func canonicalRequest(req *request.DefaultHttpRequest, signedHeaders []string, contentHeader string, hasher iHasher) (string, error) {
	hexEncode, err := getContentHash(req, contentHeader, hasher)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		req.GetMethod(),
		canonicalURI(req),
		canonicalQueryString(req),
		canonicalHeaders(req, signedHeaders),
		strings.Join(signedHeaders, ";"), hexEncode), nil
}

func getContentHash(req *request.DefaultHttpRequest, contentHeader string, hasher iHasher) (string, error) {
	if content, ok := req.GetHeaderParams()[contentHeader]; ok {
		return content, nil
	}

	buffer, err := req.GetBodyToBytes()
	if err != nil {
		return "", err
	}

	data := buffer.Bytes()
	hexEncode, err := hasher.hashHexString(data)
	if err != nil {
		return "", err
	}
	return hexEncode, nil
}

func extractTime(headers map[string]string) time.Time {
	if date, ok := headers[HeaderXDate]; ok {
		t, err := time.Parse(BasicDateFormat, date)
		if date == "" || err != nil {
			return time.Now()
		}
		return t
	}
	return time.Now()
}

// canonicalURI returns request uri
func canonicalURI(r *request.DefaultHttpRequest) string {
	pattens := strings.Split(r.GetPath(), "/")

	var uri []string
	for _, v := range pattens {
		uri = append(uri, escape(v))
	}

	urlPath := strings.Join(uri, "/")
	if len(urlPath) == 0 || urlPath[len(urlPath)-1] != '/' {
		urlPath = urlPath + "/"
	}

	return urlPath
}

func canonicalQueryString(r *request.DefaultHttpRequest) string {
	var query = make(map[string][]string, 0)
	for key, value := range r.GetQueryParams() {
		valueWithType, ok := value.(reflect.Value)
		if !ok {
			continue
		}

		if valueWithType.Kind() == reflect.Slice {
			params := r.CanonicalSliceQueryParamsToMulti(valueWithType)
			for _, param := range params {
				if _, ok := query[key]; !ok {
					query[key] = make([]string, 0)
				}
				query[key] = append(query[key], param)
			}
		} else if valueWithType.Kind() == reflect.Map {
			params := r.CanonicalMapQueryParams(key, valueWithType)
			for _, param := range params {
				for k, v := range param {
					if _, ok := query[k]; !ok {
						query[k] = make([]string, 0)
					}
					query[k] = append(query[k], v)
				}
			}
		} else {
			if _, ok := query[key]; !ok {
				query[key] = make([]string, 0)
			}
			query[key] = append(query[key], r.CanonicalStringQueryParams(valueWithType))
		}
	}

	var keys []string
	for key := range query {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var a []string
	for _, key := range keys {
		k := escape(key)
		sort.Strings(query[key])
		for _, v := range query[key] {
			kv := fmt.Sprintf("%s=%s", k, escape(v))
			a = append(a, kv)
		}
	}
	queryStr := strings.Join(a, "&")

	return queryStr
}

func canonicalHeaders(r *request.DefaultHttpRequest, signerHeaders []string) string {
	var a []string
	header := make(map[string][]string)
	userHeaders := r.GetHeaderParams()

	for k, v := range userHeaders {
		if _, ok := header[strings.ToLower(k)]; !ok {
			header[strings.ToLower(k)] = make([]string, 0)
		}
		header[strings.ToLower(k)] = append(header[strings.ToLower(k)], v)
	}

	for _, key := range signerHeaders {
		value := header[key]
		if strings.EqualFold(key, HeaderHost) {
			if u, err := url.Parse(r.GetEndpoint()); err == nil {
				header[HeaderHost] = []string{u.Host}
			}
		}

		sort.Strings(value)
		for _, v := range value {
			a = append(a, key+":"+strings.TrimSpace(v))
		}
	}

	return fmt.Sprintf("%s\n", strings.Join(a, "\n"))
}

func extractSignedHeaders(headers map[string]string) []string {
	var sh []string
	for key := range headers {
		if strings.HasPrefix(strings.ToLower(key), "content-type") {
			continue
		}
		sh = append(sh, strings.ToLower(key))
	}
	sort.Strings(sh)

	return sh
}

func shouldEscape(c byte) bool {
	if 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z' || '0' <= c && c <= '9' || c == '_' || c == '-' || c == '~' || c == '.' {
		return false
	}
	return true
}

func escape(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	t := make([]byte, len(s)+2*hexCount)
	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = "0123456789ABCDEF"[c>>4]
			t[j+2] = "0123456789ABCDEF"[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
