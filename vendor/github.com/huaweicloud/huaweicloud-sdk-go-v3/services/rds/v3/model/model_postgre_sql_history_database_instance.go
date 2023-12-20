package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistoryDatabaseInstance PostgreSQL查询可恢复库的实例信息
type PostgreSqlHistoryDatabaseInstance struct {

	// 实例ID
	Id *string `json:"id,omitempty"`

	// 实例名称
	Name *string `json:"name,omitempty"`

	// 表的个数
	TotalTables *int32 `json:"total_tables,omitempty"`

	// 数据库信息
	Databases *[]PostgreSqlHistoryDatabaseInfo `json:"databases,omitempty"`
}

func (o PostgreSqlHistoryDatabaseInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistoryDatabaseInstance struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistoryDatabaseInstance", string(data)}, " ")
}
