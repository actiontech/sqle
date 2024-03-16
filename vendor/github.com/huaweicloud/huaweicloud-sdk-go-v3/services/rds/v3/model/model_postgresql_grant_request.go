package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostgresqlGrantRequest struct {

	// 数据库名称。  数据库名称在1到63个字符之间，由字母、数字、或下划线组成，不能包含其他特殊字符，不能以“pg”和数字开头，且不能和RDS for PostgreSQL模板库重名。  RDS for PostgreSQL模板库包括postgres， template0 ，template1。
	DbName string `json:"db_name"`

	// 每个元素都是与数据库相关联的帐号。单次请求最多支持50个元素。
	Users []PostgresqlUserWithPrivilege `json:"users"`
}

func (o PostgresqlGrantRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlGrantRequest struct{}"
	}

	return strings.Join([]string{"PostgresqlGrantRequest", string(data)}, " ")
}
