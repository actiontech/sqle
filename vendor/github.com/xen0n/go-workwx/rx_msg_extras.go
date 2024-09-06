package workwx

import (
	"encoding/xml"
	"fmt"
	"io"
)

// NOTE: 这顺便就构成了一个封闭的 enum
type messageKind interface {
	formatInto(io.Writer)
}

func extractMessageExtras(common rxMessageCommon, body []byte) (messageKind, error) {
	switch common.MsgType {
	case MessageTypeText:
		var x rxTextMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case MessageTypeImage:
		var x rxImageMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case MessageTypeVoice:
		var x rxVoiceMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case MessageTypeVideo:
		var x rxVideoMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case MessageTypeLocation:
		var x rxLocationMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case MessageTypeLink:
		var x rxLinkMessageSpecifics
		err := xml.Unmarshal(body, &x)
		if err != nil {
			return nil, err
		}
		return &x, nil

	case MessageTypeEvent:
		switch common.Event {
		case EventTypeSysApprovalChange:
			var x rxEventSysApprovalChange
			err := xml.Unmarshal(body, &x)
			if err != nil {
				return nil, err
			}
			return &x, nil

		case EventTypeChangeExternalContact:
			switch common.ChangeType {
			case ChangeTypeAddExternalContact:
				var x rxEventAddExternalContact
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil

			case ChangeTypeEditExternalContact:
				var x rxEventEditExternalContact
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil

			case ChangeTypeDelExternalContact:
				var x rxEventDelExternalContact
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil

			case ChangeTypeDelFollowUser:
				var x rxEventDelFollowUser
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil

			case ChangeTypeAddHalfExternalContact:
				var x rxEventAddHalfExternalContact
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil

			case ChangeTypeTransferFail:
				var x rxEventTransferFail
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil
			case ChangeTypeCreateUser:
				var x rxEventChangeTypeCreateUser
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil

			case ChangeTypeUpdateUser:
				var x rxEventChangeTypeUpdateUser
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil
			default:
				return nil, fmt.Errorf("unknown change type '%s'", common.ChangeType)
			}
		case EventTypeChangeExternalChat:
			var x rxEventChangeExternalChat
			err := xml.Unmarshal(body, &x)
			if err != nil {
				return nil, err
			}
			return &x, nil
		case EventTypeChangeContact:
			switch common.ChangeType {
			case ChangeTypeUpdateUser:
				var x rxEventChangeTypeUpdateUser
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil
			case ChangeTypeCreateUser:
				var x rxEventChangeTypeCreateUser
				err := xml.Unmarshal(body, &x)
				if err != nil {
					return nil, err
				}
				return &x, nil
			default:
				return nil, fmt.Errorf("unknown change type '%s'", common.ChangeType)
			}
		case EventTypeAppMenuClick:
			var x rxEventAppMenuClick
			err := xml.Unmarshal(body, &x)
			if err != nil {
				return nil, err
			}
			return &x, nil
		case EventTypeAppMenuView:
			var x rxEventAppMenuView
			err := xml.Unmarshal(body, &x)
			if err != nil {
				return nil, err
			}
			return &x, nil
		case EventTypeKfMsgOrEvent:
			var x rxEventKfMsgOrEvent
			err := xml.Unmarshal(body, &x)
			if err != nil {
				return nil, err
			}
			return &x, nil

		default:
			// 返回一个未定义的事件类型
			return &rxEventUnknown{EventType: string(common.Event), Raw: string(body)}, nil
		}

	default:
		return nil, fmt.Errorf("unknown message type '%s'", common.MsgType)
	}
}

// TextMessageExtras 文本消息的参数。
type TextMessageExtras interface {
	messageKind

	// GetContent 返回文本消息的内容。
	GetContent() string
}

var _ TextMessageExtras = (*rxTextMessageSpecifics)(nil)

func (r *rxTextMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Content: %#v", r.Content)
}

func (r *rxTextMessageSpecifics) GetContent() string {
	return r.Content
}

// ImageMessageExtras 图片消息的参数。
type ImageMessageExtras interface {
	messageKind

	// GetPicURL 返回图片消息的图片链接 URL。
	GetPicURL() string

	// GetMediaID 返回图片消息的图片媒体文件 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetMediaID() string
}

