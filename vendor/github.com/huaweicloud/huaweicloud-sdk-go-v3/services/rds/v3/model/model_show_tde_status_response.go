package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowTdeStatusResponse Response Object
type ShowTdeStatusResponse struct {

	// 实例ID
	InstanceId *string `json:"instance_id,omitempty"`

	// TDE状态
	TdeStatus      *string `json:"tde_status,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowTdeStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowTdeStatusResponse struct{}"
	}

	return strings.Join([]string{"ShowTdeStatusResponse", string(data)}, " ")
}
