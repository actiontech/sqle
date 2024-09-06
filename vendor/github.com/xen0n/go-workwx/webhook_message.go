package workwx

// MentionAll 表示提醒所有人（“@所有人”）的特殊标记
const MentionAll = "@all"

// Mentions 群机器人消息的提醒设置
type Mentions struct {
	// UserIDs userid 的列表，提醒群中的指定成员（@某个成员），MentionAll 表示提醒所有人
	//
	// 如果开发者获取不到 userid，可使用 Mobiles
	UserIDs []string
	// Mobiles 手机号列表，提醒手机号对应的群成员（@某个成员），MentionAll 表示提醒所有人
	//
	// 如果开发者获取不到 userid，可使用该列表，否则可使用 UserIDs
	Mobiles []string
}

// SendTextMessage 发送文本消息
func (c *WebhookClient) SendTextMessage(
	content string,
	mentions *Mentions,
) error {
	params := map[string]interface{}{
		"content": content,
	}

	if mentions != nil {
		if len(mentions.UserIDs) > 0 {
			params["mentioned_list"] = mentions.UserIDs
		}

		if len(mentions.Mobiles) > 0 {
			params["mentioned_mobile_list"] = mentions.Mobiles
		}
	}

	return c.sendMessage("text", params)
}

// SendMarkdownMessage 发送 Markdown 消息
//
// NOTE: 使用群机器人接口发送 Markdown 消息时，不能传递 Mentions 结构体，而需要使用
// `<@userid>` 的特殊扩展语法来表示 at 给定的 userid。
func (c *WebhookClient) SendMarkdownMessage(
	content string,
) error {
	params := map[string]interface{}{
		"content": content,
	}

	return c.sendMessage("markdown", params)
}

// sendMessage 发送消息底层接口
func (c *WebhookClient) sendMessage(
	msgtype string,
	content map[string]interface{},
) error {
	req := map[string]interface{}{
		"msgtype": msgtype,
		msgtype:   content,
	}

	err := c.executeQyapiJSONPost("/cgi-bin/webhook/send", req, nil)
	if err != nil {
		return err
	}

	return nil
}
