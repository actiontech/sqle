package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CopyConfigurationRequest Request Object
type CopyConfigurationRequest struct {

	// 参数模板ID
	ConfigId string `json:"config_id"`

	Body *ConfigurationCopyRequestBody `json:"body,omitempty"`
}

func (o CopyConfigurationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CopyConfigurationRequest struct{}"
	}

	return strings.Join([]string{"CopyConfigurationRequest", string(data)}, " ")
}
