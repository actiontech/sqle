package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDbMajorVersionRequest Request Object
type UpgradeDbMajorVersionRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	Body *UpgradePgMajorVersion `json:"body,omitempty"`
}

func (o UpgradeDbMajorVersionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbMajorVersionRequest struct{}"
	}

	return strings.Join([]string{"UpgradeDbMajorVersionRequest", string(data)}, " ")
}
