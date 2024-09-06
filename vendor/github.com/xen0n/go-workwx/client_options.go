package workwx

import (
	"net/http"
)

// DefaultQYAPIHost 默认企业微信 API Host
const DefaultQYAPIHost = "https://qyapi.weixin.qq.com"

type options struct {
	QYAPIHost                      string
	HTTP                           *http.Client
	AccessTokenProvider            ITokenProvider
	JSAPITicketProvider            ITokenProvider
	JSAPITicketAgentConfigProvider ITokenProvider
}

// CtorOption 客户端对象构造参数
type CtorOption interface {
	applyTo(*options)
}

// impl Default for options
func defaultOptions() options {
	return options{
		QYAPIHost:                      DefaultQYAPIHost,
		HTTP:                           &http.Client{},
		AccessTokenProvider:            nil,
		JSAPITicketProvider:            nil,
		JSAPITicketAgentConfigProvider: nil,
	}
}

//
//
//

type withQYAPIHost struct {
	x string
}

// WithQYAPIHost 覆盖默认企业微信 API 域名
func WithQYAPIHost(host string) CtorOption {
	return &withQYAPIHost{x: host}
}

var _ CtorOption = (*withQYAPIHost)(nil)

func (x *withQYAPIHost) applyTo(y *options) {
	y.QYAPIHost = x.x
}

//
//
//

type withHTTPClient struct {
	x *http.Client
}

// WithHTTPClient 使用给定的 http.Client 作为 HTTP 客户端
func WithHTTPClient(client *http.Client) CtorOption {
	return &withHTTPClient{x: client}
}

var _ CtorOption = (*withHTTPClient)(nil)

func (x *withHTTPClient) applyTo(y *options) {
	y.HTTP = x.x
}

//
//
//

type withAccessTokenProvider struct {
	x ITokenProvider
}

func WithAccessTokenProvider(provider ITokenProvider) CtorOption {
	return &withAccessTokenProvider{x: provider}
}

var _ CtorOption = (*withAccessTokenProvider)(nil)

func (x *withAccessTokenProvider) applyTo(y *options) {
	y.AccessTokenProvider = x.x
}

//
//
//

type withJSAPITicketProvider struct {
	x ITokenProvider
}

func WithJSAPITicketProvider(provider ITokenProvider) CtorOption {
	return &withJSAPITicketProvider{x: provider}
}

var _ CtorOption = (*withJSAPITicketProvider)(nil)

func (x *withJSAPITicketProvider) applyTo(y *options) {
	y.JSAPITicketProvider = x.x
}

//
//
//

type withJSAPITicketAgentConfigProvider struct {
	x ITokenProvider
}

func WithJSAPITicketAgentConfigProvider(provider ITokenProvider) CtorOption {
	return &withJSAPITicketAgentConfigProvider{x: provider}
}

var _ CtorOption = (*withJSAPITicketAgentConfigProvider)(nil)

func (x *withJSAPITicketAgentConfigProvider) applyTo(y *options) {
	y.JSAPITicketAgentConfigProvider = x.x
}
