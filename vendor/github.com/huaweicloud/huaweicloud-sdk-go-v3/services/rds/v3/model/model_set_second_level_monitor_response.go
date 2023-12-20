package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetSecondLevelMonitorResponse Response Object
type SetSecondLevelMonitorResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetSecondLevelMonitorResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetSecondLevelMonitorResponse struct{}"
	}

	return strings.Join([]string{"SetSecondLevelMonitorResponse", string(data)}, " ")
}
