package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreatePostgresqlDbUserResponse Response Object
type CreatePostgresqlDbUserResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreatePostgresqlDbUserResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreatePostgresqlDbUserResponse struct{}"
	}

	return strings.Join([]string{"CreatePostgresqlDbUserResponse", string(data)}, " ")
}
