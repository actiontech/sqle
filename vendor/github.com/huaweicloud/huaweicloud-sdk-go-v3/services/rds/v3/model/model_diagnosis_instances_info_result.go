package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DiagnosisInstancesInfoResult struct {

	// 实例id
	Id *string `json:"id,omitempty"`
}

func (o DiagnosisInstancesInfoResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DiagnosisInstancesInfoResult struct{}"
	}

	return strings.Join([]string{"DiagnosisInstancesInfoResult", string(data)}, " ")
}
