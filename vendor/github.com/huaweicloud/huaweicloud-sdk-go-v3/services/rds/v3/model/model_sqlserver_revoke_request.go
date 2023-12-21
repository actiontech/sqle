package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SqlserverRevokeRequest struct {

	// 数据库名称。
	DbName string `json:"db_name"`

	// 每个元素都是与数据库相关联的帐号。单次请求最多支持50个元素。
	Users []SqlserverUserWithPrivilege `json:"users"`
}

func (o SqlserverRevokeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SqlserverRevokeRequest struct{}"
	}

	return strings.Join([]string{"SqlserverRevokeRequest", string(data)}, " ")
}
