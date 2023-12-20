package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifyCollationResponse Response Object
type ModifyCollationResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ModifyCollationResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyCollationResponse struct{}"
	}

	return strings.Join([]string{"ModifyCollationResponse", string(data)}, " ")
}
