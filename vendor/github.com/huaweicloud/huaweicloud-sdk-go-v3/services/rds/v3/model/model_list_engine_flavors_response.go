package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListEngineFlavorsResponse Response Object
type ListEngineFlavorsResponse struct {

	// 可用的规格列表信息
	OptionalFlavors *[]EngineFlavorData `json:"optional_flavors,omitempty"`

	// 可用的规格总数
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListEngineFlavorsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListEngineFlavorsResponse struct{}"
	}

	return strings.Join([]string{"ListEngineFlavorsResponse", string(data)}, " ")
}
