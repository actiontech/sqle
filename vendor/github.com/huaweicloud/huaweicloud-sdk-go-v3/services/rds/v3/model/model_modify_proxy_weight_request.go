package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyProxyWeightRequest struct {

	// 主实例权重，取值范围为0~1000。
	MasterWeight string `json:"master_weight"`

	// 只读实例信息。
	ReadonlyInstances []ProxyReadonlyInstances `json:"readonly_instances"`
}

func (o ModifyProxyWeightRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyProxyWeightRequest struct{}"
	}

	return strings.Join([]string{"ModifyProxyWeightRequest", string(data)}, " ")
}
