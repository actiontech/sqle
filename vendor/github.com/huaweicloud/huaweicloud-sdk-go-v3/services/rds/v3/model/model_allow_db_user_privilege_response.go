package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowDbUserPrivilegeResponse Response Object
type AllowDbUserPrivilegeResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AllowDbUserPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowDbUserPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"AllowDbUserPrivilegeResponse", string(data)}, " ")
}
