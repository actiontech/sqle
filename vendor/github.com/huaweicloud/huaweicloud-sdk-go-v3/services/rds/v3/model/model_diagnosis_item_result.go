package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DiagnosisItemResult struct {

	// 诊断项
	Name *string `json:"name,omitempty"`

	// 实例数量
	Count *int32 `json:"count,omitempty"`
}

func (o DiagnosisItemResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DiagnosisItemResult struct{}"
	}

	return strings.Join([]string{"DiagnosisItemResult", string(data)}, " ")
}
