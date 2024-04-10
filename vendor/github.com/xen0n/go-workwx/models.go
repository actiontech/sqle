package workwx

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func marshalIntoJSONBody(x interface{}) ([]byte, error) {
	y, err := json.Marshal(x)
	if err != nil {
		// should never happen unless OOM or similar bad things
		return nil, makeReqMarshalErr(err)
	}

	return y, nil
}

type reqAccessToken struct {
	CorpID     string
	CorpSecret string
}

var _ urlValuer = reqAccessToken{}

func (x reqAccessToken) intoURLValues() url.Values {
	return url.Values{
		"corpid":     {x.CorpID},
		"corpsecret": {x.CorpSecret},
	}
}

type respCommon struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

// IsOK 响应体是否为一次成功请求的响应
//
// 实现依据: https://work.weixin.qq.com/api/doc#10013
//
// > 企业微信所有接口，返回包里都有errcode、errmsg。
// > 开发者需根据errcode是否为0判断是否调用成功(errcode意义请见全局错误码)。
// > 而errmsg仅作参考，后续可能会有变动，因此不可作为是否调用成功的判据。
func (x *respCommon) IsOK() bool {
	return x.ErrCode == 0
}

func (x *respCommon) TryIntoErr() error {
	if x.IsOK() {
		return nil
	}

	return &WorkwxClientError{
		Code: x.ErrCode,
		Msg:  x.ErrMsg,
	}
}

type respAccessToken struct {
	respCommon

	AccessToken   string `json:"access_token"`
	ExpiresInSecs int64  `json:"expires_in"`
}

type reqJSAPITicketAgentConfig struct{}

var _ urlValuer = reqJSAPITicketAgentConfig{}

func (x reqJSAPITicketAgentConfig) intoURLValues() url.Values {
	return url.Values{
		"type": {"agent_config"},
	}
}

type reqJSAPITicket struct{}

var _ urlValuer = reqJSAPITicket{}

func (x reqJSAPITicket) intoURLValues() url.Values {
	return url.Values{}
}

type respJSAPITicket struct {
	respCommon

	Ticket        string `json:"ticket"`
	ExpiresInSecs int64  `json:"expires_in"`
}

// reqMessage 消息发送请求
type reqMessage struct {
	ToUser  []string
	ToParty []string
	ToTag   []string
	ChatID  string
	AgentID int64
	MsgType string
	Content map[string]interface{}
	IsSafe  bool
}

var _ bodyer = reqMessage{}

func (x reqMessage) intoBody() ([]byte, error) {
	// fuck
	safeInt := 0
	if x.IsSafe {
		safeInt = 1
	}

	obj := map[string]interface{}{
		"msgtype": x.MsgType,
		"agentid": x.AgentID,
		"safe":    safeInt,
	}

	// msgtype polymorphism
	if x.MsgType != "template_card" {
		obj[x.MsgType] = x.Content
	} else {
		obj[x.MsgType] = x.Content["template_card"]
	}

	// 复用这个结构体，因为是 package-private 的所以这么做没风险
	if x.ChatID != "" {
		obj["chatid"] = x.ChatID
	} else {
		obj["touser"] = strings.Join(x.ToUser, "|")
		obj["toparty"] = strings.Join(x.ToParty, "|")
		obj["totag"] = strings.Join(x.ToTag, "|")
	}

	return marshalIntoJSONBody(obj)
}

// respMessageSend 消息发送响应
type respMessageSend struct {
	respCommon

	InvalidUsers   string `json:"invaliduser"`
	InvalidParties string `json:"invalidparty"`
	InvalidTags    string `json:"invalidtag"`
}

type reqUserGet struct {
	UserID string
}

var _ urlValuer = reqUserGet{}

func (x reqUserGet) intoURLValues() url.Values {
	return url.Values{
		"userid": {x.UserID},
	}
}

// respUserGet 读取成员响应
type respUserGet struct {
	respCommon

	UserDetail
}

// reqUserUpdate 更新成员请求
type reqUserUpdate struct {
	UserDetail *UserDetail
}

var _ bodyer = reqUserUpdate{}

func (x reqUserUpdate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.UserDetail)
}

// respUserUpdate 更新成员响应
type respUserUpdate struct {
	respCommon
}

// reqUserList 部门成员请求
type reqUserList struct {
	DeptID     int64
	FetchChild bool
}

var _ urlValuer = reqUserList{}

func (x reqUserList) intoURLValues() url.Values {
	var fetchChild int64
	if x.FetchChild {
		fetchChild = 1
	}

	return url.Values{
		"department_id": {strconv.FormatInt(x.DeptID, 10)},
		"fetch_child":   {strconv.FormatInt(fetchChild, 10)},
	}
}

// respUsersByDeptID 部门成员详情响应
type respUserList struct {
	respCommon

	Users []*UserDetail `json:"userlist"`
}

// reqConvertUserIDToOpenID userid转openid 请求
type reqConvertUserIDToOpenID struct {
	UserID string `json:"userid"`
}

var _ bodyer = reqConvertUserIDToOpenID{}

// respConvertUserIDToOpenID userid转openid 响应
type respConvertUserIDToOpenID struct {
	respCommon

	OpenID string `json:"openid"`
}

func (x reqConvertUserIDToOpenID) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// reqConvertOpenIDToUserID openid转userid 请求
type reqConvertOpenIDToUserID struct {
	OpenID string `json:"openid"`
}

var _ bodyer = reqConvertOpenIDToUserID{}

// respConvertUserIDToOpenID openid转userid 响应
type respConvertOpenIDToUserID struct {
	respCommon

	UserID string `json:"userid"`
}

func (x reqConvertOpenIDToUserID) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// SizeType qrcode尺寸类型
//
// 1: 171 x 171; 2: 399 x 399; 3: 741 x 741; 4: 2052 x 2052
type SizeType int

const (
	// SizeTypeMini 171 x 171
	SizeTypeMini SizeType = iota + 1
	// SizeTypeSmall 399 x 399
	SizeTypeSmall
	// SizeTypeMedium 741 x 741
	SizeTypeMedium
	// SizeTypeLarge 2052 x 2052
	SizeTypeLarge
)

// reqUserJoinQrcode 获取加入企业二维码 请求
type reqUserJoinQrcode struct {
	SizeType SizeType `json:"size_type"`
}

var _ urlValuer = reqUserJoinQrcode{}

func (x reqUserJoinQrcode) intoURLValues() url.Values {
	return url.Values{
		"size_type": {strconv.Itoa(int(x.SizeType))},
	}
}

// respUserJoinQrcode 获取加入企业二维码 响应
type respUserJoinQrcode struct {
	respCommon

	JoinQrcode string `json:"join_qrcode"`
}

// reqUserIDByMobile 手机号获取 userid 请求
type reqUserIDByMobile struct {
	Mobile string `json:"mobile"`
}

var _ bodyer = reqUserIDByMobile{}

func (x reqUserIDByMobile) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respUserIDByMobile 手机号获取 userid 响应
type respUserIDByMobile struct {
	respCommon

	UserID string `json:"userid"`
}

// EmailType 用户邮箱的类型
//
// 1表示用户邮箱是企业邮箱（默认）
// 2表示用户邮箱是个人邮箱
type EmailType int

