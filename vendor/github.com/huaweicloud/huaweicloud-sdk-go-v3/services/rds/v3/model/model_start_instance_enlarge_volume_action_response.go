package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartInstanceEnlargeVolumeActionResponse Response Object
type StartInstanceEnlargeVolumeActionResponse struct {

	// 任务ID。
	JobId          *string `json:"job_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartInstanceEnlargeVolumeActionResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceEnlargeVolumeActionResponse struct{}"
	}

	return strings.Join([]string{"StartInstanceEnlargeVolumeActionResponse", string(data)}, " ")
}
