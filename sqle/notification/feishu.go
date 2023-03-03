package notification

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"
	larkIm "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

func init() {
	Notifiers = append(Notifiers, &FeishuNotifier{})
}

type FeishuNotifier struct{}

func (n *FeishuNotifier) Notify(notification Notification, users []*model.User) error {
	// workflow has been finished.
	if len(users) == 0 {
		return nil
	}

	s := model.GetStorage()
	cfg, exist, err := s.GetImConfigByType(model.ImTypeFeishu)
	if err != nil {
		return fmt.Errorf("get im config failed: %v", err)
	}
	if !exist {
		return nil
	}

	if !cfg.IsEnable {
		return nil
	}

	// 通过邮箱、手机从飞书获取用户ids
	var emails, mobiles []string
	userCount := 0
	for _, u := range users {
		if u.Email == "" && u.Phone == "" {
			continue
		}
		if u.Email != "" {
			emails = append(emails, u.Email)
		}
		if u.Phone != "" {
			mobiles = append(mobiles, u.Phone)
		}
		userCount++
		if userCount == feishu.MaxCountOfIdThatUsedToFindUser {
			break
		}
	}

	client := feishu.NewFeishuClient(cfg.AppKey, cfg.AppSecret)
	feishuUsers, err := client.GetUsersByEmailOrMobileWithLimitation(emails, mobiles)
	if err != nil {
		return fmt.Errorf("get user_ids from feishu failed: %v", err)
	}

	content, err := BuildFeishuMessageBody(notification)
	if err != nil {
		return fmt.Errorf("convert content failed: %v", err)
	}
	errMsgs := []string{}
	l := log.NewEntry()
	for id, u := range feishuUsers {
		l.Infof("send message to feishu, email=%v, phone=%v, userId=%v", u.Email, u.Mobile, id)
		if err = client.SendMessage(feishu.FeishuReceiverIdTypeUserId, id, feishu.FeishuSendMessageMsgTypePost, content); err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("send message to feishu failed: %v; email=%v; phone=%v", err, u.Email, u.Mobile))
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

func BuildFeishuMessageBody(n Notification) (string, error) {
	zhCnPostText := &larkIm.MessagePostText{Text: n.NotificationBody()}
	zhCnMessagePostContent := &larkIm.MessagePostContent{Title: n.NotificationSubject(), Content: [][]larkIm.MessagePostElement{{zhCnPostText}}}
	messagePostText := &larkIm.MessagePost{ZhCN: zhCnMessagePostContent}
	content, err := messagePostText.String()
	if err != nil {
		return "", err
	}
	return content, nil
}
