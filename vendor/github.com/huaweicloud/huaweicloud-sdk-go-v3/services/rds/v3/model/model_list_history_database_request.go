package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListHistoryDatabaseRequest Request Object
type ListHistoryDatabaseRequest struct {

	// 数据库引擎。支持的引擎如下，不区分大小写：postgresql
	DatabaseName string `json:"database_name"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *PostgreSqlHistoryDatabaseRequest `json:"body,omitempty"`
}

func (o ListHistoryDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListHistoryDatabaseRequest struct{}"
	}

	return strings.Join([]string{"ListHistoryDatabaseRequest", string(data)}, " ")
}
