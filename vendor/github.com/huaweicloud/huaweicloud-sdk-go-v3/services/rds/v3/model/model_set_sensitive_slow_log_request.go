package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetSensitiveSlowLogRequest Request Object
type SetSensitiveSlowLogRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID，可以调用“查询实例列表”接口获取。如果未申请实例，可以调用“创建实例”接口创建。
	InstanceId string `json:"instance_id"`

	// 开启或关闭慢日志敏感信息明文，取值为on或off。
	Status string `json:"status"`
}

func (o SetSensitiveSlowLogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetSensitiveSlowLogRequest struct{}"
	}

	return strings.Join([]string{"SetSensitiveSlowLogRequest", string(data)}, " ")
}
