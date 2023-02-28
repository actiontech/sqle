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

func (f *FeishuClient) GetUsersByEmailOrMobile(emails, mobiles []string) (map[string]*UserContactInfo, error) {
	//查询限制每次最多50条emails和mobiles，https://open.feishu.cn/document/uAjLw4CM/ukTMukTMukTM/reference/contact-v3/user/batch_get_id
	limitation := 50
	users := make(map[string] /*user_id*/ *UserContactInfo)
	for head, end := 0, limitation-1; ; {
		if head > len(mobiles)-1 && head > len(emails)-1 {
			break
		}

		queryBuilder := func(queryIds []string) []string {
			length := len(queryIds)
			ret := []string{}
			if end <= length-1 {
				ret = queryIds[head : end+1]
			} else if head <= length-1 {
				ret = queryIds[head:length]
			}
			return ret
		}
		queryEmails := queryBuilder(emails)
		queryMobiles := queryBuilder(mobiles)

		req := larkcontact.NewBatchGetIdUserReqBuilder().
			UserIdType(`user_id`).
			Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
				Emails(queryEmails).
				Mobiles(queryMobiles).
				Build()).
			Build()

		resp, err := f.client.Contact.User.BatchGetId(context.Background(), req)
		if err != nil {
			return nil, wrapSendRequestFailedError(err)
		}
		if !resp.Success() {
			return nil, fmt.Errorf("get user ids failed: respCode=%v, respMsg=%v", resp.Code, resp.Msg)
		}

		f.convertUsersResp(resp.Data.UserList, users)
		head = end + 1
		end += limitation
	}
	return users, nil
}

func (f *FeishuClient) convertUsersResp(raw []*larkcontact.UserContactInfo, users map[string]*UserContactInfo) {
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