const (
	// EmailTypeCorporate 企业邮箱
	EmailTypeCorporate EmailType = 1
	// EmailTypePersonal 个人邮箱
	EmailTypePersonal EmailType = 2
)

// reqUserIDByEmail 邮箱获取 userid 请求
type reqUserIDByEmail struct {
	Email     string    `json:"email"`
	EmailType EmailType `json:"email_type"`
}

var _ bodyer = reqUserIDByEmail{}

func (x reqUserIDByEmail) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respUserIDByEmail 邮箱获取 userid 响应
type respUserIDByEmail struct {
	respCommon

	UserID string `json:"userid"`
}

// reqDeptCreate 创建部门
type reqDeptCreate struct {
	DeptInfo *DeptInfo
}

var _ bodyer = reqDeptCreate{}

func (x reqDeptCreate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.DeptInfo)
}

// respDeptCreate 创建部门响应
type respDeptCreate struct {
	respCommon

	ID int64 `json:"id"`
}

// reqDeptList 获取部门列表
// 从2022年8月15日10点开始，“企业管理后台 - 管理工具 - 通讯录同步”的新增IP将不能再调用此接口，企业可通过「获取部门ID列表」接口获取部门ID列表。查看调整详情。
// https://developer.work.weixin.qq.com/document/path/96079
type reqDeptList struct {
	HaveID bool
	ID     int64
}

var _ urlValuer = reqDeptList{}

func (x reqDeptList) intoURLValues() url.Values {
	if !x.HaveID {
		return url.Values{}
	}

	return url.Values{
		"id": {strconv.FormatInt(x.ID, 10)},
	}
}

// respDeptList 部门列表响应
type respDeptList struct {
	respCommon

	// TODO: 不要懒惰，把 API 层的类型写好
	Department []*DeptInfo `json:"department"`
}

// reqDeptSimpleList 获取子部门ID列表
type reqDeptSimpleList struct {
	HaveID bool
	ID     int64
}

var _ urlValuer = reqDeptSimpleList{}

func (x reqDeptSimpleList) intoURLValues() url.Values {
	if !x.HaveID {
		return url.Values{}
	}

	return url.Values{
		"id": {strconv.FormatInt(x.ID, 10)},
	}
}

// respDeptSimpleList 部门列表响应
type respDeptSimpleList struct {
	respCommon

	DepartmentIDs []*DeptInfo `json:"department_id"`
}

// reqAppchatGet 获取群聊会话请求
type reqAppchatGet struct {
	ChatID string
}

var _ urlValuer = reqAppchatGet{}

func (x reqAppchatGet) intoURLValues() url.Values {
	return url.Values{
		"chatid": {x.ChatID},
	}
}

// respAppchatGet 获取群聊会话响应
type respAppchatGet struct {
	respCommon

	ChatInfo *ChatInfo `json:"chat_info"`
}

// reqAppchatCreate 创建群聊会话请求
type reqAppchatCreate struct {
	ChatInfo *ChatInfo
}

var _ bodyer = reqAppchatCreate{}

func (x reqAppchatCreate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.ChatInfo)
}

// respAppchatCreate 创建群聊会话响应
type respAppchatCreate struct {
	respCommon

	ChatID string `json:"chatid"`
}

// reqAppchatUpdate 修改群聊会话请求
type reqAppchatUpdate struct {
	ChatInfo
	AddMemberUserIDs []string `json:"add_user_list"`
	DelMemberUserIDs []string `json:"del_user_list"`
}

var _ bodyer = reqAppchatUpdate{}

func (x reqAppchatUpdate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respAppchatUpdate 修改群聊会话响应
type respAppchatUpdate struct {
	respCommon
}

// reqMediaUpload 临时素材上传请求
type reqMediaUpload struct {
	Type  string
	Media *Media
}

var _ urlValuer = reqMediaUpload{}
var _ mediaUploader = reqMediaUpload{}

func (x reqMediaUpload) intoURLValues() url.Values {
	return url.Values{
		"type": {x.Type},
	}
}

func (x reqMediaUpload) getMedia() *Media {
	return x.Media
}

// respMediaUpload 临时素材上传响应
type respMediaUpload struct {
	respCommon

	Type      string `json:"type"`
	MediaID   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

// reqMediaUploadImg 永久图片素材上传请求
type reqMediaUploadImg struct {
	Media *Media
}

var _ urlValuer = reqMediaUploadImg{}
var _ mediaUploader = reqMediaUploadImg{}

func (x reqMediaUploadImg) intoURLValues() url.Values {
	return url.Values{}
}

func (x reqMediaUploadImg) getMedia() *Media {
	return x.Media
}

// respMediaUploadImg 永久图片素材上传响应
type respMediaUploadImg struct {
	respCommon

	URL string `json:"url"`
}

// reqExternalContactList 获取客户列表
type reqExternalContactList struct {
	UserID string `json:"userid"`
}

var _ urlValuer = reqExternalContactList{}

func (x reqExternalContactList) intoURLValues() url.Values {
	return url.Values{
		"userid": {x.UserID},
	}
}

// respExternalContactList 获取客户列表
type respExternalContactList struct {
	respCommon

	ExternalUserID []string `json:"external_userid"`
}

// reqExternalContactGet 获取客户详情
type reqExternalContactGet struct {
	ExternalUserID string `json:"external_userid"`
}

var _ urlValuer = reqExternalContactGet{}

func (x reqExternalContactGet) intoURLValues() url.Values {
	return url.Values{
		"external_userid": {x.ExternalUserID},
	}
}

// respExternalContactGet 获取客户详情
type respExternalContactGet struct {
	respCommon
	ExternalContactInfo
}

// ExternalContactInfo 外部联系人信息
type ExternalContactInfo struct {
	ExternalContact ExternalContact `json:"external_contact"`
	FollowUser      []FollowUser    `json:"follow_user"`
}

// ExternalContactBatchInfo 外部联系人信息
type ExternalContactBatchInfo struct {
	ExternalContact ExternalContact `json:"external_contact"`
	FollowInfo      FollowInfo      `json:"follow_info"`
}

// BatchListExternalContactsResp 外部联系人信息
type BatchListExternalContactsResp struct {
	Result     []ExternalContactBatchInfo
	NextCursor string
}

// reqExternalContactBatchList 批量获取客户详情
type reqExternalContactBatchList struct {
	UserID string `json:"userid"`
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

var _ bodyer = reqExternalContactBatchList{}

func (x reqExternalContactBatchList) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respExternalContactBatchList 批量获取客户详情
type respExternalContactBatchList struct {
	respCommon
	NextCursor          string                     `json:"next_cursor"`
	ExternalContactList []ExternalContactBatchInfo `json:"external_contact_list"`
}

// reqExternalContactRemark 获取客户详情
type reqExternalContactRemark struct {
	Remark *ExternalContactRemark
}

var _ bodyer = reqExternalContactRemark{}

func (x reqExternalContactRemark) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.Remark)
}

// respExternalContactRemark 获取客户详情
type respExternalContactRemark struct {
	respCommon
}

// reqUserInfoGet 获取访问用户身份
type reqUserInfoGet struct {
	// 通过成员授权获取到的code，最大为512字节。每次成员授权带上的code将不一样，code只能使用一次，5分钟未被使用自动过期。
	Code string
}

var _ urlValuer = reqUserInfoGet{}

func (x reqUserInfoGet) intoURLValues() url.Values {
	return url.Values{
		"code": {x.Code},
	}
}

// respUserInfoGet 部门列表响应
type respUserInfoGet struct {
	respCommon
	UserIdentityInfo
}

// reqExternalContactListCorpTags 获取企业标签库
type reqExternalContactListCorpTags struct {
	// 要查询的标签id，如果不填则获取该企业的所有客户标签，目前暂不支持标签组id
	TagIDs []string `json:"tag_id"`
}

var _ bodyer = reqExternalContactListCorpTags{}

func (x reqExternalContactListCorpTags) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respExternalContactListCorpTags 获取企业标签库
type respExternalContactListCorpTags struct {
	respCommon
	// 标签组列表
	TagGroup []ExternalContactCorpTagGroup `json:"tag_group"`
}

// reqExternalContactAddCorpTag 添加企业客户标签
type reqExternalContactAddCorpTag struct {
	ExternalContactCorpTagGroup
}

var _ bodyer = reqExternalContactAddCorpTag{}

func (x reqExternalContactAddCorpTag) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.ExternalContactCorpTagGroup)
}

