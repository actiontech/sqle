package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlListHistoryTablesResponse Response Object
type ListPostgresqlListHistoryTablesResponse struct {

	// 恢复表数量限制个数
	TableLimit *int32 `json:"table_limit,omitempty"`

	// 实例信息
	Instances      *[]PostgreSqlHistoryTableInstance `json:"instances,omitempty"`
	HttpStatusCode int                               `json:"-"`
}

func (o ListPostgresqlListHistoryTablesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlListHistoryTablesResponse struct{}"
	}

	return strings.Join([]string{"ListPostgresqlListHistoryTablesResponse", string(data)}, " ")
}
