package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// GetTaskDetailListRspJobsInstance 执行任务的实例信息。
type GetTaskDetailListRspJobsInstance struct {

	// 实例ID。
	Id string `json:"id"`

	// 实例名称。
	Name string `json:"name"`
}

func (o GetTaskDetailListRspJobsInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetTaskDetailListRspJobsInstance struct{}"
	}

	return strings.Join([]string{"GetTaskDetailListRspJobsInstance", string(data)}, " ")
}
