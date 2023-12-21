package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateInstanceConfigurationAsyncResponse Response Object
type UpdateInstanceConfigurationAsyncResponse struct {

	// 任务流id
	JobId *string `json:"job_id,omitempty"`

	// 实例是否需要重启。 - “true”需要重启。 - “false”不需要重启。
	RestartRequired *bool `json:"restart_required,omitempty"`
	HttpStatusCode  int   `json:"-"`
}

func (o UpdateInstanceConfigurationAsyncResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateInstanceConfigurationAsyncResponse struct{}"
	}

	return strings.Join([]string{"UpdateInstanceConfigurationAsyncResponse", string(data)}, " ")
}
