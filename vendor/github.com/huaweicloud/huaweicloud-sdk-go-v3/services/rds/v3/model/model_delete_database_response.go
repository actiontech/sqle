package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteDatabaseResponse Response Object
type DeleteDatabaseResponse struct {

	// 操作结果。
	Resp           *string `json:"resp,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteDatabaseResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteDatabaseResponse struct{}"
	}

	return strings.Join([]string{"DeleteDatabaseResponse", string(data)}, " ")
}
