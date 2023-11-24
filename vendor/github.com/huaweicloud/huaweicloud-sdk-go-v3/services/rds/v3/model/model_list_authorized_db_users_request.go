package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAuthorizedDbUsersRequest Request Object
type ListAuthorizedDbUsersRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 数据库名。
	DbName string `json:"db-name"`

	// 分页页码，从1开始。
	Page int32 `json:"page"`

	// 每页数据条数。取值范围[1, 100]。
	Limit int32 `json:"limit"`
}

func (o ListAuthorizedDbUsersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAuthorizedDbUsersRequest struct{}"
	}

	return strings.Join([]string{"ListAuthorizedDbUsersRequest", string(data)}, " ")
}
