package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SwitchSslResponse Response Object
type SwitchSslResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o SwitchSslResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchSslResponse struct{}"
	}

	return strings.Join([]string{"SwitchSslResponse", string(data)}, " ")
}
