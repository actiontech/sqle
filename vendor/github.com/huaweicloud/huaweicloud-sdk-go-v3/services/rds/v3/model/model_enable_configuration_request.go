package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EnableConfigurationRequest Request Object
type EnableConfigurationRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 参数模板ID。
	ConfigId string `json:"config_id"`

	Body *ApplyConfigurationRequest `json:"body,omitempty"`
}

func (o EnableConfigurationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnableConfigurationRequest struct{}"
	}

	return strings.Join([]string{"EnableConfigurationRequest", string(data)}, " ")
}
