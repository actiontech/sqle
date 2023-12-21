package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ReadonlyInstances struct {

	// 只读实例ID。
	Id string `json:"id"`

	// 只读实例状态。
	Status string `json:"status"`

	// 只读实例名称。
	Name string `json:"name"`

	// 只读实例读写分离权重。
	Weight int32 `json:"weight"`

	// 可用区信息。
	AvailableZones []AvailableZone `json:"available_zones"`

	// 只读实例CPU个数。
	CpuNum int32 `json:"cpu_num"`
}

func (o ReadonlyInstances) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ReadonlyInstances struct{}"
	}

	return strings.Join([]string{"ReadonlyInstances", string(data)}, " ")
}