// respExternalContactAddCorpTag 添加企业客户标签
type respExternalContactAddCorpTag struct {
	respCommon
	// 标签组列表
	TagGroup ExternalContactCorpTagGroup `json:"tag_group"`
}

// reqExternalContactEditCorpTag 编辑企业客户标签
type reqExternalContactEditCorpTag struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order uint32 `json:"order"`
}

var _ bodyer = reqExternalContactEditCorpTag{}

func (x reqExternalContactEditCorpTag) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respExternalContactEditCorpTag 编辑企业客户标签
type respExternalContactEditCorpTag struct {
	respCommon
}

// reqExternalContactDelCorpTag 删除企业客户标签
type reqExternalContactDelCorpTag struct {
	TagID   []string `json:"tag_id"`
	GroupID []string `json:"group_id"`
}

var _ bodyer = reqExternalContactDelCorpTag{}

func (x reqExternalContactDelCorpTag) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respExternalContactDelCorpTag 删除企业客户标签
type respExternalContactDelCorpTag struct {
	respCommon
}

// reqExternalContactMarkTag 编辑企业客户标签
type reqExternalContactMarkTag struct {
	UserID         string   `json:"userid"`
	ExternalUserID string   `json:"external_userid"`
	AddTag         []string `json:"add_tag"`
	RemoveTag      []string `json:"remove_tag"`
}

var _ bodyer = reqExternalContactMarkTag{}

func (x reqExternalContactMarkTag) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respExternalContactMarkTag 编辑企业客户标签
type respExternalContactMarkTag struct {
	respCommon
}

// reqJSCode2Session 临时登录凭证校验
type reqJSCode2Session struct {
	JSCode string
}

var _ urlValuer = reqJSCode2Session{}

func (x reqJSCode2Session) intoURLValues() url.Values {
	return url.Values{
		"js_code":    {x.JSCode},
		"grant_type": {"authorization_code"},
	}
}

// respJSCode2Session 临时登录凭证校验
type respJSCode2Session struct {
	respCommon
	JSCodeSession
}

// JSCodeSession 临时登录凭证
type JSCodeSession struct {
	CorpID     string `json:"corpid"`
	UserID     string `json:"userid"`
	SessionKey string `json:"session_key"`
}

// reqAuthCode2UserInfo 获取访问用户身份
type reqAuthCode2UserInfo struct {
	Code string
}

var _ urlValuer = reqAuthCode2UserInfo{}

func (x reqAuthCode2UserInfo) intoURLValues() url.Values {
	return url.Values{
		"code": {x.Code},
	}
}

// respAuthCode2UserInfo 获取访问用户身份响应
type respAuthCode2UserInfo struct {
	respCommon
	AuthCodeUserInfo
}

// AuthCodeUserInfo 访问用户身份
type AuthCodeUserInfo struct {
	UserID         string `json:"userid,omitempty"`
	UserTicket     string `json:"user_ticket,omitempty"`
	OpenID         string `json:"openid,omitempty"`
	ExternalUserID string `json:"external_userid,omitempty"`
}

type reqMsgAuditListPermitUser struct {
	MsgAuditEdition MsgAuditEdition `json:"type"`
}

var _ bodyer = reqMsgAuditListPermitUser{}

func (x reqMsgAuditListPermitUser) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respMsgAuditListPermitUser struct {
	respCommon
	IDs []string `json:"ids"`
}

type reqMsgAuditCheckSingleAgree struct {
	Infos []CheckMsgAuditSingleAgreeUserInfo `json:"info"`
}

var _ bodyer = reqMsgAuditCheckSingleAgree{}

func (x reqMsgAuditCheckSingleAgree) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respMsgAuditCheckSingleAgree struct {
	respCommon
	AgreeInfo []struct {
		UserID           string              `json:"userid"`
		ExternalOpenID   string              `json:"exteranalopenid"`
		AgreeStatus      MsgAuditAgreeStatus `json:"agree_status"`
		StatusChangeTime int                 `json:"status_change_time"`
	} `json:"agreeinfo"`
}

func (x respMsgAuditCheckSingleAgree) intoCheckSingleAgreeInfoList() (resp []CheckMsgAuditSingleAgreeInfo) {
	for _, agreeInfo := range x.AgreeInfo {
		resp = append(resp, CheckMsgAuditSingleAgreeInfo{
			CheckMsgAuditSingleAgreeUserInfo: CheckMsgAuditSingleAgreeUserInfo{
				UserID:         agreeInfo.UserID,
				ExternalOpenID: agreeInfo.ExternalOpenID,
			},
			AgreeStatus:      agreeInfo.AgreeStatus,
			StatusChangeTime: time.Unix(int64(agreeInfo.StatusChangeTime), 0),
		})
	}
	return resp
}

type reqMsgAuditCheckRoomAgree struct {
	RoomID string `json:"roomid"`
}

var _ bodyer = reqMsgAuditCheckRoomAgree{}

func (x reqMsgAuditCheckRoomAgree) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respMsgAuditCheckRoomAgree struct {
	respCommon
	AgreeInfo []struct {
		StatusChangeTime int                 `json:"status_change_time"`
		AgreeStatus      MsgAuditAgreeStatus `json:"agree_status"`
		ExternalOpenID   string              `json:"exteranalopenid"`
	} `json:"agreeinfo"`
}

func (x respMsgAuditCheckRoomAgree) intoCheckRoomAgreeInfoList() (resp []CheckMsgAuditRoomAgreeInfo) {
	for _, agreeInfo := range x.AgreeInfo {
		resp = append(resp, CheckMsgAuditRoomAgreeInfo{
			StatusChangeTime: time.Unix(int64(agreeInfo.StatusChangeTime), 0),
			AgreeStatus:      agreeInfo.AgreeStatus,
			ExternalOpenID:   agreeInfo.ExternalOpenID,
		})
	}
	return resp
}

