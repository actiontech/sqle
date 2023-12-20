package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetPostgresqlDbUserPwdResponse Response Object
type SetPostgresqlDbUserPwdResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetPostgresqlDbUserPwdResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetPostgresqlDbUserPwdResponse struct{}"
	}

	return strings.Join([]string{"SetPostgresqlDbUserPwdResponse", string(data)}, " ")
}
