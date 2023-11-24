package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListApiVersionResponse Response Object
type ListApiVersionResponse struct {

	// API版本详细信息列表。
	Versions       *[]ApiVersion `json:"versions,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListApiVersionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListApiVersionResponse struct{}"
	}

	return strings.Join([]string{"ListApiVersionResponse", string(data)}, " ")
}
