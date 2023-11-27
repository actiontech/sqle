package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type SimplifiedInstancesRequest struct {

	// 实例id集合
	InstanceIds []string `json:"instance_ids"`
}

func (o SimplifiedInstancesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SimplifiedInstancesRequest struct{}"
	}

	return strings.Join([]string{"SimplifiedInstancesRequest", string(data)}, " ")
}
