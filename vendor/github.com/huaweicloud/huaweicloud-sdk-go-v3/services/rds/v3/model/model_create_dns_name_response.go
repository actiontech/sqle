package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// CreateDnsNameResponse Response Object
type CreateDnsNameResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateDnsNameResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDnsNameResponse struct{}"
	}

	return strings.Join([]string{"CreateDnsNameResponse", string(data)}, " ")
}
