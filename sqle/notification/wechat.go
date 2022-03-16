//go:build enterprise
// +build enterprise

package notification

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/model"

	"gopkg.in/chanxuehong/wechat.v1/corp"
	"gopkg.in/chanxuehong/wechat.v1/corp/message/send"
)

func init() {
	Notifiers = append(Notifiers, &WeChatNotifier{})
}

type WeChatNotifier struct{}

func (n *WeChatNotifier) Notify(notification Notification, users []*model.User) error {
	// workflow has been finished.
	if len(users) == 0 {
		return nil
	}
	wechatUsers := map[string] /*user name*/ string /*wechat id*/ {}
	for _, user := range users {
		if user.WeChatID != "" {
			wechatUsers[user.Name] = user.WeChatID
		}
	}

	// no user has configured email, don't send.
	if len(wechatUsers) == 0 {
		return nil
	}
	s := model.GetStorage()
	wechatC, exist, err := s.GetWeChatConfiguration()
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	if !wechatC.EnableWeChatNotify {
		return nil
	}

	client := generateWeChatClient(wechatC)
	safe := 0
	if wechatC.SafeEnabled {
		safe = 1
	}
	errs := []string{}
	for name, u := range wechatUsers {
		req := &send.Text{
			MessageHeader: send.MessageHeader{
				ToUser:  u,
				MsgType: "text",
				AgentId: int64(wechatC.AgentID),
				Safe:    &safe,
			},
		}
		req.Text.Content = fmt.Sprintf("%v \n\n %v", notification.NotificationSubject(), notification.NotificationBody())
		_, err := client.SendText(req)
		if err != nil {
			errs = append(errs, fmt.Sprintf("send message to %v failed, error: %v", name, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%v", strings.Join(errs, "\n"))
	}
	return nil
}

func generateWeChatClient(conf *model.WeChatConfiguration) *send.Client {
	proxy := http.ProxyFromEnvironment
	if conf.ProxyIP != "" {
		proxy = func(req *http.Request) (*url.URL, error) {
			return url.Parse(conf.ProxyIP)
		}
	}
	var transport http.RoundTripper = &http.Transport{
		Proxy: proxy,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	httpClient := &http.Client{
		Transport: transport,
	}
	accessTokenServer := corp.NewDefaultAccessTokenServer(conf.CorpID, conf.CorpSecret, httpClient)
	return send.NewClient(accessTokenServer, httpClient)
}
