package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type InstanceRestartRequsetBody struct {

	// 空值
	Restart *interface{} `json:"restart"`
}

func (o InstanceRestartRequsetBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceRestartRequsetBody struct{}"
	}

	return strings.Join([]string{"InstanceRestartRequsetBody", string(data)}, " ")
}
