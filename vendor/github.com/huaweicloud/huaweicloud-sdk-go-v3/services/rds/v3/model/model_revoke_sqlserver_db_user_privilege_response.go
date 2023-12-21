package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RevokeSqlserverDbUserPrivilegeResponse Response Object
type RevokeSqlserverDbUserPrivilegeResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o RevokeSqlserverDbUserPrivilegeResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeSqlserverDbUserPrivilegeResponse struct{}"
	}

	return strings.Join([]string{"RevokeSqlserverDbUserPrivilegeResponse", string(data)}, " ")
}
