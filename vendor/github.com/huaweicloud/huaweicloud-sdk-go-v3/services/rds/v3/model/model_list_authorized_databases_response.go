package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListAuthorizedDatabasesResponse Response Object
type ListAuthorizedDatabasesResponse struct {

	// 数据库及相关权限。
	Databases *[]DatabaseWithPrivilege `json:"databases,omitempty"`

	// 总数。
	TotalCount     *int32 `json:"total_count,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ListAuthorizedDatabasesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAuthorizedDatabasesResponse struct{}"
	}

	return strings.Join([]string{"ListAuthorizedDatabasesResponse", string(data)}, " ")
}
