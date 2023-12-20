package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlDatabaseResponse Response Object
type UpdatePostgresqlDatabaseResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdatePostgresqlDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlDatabaseResponse struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlDatabaseResponse", string(data)}, " ")
}
