package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSimplifiedInstancesResponse Response Object
type ListSimplifiedInstancesResponse struct {

	// 实例集合
	Instances      *[]SimplifiedInstanceEntry `json:"instances,omitempty"`
	HttpStatusCode int                        `json:"-"`
}

func (o ListSimplifiedInstancesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSimplifiedInstancesResponse struct{}"
	}

	return strings.Join([]string{"ListSimplifiedInstancesResponse", string(data)}, " ")
}
