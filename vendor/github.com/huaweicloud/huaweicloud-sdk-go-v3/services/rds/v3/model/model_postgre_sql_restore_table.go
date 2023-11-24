package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreTable 恢复表信息
type PostgreSqlRestoreTable struct {

	// 恢复前表名
	OldName *string `json:"old_name,omitempty"`

	// 恢复后表名
	NewName *string `json:"new_name,omitempty"`
}

func (o PostgreSqlRestoreTable) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreTable struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreTable", string(data)}, " ")
}
