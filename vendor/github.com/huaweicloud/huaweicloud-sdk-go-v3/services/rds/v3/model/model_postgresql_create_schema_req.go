package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PostgresqlCreateSchemaReq struct {

	// schema名称。  schema名称在1到63个字符之间，由字母、数字、或下划线组成，不能包含其他特殊字符，不能以“pg”和数字开头，且不能和RDS for PostgreSQL模板库和已存在的schema重名。 RDS for PostgreSQL模板库包括postgres， template0 ，template1。  已存在的schema包括public，information_schema。
	SchemaName string `json:"schema_name"`

	// 数据库属主用户。  数据库属主名称在1到63个字符之间，不能以“pg”和数字开头，不能和系统用户名称相同。  系统用户包括“rdsAdmin”,“ rdsMetric”, “rdsBackup”, “rdsRepl”,“ rdsProxy”, “rdsDdm”。
	Owner string `json:"owner"`
}

func (o PostgresqlCreateSchemaReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlCreateSchemaReq struct{}"
	}

	return strings.Join([]string{"PostgresqlCreateSchemaReq", string(data)}, " ")
}
