package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SqlserverUserForCreation struct {

	// 数据库用户名称。  数据库帐号名称在1到128个字符之间，不能和系统用户名称相同。  系统用户包括：rdsadmin, rdsuser, rdsbackup, rdsmirror。
	Name string `json:"name"`

	// 数据库帐号密码。  取值范围：非空，密码长度在8到128个字符之间，至少包含大写字母、小写字母、数字、特殊字符三种字符的组合。  建议您输入高强度密码，以提高安全性，防止出现密码被暴力破解等安全风险。
	Password string `json:"password"`
}

func (o SqlserverUserForCreation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SqlserverUserForCreation struct{}"
	}

	return strings.Join([]string{"SqlserverUserForCreation", string(data)}, " ")
}
