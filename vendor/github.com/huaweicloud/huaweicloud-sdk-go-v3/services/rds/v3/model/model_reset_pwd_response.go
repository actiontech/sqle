package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResetPwdResponse Response Object
type ResetPwdResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ResetPwdResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetPwdResponse struct{}"
	}

	return strings.Join([]string{"ResetPwdResponse", string(data)}, " ")
}
