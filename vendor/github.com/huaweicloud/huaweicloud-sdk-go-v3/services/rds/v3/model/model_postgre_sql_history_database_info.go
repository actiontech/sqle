package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistoryDatabaseInfo PostgreSQL查询可恢复库的数据库库信息
type PostgreSqlHistoryDatabaseInfo struct {

	// 数据库名
	Name *string `json:"name,omitempty"`

	// 表的个数
	TotalTables *int32 `json:"total_tables,omitempty"`
}

func (o PostgreSqlHistoryDatabaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistoryDatabaseInfo struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistoryDatabaseInfo", string(data)}, " ")
}
