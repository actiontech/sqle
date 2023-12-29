package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RevokePostgresqlDbPrivilegeUser struct {

	// 数据库账号名称
	Name string `json:"name"`

	// 数据库下模式名称
	SchemaName string `json:"schema_name"`
}

func (o RevokePostgresqlDbPrivilegeUser) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokePostgresqlDbPrivilegeUser struct{}"
	}

	return strings.Join([]string{"RevokePostgresqlDbPrivilegeUser", string(data)}, " ")
}
