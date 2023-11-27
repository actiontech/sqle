package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDbUserRequest Request Object
type CreateDbUserRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UserForCreation `json:"body,omitempty"`
}

func (o CreateDbUserRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDbUserRequest struct{}"
	}

	return strings.Join([]string{"CreateDbUserRequest", string(data)}, " ")
}
