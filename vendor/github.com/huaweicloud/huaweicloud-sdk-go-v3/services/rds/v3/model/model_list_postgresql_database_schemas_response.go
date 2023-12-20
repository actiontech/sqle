package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlDatabaseSchemasResponse Response Object
type ListPostgresqlDatabaseSchemasResponse struct {

	// 列表中每个元素表示一个数据库schema。
	DatabaseSchemas *[]PostgresqlDatabaseForListSchema `json:"database_schemas,omitempty"`

	// 数据库schema总数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListPostgresqlDatabaseSchemasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlDatabaseSchemasResponse struct{}"
	}

	return strings.Join([]string{"ListPostgresqlDatabaseSchemasResponse", string(data)}, " ")
}
