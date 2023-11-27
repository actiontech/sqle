package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserWithPrivilege 用户及其权限。
type UserWithPrivilege struct {

	// 用户名。
	Name string `json:"name"`

	// 是否为只读权限。
	Readonly bool `json:"readonly"`
}

func (o UserWithPrivilege) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserWithPrivilege struct{}"
	}

	return strings.Join([]string{"UserWithPrivilege", string(data)}, " ")
}
