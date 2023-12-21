package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type OpenProxyRequest struct {

	// 规格ID。
	FlavorId *string `json:"flavor_id,omitempty"`

	// 节点数量。
	NodeNum *int32 `json:"node_num,omitempty"`
}

func (o OpenProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OpenProxyRequest struct{}"
	}

	return strings.Join([]string{"OpenProxyRequest", string(data)}, " ")
}
