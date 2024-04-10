package workwx

import (
	"context"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
)

// ITokenProvider 是鉴权 token 的外部提供者需要实现的 interface。可用于官方所谓
// 使用“中控服务”集中提供、刷新 token 的场景。
//
// 不同类型的 tokens（如 access token、JSAPI token 等）都是这个 interface 提供，
// 实现方需要自行掌握 token 的类别，避免在 client 构造函数的选项中传入错误的种类。
type ITokenProvider interface {
	// GetToken 取回一个 token。有可能被并发调用。
	GetToken(context.Context) (string, error)
}

type tokenInfo struct {
	token     string
	expiresIn time.Duration
}

type token struct {
	mutex *sync.RWMutex
	tokenInfo
	lastRefresh      time.Time
	getTokenFunc     func() (tokenInfo, error)
	externalProvider ITokenProvider
}

func newToken(
	externalProvider ITokenProvider,
	refresher func() (tokenInfo, error),
) *token {
	if externalProvider != nil {
		return &token{
			externalProvider: externalProvider,
		}
	}

	return &token{
		mutex:        &sync.RWMutex{},
		getTokenFunc: refresher,
	}
}

func (t *token) usingExternalProvider() bool {
	return t.externalProvider != nil
}

// getAccessToken 获取 access token
func (c *WorkwxApp) getAccessToken() (tokenInfo, error) {
	get, err := c.execGetAccessToken(reqAccessToken{
		CorpID:     c.CorpID,
		CorpSecret: c.CorpSecret,
	})
	if err != nil {
		return tokenInfo{}, err
	}
	return tokenInfo{token: get.AccessToken, expiresIn: time.Duration(get.ExpiresInSecs)}, nil
}

// SpawnAccessTokenRefresher 启动该 app 的 access token 刷新 goroutine
//
// 如果使用了外部 token provider 提供 access token 则没有必要调用此方法：调用效果为空操作。
//
// NOTE: 该 goroutine 本身没有 keep-alive 逻辑，需要自助保活
func (c *WorkwxApp) SpawnAccessTokenRefresher() {
	ctx := context.Background()
	c.SpawnAccessTokenRefresherWithContext(ctx)
}

// SpawnAccessTokenRefresherWithContext 启动该 app 的 access token 刷新 goroutine
// 可以通过 context cancellation 停止此 goroutine
//
// 如果使用了外部 token provider 提供 access token 则没有必要调用此方法：调用效果为空操作。
//
// NOTE: 该 goroutine 本身没有 keep-alive 逻辑，需要自助保活
func (c *WorkwxApp) SpawnAccessTokenRefresherWithContext(ctx context.Context) {
	if c.accessToken.usingExternalProvider() {
		return
	}

	go c.accessToken.tokenRefresher(ctx)
}

// GetJSAPITicket 获取 JSAPI_ticket
func (c *WorkwxApp) GetJSAPITicket() (string, error) {
	return c.jsapiTicket.getToken()
}

// getJSAPITicket 获取 JSAPI_ticket
func (c *WorkwxApp) getJSAPITicket() (tokenInfo, error) {
	get, err := c.execGetJSAPITicket(reqJSAPITicket{})
	if err != nil {
		return tokenInfo{}, err
	}
	return tokenInfo{token: get.Ticket, expiresIn: time.Duration(get.ExpiresInSecs)}, nil
}

// SpawnJSAPITicketRefresher 启动该 app 的 JSAPI_ticket 刷新 goroutine
//
// 如果使用了外部 token provider 提供 JSAPI ticket 则没有必要调用此方法：调用效果为空操作。
//
// NOTE: 该 goroutine 本身没有 keep-alive 逻辑，需要自助保活
func (c *WorkwxApp) SpawnJSAPITicketRefresher() {
	ctx := context.Background()
	c.SpawnJSAPITicketRefresherWithContext(ctx)
}

// SpawnJSAPITicketRefresherWithContext 启动该 app 的 JSAPI_ticket 刷新 goroutine
// 可以通过 context cancellation 停止此 goroutine
//
// 如果使用了外部 token provider 提供 JSAPI ticket 则没有必要调用此方法：调用效果为空操作。
//
// NOTE: 该 goroutine 本身没有 keep-alive 逻辑，需要自助保活
func (c *WorkwxApp) SpawnJSAPITicketRefresherWithContext(ctx context.Context) {
	if c.jsapiTicket.usingExternalProvider() {
		return
	}

	go c.jsapiTicket.tokenRefresher(ctx)
}

