package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteSqlserverDbUserResponse Response Object
type DeleteSqlserverDbUserResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteSqlserverDbUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteSqlserverDbUserResponse struct{}"
	}

	return strings.Join([]string{"DeleteSqlserverDbUserResponse", string(data)}, " ")
}
