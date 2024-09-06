package workwx

// CreateAppchat 创建群聊会话
func (c *WorkwxApp) CreateAppchat(chatInfo *ChatInfo) (chatID string, err error) {
	resp, err := c.execAppchatCreate(reqAppchatCreate{
		ChatInfo: chatInfo,
	})
	if err != nil {
		return "", err
	}
	return resp.ChatID, nil
}

// UpdateAppchat 修改群聊会话
func (c *WorkwxApp) UpdateAppchat(chatInfo ChatInfo, addMemberUserIDs, delMemberUserIDs []string) (err error) {
	_, err = c.execAppchatUpdate(reqAppchatUpdate{
		ChatInfo:         chatInfo,
		AddMemberUserIDs: addMemberUserIDs,
		DelMemberUserIDs: delMemberUserIDs,
	})
	if err != nil {
		return err
	}
	return nil
}

// GetAppchat 获取群聊会话
func (c *WorkwxApp) GetAppchat(chatID string) (*ChatInfo, error) {
	resp, err := c.execAppchatGet(reqAppchatGet{
		ChatID: chatID,
	})
	if err != nil {
		return nil, err
	}

	// TODO: return bare T instead of &T?
	obj := resp.ChatInfo
	return obj, nil
}

// GetAppChatList 获取客户群列表 企业微信接口调整 此API同GetGroupChatList 兼容处理
func (c *WorkwxApp) GetAppChatList(req ReqChatList) (*RespAppchatList, error) {
	resp, err := c.execGroupChatListGet(reqGroupChatList{
		ReqChatList: req,
	})
	if err != nil {
		return nil, err
	}
	return resp.RespGroupChatList, nil
}

// GetAppChatInfo 获取客户群详细信息 企业微信接口调整 此API同GetGroupChatInfo 兼容处理
func (c *WorkwxApp) GetAppChatInfo(chatID string) (*RespAppChatInfo, error) {
	resp, err := c.execGroupChatInfoGet(reqGroupChatInfo{
		ChatID:   chatID,
		NeedName: ChatNeedName,
	})
	if err != nil {
		return nil, err
	}
	return resp.GroupChat, nil
}
