package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlDbUserCommentRequest Request Object
type UpdatePostgresqlDbUserCommentRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 数据库用户名。
	UserName string `json:"user_name"`

	Body *UpdateDbUserReq `json:"body,omitempty"`
}

func (o UpdatePostgresqlDbUserCommentRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlDbUserCommentRequest struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlDbUserCommentRequest", string(data)}, " ")
}
