package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetDatabaseUserPrivilegeResponse Response Object
type SetDatabaseUserPrivilegeResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SetDatabaseUserPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetDatabaseUserPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"SetDatabaseUserPrivilegeResponse", string(data)}, " ")
}
