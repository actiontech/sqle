package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlDatabaseRequest Request Object
type DeletePostgresqlDatabaseRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 需要删除的数据库名。
	DbName string `json:"db_name"`
}

func (o DeletePostgresqlDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlDatabaseRequest struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlDatabaseRequest", string(data)}, " ")
}
