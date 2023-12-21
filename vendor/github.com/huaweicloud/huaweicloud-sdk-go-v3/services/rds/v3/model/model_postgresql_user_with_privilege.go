package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgresqlUserWithPrivilege 用户及其权限。
type PostgresqlUserWithPrivilege struct {

	// 数据库帐号名称。  数据库帐号名称在1到63个字符之间，由字母、数字、或下划线组成，不能包含其他特殊字符，不能以“pg”和数字开头，不能和系统用户名称相同且帐号名称必须存在。  系统用户包括“rdsAdmin”,“ rdsMetric”, “rdsBackup”, “rdsRepl”,“ rdsProxy”, “rdsDdm”。
	Name string `json:"name"`

	// 数据库帐号权限。 - true：只读。 - false：可读可写。
	Readonly bool `json:"readonly"`

	// schema名称。  schema名称在1到63个字符之间，由字母、数字、或下划线组成，不能包含其他特殊字符，不能以“pg”和数字开头，不能和RDS for PostgreSQL模板库重名，且schema名称必须存在。  RDS for PostgreSQL模板库包括postgres， template0 ，template1。
	SchemaName string `json:"schema_name"`
}

func (o PostgresqlUserWithPrivilege) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlUserWithPrivilege struct{}"
	}

	return strings.Join([]string{"PostgresqlUserWithPrivilege", string(data)}, " ")
}
