package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlDatabaseResponse Response Object
type DeletePostgresqlDatabaseResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeletePostgresqlDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlDatabaseResponse struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlDatabaseResponse", string(data)}, " ")
}
