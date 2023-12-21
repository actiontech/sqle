package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// BatchAddMsdtcsResponse Response Object
type BatchAddMsdtcsResponse struct {

	// 任务流id
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o BatchAddMsdtcsResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchAddMsdtcsResponse struct{}"
	}

	return strings.Join([]string{"BatchAddMsdtcsResponse", string(data)}, " ")
}
