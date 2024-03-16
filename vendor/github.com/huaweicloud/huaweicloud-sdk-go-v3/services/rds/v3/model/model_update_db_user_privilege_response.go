package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDbUserPrivilegeResponse Response Object
type UpdateDbUserPrivilegeResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateDbUserPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDbUserPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"UpdateDbUserPrivilegeResponse", string(data)}, " ")
}
