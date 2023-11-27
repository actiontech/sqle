package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePostgresqlDatabaseSchemaResponse Response Object
type CreatePostgresqlDatabaseSchemaResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreatePostgresqlDatabaseSchemaResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePostgresqlDatabaseSchemaResponse struct{}"
	}

	return strings.Join([]string{"CreatePostgresqlDatabaseSchemaResponse", string(data)}, " ")
}
