package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ScaleProxyRequestBody struct {

	// 需要变更的新规格ID。
	FlavorRef string `json:"flavor_ref"`

	// 是否延迟变更。  - true：延迟变更，将在运维时间窗内自动变更。 - false：立即变更。
	Delay bool `json:"delay"`
}

func (o ScaleProxyRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleProxyRequestBody struct{}"
	}

	return strings.Join([]string{"ScaleProxyRequestBody", string(data)}, " ")
}
