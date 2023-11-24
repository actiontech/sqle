package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestoreDatabasesInfoNew struct {

	// 库名
	Database string `json:"database"`

	// 表信息
	Tables []RestoreTableInfoNew `json:"tables"`
}

func (o RestoreDatabasesInfoNew) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreDatabasesInfoNew struct{}"
	}

	return strings.Join([]string{"RestoreDatabasesInfoNew", string(data)}, " ")
}
