package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RevokePostgresqlDbPrivilegeRequestBody struct {

	// 数据库名称
	DbName string `json:"db_name"`

	// 用户信息，最大值50个
	Users []RevokePostgresqlDbPrivilegeUser `json:"users"`
}

func (o RevokePostgresqlDbPrivilegeRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokePostgresqlDbPrivilegeRequestBody struct{}"
	}

	return strings.Join([]string{"RevokePostgresqlDbPrivilegeRequestBody", string(data)}, " ")
}
