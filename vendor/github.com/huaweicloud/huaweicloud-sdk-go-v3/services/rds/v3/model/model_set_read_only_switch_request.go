package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetReadOnlySwitchRequest Request Object
type SetReadOnlySwitchRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *MysqlReadOnlySwitch `json:"body,omitempty"`
}

func (o SetReadOnlySwitchRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetReadOnlySwitchRequest struct{}"
	}

	return strings.Join([]string{"SetReadOnlySwitchRequest", string(data)}, " ")
}
