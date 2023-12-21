package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListApiVersionNewResponse Response Object
type ListApiVersionNewResponse struct {

	// API版本详细信息列表。
	Versions       *[]ApiVersion `json:"versions,omitempty"`
	HttpStatusCode int           `json:"-"`
}

func (o ListApiVersionNewResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListApiVersionNewResponse struct{}"
	}

	return strings.Join([]string{"ListApiVersionNewResponse", string(data)}, " ")
}
