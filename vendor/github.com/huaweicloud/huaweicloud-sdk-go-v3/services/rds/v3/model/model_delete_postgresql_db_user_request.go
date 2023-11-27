package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlDbUserRequest Request Object
type DeletePostgresqlDbUserRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 需要删除的帐号名。
	UserName string `json:"user_name"`
}

func (o DeletePostgresqlDbUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlDbUserRequest struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlDbUserRequest", string(data)}, " ")
}
