package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistorySchema PostgreSQL查询可恢复表的模式信息
type PostgreSqlHistorySchema struct {

	// 模式名
	Name *string `json:"name,omitempty"`

	// 可恢复表的数量
	TotalTables *int32 `json:"total_tables,omitempty"`

	// 表信息
	Tables *[]PostgreSqlHistoryTable `json:"tables,omitempty"`
}

func (o PostgreSqlHistorySchema) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistorySchema struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistorySchema", string(data)}, " ")
}
