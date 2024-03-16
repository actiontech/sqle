package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DbUserPwdRequest struct {

	// 数据库帐号名称。
	Name string `json:"name"`

	// 数据库帐号密码。  取值范围：  非空，至少包含以下字符中的三种：大写字母、小写字母、数字和特殊符号~!@#%^*-_=+?,组成，长度8~32个字符，不能和数据库帐号“name”或“name”的逆序相同。  建议您输入高强度密码，以提高安全性，防止出现密码被暴力破解等安全风险。
	Password string `json:"password"`
}

func (o DbUserPwdRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DbUserPwdRequest struct{}"
	}

	return strings.Join([]string{"DbUserPwdRequest", string(data)}, " ")
}
