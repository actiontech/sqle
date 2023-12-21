package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistoryDatabaseRequest 查询可恢复库的请求信息
type PostgreSqlHistoryDatabaseRequest struct {

	// 实例ID集合
	InstanceIds []string `json:"instance_ids"`

	// 恢复时间点
	RestoreTime int64 `json:"restore_time"`

	// 数据库名，模糊查询
	DatabaseNameLike *string `json:"database_name_like,omitempty"`

	// 实例名称，模糊查询
	InstanceNameLike *string `json:"instance_name_like,omitempty"`
}

func (o PostgreSqlHistoryDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistoryDatabaseRequest struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistoryDatabaseRequest", string(data)}, " ")
}
