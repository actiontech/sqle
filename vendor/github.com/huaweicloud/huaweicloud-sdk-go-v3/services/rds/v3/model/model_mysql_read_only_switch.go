package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// MysqlReadOnlySwitch 设置实例只读参数。
type MysqlReadOnlySwitch struct {

	// 是否设置为只读权限 - true，表示设置为只读权限 - false，表示解除已设置的只读权限
	Readonly bool `json:"readonly"`
}

func (o MysqlReadOnlySwitch) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlReadOnlySwitch struct{}"
	}

	return strings.Join([]string{"MysqlReadOnlySwitch", string(data)}, " ")
}
