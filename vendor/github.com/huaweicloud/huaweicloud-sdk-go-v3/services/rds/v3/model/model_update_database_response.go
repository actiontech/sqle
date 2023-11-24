package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDatabaseResponse Response Object
type UpdateDatabaseResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDatabaseResponse struct{}"
	}

	return strings.Join([]string{"UpdateDatabaseResponse", string(data)}, " ")
}
