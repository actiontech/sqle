package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CustomerUpgradeDatabaseVersionReqNew struct {

	// 是否延迟至可维护时间段内升级。 取值范围： - true：延迟升级。表示实例将在设置的可维护时间段内升级。 - false：立即升级，默认该方式。
	IsDelayed *bool `json:"is_delayed,omitempty"`
}

func (o CustomerUpgradeDatabaseVersionReqNew) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CustomerUpgradeDatabaseVersionReqNew struct{}"
	}

	return strings.Join([]string{"CustomerUpgradeDatabaseVersionReqNew", string(data)}, " ")
}
