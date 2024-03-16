package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ErrorLogItem struct {

	// 日期时间UTC时间。
	Time *string `json:"time,omitempty"`

	// 日志级别。
	Level *string `json:"level,omitempty"`

	// 错误日志内容。
	Content *string `json:"content,omitempty"`

	// 日志单行序列号。
	LineNum *string `json:"line_num,omitempty"`
}

func (o ErrorLogItem) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrorLogItem struct{}"
	}

	return strings.Join([]string{"ErrorLogItem", string(data)}, " ")
}
