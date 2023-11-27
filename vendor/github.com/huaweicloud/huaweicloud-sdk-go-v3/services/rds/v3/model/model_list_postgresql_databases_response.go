package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlDatabasesResponse Response Object
type ListPostgresqlDatabasesResponse struct {

	// 列表中每个元素表示一个数据库。
	Databases *[]PostgresqlListDatabase `json:"databases,omitempty"`

	// 数据库总数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListPostgresqlDatabasesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlDatabasesResponse struct{}"
	}

	return strings.Join([]string{"ListPostgresqlDatabasesResponse", string(data)}, " ")
}
