package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchTagDelActionResponse Response Object
type BatchTagDelActionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchTagDelActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchTagDelActionResponse struct{}"
	}

	return strings.Join([]string{"BatchTagDelActionResponse", string(data)}, " ")
}
