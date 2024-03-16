package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchRestorePostgreSqlTablesRequest Request Object
type BatchRestorePostgreSqlTablesRequest struct {
	Body *PostgreSqlRestoreTableRequest `json:"body,omitempty"`
}

func (o BatchRestorePostgreSqlTablesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchRestorePostgreSqlTablesRequest struct{}"
	}

	return strings.Join([]string{"BatchRestorePostgreSqlTablesRequest", string(data)}, " ")
}
