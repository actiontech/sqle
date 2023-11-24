package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AllowDbPrivilegeResponse Response Object
type AllowDbPrivilegeResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o AllowDbPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AllowDbPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"AllowDbPrivilegeResponse", string(data)}, " ")
}
