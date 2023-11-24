package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreDatabase 恢复库信息
type PostgreSqlRestoreDatabase struct {

	// 数据库名
	Database *string `json:"database,omitempty"`

	// 模式信息
	Schemas *[]PostgreSqlRestoreSchema `json:"schemas,omitempty"`
}

func (o PostgreSqlRestoreDatabase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreDatabase struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreDatabase", string(data)}, " ")
}
