package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RecyclePolicy 回收策略。
type RecyclePolicy struct {

	// 保留时间，设置已删除实例保留天数，支持整数，可设置范围为1~7天。  当“retention_period_in_days”为空，保留时间默认是7天。
	RetentionPeriodInDays *string `json:"retention_period_in_days,omitempty"`
}

func (o RecyclePolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecyclePolicy struct{}"
	}

	return strings.Join([]string{"RecyclePolicy", string(data)}, " ")
}
