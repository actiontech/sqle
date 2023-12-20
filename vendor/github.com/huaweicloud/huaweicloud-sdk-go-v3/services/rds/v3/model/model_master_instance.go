package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MasterInstance 主实例信息。
type MasterInstance struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 区域。
	Region string `json:"region"`

	// 项目ID。
	ProjectId string `json:"project_id"`

	// 项目名称。
	ProjectName string `json:"project_name"`
}

func (o MasterInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MasterInstance struct{}"
	}

	return strings.Join([]string{"MasterInstance", string(data)}, " ")
}
