package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UserForCreation struct {

	// 数据库用户名称。 数据库帐号名称在1到32个字符之间，由字母、数字、中划线或下划线组成，不能包含其他特殊字符。 - 若数据库版本为MySQL5.6，帐号长度为1～16个字符。 - 若数据库版本为MySQL5.7和8.0，帐号长度为1～32个字符。
	Name string `json:"name"`

	// 数据库帐号密码。  取值范围：  非空，由大小写字母、数字和特殊符号~!@#$%^*-_=+?,()&组成，长度8~32个字符，不能和数据库帐号“name”或“name”的逆序相同。  建议您输入高强度密码，以提高安全性，防止出现密码被暴力破解等安全风险。
	Password string `json:"password"`

	// 数据库用户备注。 取值范围：长度1~512个字符。目前仅支持MySQL 8.0.25及以上版本。
	Comment *string `json:"comment,omitempty"`

	// 授权用户登录主机IP列表 • 若IP地址为%，则表示允许所有地址访问MySQL实例。 • 若IP地址为“10.10.10.%”，则表示10.10.10.X的IP地址都可以访问该MySQL实例。 • 支持添加多个IP地址。
	Hosts *[]string `json:"hosts,omitempty"`

	// 授权用户数据库权限
	Databases *[]DatabaseWithPrivilegeObject `json:"databases,omitempty"`
}

func (o UserForCreation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserForCreation struct{}"
	}

	return strings.Join([]string{"UserForCreation", string(data)}, " ")
}
