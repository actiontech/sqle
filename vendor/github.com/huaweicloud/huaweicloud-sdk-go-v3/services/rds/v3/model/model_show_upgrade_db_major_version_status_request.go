package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowUpgradeDbMajorVersionStatusRequest Request Object
type ShowUpgradeDbMajorVersionStatusRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 要查询的状态 check：查询升级预检查的状态。 upgrade：查询大板本升级的状态。
	Action string `json:"action"`

	// 语言。默认en-us。
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowUpgradeDbMajorVersionStatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUpgradeDbMajorVersionStatusRequest struct{}"
	}

	return strings.Join([]string{"ShowUpgradeDbMajorVersionStatusRequest", string(data)}, " ")
}
