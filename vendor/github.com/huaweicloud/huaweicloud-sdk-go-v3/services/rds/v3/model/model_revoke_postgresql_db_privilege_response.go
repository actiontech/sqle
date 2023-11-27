package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokePostgresqlDbPrivilegeResponse Response Object
type RevokePostgresqlDbPrivilegeResponse struct {

	// 调用正常时，返回“successful”。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RevokePostgresqlDbPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokePostgresqlDbPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"RevokePostgresqlDbPrivilegeResponse", string(data)}, " ")
}
