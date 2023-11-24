package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgreSqlRestoreDatabaseInfo 库级恢复数据库信息
type PostgreSqlRestoreDatabaseInfo struct {

	// 恢复前库名
	OldName *string `json:"old_name,omitempty"`

	// 恢复后库名
	NewName *string `json:"new_name,omitempty"`
}

func (o PostgreSqlRestoreDatabaseInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgreSqlRestoreDatabaseInfo struct{}"
	}

	return strings.Join([]string{"PostgreSqlRestoreDatabaseInfo", string(data)}, " ")
}
