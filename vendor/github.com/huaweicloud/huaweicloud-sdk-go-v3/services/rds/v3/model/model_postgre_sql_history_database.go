package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistoryDatabase PostgreSQL查询可恢复表的数据库信息
type PostgreSqlHistoryDatabase struct {

	// 数据库名
	Name *string `json:"name,omitempty"`

	// 可恢复表的数量
	TotalTables *int32 `json:"total_tables,omitempty"`

	// 模式信息
	Schemas *[]PostgreSqlHistorySchema `json:"schemas,omitempty"`
}

func (o PostgreSqlHistoryDatabase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistoryDatabase struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistoryDatabase", string(data)}, " ")
}
