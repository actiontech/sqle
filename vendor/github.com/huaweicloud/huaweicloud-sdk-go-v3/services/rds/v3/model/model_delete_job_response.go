package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeleteJobResponse Response Object
type DeleteJobResponse struct {
	Body           *string `json:"body,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o DeleteJobResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteJobResponse struct{}"
	}

	return strings.Join([]string{"DeleteJobResponse", string(data)}, " ")
}