var _ ImageMessageExtras = (*rxImageMessageSpecifics)(nil)

func (r *rxImageMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "PicURL: %#v, MediaID: %#v", r.PicURL, r.MediaID)
}

func (r *rxImageMessageSpecifics) GetPicURL() string {
	return r.PicURL
}

func (r *rxImageMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

// VoiceMessageExtras 语音消息的参数。
type VoiceMessageExtras interface {
	messageKind

	// GetMediaID 返回语音消息的语音媒体文件 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetMediaID() string

	// GetFormat 返回语音消息的语音格式，如 "amr"、"speex" 等。
	GetFormat() string
}

var _ VoiceMessageExtras = (*rxVoiceMessageSpecifics)(nil)

func (r *rxVoiceMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v, Format: %#v", r.MediaID, r.Format)
}

func (r *rxVoiceMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

func (r *rxVoiceMessageSpecifics) GetFormat() string {
	return r.Format
}

// VideoMessageExtras 视频消息的参数。
type VideoMessageExtras interface {
	messageKind

	// GetMediaID 返回视频消息的视频媒体文件 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetMediaID() string

	// GetThumbMediaID 返回视频消息缩略图的媒体 ID。
	//
	// 可以调用【获取媒体文件】接口拉取，仅三天内有效。
	GetThumbMediaID() string
}

var _ VideoMessageExtras = (*rxVideoMessageSpecifics)(nil)

func (r *rxVideoMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "MediaID: %#v, ThumbMediaID: %#v", r.MediaID, r.ThumbMediaID)
}

func (r *rxVideoMessageSpecifics) GetMediaID() string {
	return r.MediaID
}

func (r *rxVideoMessageSpecifics) GetThumbMediaID() string {
	return r.ThumbMediaID
}

// LocationMessageExtras 位置消息的参数。
type LocationMessageExtras interface {
	messageKind

	// GetLatitude 返回位置消息的纬度（角度值；北纬为正）。
	GetLatitude() float64

	// GetLongitude 返回位置消息的经度（角度值；东经为正）。
	GetLongitude() float64

	// GetScale 返回位置消息的地图缩放大小。
	GetScale() int

	// GetLabel 返回位置消息的地理位置信息。
	GetLabel() string

	// 不知道这个有啥用，先不暴露
	// GetAppType() string
}

var _ LocationMessageExtras = (*rxLocationMessageSpecifics)(nil)

func (r *rxLocationMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"Latitude: %#v, Longitude: %#v, Scale: %d, Label: %#v",
		r.Lat,
		r.Lon,
		r.Scale,
		r.Label,
	)
}

func (r *rxLocationMessageSpecifics) GetLatitude() float64 {
	return r.Lat
}

func (r *rxLocationMessageSpecifics) GetLongitude() float64 {
	return r.Lon
}

func (r *rxLocationMessageSpecifics) GetScale() int {
	return r.Scale
}

func (r *rxLocationMessageSpecifics) GetLabel() string {
	return r.Label
}

// LinkMessageExtras 链接消息的参数。
type LinkMessageExtras interface {
	messageKind

	// GetTitle 返回链接消息的标题。
	GetTitle() string

	// GetDescription 返回链接消息的描述。
	GetDescription() string

	// GetURL 返回链接消息的跳转 URL。
	GetURL() string

	// GetPicURL 返回链接消息的封面缩略图 URL。
	GetPicURL() string
}

var _ LinkMessageExtras = (*rxLinkMessageSpecifics)(nil)

func (r *rxLinkMessageSpecifics) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"Title: %#v, Description: %#v, URL: %#v, PicURL: %#v",
		r.Title,
		r.Description,
		r.URL,
		r.PicURL,
	)
}

func (r *rxLinkMessageSpecifics) GetTitle() string {
	return r.Title
}

func (r *rxLinkMessageSpecifics) GetDescription() string {
	return r.Description
}

func (r *rxLinkMessageSpecifics) GetURL() string {
	return r.URL
}

func (r *rxLinkMessageSpecifics) GetPicURL() string {
	return r.PicURL
}