type reqMsgAuditGetGroupChat struct {
	RoomID string `json:"roomid"`
}

var _ bodyer = reqMsgAuditGetGroupChat{}

func (x reqMsgAuditGetGroupChat) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respMsgAuditGetGroupChat struct {
	respCommon
	Members []struct {
		MemberID string `json:"memberid"`
		JoinTime int    `json:"jointime"`
	} `json:"members"`
	RoomName       string `json:"roomname"`
	Creator        string `json:"creator"`
	RoomCreateTime int    `json:"room_create_time"`
	Notice         string `json:"notice"`
}

func (x respMsgAuditGetGroupChat) intoGroupChat() (resp MsgAuditGroupChat) {
	resp.Creator = x.Creator
	resp.Notice = x.Notice
	resp.RoomName = x.RoomName
	resp.RoomCreateTime = time.Unix(int64(x.RoomCreateTime), 0)
	for _, member := range x.Members {
		resp.Members = append(resp.Members, MsgAuditGroupChatMember{
			MemberID: member.MemberID,
			JoinTime: time.Unix(int64(member.JoinTime), 0),
		})
	}
	return resp
}

type reqListUnassignedExternalContact struct {
	// PageID 分页查询，要查询页号，从0开始
	PageID uint32 `json:"page_id"`
	// PageSize 每次返回的最大记录数，默认为1000，最大值为1000
	PageSize uint32 `json:"page_size"`
	// Cursor 分页查询游标，字符串类型，适用于数据量较大的情况，如果使用该参数则无需填写page_id，该参数由上一次调用返回
	Cursor string `json:"cursor"`
}

var _ bodyer = reqListUnassignedExternalContact{}

