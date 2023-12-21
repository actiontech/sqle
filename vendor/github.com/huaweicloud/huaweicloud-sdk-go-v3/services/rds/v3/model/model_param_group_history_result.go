package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ParamGroupHistoryResult struct {

	// 参数名称
	ParameterName *string `json:"parameter_name,omitempty"`

	// 旧值
	OldValue *string `json:"old_value,omitempty"`

	// 新值
	NewValue *string `json:"new_value,omitempty"`

	// 更新结果 成功：SUCCESS 失败： FAILED
	UpdateResult *string `json:"update_result,omitempty"`

	// 是否已应用 true：已应用 false：未应用
	Applied *bool `json:"applied,omitempty"`

	// 修改时间
	UpdateTime *string `json:"update_time,omitempty"`

	// 应用时间
	ApplyTime *string `json:"apply_time,omitempty"`
}

func (o ParamGroupHistoryResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ParamGroupHistoryResult struct{}"
	}

	return strings.Join([]string{"ParamGroupHistoryResult", string(data)}, " ")
}
