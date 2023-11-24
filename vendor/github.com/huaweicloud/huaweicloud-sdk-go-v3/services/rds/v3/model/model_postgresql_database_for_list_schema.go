package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// PostgresqlDatabaseForListSchema 数据库schema信息。
type PostgresqlDatabaseForListSchema struct {

	// schema名称。
	SchemaName string `json:"schema_name"`

	// schema所属用户。
	Owner string `json:"owner"`
}

func (o PostgresqlDatabaseForListSchema) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PostgresqlDatabaseForListSchema struct{}"
	}

	return strings.Join([]string{"PostgresqlDatabaseForListSchema", string(data)}, " ")
}
