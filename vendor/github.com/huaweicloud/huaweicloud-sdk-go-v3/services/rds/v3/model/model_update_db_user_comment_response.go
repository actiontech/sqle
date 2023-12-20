package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDbUserCommentResponse Response Object
type UpdateDbUserCommentResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDbUserCommentResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDbUserCommentResponse struct{}"
	}

	return strings.Join([]string{"UpdateDbUserCommentResponse", string(data)}, " ")
}
