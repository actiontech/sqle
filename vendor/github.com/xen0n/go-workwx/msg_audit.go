package workwx

import (
	"time"
)

// MsgAuditAgreeStatus 会话中外部成员的同意状态
type MsgAuditAgreeStatus string

const (
	// MsgAuditAgreeStatusAgree 同意
	MsgAuditAgreeStatusAgree = "Agree"
	// MsgAuditAgreeStatusDisagree 不同意
	MsgAuditAgreeStatusDisagree = "Disagree"
	// MsgAuditAgreeStatusDefaultAgree 默认同意
	MsgAuditAgreeStatusDefaultAgree = "Default_Agree"
)

// CheckMsgAuditSingleAgreeUserInfo 获取会话同意情况（单聊）内外成员
type CheckMsgAuditSingleAgreeUserInfo struct {
	// UserID 内部成员的userid
	UserID string `json:"userid"`
	// ExternalOpenID 外部成员的externalopenid
	ExternalOpenID string `json:"exteranalopenid"`
}

// CheckMsgAuditSingleAgreeInfo 获取会话同意情况（单聊）同意信息
type CheckMsgAuditSingleAgreeInfo struct {
	CheckMsgAuditSingleAgreeUserInfo
	// AgreeStatus 同意:”Agree”，不同意:”Disagree”，默认同意:”Default_Agree”
	AgreeStatus MsgAuditAgreeStatus
	// StatusChangeTime 同意状态改变的具体时间
	StatusChangeTime time.Time
}

// CheckMsgAuditSingleAgree 获取会话同意情况（单聊）
func (c *WorkwxApp) CheckMsgAuditSingleAgree(infos []CheckMsgAuditSingleAgreeUserInfo) ([]CheckMsgAuditSingleAgreeInfo, error) {
	resp, err := c.execMsgAuditCheckSingleAgree(reqMsgAuditCheckSingleAgree{
		Infos: infos,
	})
	if err != nil {
		return nil, err
	}
	return resp.intoCheckSingleAgreeInfoList(), nil
}

// CheckMsgAuditRoomAgreeInfo 获取会话同意情况（群聊）同意信息
type CheckMsgAuditRoomAgreeInfo struct {
	// StatusChangeTime 同意状态改变的具体时间
	StatusChangeTime time.Time
	// AgreeStatus 同意:”Agree”，不同意:”Disagree”，默认同意:”Default_Agree”
	AgreeStatus MsgAuditAgreeStatus
	// ExternalOpenID 群内外部联系人的externalopenid
	ExternalOpenID string
}

// CheckMsgAuditRoomAgree 获取会话同意情况（群聊）
func (c *WorkwxApp) CheckMsgAuditRoomAgree(roomID string) ([]CheckMsgAuditRoomAgreeInfo, error) {
	resp, err := c.execMsgAuditCheckRoomAgree(reqMsgAuditCheckRoomAgree{
		RoomID: roomID,
	})
	if err != nil {
		return nil, err
	}
	return resp.intoCheckRoomAgreeInfoList(), nil
}

// MsgAuditEdition 会话内容存档版本
type MsgAuditEdition uint8

const (
	// MsgAuditEditionOffice 会话内容存档办公版
	MsgAuditEditionOffice MsgAuditEdition = 1
	// MsgAuditEditionService 会话内容存档服务版
	MsgAuditEditionService MsgAuditEdition = 2
	// MsgAuditEditionEnterprise 会话内容存档企业版
	MsgAuditEditionEnterprise MsgAuditEdition = 3
)

// ListMsgAuditPermitUser 获取会话内容存档开启成员列表
func (c *WorkwxApp) ListMsgAuditPermitUser(msgAuditEdition MsgAuditEdition) ([]string, error) {
	resp, err := c.execMsgAuditListPermitUser(reqMsgAuditListPermitUser{
		MsgAuditEdition: msgAuditEdition,
	})
	if err != nil {
		return nil, err
	}
	return resp.IDs, nil
}

// MsgAuditGroupChatMember 获取会话内容存档内部群成员
type MsgAuditGroupChatMember struct {
	// MemberID roomid群成员的id，userid
	MemberID string
	// JoinTime roomid群成员的入群时间
	JoinTime time.Time
}

// MsgAuditGroupChat 获取会话内容存档内部群信息
type MsgAuditGroupChat struct {
	// Members roomid对应的群成员列表
	Members []MsgAuditGroupChatMember
	// RoomName roomid对应的群名称
	RoomName string
	// Creator roomid对应的群创建者，userid
	Creator string
	// RoomCreateTime roomid对应的群创建时间
	RoomCreateTime time.Time
	// Notice roomid对应的群公告
	Notice string
}

// GetMsgAuditGroupChat 获取会话内容存档内部群信息
func (c *WorkwxApp) GetMsgAuditGroupChat(roomID string) (*MsgAuditGroupChat, error) {
	resp, err := c.execMsgAuditGetGroupChat(reqMsgAuditGetGroupChat{
		RoomID: roomID,
	})
	if err != nil {
		return nil, err
	}
	groupChat := resp.intoGroupChat()
	return &groupChat, nil
}
