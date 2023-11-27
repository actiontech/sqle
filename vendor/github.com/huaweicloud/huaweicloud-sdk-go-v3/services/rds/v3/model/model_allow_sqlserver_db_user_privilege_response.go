package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowSqlserverDbUserPrivilegeResponse Response Object
type AllowSqlserverDbUserPrivilegeResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AllowSqlserverDbUserPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowSqlserverDbUserPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"AllowSqlserverDbUserPrivilegeResponse", string(data)}, " ")
}
