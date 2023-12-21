package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHistoryDatabaseResponse Response Object
type ListHistoryDatabaseResponse struct {

	// 恢复库数量限制个数
	DatabaseLimit *int32 `json:"database_limit,omitempty"`

	// 恢复表数量限制个数
	TableLimit *int32 `json:"table_limit,omitempty"`

	// 实例信息
	Instances      *[]PostgreSqlHistoryDatabaseInstance `json:"instances,omitempty"`
	HttpStatusCode int                                  `json:"-"`
}

func (o ListHistoryDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHistoryDatabaseResponse struct{}"
	}

	return strings.Join([]string{"ListHistoryDatabaseResponse", string(data)}, " ")
}
