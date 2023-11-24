package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlHistoryTable PostgreSQL查询可恢复表的表信息
type PostgreSqlHistoryTable struct {

	// 表名
	Name *string `json:"name,omitempty"`
}

func (o PostgreSqlHistoryTable) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlHistoryTable struct{}"
	}

	return strings.Join([]string{"PostgreSqlHistoryTable", string(data)}, " ")
}
