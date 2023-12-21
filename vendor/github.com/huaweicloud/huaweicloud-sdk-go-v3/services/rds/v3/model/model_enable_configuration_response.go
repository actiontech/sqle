package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EnableConfigurationResponse Response Object
type EnableConfigurationResponse struct {

	// 参数组ID。
	ConfigurationId *string `json:"configuration_id,omitempty"`

	// 参数组名称。
	ConfigurationName *string `json:"configuration_name,omitempty"`

	// 参数模板是否都应用成功。  - “true”表示参数模板都应用成功。 - “false”表示存在应用失败的参数模板。
	Success *bool `json:"success,omitempty"`

	// 对每个实例的应用结果。
	ApplyResults   *[]ApplyConfigurationResponseApplyResults `json:"apply_results,omitempty"`
	HttpStatusCode int                                       `json:"-"`
}

func (o EnableConfigurationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnableConfigurationResponse struct{}"
	}

	return strings.Join([]string{"EnableConfigurationResponse", string(data)}, " ")
}
