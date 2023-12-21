package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// DeletePostgresqlExtensionRequest Request Object
type DeletePostgresqlExtensionRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *ExtensionRequest `json:"body,omitempty"`
}

func (o DeletePostgresqlExtensionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeletePostgresqlExtensionRequest struct{}"
	}

	return strings.Join([]string{"DeletePostgresqlExtensionRequest", string(data)}, " ")
}
