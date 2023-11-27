package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ScaleFlavors struct {

	// 规格码。
	Code *string `json:"code,omitempty"`

	// CPU个数。
	Cpu *string `json:"cpu,omitempty"`

	// 内存大小（单位：GB）。
	Mem *string `json:"mem,omitempty"`
}

func (o ScaleFlavors) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ScaleFlavors struct{}"
	}

	return strings.Join([]string{"ScaleFlavors", string(data)}, " ")
}
