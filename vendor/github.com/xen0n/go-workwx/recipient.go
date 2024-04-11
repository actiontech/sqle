package workwx

// Recipient 消息收件人定义
type Recipient struct {
	// UserIDs 成员ID列表（消息接收者），最多支持1000个
	UserIDs []string
	// PartyIDs 部门ID列表，最多支持100个。
	PartyIDs []string
	// TagIDs 标签ID列表，最多支持100个
	TagIDs []string
	// ChatID 应用关联群聊ID，仅用于【发送消息到群聊会话】
	ChatID string
	// OpenKfID 应用关联客服ID，仅用于【客服发送消息】
	OpenKfID string
	// Code 仅用于【客服发送欢迎语等事件响应消息】
	Code string
}

// isIndividualTargetsEmpty 对非群发收件人字段而言，是否全为空
//
// 文档注释摘抄:
//
// > touser、toparty、totag不能同时为空，后面不再强调。
func (x *Recipient) isIndividualTargetsEmpty() bool {
	return len(x.UserIDs) == 0 && len(x.PartyIDs) == 0 && len(x.TagIDs) == 0
}

// isValidForMessageSend 本结构体是否对【发送应用消息】请求有效
func (x *Recipient) isValidForMessageSend() bool {
	if x.OpenKfID != "" {
		// 这时候你应该用 KfSend 接口
		return false
	}

	if x.Code != "" {
		// 这时候你应该用 KfOnEventSend 接口
		return false
	}

	if x.ChatID != "" {
		// 这时候你应该用 AppchatSend 接口
		return false
	}

	if x.isIndividualTargetsEmpty() {
		// 见这个方法的注释
		return false
	}

	if len(x.UserIDs) > 1000 || len(x.PartyIDs) > 100 || len(x.TagIDs) > 100 {
		// 见字段注释
		return false
	}

	return true
}

// isValidForAppchatSend 本结构体是否对【发送消息到群聊会话】请求有效
func (x *Recipient) isValidForAppchatSend() bool {
	if x.OpenKfID != "" {
		// 这时候你应该用 KfSend 接口
		return false
	}

	if x.Code != "" {
		// 这时候你应该用 KfOnEventSend 接口
		return false
	}

	if !x.isIndividualTargetsEmpty() {
		return false
	}

	return x.ChatID != ""
}

// isValidForKfSend 本结构体是否对【客服发送消息】请求有效
func (x *Recipient) isValidForKfSend() bool {
	if x.Code != "" {
		// 这时候你应该用 KfOnEventSend 接口
		return false
	}

	if !x.isIndividualTargetsEmpty() {
		return false
	}

	return x.OpenKfID != ""
}

// isValidForKfOnEventSend 本结构体是否对【客服发送欢迎语等事件响应消息】请求有效
func (x *Recipient) isValidForKfOnEventSend() bool {
	return x.Code != ""
}
