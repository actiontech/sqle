package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreTableRequest 表级时间点恢复的请求信息
type PostgreSqlRestoreTableRequest struct {

	// 表信息
	Instances *[]PostgreSqlRestoreTableInstance `json:"instances,omitempty"`
}

func (o PostgreSqlRestoreTableRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreTableRequest struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreTableRequest", string(data)}, " ")
}
