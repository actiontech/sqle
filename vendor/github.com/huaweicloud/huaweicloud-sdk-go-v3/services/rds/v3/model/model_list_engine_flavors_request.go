package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListEngineFlavorsRequest Request Object
type ListEngineFlavorsRequest struct {

	// 实例ID
	InstanceId string `json:"instance_id"`

	// 可用区，多个用\",\"分割，如cn-southwest-244a,cn-southwest-244b
	AvailabilityZoneIds string `json:"availability_zone_ids"`

	// 模式，包括如下类型： ha：主备实例。 replica：只读实例。 single：单实例。
	HaMode string `json:"ha_mode"`

	// 性能规格,如rds.dec.pg.s1.medium，模糊匹配该规格类型
	SpecCodeLike *string `json:"spec_code_like,omitempty"`

	// 规格类型，包括如下类型：simple、dec
	FlavorCategoryType *string `json:"flavor_category_type,omitempty"`

	// 是否显示高可用只读类型
	IsRhaFlavor *bool `json:"is_rha_flavor,omitempty"`

	// 索引位置，偏移量。 从第一条数据偏移offset条数据后开始查询，默认为0。 取值必须为数字，且不能为负数。
	Offset *int32 `json:"offset,omitempty"`

	// 查询个数上限值。 取值范围：1~100。 不传该参数时，默认查询前100条信息。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListEngineFlavorsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEngineFlavorsRequest struct{}"
	}

	return strings.Join([]string{"ListEngineFlavorsRequest", string(data)}, " ")
}
