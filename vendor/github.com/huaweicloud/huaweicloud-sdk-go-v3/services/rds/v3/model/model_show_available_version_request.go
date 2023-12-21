package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAvailableVersionRequest Request Object
type ShowAvailableVersionRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 语言。默认en-us。
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowAvailableVersionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAvailableVersionRequest struct{}"
	}

	return strings.Join([]string{"ShowAvailableVersionRequest", string(data)}, " ")
}
