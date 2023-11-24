package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartInstanceRestartActionResponse Response Object
type StartInstanceRestartActionResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartInstanceRestartActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceRestartActionResponse struct{}"
	}

	return strings.Join([]string{"StartInstanceRestartActionResponse", string(data)}, " ")
}
