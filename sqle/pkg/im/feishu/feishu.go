package feishu

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/utils"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var wrapSendRequestFailedError = func(err error) error { return fmt.Errorf("send message failed: %v", err) }

type FeishuClient struct {
	client *lark.Client
}

func NewFeishuClient(appId, appSecret string) *FeishuClient {
	return &FeishuClient{client: lark.NewClient(appId, appSecret)}
}

type UserContactInfo struct {
	Email  string
	Mobile string
}

func (f *FeishuClient) GetUserIdsByEmailOrMobile(emails, mobiles []string) (map[string]*UserContactInfo, error) {
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`user_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails(emails).
			Mobiles(mobiles).
			Build()).
		Build()

	resp, err := f.client.Contact.User.BatchGetId(context.Background(), req)
	if err != nil {
		return nil, wrapSendRequestFailedError(err)
	}
	if !resp.Success() {
		return nil, fmt.Errorf("get user ids failed: respCode=%v, respMsg=%v", resp.Code, resp.Msg)
	}

	users := make(map[string] /*user_id*/ *UserContactInfo)
	for _, user := range resp.Data.UserList {
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
	return users, nil
}

const (
	FeishuRceiveIdTypeUserId = "user_id"

	FeishuSendMessageMsgTypePost = "post"
)

func (f FeishuClient) SendMessage(receiveIdType, receiveId, msgType, content string) error {
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(receiveIdType).
		Body(larkim.NewCreateMessageReqBodyBuilder().
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