func (x reqListUnassignedExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respListUnassignedExternalContact struct {
	respCommon
	Info []struct {
		HandoverUserID string `json:"handover_userid"`
		ExternalUserID string `json:"external_userid"`
		DemissionTime  int    `json:"dimission_time"`
	} `json:"info"`
	IsLast     bool   `json:"is_last"`
	NextCursor string `json:"next_cursor"`
}

func (x respListUnassignedExternalContact) intoExternalContactUnassignedList() (resp ExternalContactUnassignedList) {
	list := make([]ExternalContactUnassigned, 0, len(x.Info))
	for _, info := range x.Info {
		list = append(list, ExternalContactUnassigned{
			HandoverUserID: info.HandoverUserID,
			ExternalUserID: info.ExternalUserID,
			DemissionTime:  time.Unix(int64(info.DemissionTime), 0),
		})
	}
	resp.Info = list
	resp.IsLast = x.IsLast
	resp.NextCursor = x.NextCursor
	return resp
}

type reqTransferExternalContact struct {
	// ExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	ExternalUserID string `json:"external_userid"`
	// HandoverUserID 原跟进成员的userid
	HandoverUserID string `json:"handover_userid"`
	// TakeoverUserID 接替成员的userid
	TakeoverUserID string `json:"takeover_userid"`
	// TransferSuccessMsg 转移成功后发给客户的消息，最多200个字符，不填则使用默认文案，目前只对在职成员分配客户的情况生效
	TransferSuccessMsg string `json:"transfer_success_msg"`
}

var _ bodyer = reqTransferExternalContact{}

func (x reqTransferExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respTransferExternalContact struct {
	respCommon
}

type reqGetTransferExternalContactResult struct {
	// ExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	ExternalUserID string `json:"external_userid"`
	// HandoverUserID 原跟进成员的userid
	HandoverUserID string `json:"handover_userid"`
	// TakeoverUserID 接替成员的userid
	TakeoverUserID string `json:"takeover_userid"`
}

var _ bodyer = reqGetTransferExternalContactResult{}

func (x reqGetTransferExternalContactResult) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respGetTransferExternalContactResult struct {
	respCommon
	Status       uint8 `json:"status"`
	TakeoverTime int   `json:"takeover_time"`
}

func (x respGetTransferExternalContactResult) intoExternalContactTransferResult() ExternalContactTransferResult {
	return ExternalContactTransferResult{
		Status:       ExternalContactTransferStatus(x.Status),
		TakeoverTime: time.Unix(int64(x.TakeoverTime), 0),
	}
}

type reqTransferGroupChatExternalContact struct {
	// ChatIDList 需要转群主的客户群ID列表。取值范围： 1 ~ 100
	ChatIDList []string `json:"chat_id_list"`
	// NewOwner 新群主ID
	NewOwner string `json:"new_owner"`
}

var _ bodyer = reqTransferGroupChatExternalContact{}

func (x reqTransferGroupChatExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respTransferGroupChatExternalContact struct {
	respCommon
	FailedChatList []ExternalContactGroupChatTransferFailed `json:"failed_chat_list"`
}

type reqOAGetTemplateDetail struct {
	TemplateID string `json:"template_id"`
}

var _ bodyer = reqOAGetTemplateDetail{}

func (x reqOAGetTemplateDetail) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respOAGetTemplateDetail struct {
	respCommon
	OATemplateDetail
}

type reqOAApplyEvent struct {
	OAApplyEvent
}

var _ bodyer = reqOAApplyEvent{}

func (x reqOAApplyEvent) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respOAApplyEvent struct {
	respCommon
	// SpNo 表单提交成功后，返回的表单编号
	SpNo string `json:"sp_no"`
}

type reqOAGetApprovalInfo struct {
	StartTime string                 `json:"starttime"`
	EndTime   string                 `json:"endtime"`
	Cursor    int                    `json:"cursor"`
	Size      uint32                 `json:"size"`
	Filters   []OAApprovalInfoFilter `json:"filters"`
}

var _ bodyer = reqOAGetApprovalInfo{}

func (x reqOAGetApprovalInfo) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respOAGetApprovalInfo struct {
	respCommon
	// SpNoList 审批单号列表，包含满足条件的审批申请
	SpNoList []string `json:"sp_no_list"`
}

type reqOAGetApprovalDetail struct {
	// SpNo 审批单编号。
	SpNo string `json:"sp_no"`
}

var _ bodyer = reqOAGetApprovalDetail{}

func (x reqOAGetApprovalDetail) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respOAGetApprovalDetail struct {
	respCommon
	// Info 审批申请详情
	Info OAApprovalDetail `json:"info"`
}

// TaskCardBtn 任务卡片消息按钮
type TaskCardBtn struct {
	// Key 按钮key值，用户点击后，会产生任务卡片回调事件，回调事件会带上该key值，只能由数字、字母和“_-@”组成，最长支持128字节
	Key string `json:"key"`
	// Name 按钮名称
	Name string `json:"name"`
	// ReplaceName 点击按钮后显示的名称，默认为“已处理”
	ReplaceName string `json:"replace_name"`
	// Color 按钮字体颜色，可选“red”或者“blue”,默认为“blue”
	Color string `json:"color"`
	// IsBold 按钮字体是否加粗，默认false
	IsBold bool `json:"is_bold"`
}

// Article news 类型的文章
type Article struct {
	// 标题，不超过128个字节，超过会自动截断（支持id转译）
	Title string `json:"title"`
	// 描述，不超过512个字节，超过会自动截断（支持id转译）
	Description string `json:"description"`
	// 点击后跳转的链接。 最长2048字节，请确保包含了协议头(http/https)，小程序或者url必须填写一个
	URL string `json:"url"`
	// 图文消息的图片链接，最长2048字节，支持JPG、PNG格式，较好的效果为大图 1068*455，小图150*150
	PicURL string `json:"picurl"`
	// 小程序appid，必须是与当前应用关联的小程序，appid和pagepath必须同时填写，填写后会忽略url字段
	AppID string `json:"appid"`
	// 点击消息卡片后的小程序页面，最长128字节，仅限本小程序内的页面。appid和pagepath必须同时填写，填写后会忽略url字段
	PagePath string `json:"pagepath"`
}

// MPArticle mpnews 类型的文章
type MPArticle struct {
	// 标题，不超过128个字节，超过会自动截断（支持id转译）
	Title string `json:"title"`
	// 图文消息缩略图的media_id, 可以通过素材管理接口获得。此处thumb_media_id即上传接口返回的media_id
	ThumbMediaID string `json:"thumb_media_id"`
	// 图文消息的作者，不超过64个字节
	Author string `json:"author"`
	// 图文消息点击“阅读原文”之后的页面链接
	ContentSourceURL string `json:"content_source_url"`
	// 图文消息的内容，支持html标签，不超过666 K个字节（支持id转译）
	Content string `json:"content"`
	// 图文消息的描述，不超过512个字节，超过会自动截断（支持id转译）
	Digest string `json:"digest"`
}

// Source 卡片来源样式信息，不需要来源样式可不填写
type Source struct {
	// 来源图片的url，来源图片的尺寸建议为72*72
	IconURL string `json:"icon_url"`
	// 来源图片的描述，建议不超过20个字，（支持id转译）
	Desc string `json:"desc"`
	// 来源文字的颜色，目前支持：0(默认) 灰色，1 黑色，2 红色，3 绿色
	DescColor int `json:"desc_color"`
}

// ActionList 操作列表，列表长度取值范围为 [1, 3]
type ActionList struct {
	// 操作的描述文案
	Text string `json:"text"`
	// 操作key值，用户点击后，会产生回调事件将本参数作为EventKey返回，回调事件会带上该key值，最长支持1024字节，不可重复
	Key string `json:"key"`
}

// ActionMenu 卡片右上角更多操作按钮
type ActionMenu struct {
	// 更多操作界面的描述
	Desc       string       `json:"desc"`
	ActionList []ActionList `json:"action_list"`
}

// MainTitle 一级标题
type MainTitle struct {
	// 一级标题，建议不超过36个字，文本通知型卡片本字段非必填，但不可本字段和sub_title_text都不填，（支持id转译）
	Title string `json:"title"`
	// 标题辅助信息，建议不超过160个字，（支持id转译）
	Desc string `json:"desc"`
}

// QuoteArea 引用文献样式
type QuoteArea struct {
	// 引用文献样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序
	Type int `json:"type"`
	// 点击跳转的url，quote_area.type是1时必填
	URL string `json:"url"`
	// 引用文献样式的标题
	Title string `json:"title"`
	// 引用文献样式的引用文案
	QuoteText string `json:"quote_text"`
	// 小程序appid，必须是与当前应用关联的小程序，appid和pagepath必须同时填写，填写后会忽略url字段
	AppID string `json:"appid"`
	// 点击消息卡片后的小程序页面，最长128字节，仅限本小程序内的页面。appid和pagepath必须同时填写，填写后会忽略url字段
	PagePath string `json:"pagepath"`
}

// EmphasisContent 关键数据样式
type EmphasisContent struct {
	// 关键数据样式的数据内容，建议不超过14个字
	Title string `json:"title"`
	// 关键数据样式的数据描述内容，建议不超过22个字
	Desc string `json:"desc"`
}

// HorizontalContentList 二级标题+文本列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过6
type HorizontalContentList struct {
	// 二级标题，建议不超过5个字
	KeyName string `json:"keyname"`
	// 二级文本，如果horizontal_content_list.type是2，该字段代表文件名称（要包含文件类型），建议不超过30个字，（支持id转译）
	Value string `json:"value"`
	// 链接类型，0或不填代表不是链接，1 代表跳转url，2 代表下载附件，3 代表点击跳转成员详情
	Type int `json:"type,omitempty"`
	// 链接跳转的url，horizontal_content_list.type是1时必填
	URL string `json:"url,omitempty"`
	// 附件的media_id，horizontal_content_list.type是2时必填
	MediaID string `json:"media_id,omitempty"`
	// 成员详情的userid，horizontal_content_list.type是3时必填
	Userid string `json:"userid,omitempty"`
}

// JumpList 跳转指引样式的列表，该字段可为空数组，但有数据的话需确认对应字段是否必填，列表长度不超过3
type JumpList struct {
	// 跳转链接类型，0或不填代表不是链接，1 代表跳转url，2 代表跳转小程序
	Type int `json:"type"`
	// 跳转链接样式的文案内容，建议不超过18个字
	Title string `json:"title"`
	// 跳转链接的url，jump_list.type是1时必填
	URL string `json:"url,omitempty"`
	// 跳转链接的小程序的appid，必须是与当前应用关联的小程序，jump_list.type是2时必填
	Appid string `json:"appid,omitempty"`
	// 跳转链接的小程序的pagepath，jump_list.type是2时选填
	PagePath string `json:"pagepath,omitempty"`
}

// CardAction 整体卡片的点击跳转事件，text_notice必填本字段
type CardAction struct {
	// 跳转事件类型，1 代表跳转url，2 代表打开小程序。text_notice卡片模版中该字段取值范围为[1,2]
	Type int `json:"type"`
	// 跳转事件的url，card_action.type是1时必填
	URL string `json:"url"`
	// 跳转事件的小程序的appid，必须是与当前应用关联的小程序，card_action.type是2时必填
	Appid string `json:"appid"`
	// 跳转事件的小程序的pagepath，card_action.type是2时选填
	Pagepath string `json:"pagepath"`
}

// ImageTextArea 左图右文样式，news_notice类型的卡片，card_image和image_text_area两者必填一个字段，不可都不填
type ImageTextArea struct {
	// 左图右文样式区域点击事件，0或不填代表没有点击事件，1 代表跳转url，2 代表跳转小程序
	Type int `json:"type"`
	// 点击跳转的url，image_text_area.type是1时必填
	URL string `json:"url"`
	// 点击跳转的小程序的appid，必须是与当前应用关联的小程序，image_text_area.type是2时必填
	AppID string `json:"appid,omitempty"`
	// 点击跳转的小程序的pagepath，image_text_area.type是2时选填
	PagePath string `json:"pagepath,omitempty"`
	// 左图右文样式的标题
	Title string `json:"title"`
	// 左图右文样式的描述
	Desc string `json:"desc"`
	// 左图右文样式的图片url
	ImageURL string `json:"image_url"`
}

// CardImage 图片样式，news_notice类型的卡片，card_image和image_text_area两者必填一个字段，不可都不填
type CardImage struct {
	// 图片的url
	URL string `json:"url"`
	// 图片的宽高比，宽高比要小于2.25，大于1.3，不填该参数默认1.3
	AspectRatio float32 `json:"aspect_ratio"`
}

// ButtonSelection 按钮交互型
type ButtonSelection struct {
	// 下拉式的选择器的key，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节
	QuestionKey string `json:"question_key"`
	// 下拉式的选择器的key，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节
	Title string `json:"title"`
	// 选项列表，下拉选项不超过 10 个，最少1个
	OptionList []struct {
		// 下拉式的选择器选项的id，用户提交后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复
		ID string `json:"id"`
		// 下拉式的选择器选项的文案，建议不超过16个字
		Text string `json:"text"`
	} `json:"option_list"`
	// 默认选定的id，不填或错填默认第一个
	SelectedID string `json:"selected_id"`
}

type Button struct {
	// 按钮点击事件类型，0 或不填代表回调点击事件，1 代表跳转url
	Type int `json:"type,omitempty"`
	// 按钮文案，建议不超过10个字
	Text string `json:"text"`
	// 按钮样式，目前可填1~4，不填或错填默认1
	Style int `json:"style,omitempty"`
	// 按钮key值，用户点击后，会产生回调事件将本参数作为EventKey返回，回调事件会带上该key值，最长支持1024字节，不可重复，button_list.type是0时必填
	Key string `json:"key,omitempty"`
	// 跳转事件的url，button_list.type是1时必填
	URL string `json:"url,omitempty"`
}

// CheckBox 选择题样式
type CheckBox struct {
	// 选择题key值，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节
	QuestionKey string `json:"question_key"`
	// 选项list，选项个数不超过 20 个，最少1个
	OptionList []struct {
		// 选项id，用户提交选项后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复
		ID string `json:"id"`
		// 选项文案描述，建议不超过17个字
		Text string `json:"text"`
		// 该选项是否要默认选中
		IsChecked bool `json:"is_checked"`
	} `json:"option_list" validate:"required,min=1,max=20"`
	// 选择题模式，单选：0，多选：1，不填默认0
	Mode int `json:"mode" validate:"omitempty,oneof=0 1"`
}

// SubmitButton 提交按钮样式
type SubmitButton struct {
	// 按钮文案，建议不超过10个字，不填默认为提交
	Text string `json:"text"`
	// 提交按钮的key，会产生回调事件将本参数作为EventKey返回，最长支持1024字节
	Key string `json:"key"`
}

// SelectList 下拉式的选择器列表，multiple_interaction类型的卡片该字段不可为空，一个消息最多支持 3 个选择器
type SelectList struct {
	// 下拉式的选择器题目的key，用户提交选项后，会产生回调事件，回调事件会带上该key值表示该题，最长支持1024字节，不可重复
	QuestionKey string `json:"question_key"`
	// 下拉式的选择器上面的title
	Title string `json:"title,omitempty"`
	// 默认选定的id，不填或错填默认第一个
	SelectedID string       `json:"selected_id,omitempty"`
	OptionList []OptionList `json:"option_list"`
}

// 项列表，下拉选项不超过 10 个，最少1个
type OptionList struct {
	// 下拉式的选择器选项的id，用户提交选项后，会产生回调事件，回调事件会带上该id值表示该选项，最长支持128字节，不可重复
	ID string `json:"id"`
	// 下拉式的选择器选项的文案，建议不超过16个字
	Text string `json:"text"`
}

// TemplateCardType 模板卡片的类型
type TemplateCardType string

const (
	CardTypeTextNotice          TemplateCardType = "text_notice"          // 文本通知型
	CardTypeNewsNotice          TemplateCardType = "news_notice"          // 图文展示型
	CardTypeButtonInteraction   TemplateCardType = "button_interaction"   // 按钮交互型
	CardTypeVoteInteraction     TemplateCardType = "vote_interaction"     // 投票选择型
	CardTypeMultipleInteraction TemplateCardType = "multiple_interaction" // 多项选择型
)

type TemplateCard struct {
	CardType   TemplateCardType `json:"card_type"`
	Source     Source           `json:"source"`
	ActionMenu *ActionMenu      `json:"action_menu,omitempty" validate:"required_with=TaskID"`
	TaskID     string           `json:"task_id,omitempty" validate:"required_with=ActionMenu"`
	MainTitle  *MainTitle       `json:"main_title"`
	QuoteArea  *QuoteArea       `json:"quote_area,omitempty"`
	// 文本通知型
	EmphasisContent *EmphasisContent `json:"emphasis_content,omitempty"`
	SubTitleText    string           `json:"sub_title_text,omitempty"`
	// 图文展示型
	ImageTextArea         *ImageTextArea          `json:"image_text_area,omitempty"`
	CardImage             *CardImage              `json:"card_image,omitempty"`
	HorizontalContentList []HorizontalContentList `json:"horizontal_content_list"`
	JumpList              []JumpList              `json:"jump_list"`
	CardAction            *CardAction             `json:"card_action,omitempty"`
	// 按钮交互型
	ButtonSelection *ButtonSelection `json:"button_selection,omitempty"`
	ButtonList      []Button         `json:"button_list,omitempty" validate:"omitempty,max=6"`
	// 投票选择型
	CheckBox     *CheckBox     `json:"checkbox,omitempty"`
	SelectList   []SelectList  `json:"select_list,omitempty" validate:"max=3"`
	SubmitButton *SubmitButton `json:"submit_button,omitempty"`
}

type TemplateCardUpdateMessage struct {
	UserIds      []string `json:"userids" validate:"omitempty,max=100"`
	PartyIds     []int64  `json:"partyids" validate:"omitempty,max=100"`
	TagIds       []int32  `json:"tagids" validate:"omitempty,max=100"`
	AtAll        int      `json:"atall,omitempty"`
	ResponseCode string   `json:"response_code"`
	Button       struct {
		ReplaceName string `json:"replace_name"`
	} `json:"button" validate:"required_without=TemplateCard"`
	TemplateCard TemplateCard `json:"template_card" validate:"required_without=Button"`
	ReplaceText  string       `json:"replace_text,omitempty"`
}

type reqTransferCustomer struct {
	// HandoverUserID 原跟进成员的userid
	HandoverUserID string `json:"handover_userid"`
	// TakeoverUserID 接替成员的userid
	TakeoverUserID string `json:"takeover_userid"`
	// ExternalUserID 客户的external_userid列表，每次最多分配100个客户
	ExternalUserID []string `json:"external_userid"`
	// TransferSuccessMsg 转移成功后发给客户的消息，最多200个字符，不填则使用默认文案
	TransferSuccessMsg string `json:"transfer_success_msg"`
}

var _ bodyer = reqTransferCustomer{}

func (x reqTransferCustomer) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respTransferCustomer struct {
	respCommon
	Customer []struct {
		// ExternalUserID 转接客户的外部联系人userid
		ExternalUserID string `json:"external_userid"`
		// Errcode 对此客户进行分配的结果, 具体可参考全局错误码(https://developer.work.weixin.qq.com/document/path/90475), 0表示成功发起接替,待24小时后自动接替,并不代表最终接替成功
		Errcode int `json:"errcode"`
	} `json:"customer"`
}

func (x respTransferCustomer) intoTransferCustomerResult() TransferCustomerResult {
	return x.Customer
}

type reqGetTransferCustomerResult struct {
	// HandoverUserID 原跟进成员的userid
	HandoverUserID string `json:"handover_userid"`
	// TakeoverUserID 接替成员的userid
	TakeoverUserID string `json:"takeover_userid"`
	// Cursor 分页查询的cursor，每个分页返回的数据不会超过1000条；不填或为空表示获取第一个分页
	Cursor string `json:"cursor"`
}

var _ bodyer = reqGetTransferCustomerResult{}

func (x reqGetTransferCustomerResult) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respGetTransferCustomerResult struct {
	respCommon
	Customer []struct {
		// ExternalUserID 转接客户的外部联系人userid
		ExternalUserID string `json:"external_userid"`
		// Status 接替状态， 1-接替完毕 2-等待接替 3-客户拒绝 4-接替成员客户达到上限 5-无接替记录
		Status int `json:"status"`
		// TakeoverTime 接替客户的时间，如果是等待接替状态，则为未来的自动接替时间
		TakeoverTime int `json:"takeover_time"`
	} `json:"customer"`
	// NextCursor 下个分页的起始cursor
	NextCursor string `json:"next_cursor"`
}

func (x respGetTransferCustomerResult) intoCustomerTransferResult() CustomerTransferResult {
	return CustomerTransferResult{
		Customer:   x.Customer,
		NextCursor: x.NextCursor,
	}
}

type reqListFollowUserExternalContact struct {
}

var _ urlValuer = reqListFollowUserExternalContact{}

func (x reqListFollowUserExternalContact) intoURLValues() url.Values {
	return url.Values{}
}

type respListFollowUserExternalContact struct {
	respCommon
	ExternalContactFollowUserList
}

type reqAddContactExternalContact struct {
	ExternalContactWay
}

var _ bodyer = reqAddContactExternalContact{}

func (x reqAddContactExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respAddContactExternalContact struct {
	respCommon
	ExternalContactAddContact
}

type ExternalContactAddContact struct {
	ConfigID string `json:"config_id"`
	QRCode   string `json:"qr_code"`
}

type reqGetContactWayExternalContact struct {
	ConfigID string `json:"config_id"`
}

var _ bodyer = reqGetContactWayExternalContact{}

func (x reqGetContactWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respGetContactWayExternalContact struct {
	respCommon
	ContactWay ExternalContactContactWay `json:"contact_way"`
}

type ExternalContactContactWay struct {
	ConfigID string `json:"config_id"`
	QRCode   string `json:"qr_code"`
	ExternalContactWay
}

var _ bodyer = reqListContactWayExternalContact{}

func (x reqListContactWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respListContactWayChatExternalContact struct {
	respCommon
	ExternalContactListContactWayChat
}

type ExternalContactListContactWayChat struct {
	NextCursor string       `json:"next_cursor"`
	ContactWay []contactWay `json:"contact_way"`
}

type contactWay struct {
	ConfigID string `json:"config_id"`
}

var _ bodyer = reqUpdateContactWayExternalContact{}

func (x reqUpdateContactWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respUpdateContactWayExternalContact struct {
	respCommon
}

type reqDelContactWayExternalContact struct {
	ConfigID string `json:"config_id"`
}

var _ bodyer = reqDelContactWayExternalContact{}

func (x reqDelContactWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respDelContactWayExternalContact struct {
	respCommon
}

type reqGroupChatList struct {
	ReqChatList ReqChatList
}

func (x reqGroupChatList) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.ReqChatList)
}

var _ bodyer = reqGroupChatList{}

type respGroupChatList struct {
	respCommon
	*RespGroupChatList
}

type reqGroupChatInfo struct {
	ChatID   string `json:"chat_id"`
	NeedName int64  `json:"need_name"`
}

func (x reqGroupChatInfo) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

var _ bodyer = reqGroupChatInfo{}

type respGroupChatInfo struct {
	respCommon
	GroupChat *RespGroupChatInfo `json:"group_chat"`
}

// reqConvertOpenGIDToChatID 客户群opengid转换 请求
type reqConvertOpenGIDToChatID struct {
	OpenGID string `json:"opengid"`
}

var _ bodyer = reqConvertOpenGIDToChatID{}

func (x reqConvertOpenGIDToChatID) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respConvertOpenGIDToChatID 客户群opengid转换 响应
type respConvertOpenGIDToChatID struct {
	respCommon

	ChatID string `json:"chat_id"`
}

type reqAddGroupChatJoinWayExternalContact struct {
	ExternalGroupChatJoinWay
}

var _ bodyer = reqAddGroupChatJoinWayExternalContact{}

func (x reqAddGroupChatJoinWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respAddGroupChatJoinWayExternalContact struct {
	respCommon

	ConfigID string `json:"config_id"`
}

type reqGetGroupChatJoinWayExternalContact struct {
	ConfigID string `json:"config_id"`
}

var _ bodyer = reqGetGroupChatJoinWayExternalContact{}

func (x reqGetGroupChatJoinWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respGetGroupChatJoinWayExternalContact struct {
	respCommon
	JoinWay ExternalContactGroupChatJoinWay `json:"join_way"`
}

type ExternalContactGroupChatJoinWay struct {
	ConfigID string `json:"config_id"`
	QRCode   string `json:"qr_code"`
	ExternalGroupChatJoinWay
}

type reqUpdateGroupChatJoinWayExternalContact struct {
	ConfigID string `json:"config_id"`
	ExternalGroupChatJoinWay
}

var _ bodyer = reqUpdateGroupChatJoinWayExternalContact{}

func (x reqUpdateGroupChatJoinWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respUpdateGroupChatJoinWayExternalContact struct {
	respCommon
}

type reqDelGroupChatJoinWayExternalContact struct {
	ConfigID string `json:"config_id"`
}

var _ bodyer = reqDelGroupChatJoinWayExternalContact{}

func (x reqDelGroupChatJoinWayExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respDelGroupChatJoinWayExternalContact struct {
	respCommon
}

type reqCloseTempChatExternalContact struct {
	UserID         string `json:"userid"`
	ExternalUserID string `json:"external_userid"`
}

var _ bodyer = reqCloseTempChatExternalContact{}

func (x reqCloseTempChatExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respCloseTempChatExternalContact struct {
	respCommon
}

type reqAddMsgTemplateExternalContact struct {
	AddMsgTemplateExternalContact
}

var _ bodyer = reqAddMsgTemplateExternalContact{}

func (x reqAddMsgTemplateExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respAddMsgTemplateExternalContact struct {
	respCommon
	AddMsgTemplateDetail
}

type AddMsgTemplateDetail struct {
	FailList []string `json:"fail_list"`
	MsgID    string   `json:"msgid"`
}

// reqSendWelcomeMsgExternalContact 发送新客户欢迎语
type reqSendWelcomeMsgExternalContact struct {
	SendWelcomeMsgExternalContact
}

var _ bodyer = reqSendWelcomeMsgExternalContact{}

func (x reqSendWelcomeMsgExternalContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

type respSendWelcomeMsgExternalContact struct {
	respCommon
}

// reqExternalContactAddCorpTag 添加企业客户标签
type reqExternalContactAddCorpTagGroup struct {
	ExternalContactAddCorpTagGroup
}

var _ bodyer = reqExternalContactAddCorpTagGroup{}

func (x reqExternalContactAddCorpTagGroup) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x.ExternalContactAddCorpTagGroup)
}

// reqKfAccountCreate 创建客服账号
type reqKfAccountCreate struct {
	Name    string `json:"name"`
	MediaID string `json:"media_id"`
}

var _ bodyer = reqKfAccountCreate{}

func (x reqKfAccountCreate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfAccountCreate 创建客服账号 响应
type respKfAccountCreate struct {
	respCommon

	OpenKfID string `json:"open_kfid"`
}

// reqKfAccountDelete 删除客服账号
type reqKfAccountDelete struct {
	OpenKfID string `json:"open_kfid"`
}

var _ bodyer = reqKfAccountDelete{}

func (x reqKfAccountDelete) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfAccountDelete 删除客服账号 响应
type respKfAccountDelete struct {
	respCommon
}

// reqKfAccountUpdate 修改客服账号
type reqKfAccountUpdate struct {
	OpenKfID string `json:"open_kfid"`
	Name     string `json:"name"`
	MediaID  string `json:"media_id"`
}

var _ bodyer = reqKfAccountUpdate{}

func (x reqKfAccountUpdate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfAccountUpdate 修改客服账号 响应
type respKfAccountUpdate struct {
	respCommon
}

// reqKfAccountList 获取客服账号列表
type reqKfAccountList struct {
	Offset int64 `json:"offset"`
	Limit  int64 `json:"limit"`
}

var _ urlValuer = reqKfAccountList{}

func (x reqKfAccountList) intoURLValues() url.Values {
	return url.Values{
		"offset": {strconv.FormatInt(x.Offset, 10)},
		"limit":  {strconv.FormatInt(x.Limit, 10)},
	}
}

// respKfAccountList 客服账号列表 响应
type respKfAccountList struct {
	respCommon

	AccountList []*KfAccount `json:"account_list"`
}

// reqAddKfContact 获取客服账号链接
type reqAddKfContact struct {
	OpenKfID string `json:"open_kfid"`
	Scene    string `json:"scene"`
}

var _ bodyer = reqAddKfContact{}

func (x reqAddKfContact) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respAddKfContact 获取客服账号链接 响应
type respAddKfContact struct {
	respCommon

	URL string `json:"url"`
}

// reqKfServicerCreate 添加接待人员
type reqKfServicerCreate struct {
	OpenKfID      string   `json:"open_kfid"`
	UserIDs       []string `json:"userid_list"`
	DepartmentIDs []int64  `json:"department_id_list"`
}

var _ bodyer = reqKfServicerCreate{}

func (x reqKfServicerCreate) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfServicerCreate 添加接待人员 响应
type respKfServicerCreate struct {
	respCommon

	ResultList []*KfServicerResult `json:"result_list"`
}

// reqKfServicerDelete 删除接待人员
type reqKfServicerDelete struct {
	OpenKfID      string   `json:"open_kfid"`
	UserIDs       []string `json:"userid_list"`
	DepartmentIDs []int64  `json:"department_id_list"`
}

var _ bodyer = reqKfServicerDelete{}

func (x reqKfServicerDelete) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfServicerDelete 删除接待人员 响应
type respKfServicerDelete struct {
	respCommon

	ResultList []*KfServicerResult `json:"result_list"`
}

// reqKfServicerList 获取接待人员列表
type reqKfServicerList struct {
	OpenKfID string `json:"open_kfid"`
}

var _ urlValuer = reqKfServicerList{}

func (x reqKfServicerList) intoURLValues() url.Values {
	return url.Values{
		"open_kfid": {x.OpenKfID},
	}
}

// respKfServicerList 接待人员列表 响应
type respKfServicerList struct {
	respCommon

	ServicerList []*KfServicer `json:"servicer_list"`
}

// reqKfServiceStateGet 获取会话状态
type reqKfServiceStateGet struct {
	OpenKfID       string `json:"open_kfid"`
	ExternalUserID string `json:"external_userid"`
}

var _ bodyer = reqKfServiceStateGet{}

func (x reqKfServiceStateGet) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfServiceStateGet 获取会话状态 响应
type respKfServiceStateGet struct {
	respCommon

	ServiceState   KfServiceState `json:"service_state"`
	ServicerUserID string         `json:"servicer_userid"`
}

// reqKfServiceStateTrans 变更会话状态
type reqKfServiceStateTrans struct {
	OpenKfID       string         `json:"open_kfid"`
	ExternalUserID string         `json:"external_userid"`
	ServiceState   KfServiceState `json:"service_state"`
	ServicerUserID string         `json:"servicer_userid"`
}

var _ bodyer = reqKfServiceStateTrans{}

func (x reqKfServiceStateTrans) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfServiceStateTrans 变更会话状态 响应
type respKfServiceStateTrans struct {
	respCommon

	MsgCode string `json:"msg_code"`
}

// reqKfSyncMsg 读取消息
type reqKfSyncMsg struct {
	OpenKfID    string `json:"open_kfid"`
	Cursor      string `json:"cursor"`
	Token       string `json:"token"`
	Limit       int64  `json:"limit"`
	VoiceFormat int    `json:"voice_format"`
}

var _ bodyer = reqKfSyncMsg{}

func (x reqKfSyncMsg) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respKfSyncMsg 读取消息 响应
type respKfSyncMsg struct {
	respCommon
	NextCursor string  `json:"next_cursor"`
	HasMore    int     `json:"has_more"`
	MsgList    []KfMsg `json:"msg_list"`
}

// reqOAGetCorpVacationConf 获取企业假期管理配置
type reqOAGetCorpVacationConf struct {
}

var _ urlValuer = reqOAGetCorpVacationConf{}

func (x reqOAGetCorpVacationConf) intoURLValues() url.Values {
	return url.Values{}
}

// respOAGetCorpVacationConf 获取企业假期管理配置 响应
type respOAGetCorpVacationConf struct {
	respCommon
	// Lists 假期列表
	Lists []CorpVacationConf `json:"lists"`
}

// reqOAGetUserVacationQuota 获取成员假期余额
type reqOAGetUserVacationQuota struct {
	UserID string `json:"userid"`
}

var _ bodyer = reqOAGetUserVacationQuota{}

func (x reqOAGetUserVacationQuota) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respOAGetUserVacationQuota 获取成员假期余额 响应
type respOAGetUserVacationQuota struct {
	respCommon
	// Lists 假期列表
	Lists []UserVacationQuota `json:"lists"`
}

// reqOASetOneUserVacationQuota 修改成员假期余额
type reqOASetOneUserVacationQuota struct {
	UserID       string `json:"userid"`
	VacationID   string `json:"vacation_id"`
	LeftDuration string `json:"leftduration"`
	TimeAttr     int64  `json:"time_attr"`
	Remarks      string `json:"remarks"`
}

var _ bodyer = reqOASetOneUserVacationQuota{}

func (x reqOASetOneUserVacationQuota) intoBody() ([]byte, error) {
	return marshalIntoJSONBody(x)
}

// respOASetOneUserVacationQuota 修改成员假期余额 响应
type respOASetOneUserVacationQuota struct {
	respCommon
}
