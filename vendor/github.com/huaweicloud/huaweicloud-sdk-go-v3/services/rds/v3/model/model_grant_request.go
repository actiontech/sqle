package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GrantRequest struct {

	// 数据库名称。
	DbName string `json:"db_name"`

	// 每个元素都是与数据库相关联的帐号。单次请求最多支持50个元素。
	Users []UserWithPrivilege `json:"users"`
}

func (o GrantRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GrantRequest struct{}"
	}

	return strings.Join([]string{"GrantRequest", string(data)}, " ")
}
