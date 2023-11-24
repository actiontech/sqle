package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateInstanceConfigurationRequestBody struct {

	// 参数值对象，用户基于默认参数模板自定义的参数值。为空时不修改参数值。  - key：参数名称，\"max_connections\":\"10\"。为空时不修改参数值，key不为空时value也不可为空。 - value：参数值，\"max_connections\":\"10\"。
	Values map[string]string `json:"values"`
}

func (o UpdateInstanceConfigurationRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceConfigurationRequestBody struct{}"
	}

	return strings.Join([]string{"UpdateInstanceConfigurationRequestBody", string(data)}, " ")
}
