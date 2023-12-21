package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateTdeStatusResponse Response Object
type UpdateTdeStatusResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o UpdateTdeStatusResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTdeStatusResponse struct{}"
	}

	return strings.Join([]string{"UpdateTdeStatusResponse", string(data)}, " ")
}
