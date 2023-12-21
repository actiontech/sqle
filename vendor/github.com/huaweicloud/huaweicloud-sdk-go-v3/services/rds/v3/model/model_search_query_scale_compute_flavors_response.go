package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SearchQueryScaleComputeFlavorsResponse Response Object
type SearchQueryScaleComputeFlavorsResponse struct {

	// 计算规格列表对象。
	ComputeFlavorGroups *[]Computes `json:"compute_flavor_groups,omitempty"`
	HttpStatusCode      int         `json:"-"`
}

func (o SearchQueryScaleComputeFlavorsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SearchQueryScaleComputeFlavorsResponse struct{}"
	}

	return strings.Join([]string{"SearchQueryScaleComputeFlavorsResponse", string(data)}, " ")
}
