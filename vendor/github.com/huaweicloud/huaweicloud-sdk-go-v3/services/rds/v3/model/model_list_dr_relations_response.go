package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDrRelationsResponse Response Object
type ListDrRelationsResponse struct {
	InstanceDrRelations *[]InstanceDrRelation `json:"instance_dr_relations,omitempty"`
	HttpStatusCode      int                   `json:"-"`
}

func (o ListDrRelationsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDrRelationsResponse struct{}"
	}

	return strings.Join([]string{"ListDrRelationsResponse", string(data)}, " ")
}
