package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// StartDatabaseProxyResponse Response Object
type StartDatabaseProxyResponse struct {

	// 工作流ID
	WorkflowId     *string `json:"workflow_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o StartDatabaseProxyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartDatabaseProxyResponse struct{}"
	}

	return strings.Join([]string{"StartDatabaseProxyResponse", string(data)}, " ")
}
