package workwx

import (
	"errors"
)

// SendTextMessage 发送文本消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendTextMessage(
	recipient *Recipient,
	content string,
	isSafe bool,
) error {
	return c.sendMessage(recipient, "text", map[string]interface{}{"content": content}, isSafe)
}

// SendImageMessage 发送图片消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendImageMessage(
	recipient *Recipient,
	mediaID string,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"image",
		map[string]interface{}{
			"media_id": mediaID,
		}, isSafe,
	)
}

// SendVoiceMessage 发送语音消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendVoiceMessage(
	recipient *Recipient,
	mediaID string,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"voice",
		map[string]interface{}{
			"media_id": mediaID,
		}, isSafe,
	)
}

// SendVideoMessage 发送视频消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendVideoMessage(
	recipient *Recipient,
	mediaID string,
	description string,
	title string,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"video",
		map[string]interface{}{
			"media_id":    mediaID,
			"description": description, // TODO: 零值
			"title":       title,       // TODO: 零值
		}, isSafe,
	)
}

// SendFileMessage 发送文件消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendFileMessage(
	recipient *Recipient,
	mediaID string,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"file",
		map[string]interface{}{
			"media_id": mediaID,
		}, isSafe,
	)
}

// SendTextCardMessage 发送文本卡片消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendTextCardMessage(
	recipient *Recipient,
	title string,
	description string,
	url string,
	buttonText string,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"textcard",
		map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
			"btntxt":      buttonText, // TODO: 零值
		}, isSafe,
	)
}

// SendNewsMessage 发送图文消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendNewsMessage(
	recipient *Recipient,
	articles []Article,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"news",
		map[string]interface{}{
			"articles": articles,
		}, isSafe,
	)
}

// SendMPNewsMessage 发送 mpnews 类型的图文消息
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendMPNewsMessage(
	recipient *Recipient,
	mparticles []MPArticle,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"mpnews",
		map[string]interface{}{
			"articles": mparticles,
		}, isSafe,
	)
}

// SendMarkdownMessage 发送 Markdown 消息
//
// 仅支持 Markdown 的子集，详见[官方文档](https://work.weixin.qq.com/api/doc#90002/90151/90854/%E6%94%AF%E6%8C%81%E7%9A%84markdown%E8%AF%AD%E6%B3%95)。
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) SendMarkdownMessage(
	recipient *Recipient,
	content string,
	isSafe bool,
) error {
	return c.sendMessage(recipient, "markdown", map[string]interface{}{"content": content}, isSafe)
}

// SendTaskCardMessage 发送 任务卡片 消息
func (c *WorkwxApp) SendTaskCardMessage(
	recipient *Recipient,
	title string,
	description string,
	url string,
	taskid string,
	btn []TaskCardBtn,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"taskcard",
		map[string]interface{}{
			"title":       title,
			"description": description,
			"url":         url,
			"task_id":     taskid,
			"btn":         btn,
		}, isSafe,
	)
}

// SendTemplateCardMessage 发送卡片模板消息
func (c *WorkwxApp) SendTemplateCardMessage(
	recipient *Recipient,
	templateCard TemplateCard,
	isSafe bool,
) error {
	return c.sendMessage(
		recipient,
		"template_card",
		map[string]interface{}{
			"template_card": templateCard,
		}, isSafe,
	)
}

// sendMessage 发送消息底层接口
//
// 收件人参数如果仅设置了 `ChatID` 字段，则为【发送消息到群聊会话】接口调用；
// 收件人参数如果仅设置了 `OpenKfID` 字段，则为【客服发送消息】接口调用；
// 收件人参数如果仅设置了 `Code` 字段，则为【发送欢迎语等事件响应消息】接口调用；
// 否则为单纯的【发送应用消息】接口调用。
func (c *WorkwxApp) sendMessage(
	recipient *Recipient,
	msgtype string,
	content map[string]interface{},
	isSafe bool,
) error {
	sendRequestFunc := c.execMessageSend
	if !recipient.isValidForMessageSend() {
		if recipient.isValidForAppchatSend() {
			sendRequestFunc = c.execAppchatSend
		} else if recipient.isValidForKfSend() {
			sendRequestFunc = c.execKfSend
		} else if recipient.isValidForKfOnEventSend() {
			sendRequestFunc = c.execKfOnEventSend
		} else {
			// TODO: better error
			return errors.New("recipient invalid for message sending")
		}
	}

	req := reqMessage{
		ToUser:  recipient.UserIDs,
		ToParty: recipient.PartyIDs,
		ToTag:   recipient.TagIDs,
		ChatID:  recipient.ChatID,
		AgentID: c.AgentID,
		MsgType: msgtype,
		Content: content,
		IsSafe:  isSafe,
	}

	resp, err := sendRequestFunc(req)

	if err != nil {
		return err
	}

	// TODO: what to do with resp?
	_ = resp
	return nil
}
