package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlHbaConfResponse Response Object
type DeletePostgresqlHbaConfResponse struct {

	// 结果码
	Code *string `json:"code,omitempty"`

	// 结果描述
	Message        *string `json:"message,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeletePostgresqlHbaConfResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlHbaConfResponse struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlHbaConfResponse", string(data)}, " ")
}
