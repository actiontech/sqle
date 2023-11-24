package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlDbUserResponse Response Object
type DeletePostgresqlDbUserResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeletePostgresqlDbUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlDbUserResponse struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlDbUserResponse", string(data)}, " ")
}
