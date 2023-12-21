package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SupportFastRestoreList struct {

	// 实例id。
	InstanceId *string `json:"instance_id,omitempty"`

	// 表级恢复是否支持极速恢复。
	IsSupportFastTableRestore *bool `json:"is_support_fast_table_restore,omitempty"`

	// 库级恢复是否支持极速恢复。
	IsSupportFastDatabaseRestore *bool `json:"is_support_fast_database_restore,omitempty"`
}

func (o SupportFastRestoreList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SupportFastRestoreList struct{}"
	}

	return strings.Join([]string{"SupportFastRestoreList", string(data)}, " ")
}
