package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowSecondLevelMonitoringRequest Request Object
type ShowSecondLevelMonitoringRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowSecondLevelMonitoringRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSecondLevelMonitoringRequest struct{}"
	}

	return strings.Join([]string{"ShowSecondLevelMonitoringRequest", string(data)}, " ")
}