// EventAddExternalContact 添加企业客户事件的参数。
type EventAddExternalContact interface {
	messageKind

	// GetUserID 企业服务人员的UserID
	GetUserID() string

	// GetExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	GetExternalUserID() string

	// GetState 添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
	GetState() string

	// GetWelcomeCode 欢迎语code，可用于发送欢迎语
	GetWelcomeCode() string
}

var _ EventAddExternalContact = (*rxEventAddExternalContact)(nil)

func (r *rxEventAddExternalContact) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UserID: %#v, ExternalUserID: %#v, State: %#v, WelcomeCode: %#v",
		r.UserID,
		r.ExternalUserID,
		r.State,
		r.WelcomeCode,
	)
}

func (r *rxEventAddExternalContact) GetUserID() string {
	return r.UserID
}

func (r *rxEventAddExternalContact) GetExternalUserID() string {
	return r.ExternalUserID
}

func (r *rxEventAddExternalContact) GetState() string {
	return r.State
}

func (r *rxEventAddExternalContact) GetWelcomeCode() string {
	return r.WelcomeCode
}

// EventEditExternalContact 编辑企业客户事件的参数。
type EventEditExternalContact interface {
	messageKind

	// GetUserID 企业服务人员的UserID
	GetUserID() string

	// GetExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	GetExternalUserID() string

	// GetState 添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
	GetState() string
}

var _ EventEditExternalContact = (*rxEventEditExternalContact)(nil)

func (r *rxEventEditExternalContact) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UserID: %#v, ExternalUserID: %#v, State: %#v",
		r.UserID,
		r.ExternalUserID,
		r.State,
	)
}

func (r *rxEventEditExternalContact) GetUserID() string {
	return r.UserID
}

func (r *rxEventEditExternalContact) GetExternalUserID() string {
	return r.ExternalUserID
}

func (r *rxEventEditExternalContact) GetState() string {
	return r.State
}

// EventAddHalfExternalContact 外部联系人免验证添加成员事件。
type EventAddHalfExternalContact interface {
	messageKind

	// GetUserID 企业服务人员的UserID
	GetUserID() string

	// GetExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	GetExternalUserID() string

	// GetState 添加此用户的「联系我」方式配置的state参数，可用于识别添加此用户的渠道
	GetState() string
}

var _ EventAddHalfExternalContact = (*rxEventAddHalfExternalContact)(nil)

func (r *rxEventAddHalfExternalContact) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UserID: %#v, ExternalUserID: %#v, State: %#v, WelcomeCode: %#v",
		r.UserID,
		r.ExternalUserID,
		r.State,
		r.WelcomeCode,
	)
}

func (r *rxEventAddHalfExternalContact) GetUserID() string {
	return r.UserID
}

func (r *rxEventAddHalfExternalContact) GetExternalUserID() string {
	return r.ExternalUserID
}

func (r *rxEventAddHalfExternalContact) GetState() string {
	return r.State
}

func (r *rxEventAddHalfExternalContact) GetWelcomeCode() string {
	return r.WelcomeCode
}

// EventDelExternalContact 删除企业客户事件
type EventDelExternalContact interface {
	messageKind

	// GetUserID 企业服务人员的UserID
	GetUserID() string

	// GetExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	GetExternalUserID() string
}

var _ EventDelExternalContact = (*rxEventDelExternalContact)(nil)

func (r *rxEventDelExternalContact) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UserID: %#v, ExternalUserID: %#v",
		r.UserID,
		r.ExternalUserID,
	)
}

func (r *rxEventDelExternalContact) GetUserID() string {
	return r.UserID
}

func (r *rxEventDelExternalContact) GetExternalUserID() string {
	return r.ExternalUserID
}

// EventDelFollowUser 删除跟进成员事件
type EventDelFollowUser interface {
	messageKind

	// GetUserID 企业服务人员的UserID
	GetUserID() string

	// GetExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	GetExternalUserID() string
}

var _ EventDelFollowUser = (*rxEventDelFollowUser)(nil)

func (r *rxEventDelFollowUser) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UserID: %#v, ExternalUserID: %#v",
		r.UserID,
		r.ExternalUserID,
	)
}

func (r *rxEventDelFollowUser) GetUserID() string {
	return r.UserID
}

func (r *rxEventDelFollowUser) GetExternalUserID() string {
	return r.ExternalUserID
}

