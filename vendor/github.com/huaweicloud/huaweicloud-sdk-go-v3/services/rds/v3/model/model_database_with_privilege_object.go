package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DatabaseWithPrivilegeObject 数据库及其权限。
type DatabaseWithPrivilegeObject struct {

	// 已有数据库名称。
	Name string `json:"name"`

	// 是否为只读权限。
	Readonly *bool `json:"readonly,omitempty"`
}

func (o DatabaseWithPrivilegeObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DatabaseWithPrivilegeObject struct{}"
	}

	return strings.Join([]string{"DatabaseWithPrivilegeObject", string(data)}, " ")
}
