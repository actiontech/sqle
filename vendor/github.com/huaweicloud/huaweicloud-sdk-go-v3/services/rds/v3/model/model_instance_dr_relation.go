package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// InstanceDrRelation 容灾实例信息。
type InstanceDrRelation struct {

	// 当前区域实例ID。
	InstanceId *string `json:"instance_id,omitempty"`

	MasterInstance *MasterInstance `json:"master_instance,omitempty"`

	// 容灾实例信息列表。
	SlaveInstances *[]SlaveInstance `json:"slave_instances,omitempty"`
}

func (o InstanceDrRelation) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceDrRelation struct{}"
	}

	return strings.Join([]string{"InstanceDrRelation", string(data)}, " ")
}
