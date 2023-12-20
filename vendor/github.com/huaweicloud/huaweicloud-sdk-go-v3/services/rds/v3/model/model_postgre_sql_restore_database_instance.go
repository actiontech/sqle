package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreDatabaseInstance 查询可恢复库的响应信息
type PostgreSqlRestoreDatabaseInstance struct {

	// 恢复时间
	RestoreTime *int64 `json:"restore_time,omitempty"`

	// 实例ID
	InstanceId *string `json:"instance_id,omitempty"`

	// 库信息
	Databases *[]PostgreSqlRestoreDatabaseInfo `json:"databases,omitempty"`
}

func (o PostgreSqlRestoreDatabaseInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreDatabaseInstance struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreDatabaseInstance", string(data)}, " ")
}
