package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPostgresqlParamValueRequest Request Object
type ShowPostgresqlParamValueRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 参数名称。
	Name string `json:"name"`
}

func (o ShowPostgresqlParamValueRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPostgresqlParamValueRequest struct{}"
	}

	return strings.Join([]string{"ShowPostgresqlParamValueRequest", string(data)}, " ")
}
