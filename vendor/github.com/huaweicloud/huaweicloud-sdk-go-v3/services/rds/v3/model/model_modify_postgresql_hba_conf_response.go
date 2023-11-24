package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyPostgresqlHbaConfResponse Response Object
type ModifyPostgresqlHbaConfResponse struct {

	// 结果码
	Code *string `json:"code,omitempty"`

	// 结果描述
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ModifyPostgresqlHbaConfResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyPostgresqlHbaConfResponse struct{}"
	}

	return strings.Join([]string{"ModifyPostgresqlHbaConfResponse", string(data)}, " ")
}
