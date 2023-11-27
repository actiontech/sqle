package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ChangeOpsWindowResponse Response Object
type ChangeOpsWindowResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o ChangeOpsWindowResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeOpsWindowResponse struct{}"
	}

	return strings.Join([]string{"ChangeOpsWindowResponse", string(data)}, " ")
}
