package feishu

import (
	"context"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkContact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkIm "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var wrapSendRequestFailedError = func(err error) error { return fmt.Errorf("send message failed: %v", err) }

type FeishuClient struct {
	client *lark.Client
}

func NewFeishuClient(appId, appSecret string) *FeishuClient {
	return &FeishuClient{client: lark.NewClient(appId, appSecret, lark.WithReqTimeout(30*time.Second))}
}

type UserContactInfo struct {
	Email  string
	Mobile string
}

const MaxCountOfIdThatUsedToFindUser = 50

// 查询限制每次最多50条emails和mobiles，https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/user/batch_get_id
// 每次最多查询50个邮箱和50个手机号，如果超出50个，只查询前50个
func (f *FeishuClient) GetUsersByEmailOrMobileWithLimitation(emails, mobiles []string, userType string) (map[string]*UserContactInfo, error) {
	tempEmails, tempMobiles := emails, mobiles
	if len(emails) > MaxCountOfIdThatUsedToFindUser {
		tempEmails = emails[:MaxCountOfIdThatUsedToFindUser]
	}
	if len(mobiles) > MaxCountOfIdThatUsedToFindUser {
		tempMobiles = mobiles[:MaxCountOfIdThatUsedToFindUser]
	}

	req := larkContact.NewBatchGetIdUserReqBuilder().
		UserIdType(userType).
		Body(larkContact.NewBatchGetIdUserReqBodyBuilder().
			Emails(tempEmails).
			Mobiles(tempMobiles).
			Build()).
		Build()

	resp, err := f.client.Contact.User.BatchGetId(context.Background(), req)
	if err != nil {
		return nil, wrapSendRequestFailedError(err)
	}
	if !resp.Success() {
		return nil, fmt.Errorf("get user ids failed: respCode=%v, respMsg=%v", resp.Code, resp.Msg)
	}

	users := make(map[string]*UserContactInfo)
	f.convertUsersResp(resp.Data.UserList, users)
	return users, nil
}

func (f *FeishuClient) convertUsersResp(raw []*larkContact.UserContactInfo, users map[string]*UserContactInfo) {
	for _, user := range raw {
		id := utils.NvlString(user.UserId)
		if id == "" {
			continue
		}

		_, ok := users[id]
		if !ok {
			users[id] = &UserContactInfo{}
		}
		info := users[id]

		// 飞书接口的响应结构里不会同时有email和mobile
		if email := utils.NvlString(user.Email); email != "" {
			info.Email = email
			continue
		}
		if mobile := utils.NvlString(user.Mobile); mobile != "" {
			info.Mobile = mobile
			continue
		}
	}
}

const (
	FeishuReceiverIdTypeUserId = "user_id"

	FeishuSendMessageMsgTypePost = "post"
)

func (f FeishuClient) SendMessage(receiveIdType, receiveId, msgType, content string) error {
	req := larkIm.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkIm.NewCreateMessageReqBodyBuilder().
			ReceiveId(receiveId).
			MsgType(msgType).
			Content(content).
			Build()).
		Build()

	resp, err := f.client.Im.Message.Create(context.Background(), req)
	if err != nil {
		return wrapSendRequestFailedError(err)
	}

	if !resp.Success() {
		return fmt.Errorf("send message to user failed: respCode=%v, respMsg=%v", resp.Code, resp.Msg)
	}

	return nil
}

func (f FeishuClient) GetFeishuUserIdList(users []*model.User, userType string) ([]string, error) {
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
		if userCount == MaxCountOfIdThatUsedToFindUser {
			log.NewEntry().Infof("user %v exceed max count %v", u.Name, MaxCountOfIdThatUsedToFindUser)
			break
		}
	}

	feishuUserMap, err := f.GetUsersByEmailOrMobileWithLimitation(emails, mobiles, userType)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(feishuUserMap))
	for feishuUserID := range feishuUserMap {
		userIDs = append(userIDs, feishuUserID)
	}

	return userIDs, nil
}

func (f FeishuClient) GetFeishuUserInfo(userID string, userType string) (*larkContact.User, error) {
	req := larkContact.NewGetUserReqBuilder().
		UserId(userID).
		UserIdType(userType).
		Build()

	resp, err := f.client.Contact.User.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("get user failed: respCode=%v, respMsg=%v", resp.Code, resp.Msg)
	}

	return resp.Data.User, nil
}
