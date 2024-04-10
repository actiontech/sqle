package workwx

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// RxMessage 一条接收到的消息
type RxMessage struct {
	// FromUserID 发送者的 UserID
	FromUserID string
	// SendTime 消息发送时间
	SendTime time.Time
	// MsgType 消息类型
	MsgType MessageType
	// MsgID 消息 ID
	MsgID int64
	// AgentID 企业应用 ID，可在应用的设置页面查看
	AgentID int64
	// Event 事件类型 MsgType为event存在
	Event EventType
	// ChangeType 变更类型 Event为change_external_contact存在
	ChangeType ChangeType

	extras messageKind
}

func fromEnvelope(body []byte) (*RxMessage, error) {
	// extract common part
	var common rxMessageCommon
	err := xml.Unmarshal(body, &common)
	if err != nil {
		return nil, err
	}

	// deal with polymorphic message types
	extras, err := extractMessageExtras(common, body)
	if err != nil {
		return nil, err
	}

	// assemble message object
	var obj RxMessage
	{
		// let's force people to think about timezones okay?
		// -- let's not
		sendTime := time.Unix(common.CreateTime, 0) // in time.Local

		obj = RxMessage{
			FromUserID: common.FromUserName,
			SendTime:   sendTime,
			MsgType:    common.MsgType,
			MsgID:      common.MsgID,
			AgentID:    common.AgentID,
			Event:      common.Event,
			ChangeType: common.ChangeType,

			extras: extras,
		}
	}

	return &obj, nil
}

func (m *RxMessage) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(
		&sb,
		"RxMessage { FromUserID: %#v, SendTime: %d, MsgType: %#v, MsgID: %d, AgentID: %d, Event: %#v, ChangeType: %#v, ",
		m.FromUserID,
		m.SendTime.UnixNano(),
		m.MsgType,
		m.MsgID,
		m.AgentID,
		m.Event,
		m.ChangeType,
	)

	m.extras.formatInto(&sb)

	sb.WriteString(" }")

	return sb.String()
}

// Text 如果消息为文本类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Text() (TextMessageExtras, bool) {
	y, ok := m.extras.(TextMessageExtras)
	return y, ok
}

// Image 如果消息为图片类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Image() (ImageMessageExtras, bool) {
	y, ok := m.extras.(ImageMessageExtras)
	return y, ok
}

// Voice 如果消息为语音类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Voice() (VoiceMessageExtras, bool) {
	y, ok := m.extras.(VoiceMessageExtras)
	return y, ok
}

// Video 如果消息为视频类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Video() (VideoMessageExtras, bool) {
	y, ok := m.extras.(VideoMessageExtras)
	return y, ok
}

// Location 如果消息为位置类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Location() (LocationMessageExtras, bool) {
	y, ok := m.extras.(LocationMessageExtras)
	return y, ok
}

// Link 如果消息为链接类型，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) Link() (LinkMessageExtras, bool) {
	y, ok := m.extras.(LinkMessageExtras)
	return y, ok
}

// EventAddExternalContact 如果消息为添加企业客户事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventAddExternalContact() (EventAddExternalContact, bool) {
	y, ok := m.extras.(EventAddExternalContact)
	return y, ok
}

// EventEditExternalContact 如果消息为编辑企业客户事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventEditExternalContact() (EventEditExternalContact, bool) {
	y, ok := m.extras.(EventEditExternalContact)
	return y, ok
}

// EventDelExternalContact 如果消息为删除企业客户事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventDelExternalContact() (EventDelExternalContact, bool) {
	y, ok := m.extras.(EventDelExternalContact)
	return y, ok
}

// EventDelFollowUser 如果消息为删除跟进成员事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventDelFollowUser() (EventDelFollowUser, bool) {
	y, ok := m.extras.(EventDelFollowUser)
	return y, ok
}

// EventAddHalfExternalContact 如果消息为外部联系人免验证添加成员事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventAddHalfExternalContact() (EventAddHalfExternalContact, bool) {
	y, ok := m.extras.(EventAddHalfExternalContact)
	return y, ok
}

// EventTransferFail 如果消息为客户接替失败事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventTransferFail() (EventTransferFail, bool) {
	y, ok := m.extras.(EventTransferFail)
	return y, ok
}

// EventChangeExternalChat 如果消息为客户群变更事件，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventChangeExternalChat() (EventChangeExternalChat, bool) {
	y, ok := m.extras.(EventChangeExternalChat)
	return y, ok
}

// EventSysApprovalChange 如果消息为审批申请状态变化回调通知，则拿出相应的消息参数，否则返回 nil, false
func (m *RxMessage) EventSysApprovalChange() (EventSysApprovalChange, bool) {
	y, ok := m.extras.(EventSysApprovalChange)
	return y, ok
}

// EventChangeTypeUpdateUser 如果消息为更新成员事件通知，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventChangeTypeUpdateUser() (*rxEventChangeTypeUpdateUser, bool) {
	y, ok := m.extras.(*rxEventChangeTypeUpdateUser)
	return y, ok
}

// EventChangeTypeCreateUser  如果消息为创建成员事件通知，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventChangeTypeCreateUser() (*rxEventChangeTypeCreateUser, bool) {
	y, ok := m.extras.(*rxEventChangeTypeCreateUser)
	return y, ok
}

// EventAppMenuClick  如果消息为应用菜单点击事件通知，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventAppMenuClick() (*rxEventAppMenuClick, bool) {
	y, ok := m.extras.(*rxEventAppMenuClick)
	return y, ok
}

// EventAppMenuView  如果消息为应用菜单链接点击事件通知，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventAppMenuView() (*rxEventAppMenuView, bool) {
	y, ok := m.extras.(*rxEventAppMenuView)
	return y, ok
}

// EventAppSubscribe  如果消息为应用订阅事件通知，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventAppSubscribe() (*rxEventAppSubscribe, bool) {
	y, ok := m.extras.(*rxEventAppSubscribe)
	return y, ok
}

// EventAppUnsubscribe  如果消息为应用订阅取消事件通知，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventAppUnsubscribe() (*rxEventAppUnsubscribe, bool) {
	y, ok := m.extras.(*rxEventAppUnsubscribe)
	return y, ok
}

// EventKfMsgOrEvent  如果消息为客服接收消息和事件，则拿出相应消息参数，否则返回 nil, false
func (m *RxMessage) EventKfMsgOrEvent() (*rxEventKfMsgOrEvent, bool) {
	y, ok := m.extras.(*rxEventKfMsgOrEvent)
	return y, ok
}

// EventUnknown  未定义的event类型
func (m *RxMessage) EventUnknown() (*rxEventUnknown, bool) {
	y, ok := m.extras.(*rxEventUnknown)
	return y, ok
}
