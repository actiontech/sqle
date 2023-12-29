package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartFailoverResponse Response Object
type StartFailoverResponse struct {

	// 实例Id
	InstanceId *string `json:"instanceId,omitempty"`

	// 节点Id
	NodeId *string `json:"nodeId,omitempty"`

	// 任务Id
	WorkflowId     *string `json:"workflowId,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartFailoverResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartFailoverResponse struct{}"
	}

	return strings.Join([]string{"StartFailoverResponse", string(data)}, " ")
}
