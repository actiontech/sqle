package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgresqlUserForList 数据库用户信息。
type PostgresqlUserForList struct {

	// 帐号名。
	Name string `json:"name"`

	// 用户的权限属性。
	Attributes *interface{} `json:"attributes,omitempty"`

	// 用户的默认权限。
	Memberof *[]string `json:"memberof,omitempty"`

	// 数据库用户备注。
	Comment *string `json:"comment,omitempty"`
}

func (o PostgresqlUserForList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlUserForList struct{}"
	}

	return strings.Join([]string{"PostgresqlUserForList", string(data)}, " ")
}
