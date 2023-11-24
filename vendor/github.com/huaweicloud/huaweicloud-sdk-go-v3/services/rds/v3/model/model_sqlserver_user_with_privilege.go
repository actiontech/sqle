package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SqlserverUserWithPrivilege 数据库帐号名称。  数据库帐号名称在1到128个字符之间，不能和系统用户名称相同。  系统用户包括：rdsadmin, rdsuser, rdsbackup, rdsmirror。
type SqlserverUserWithPrivilege struct {

	// 数据库帐号名称。
	Name string `json:"name"`

	// 是否为只读权限。
	Readonly *bool `json:"readonly,omitempty"`
}

func (o SqlserverUserWithPrivilege) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SqlserverUserWithPrivilege struct{}"
	}

	return strings.Join([]string{"SqlserverUserWithPrivilege", string(data)}, " ")
}
