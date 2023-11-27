package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlInstanceAliasResponse Response Object
type UpdatePostgresqlInstanceAliasResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdatePostgresqlInstanceAliasResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlInstanceAliasResponse struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlInstanceAliasResponse", string(data)}, " ")
}
