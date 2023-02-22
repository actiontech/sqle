package feishu

import (
	"context"
	"fmt"

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

func (f *FeishuClient) GetUserIdsByEmailOrMobile(emails, mobiles []string) ([]*larkcontact.UserContactInfo, error) {
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
	return resp.Data.UserList, nil
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
