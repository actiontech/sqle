package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AdDomainInfo 域信息，加域实例单转主备时必填，非加域实例不需要填写
type AdDomainInfo struct {

	// 域管理员账号名
	DomainAdminAccountName string `json:"domain_admin_account_name"`

	// 域管理员密码
	DomainAdminPwd string `json:"domain_admin_pwd"`
}

func (o AdDomainInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AdDomainInfo struct{}"
	}

	return strings.Join([]string{"AdDomainInfo", string(data)}, " ")
}
