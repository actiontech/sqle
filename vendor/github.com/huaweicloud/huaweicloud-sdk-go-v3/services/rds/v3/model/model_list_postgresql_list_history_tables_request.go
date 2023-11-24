package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlListHistoryTablesRequest Request Object
type ListPostgresqlListHistoryTablesRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 数据库引擎。支持的引擎如下，不区分大小写：postgresql
	DatabaseName string `json:"database_name"`

	Body *PostgreSqlHistoryTableRequest `json:"body,omitempty"`
}

func (o ListPostgresqlListHistoryTablesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlListHistoryTablesRequest struct{}"
	}

	return strings.Join([]string{"ListPostgresqlListHistoryTablesRequest", string(data)}, " ")
}
