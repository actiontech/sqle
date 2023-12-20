package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// SlaveInstance 容灾实例信息。
type SlaveInstance struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 区域。
	Region string `json:"region"`

	// 项目ID。
	ProjectId string `json:"project_id"`

	// 项目名称。
	ProjectName string `json:"project_name"`
}

func (o SlaveInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SlaveInstance struct{}"
	}

	return strings.Join([]string{"SlaveInstance", string(data)}, " ")
}
