package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlDbUserCommentResponse Response Object
type UpdatePostgresqlDbUserCommentResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdatePostgresqlDbUserCommentResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlDbUserCommentResponse struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlDbUserCommentResponse", string(data)}, " ")
}
