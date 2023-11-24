package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAuthorizedSqlserverDbUsersResponse Response Object
type ListAuthorizedSqlserverDbUsersResponse struct {

	// 用户及相关权限。
	Users *[]UserWithPrivilege `json:"users,omitempty"`

	// 总数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListAuthorizedSqlserverDbUsersResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAuthorizedSqlserverDbUsersResponse struct{}"
	}

	return strings.Join([]string{"ListAuthorizedSqlserverDbUsersResponse", string(data)}, " ")
}
