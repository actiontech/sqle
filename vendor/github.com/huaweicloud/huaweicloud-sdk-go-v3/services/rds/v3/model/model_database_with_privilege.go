package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DatabaseWithPrivilege 数据库及其权限。
type DatabaseWithPrivilege struct {

	// 数据库名称。
	Name string `json:"name"`

	// 是否为只读权限。
	Readonly bool `json:"readonly"`
}

func (o DatabaseWithPrivilege) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DatabaseWithPrivilege struct{}"
	}

	return strings.Join([]string{"DatabaseWithPrivilege", string(data)}, " ")
}
