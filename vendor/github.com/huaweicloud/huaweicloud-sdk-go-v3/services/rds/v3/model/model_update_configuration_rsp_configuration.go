package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateConfigurationRspConfiguration 参数模板信息
type UpdateConfigurationRspConfiguration struct {

	// 参数模板ID。
	Id *string `json:"id,omitempty"`

	// 参数模板名称。
	Name *string `json:"name,omitempty"`

	// 请求参数“values”中被忽略掉，没有生效的参数名称列表。 当参数不存在时，参数修改不会下发，并通过此参数返回所有被忽略的参数名称。
	IgnoredParams *[]string `json:"ignored_params,omitempty"`
}

func (o UpdateConfigurationRspConfiguration) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateConfigurationRspConfiguration struct{}"
	}

	return strings.Join([]string{"UpdateConfigurationRspConfiguration", string(data)}, " ")
}
