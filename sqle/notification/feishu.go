package notification

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/im/feishu"
	"github.com/actiontech/sqle/sqle/utils"
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
	for _, u := range users {
		if u.Email != "" {
			emails = append(emails, u.Email)
		}
		if u.Phone != "" {
			mobiles = append(mobiles, u.Phone)
		}
	}
	emails = utils.RemoveDuplicate(emails)
	mobiles = utils.RemoveDuplicate(mobiles)

	client := feishu.NewFeishuClient(cfg.AppKey, cfg.AppSecret)
	feishuUsers, err := client.GetUserIdsByEmailOrMobile(emails, mobiles)
	if err != nil {
		return fmt.Errorf("get user_ids from feishu failed: %v", err)
	}

	// content是要作为json文本的一个value传给飞书，需要转换为显式的换行符
	content := strings.Replace(notification.NotificationBody(), "\n", "\\n", -1)
	errMsgs := []string{}
	l := log.NewEntry()
	sentUserIds := make(map[string]struct{})
	for _, u := range feishuUsers {
		id := utils.NvlString(u.UserId)
		if id == "" {
			continue
		}
		// feishuUsers是通过email和phone获取的userId，可能会有重复的userId
		if _, ok := sentUserIds[id]; ok {
			continue
		} else {
			sentUserIds[id] = struct{}{}
		}
		email := utils.NvlString(u.Email)
		phone := utils.NvlString(u.Mobile)
		l.Infof("send message to feishu, email=%v, phone=%v, userId=%v", email, phone, id)
		if err = client.SendMessage(feishu.FeishuRceiveIdTypeUserId, id, feishu.FeishuSendMessageMsgTypePost, fmt.Sprintf(FeishuContentPattern, notification.NotificationSubject(), content)); err != nil {
			errMsgs = append(errMsgs, fmt.Sprintf("send message to feishu failed: %v; email=%v; phone=%v", err, email, phone))
		}
	}
	if len(errMsgs) > 0 {
		return fmt.Errorf(strings.Join(errMsgs, "\n"))
	}
	return nil
}

var FeishuContentPattern = `
{
  "zh_cn": {
    "title": "%v",
    "content": [
      [
        {
          "tag": "text",
          "text": "%v"
        }
      ]
    ]
  }
}`
