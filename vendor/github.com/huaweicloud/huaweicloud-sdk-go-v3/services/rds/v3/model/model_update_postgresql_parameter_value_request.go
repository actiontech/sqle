package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlParameterValueRequest Request Object
type UpdatePostgresqlParameterValueRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 参数名称。
	Name string `json:"name"`

	Body *ModifyParamRequest `json:"body,omitempty"`
}

func (o UpdatePostgresqlParameterValueRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlParameterValueRequest struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlParameterValueRequest", string(data)}, " ")
}
