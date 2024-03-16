package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateSqlserverDbUserResponse Response Object
type CreateSqlserverDbUserResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateSqlserverDbUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateSqlserverDbUserResponse struct{}"
	}

	return strings.Join([]string{"CreateSqlserverDbUserResponse", string(data)}, " ")
}
