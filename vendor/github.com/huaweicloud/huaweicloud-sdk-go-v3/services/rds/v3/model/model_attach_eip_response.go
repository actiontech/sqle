package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// AttachEipResponse Response Object
type AttachEipResponse struct {
	HttpStatusCode int `json:"-"`
}

func (o AttachEipResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AttachEipResponse struct{}"
	}

	return strings.Join([]string{"AttachEipResponse", string(data)}, " ")
}
