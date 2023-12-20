package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowUpgradeDbMajorVersionStatusResponse Response Object
type ShowUpgradeDbMajorVersionStatusResponse struct {

	// 实例大版本升级状态 \" running\"：预检查或大版本升级进行中 \" success\"：预检查或大版本升级成功 \" failed\"：预检查或大版本升级失败
	Status *string `json:"status,omitempty"`

	// 目标版本。
	TargetVersion *string `json:"target_version,omitempty"`

	// 开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	StartTime *string `json:"start_time,omitempty"`

	// 检查成功时，检查报告到期时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。 该字段仅在action为check时返回。
	ReportExpirationTime *string `json:"report_expiration_time,omitempty"`

	// 预检查或升级报告信息。
	Detail         *string `json:"detail,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowUpgradeDbMajorVersionStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowUpgradeDbMajorVersionStatusResponse struct{}"
	}

	return strings.Join([]string{"ShowUpgradeDbMajorVersionStatusResponse", string(data)}, " ")
}
