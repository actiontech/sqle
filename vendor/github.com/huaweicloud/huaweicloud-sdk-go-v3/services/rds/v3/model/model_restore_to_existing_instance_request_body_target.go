package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreToExistingInstanceRequestBodyTarget 恢复目标对象。
type RestoreToExistingInstanceRequestBodyTarget struct {

	// 恢复目标实例ID。
	InstanceId string `json:"instance_id"`
}

func (o RestoreToExistingInstanceRequestBodyTarget) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreToExistingInstanceRequestBodyTarget struct{}"
	}

	return strings.Join([]string{"RestoreToExistingInstanceRequestBodyTarget", string(data)}, " ")
}
