package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ProxyReadonlyInstances struct {

	// 只读实例ID。
	Id string `json:"id"`

	// 只读实例权重，取值范围为0~1000。
	Weight int32 `json:"weight"`
}

func (o ProxyReadonlyInstances) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProxyReadonlyInstances struct{}"
	}

	return strings.Join([]string{"ProxyReadonlyInstances", string(data)}, " ")
}
