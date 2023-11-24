package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DbUserPrivilegeRequest 入参对象
type DbUserPrivilegeRequest struct {

	// 账号，数据库相关联的帐号
	UserName string `json:"user_name"`

	// 授权SQL类型,枚举： 1、ROLE(支持对指定用户设置以下权限) 2、RECYCLING_ROLE(支持对指定用户回收以下权限) 3、SYSTEM_ROLE(支持对指定用户授予以下系统角色) 4、RECYCLING_SYSTEM_ROLE(支持对指定用户回收以下系统角色)
	AuthorizationType string `json:"authorization_type"`

	// 支持用户设置的权限集合。 1、privilege_type为ROLE时需要使用,枚举： CREATEDB CREATEROLE LOGIN REPLICATION 2、privilege_type为RECYCLING_ROLE时需要使用,枚举： NOCREATEDB NOCREATEROLE NOLOGIN NOREPLICATION 3、privilege_type为SYSTEM_ROLE /RECYCLING_ SYSTEM_ROLE时需要使用,枚举： pg_monitor pg_signal_backend root
	Privileges []string `json:"privileges"`
}

func (o DbUserPrivilegeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DbUserPrivilegeRequest struct{}"
	}

	return strings.Join([]string{"DbUserPrivilegeRequest", string(data)}, " ")
}
