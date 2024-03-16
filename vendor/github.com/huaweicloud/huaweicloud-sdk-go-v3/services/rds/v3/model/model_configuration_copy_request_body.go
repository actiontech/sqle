package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ConfigurationCopyRequestBody struct {

	// 只支持a-zA-Z0-9._- 以上字符，长度限制1-64字符
	Name string `json:"name"`

	// 不支持 !<>=&\" ' 字符，长度限制0-256字符
	Description *string `json:"description,omitempty"`
}

func (o ConfigurationCopyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationCopyRequestBody struct{}"
	}

	return strings.Join([]string{"ConfigurationCopyRequestBody", string(data)}, " ")
}
