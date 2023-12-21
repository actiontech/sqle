package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestoreTableInfo struct {

	// 旧表名
	OldName string `json:"oldName"`

	// 新表名
	NewName string `json:"newName"`
}

func (o RestoreTableInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreTableInfo struct{}"
	}

	return strings.Join([]string{"RestoreTableInfo", string(data)}, " ")
}
