package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InspectionReports struct {

	// 检查报告ID。
	Id string `json:"id"`

	// 检查时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	CheckTime string `json:"check_time"`

	// 到期时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如偏移1个小时显示为+0100。
	ExpirationTime string `json:"expiration_time"`

	// 目标版本。
	TargetVersion string `json:"target_version"`

	// 检查结果。 success，表示成功。 failed，表示失败。 running， 表示检查中。
	Result string `json:"result"`

	// 检查报告详情。
	Detail string `json:"detail"`
}

func (o InspectionReports) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InspectionReports struct{}"
	}

	return strings.Join([]string{"InspectionReports", string(data)}, " ")
}
