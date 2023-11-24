package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDbMajorVersionPreCheckRequest Request Object
type UpgradeDbMajorVersionPreCheckRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 语言。默认en-us。
	XLanguage *string `json:"X-Language,omitempty"`

	Body *PostgresqlPreCheckUpgradeMajorVersionReq `json:"body,omitempty"`
}

func (o UpgradeDbMajorVersionPreCheckRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbMajorVersionPreCheckRequest struct{}"
	}

	return strings.Join([]string{"UpgradeDbMajorVersionPreCheckRequest", string(data)}, " ")
}
