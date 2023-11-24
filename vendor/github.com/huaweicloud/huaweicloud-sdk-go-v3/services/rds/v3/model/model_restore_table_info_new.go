package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestoreTableInfoNew struct {

	// 旧表名
	OldName string `json:"old_name"`

	// 新表名
	NewName string `json:"new_name"`
}

func (o RestoreTableInfoNew) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreTableInfoNew struct{}"
	}

	return strings.Join([]string{"RestoreTableInfoNew", string(data)}, " ")
}
