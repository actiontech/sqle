package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Computes struct {

	// 群组类型。  - X86：X86架构。 - ARM：ARM架构。
	GroupType *string `json:"group_type,omitempty"`

	// 计算规格信息。
	ComputeFlavors *[]ScaleFlavors `json:"compute_flavors,omitempty"`
}

func (o Computes) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Computes struct{}"
	}

	return strings.Join([]string{"Computes", string(data)}, " ")
}
