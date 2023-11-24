package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BackupDatabase struct {

	// SqlServer引擎指定备份的数据库。
	Name string `json:"name"`
}

func (o BackupDatabase) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupDatabase struct{}"
	}

	return strings.Join([]string{"BackupDatabase", string(data)}, " ")
}
