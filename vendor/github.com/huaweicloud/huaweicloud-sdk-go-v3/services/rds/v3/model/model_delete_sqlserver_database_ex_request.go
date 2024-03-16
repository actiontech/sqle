package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSqlserverDatabaseExRequest Request Object
type DeleteSqlserverDatabaseExRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 需要删除的数据库名。
	DbName string `json:"db_name"`

	Body *DropDatabaseV3Req `json:"body,omitempty"`
}

func (o DeleteSqlserverDatabaseExRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlserverDatabaseExRequest struct{}"
	}

	return strings.Join([]string{"DeleteSqlserverDatabaseExRequest", string(data)}, " ")
}