// GetJSAPITicketAgentConfig 获取 JSAPI_ticket_agent_config
func (c *WorkwxApp) GetJSAPITicketAgentConfig() (string, error) {
	return c.jsapiTicketAgentConfig.getToken()
}

// getJSAPITicketAgentConfig 获取 JSAPI_ticket_agent_config
func (c *WorkwxApp) getJSAPITicketAgentConfig() (tokenInfo, error) {
	get, err := c.execGetJSAPITicketAgentConfig(reqJSAPITicketAgentConfig{})
	if err != nil {
		return tokenInfo{}, err
	}
	return tokenInfo{token: get.Ticket, expiresIn: time.Duration(get.ExpiresInSecs)}, nil
}

// SpawnJSAPITicketAgentConfigRefresher 启动该 app 的 JSAPI_ticket_agent_config 刷新 goroutine
//
// 如果使用了外部 token provider 提供 JSAPI ticket agent config 则没有必要调用此方法：调用效果为空操作。
//
// NOTE: 该 goroutine 本身没有 keep-alive 逻辑，需要自助保活
func (c *WorkwxApp) SpawnJSAPITicketAgentConfigRefresher() {
	ctx := context.Background()
	c.SpawnJSAPITicketAgentConfigRefresherWithContext(ctx)
}

// SpawnJSAPITicketAgentConfigRefresherWithContext 启动该 app 的 JSAPI_ticket_agent_config 刷新 goroutine
// 可以通过 context cancellation 停止此 goroutine
//
// 如果使用了外部 token provider 提供 JSAPI ticket agent config 则没有必要调用此方法：调用效果为空操作。
//
// NOTE: 该 goroutine 本身没有 keep-alive 逻辑，需要自助保活
func (c *WorkwxApp) SpawnJSAPITicketAgentConfigRefresherWithContext(ctx context.Context) {
	if c.jsapiTicketAgentConfig.usingExternalProvider() {
		return
	}

	go c.jsapiTicketAgentConfig.tokenRefresher(ctx)
}

func (t *token) getToken() (string, error) {
	if t.externalProvider != nil {
		tok, err := t.externalProvider.GetToken(context.TODO())
		if err != nil {
			return "", err
		}
		return tok, nil
	}

	// intensive mutex juggling action
	t.mutex.RLock()
	if t.token == "" {
		t.mutex.RUnlock() // RWMutex doesn't like recursive locking
		err := t.syncToken()
		if err != nil {
			return "", err
		}
		t.mutex.RLock()
	}
	tokenToUse := t.token
	t.mutex.RUnlock()
	return tokenToUse, nil
}

func (t *token) syncToken() error {
	get, err := t.getTokenFunc()
	if err != nil {
		return err
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.token = get.token
	t.expiresIn = get.expiresIn * time.Second
	t.lastRefresh = time.Now()
	return nil
}

func (t *token) tokenRefresher(ctx context.Context) {
	const refreshTimeWindow = 30 * time.Minute
	const minRefreshDuration = 5 * time.Second

	var waitDuration time.Duration
	for {
		select {
		case <-time.After(waitDuration):
			retryer := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
			if err := backoff.Retry(t.syncToken, retryer); err != nil {
				// TODO: logging
				_ = err
			}

			waitUntilTime := t.lastRefresh.Add(t.expiresIn).Add(-refreshTimeWindow)
			waitDuration = time.Until(waitUntilTime)
			if waitDuration < minRefreshDuration {
				waitDuration = minRefreshDuration
			}
		case <-ctx.Done():
			return
		}
	}
}

// JSCode2Session 临时登录凭证校验
func (c *WorkwxApp) JSCode2Session(jscode string) (*JSCodeSession, error) {
	resp, err := c.execJSCode2Session(reqJSCode2Session{JSCode: jscode})
	if err != nil {
		return nil, err
	}
	return &resp.JSCodeSession, nil
}

// AuthCode2UserInfo 获取访问用户身份
func (c *WorkwxApp) AuthCode2UserInfo(code string) (*AuthCodeUserInfo, error) {
	resp, err := c.execAuthCode2UserInfo(reqAuthCode2UserInfo{Code: code})
	if err != nil {
		return nil, err
	}
	return &resp.AuthCodeUserInfo, nil
}
