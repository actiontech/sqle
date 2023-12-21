package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SetSecurityGroupResponse Response Object
type SetSecurityGroupResponse struct {

	// 任务ID
	WorkflowId     *string `json:"workflowId,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o SetSecurityGroupResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetSecurityGroupResponse struct{}"
	}

	return strings.Join([]string{"SetSecurityGroupResponse", string(data)}, " ")
}
