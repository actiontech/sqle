package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowPostgresqlParamValueResponse Response Object
type ShowPostgresqlParamValueResponse struct {

	// 参数名称。
	Name *string `json:"name,omitempty"`

	// 参数值。
	Value *string `json:"value,omitempty"`

	// 是否需要重启。 - \"false\"表示否 - \"true\"表示是
	RestartRequired *bool `json:"restart_required,omitempty"`

	// 参数值范围，如Integer取值0-1、Boolean取值true|false等。
	ValueRange *string `json:"value_range,omitempty"`

	// 参数类型，取值为“string”、“integer”、“boolean”、“list”或“float”之一。
	Type *string `json:"type,omitempty"`

	// 参数描述。
	Description    *string `json:"description,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowPostgresqlParamValueResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowPostgresqlParamValueResponse struct{}"
	}

	return strings.Join([]string{"ShowPostgresqlParamValueResponse", string(data)}, " ")
}
