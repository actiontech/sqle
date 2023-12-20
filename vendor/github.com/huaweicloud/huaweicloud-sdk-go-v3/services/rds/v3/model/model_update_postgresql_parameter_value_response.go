package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdatePostgresqlParameterValueResponse Response Object
type UpdatePostgresqlParameterValueResponse struct {

	// 任务ID。
	JobId *string `json:"job_id,omitempty"`

	// 实例是否需要重启。 - “true”需要重启。 - “false”不需要重启。
	RestartRequired *bool `json:"restart_required,omitempty"`
	HttpStatusCode  int   `json:"-"`
}

func (o UpdatePostgresqlParameterValueResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlParameterValueResponse struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlParameterValueResponse", string(data)}, " ")
}
