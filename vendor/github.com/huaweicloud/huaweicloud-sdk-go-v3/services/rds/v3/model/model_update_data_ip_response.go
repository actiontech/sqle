package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpdateDataIpResponse Response Object
type UpdateDataIpResponse struct {

	// 任务ID
	WorkflowId     *string `json:"workflowId,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpdateDataIpResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDataIpResponse struct{}"
	}

	return strings.Join([]string{"UpdateDataIpResponse", string(data)}, " ")
}
