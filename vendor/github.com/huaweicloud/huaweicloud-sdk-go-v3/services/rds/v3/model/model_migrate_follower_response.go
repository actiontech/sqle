package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MigrateFollowerResponse Response Object
type MigrateFollowerResponse struct {

	// 任务ID
	WorkflowId     *string `json:"workflowId,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o MigrateFollowerResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MigrateFollowerResponse struct{}"
	}

	return strings.Join([]string{"MigrateFollowerResponse", string(data)}, " ")
}
