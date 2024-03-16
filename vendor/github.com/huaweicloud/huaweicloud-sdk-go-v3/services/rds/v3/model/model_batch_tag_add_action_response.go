package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchTagAddActionResponse Response Object
type BatchTagAddActionResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o BatchTagAddActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchTagAddActionResponse struct{}"
	}

	return strings.Join([]string{"BatchTagAddActionResponse", string(data)}, " ")
}
