package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSqlserverDbUserRequest Request Object
type DeleteSqlserverDbUserRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 需要删除的帐号名。
	UserName string `json:"user_name"`
}

func (o DeleteSqlserverDbUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlserverDbUserRequest struct{}"
	}

	return strings.Join([]string{"DeleteSqlserverDbUserRequest", string(data)}, " ")
}
