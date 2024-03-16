package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSqlserverDatabaseResponse Response Object
type CreateSqlserverDatabaseResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateSqlserverDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSqlserverDatabaseResponse struct{}"
	}

	return strings.Join([]string{"CreateSqlserverDatabaseResponse", string(data)}, " ")
}
