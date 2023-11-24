package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestoreTablesRequestBody struct {

	// 恢复时间戳
	RestoreTime int64 `json:"restoreTime"`

	// 表信息
	RestoreTables []RestoreDatabasesInfo `json:"restoreTables"`

	// 是否使用极速恢复，可先根据”获取实例是否能使用极速恢复“接口判断本次恢复是否能使用急速恢复。 如果实例使用了XA事务，则不可使用极速恢复！使用恢复会导致恢复失败！
	IsFastRestore *bool `json:"is_fast_restore,omitempty"`
}

func (o RestoreTablesRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreTablesRequestBody struct{}"
	}

	return strings.Join([]string{"RestoreTablesRequestBody", string(data)}, " ")
}
