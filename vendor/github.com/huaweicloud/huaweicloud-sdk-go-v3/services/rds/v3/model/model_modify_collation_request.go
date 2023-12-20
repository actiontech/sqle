package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyCollationRequest Request Object
type ModifyCollationRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ModifyCollationRequestBody `json:"body,omitempty"`
}

func (o ModifyCollationRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyCollationRequest struct{}"
	}

	return strings.Join([]string{"ModifyCollationRequest", string(data)}, " ")
}
