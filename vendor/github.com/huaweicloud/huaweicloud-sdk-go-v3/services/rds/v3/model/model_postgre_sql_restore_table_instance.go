package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreTableInstance 表级恢复实例信息
type PostgreSqlRestoreTableInstance struct {

	// 恢复时间
	RestoreTime *int64 `json:"restore_time,omitempty"`

	// 实例ID
	InstanceId *string `json:"instance_id,omitempty"`

	// 数据库信息
	Databases *[]PostgreSqlRestoreDatabase `json:"databases,omitempty"`
}

func (o PostgreSqlRestoreTableInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreTableInstance struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreTableInstance", string(data)}, " ")
}