// EventTransferFail 客户接替失败事件
type EventTransferFail interface {
	messageKind

	// GetUserID 企业服务人员的UserID
	GetUserID() string

	// GetExternalUserID 外部联系人的userid，注意不是企业成员的帐号
	GetExternalUserID() string

	// GetFailReason 接替失败的原因, customer_refused-客户拒绝， customer_limit_exceed-接替成员的客户数达到上限
	GetFailReason() string
}

var _ EventTransferFail = (*rxEventTransferFail)(nil)

func (r *rxEventTransferFail) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UserID: %#v, ExternalUserID: %#v, FailReason: %#v",
		r.UserID,
		r.ExternalUserID,
		r.FailReason,
	)
}

func (r *rxEventTransferFail) GetUserID() string {
	return r.UserID
}

func (r *rxEventTransferFail) GetExternalUserID() string {
	return r.ExternalUserID
}

func (r *rxEventTransferFail) GetFailReason() string {
	return r.FailReason
}

// EventChangeExternalChat 客户群变更事件
type EventChangeExternalChat interface {
	messageKind

	// GetChatID 群ID
	GetChatID() string

	// GetToUserName 企业微信CorpID
	GetToUserName() string

	// GetFromUserName 此事件该值固定为sys，表示该消息由系统生成
	GetFromUserName() string

	// GetFailReason 接替失败的原因, customer_refused-客户拒绝， customer_limit_exceed-接替成员的客户数达到上限
	GetFailReason() string
}

var _ EventChangeExternalChat = (*rxEventChangeExternalChat)(nil)

func (r *rxEventChangeExternalChat) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"ChatID: %#v, ToUserName: %#v, FromUserName: %#v, FailReason: %#v",
		r.ChatID,
		r.ToUserName,
		r.FromUserName,
		r.FailReason,
	)
}

func (r *rxEventChangeExternalChat) GetChatID() string {
	return r.ChatID
}

func (r *rxEventChangeExternalChat) GetToUserName() string {
	return r.ToUserName
}

func (r *rxEventChangeExternalChat) GetFromUserName() string {
	return r.FromUserName
}

func (r *rxEventChangeExternalChat) GetFailReason() string {
	return r.FailReason
}

// EventSysApprovalChange 审批申请状态变化回调通知
type EventSysApprovalChange interface {
	messageKind

	// GetApprovalInfo 获取审批模板详情
	GetApprovalInfo() OAApprovalInfo
}

func (r rxEventSysApprovalChange) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"ApprovalInfo: %#v",
		r.ApprovalInfo,
	)
}

func (r rxEventSysApprovalChange) GetApprovalInfo() OAApprovalInfo {
	return r.ApprovalInfo
}

func (r rxEventChangeTypeUpdateUser) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"UpdateUser: %#v",
		r,
	)
}

func (r rxEventChangeTypeCreateUser) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"CreateUser: %#v",
		r,
	)
}

func (r rxEventAppMenuClick) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "EventKey: %#v", r.EventKey)
}

func (r rxEventAppMenuView) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "EventKey: %#v", r.EventKey)
}

func (r rxEventAppSubscribe) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "EventKey: %#v", r.EventKey)
}

func (r rxEventAppUnsubscribe) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "EventKey: %#v", r.EventKey)
}

// EventKfMsgOrEvent 客服接收消息和事件
type EventKfMsgOrEvent interface {
	messageKind

	// GetOpenKfID 客服账号ID
	GetOpenKfID() string

	// GetToken 调用拉取消息接口时，需要传此token，用于校验请求的合法性
	GetToken() string
}

var _ EventKfMsgOrEvent = (*rxEventKfMsgOrEvent)(nil)

func (r *rxEventKfMsgOrEvent) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(
		w,
		"OpenKfID: %#v, Token: %#v",
		r.OpenKfID,
		r.Token,
	)
}

func (r *rxEventKfMsgOrEvent) GetOpenKfID() string {
	return r.OpenKfID
}

func (r *rxEventKfMsgOrEvent) GetToken() string {
	return r.Token
}

func (r rxEventUnknown) formatInto(w io.Writer) {
	_, _ = fmt.Fprintf(w, "Raw: %#v", r.Raw)
}
