package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetSecondLevelMonitorRequest Request Object
type SetSecondLevelMonitorRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *SecondMonitor `json:"body,omitempty"`
}

func (o SetSecondLevelMonitorRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetSecondLevelMonitorRequest struct{}"
	}

	return strings.Join([]string{"SetSecondLevelMonitorRequest", string(data)}, " ")
}
