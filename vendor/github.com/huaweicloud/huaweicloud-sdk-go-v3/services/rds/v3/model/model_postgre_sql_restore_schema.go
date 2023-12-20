package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreSchema 恢复模式信息
type PostgreSqlRestoreSchema struct {

	// 模式信息
	Schema *string `json:"schema,omitempty"`

	// 表信息
	Tables *[]PostgreSqlRestoreTable `json:"tables,omitempty"`
}

func (o PostgreSqlRestoreSchema) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreSchema struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreSchema", string(data)}, " ")
}
