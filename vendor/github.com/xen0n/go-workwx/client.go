package workwx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/url"
)

// Workwx 企业微信客户端
type Workwx struct {
	opts options

	// CorpID 企业 ID，必填
	CorpID string
}

// WorkwxApp 企业微信客户端（分应用）
//
//nolint:revive // The (stuttering) name is part of public API, so cannot be fixed without a v2 bump
type WorkwxApp struct {
	*Workwx

	// CorpSecret 应用的凭证密钥，必填
	CorpSecret string
	// AgentID 应用 ID，必填
	AgentID int64

	accessToken            *token
	jsapiTicket            *token
	jsapiTicketAgentConfig *token
}

// New 构造一个 Workwx 客户端对象，需要提供企业 ID
func New(corpID string, opts ...CtorOption) *Workwx {
	optionsObj := defaultOptions()

	for _, o := range opts {
		o.applyTo(&optionsObj)
	}

	return &Workwx{
		opts: optionsObj,

		CorpID: corpID,
	}
}

// WithApp 构造本企业下某自建 app 的客户端
func (c *Workwx) WithApp(corpSecret string, agentID int64) *WorkwxApp {
	app := WorkwxApp{
		Workwx: c,

		CorpSecret: corpSecret,
		AgentID:    agentID,
	}

	app.accessToken = newToken(c.opts.AccessTokenProvider, app.getAccessToken)
	app.jsapiTicket = newToken(c.opts.JSAPITicketProvider, app.getJSAPITicket)
	app.jsapiTicketAgentConfig = newToken(c.opts.JSAPITicketAgentConfigProvider, app.getJSAPITicketAgentConfig)

	return &app
}

func (c *WorkwxApp) composeQyapiURL(path string, req interface{}) (*url.URL, error) {
	values := url.Values{}
	if valuer, ok := req.(urlValuer); ok {
		values = valuer.intoURLValues()
	}

	// TODO: refactor
	base, err := url.Parse(c.opts.QYAPIHost)
	if err != nil {
		return nil, fmt.Errorf("qyapiHost invalid: host=%s err=%w", c.opts.QYAPIHost, err)
	}

	base.Path = path
	base.RawQuery = values.Encode()

	return base, nil
}

func (c *WorkwxApp) composeQyapiURLWithToken(path string, req interface{}, withAccessToken bool) (*url.URL, error) {
	url, err := c.composeQyapiURL(path, req)
	if err != nil {
		return nil, err
	}

	if !withAccessToken {
		return url, nil
	}

	tok, err := c.accessToken.getToken()
	if err != nil {
		return nil, err
	}

	q := url.Query()
	q.Set("access_token", tok)
	url.RawQuery = q.Encode()

	return url, nil
}

func (c *WorkwxApp) executeQyapiGet(path string, req urlValuer, respObj interface{}, withAccessToken bool) error {
	url, err := c.composeQyapiURLWithToken(path, req, withAccessToken)
	if err != nil {
		return err
	}
	urlStr := url.String()

	resp, err := c.opts.HTTP.Get(urlStr)
	if err != nil {
		return makeRequestErr(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(respObj)
	if err != nil {
		return makeRespUnmarshalErr(err)
	}

	return nil
}

func (c *WorkwxApp) executeQyapiJSONPost(path string, req bodyer, respObj interface{}, withAccessToken bool) error {
	url, err := c.composeQyapiURLWithToken(path, req, withAccessToken)
	if err != nil {
		return err
	}
	urlStr := url.String()

	body, err := req.intoBody()
	if err != nil {
		return makeReqMarshalErr(err)
	}

	resp, err := c.opts.HTTP.Post(urlStr, "application/json", bytes.NewReader(body))
	if err != nil {
		return makeRequestErr(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(respObj)
	if err != nil {
		return makeRespUnmarshalErr(err)
	}

	return nil
}

func (c *WorkwxApp) executeQyapiMediaUpload(
	path string,
	req mediaUploader,
	respObj interface{},
	withAccessToken bool,
) error {
	url, err := c.composeQyapiURLWithToken(path, req, withAccessToken)
	if err != nil {
		return err
	}
	urlStr := url.String()

	m := req.getMedia()

	// FIXME: use streaming upload to conserve memory!
	buf := bytes.Buffer{}
	mw := multipart.NewWriter(&buf)

	err = m.writeTo(mw)
	if err != nil {
		return err
	}

	err = mw.Close()
	if err != nil {
		return err
	}

	resp, err := c.opts.HTTP.Post(urlStr, mw.FormDataContentType(), &buf)
	if err != nil {
		return makeRequestErr(err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(respObj)
	if err != nil {
		return makeRespUnmarshalErr(err)
	}

	return nil
}
