package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreResult 表级时间点恢复的请求信息
type PostgreSqlRestoreResult struct {

	// 实例ID
	InstanceId *string `json:"instance_id,omitempty"`

	// 工作流id
	JobId *string `json:"job_id,omitempty"`
}

func (o PostgreSqlRestoreResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreResult struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreResult", string(data)}, " ")
}
