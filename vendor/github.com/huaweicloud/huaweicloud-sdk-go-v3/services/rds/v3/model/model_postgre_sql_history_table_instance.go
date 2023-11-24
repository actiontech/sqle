package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistoryTableInstance PostgreSQL查询可恢复表的实例信息
type PostgreSqlHistoryTableInstance struct {

	// 实例ID
	Id *string `json:"id,omitempty"`

	// 实例名称
	Name *string `json:"name,omitempty"`

	// 可恢复表的数量
	TotalTables *int32 `json:"total_tables,omitempty"`

	// 数据库信息
	Databases *[]PostgreSqlHistoryDatabase `json:"databases,omitempty"`
}

func (o PostgreSqlHistoryTableInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistoryTableInstance struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistoryTableInstance", string(data)}, " ")
}
