package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// sys用户长有效期token，有限期至2073年
var DefaultDMSToken = "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOiI3MDAyMDEiLCJleHAiOjMyNTkxMzU3ODIsImlzcyI6ImFjdGlvbnRlY2ggZG1zIn0.pdCLmGM-lZwcBsOnwfxM2m5xbUzGpEqAiRurkCj-8YY"

func Get(ctx context.Context, url string, headers map[string]string, body, out interface{}) error {
	return Call(ctx, http.MethodGet, url, headers, body, out)
}

func POST(ctx context.Context, url string, headers map[string]string, body, out interface{}) error {
	return Call(ctx, http.MethodPost, url, headers, body, out)
}

func Call(ctx context.Context, method, url string, headers map[string]string, body, out interface{}) error {
	var bodyReader io.Reader
	if body != nil {
		bodyJson, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal error: %v", err)
		}
		bodyReader = bytes.NewReader(bodyJson)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("new request error: %v", err)
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: time.Second * 15,
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("query data error: %v", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read data error: %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("query data error: %v;%v", resp.Status, string(result))
	}

	err = json.Unmarshal(result, &out)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}
	return nil
}
