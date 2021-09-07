package httpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

const (
	defaultTimeout = time.Second * 10
)

// HTTPClient is a wrap of http.Client
type Client struct {
	*http.Client
}

// NewHTTPClient returns a new HTTP client with timeout and HTTPS support
func NewHTTPClient(timeout time.Duration, tlsCfg *tls.Config) *Client {
	if timeout < time.Second {
		timeout = defaultTimeout
	}
	tp := &http.Transport{
		TLSClientConfig: tlsCfg,
		Dial:            (&net.Dialer{Timeout: 3 * time.Second}).Dial,
	}
	return &Client{&http.Client{
		Timeout:   timeout,
		Transport: tp,
	}}
}

func (c *Client) SendRequest(ctx context.Context, url, method, token string, body io.Reader) ([]byte, error) {
	defer c.CloseIdleConnections()
	switch method {
	case http.MethodGet:
		return c.Get(ctx, url, token)
	case http.MethodPost:
		return c.Post(ctx, url, token, body)
	default:
		return nil, fmt.Errorf("invalid request method")
	}
}

// Get fetch a URL with GET method and returns the response
func (c *Client) Get(ctx context.Context, url, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return checkHTTPResponse(res)
}

// Post send a POST request to the url and returns the response
func (c *Client) Post(ctx context.Context, url, token string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return checkHTTPResponse(res)
}

// checkHTTPResponse checks if an HTTP response is with normal status codes
func checkHTTPResponse(res *http.Response) ([]byte, error) {
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return body, fmt.Errorf("error requesting %s, response: %s, code %d", res.Request.URL, string(body), res.StatusCode)
	}
	return body, nil
}
