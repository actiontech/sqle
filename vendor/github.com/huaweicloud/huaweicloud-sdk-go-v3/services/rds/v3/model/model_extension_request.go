package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ExtensionRequest struct {

	// 数据库名称。
	DatabaseName string `json:"database_name"`

	// 插件名称。
	ExtensionName string `json:"extension_name"`
}

func (o ExtensionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ExtensionRequest struct{}"
	}

	return strings.Join([]string{"ExtensionRequest", string(data)}, " ")
}
