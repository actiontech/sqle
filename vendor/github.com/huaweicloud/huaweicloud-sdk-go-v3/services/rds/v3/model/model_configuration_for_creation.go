package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ConfigurationForCreation struct {

	// 参数模板名称。最长64个字符，只允许大写字母、小写字母、数字、和“-_.”特殊字符。
	Name string `json:"name"`

	// 参数模板描述。最长256个字符，不支持>!<\"&'=特殊字符。默认为空。
	Description *string `json:"description,omitempty"`

	Datastore *ParaGroupDatastore `json:"datastore"`

	// 参数值对象，用户基于默认参数模板自定义的参数值。为空时不修改参数值。  - key：参数名称，\"max_connections\":\"10\"。为空时不修改参数值，key不为空时value也不可为空。 - value：参数值，\"max_connections\":\"10\"。
	Values map[string]string `json:"values,omitempty"`
}

func (o ConfigurationForCreation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ConfigurationForCreation struct{}"
	}

	return strings.Join([]string{"ConfigurationForCreation", string(data)}, " ")
}
