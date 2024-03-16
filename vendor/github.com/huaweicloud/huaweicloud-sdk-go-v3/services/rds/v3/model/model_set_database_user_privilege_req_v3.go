package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SetDatabaseUserPrivilegeReqV3 struct {

	// 是否设置所有用户。
	AllUsers bool `json:"all_users"`

	// 数据库用户名。
	UserName *string `json:"user_name,omitempty"`

	// 是否为只读权限。
	Readonly bool `json:"readonly"`
}

func (o SetDatabaseUserPrivilegeReqV3) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetDatabaseUserPrivilegeReqV3 struct{}"
	}

	return strings.Join([]string{"SetDatabaseUserPrivilegeReqV3", string(data)}, " ")
}
