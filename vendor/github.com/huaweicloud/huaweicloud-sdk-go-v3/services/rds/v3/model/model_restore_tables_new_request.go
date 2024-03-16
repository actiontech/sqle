package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// RestoreTablesNewRequest Request Object
type RestoreTablesNewRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *RestoreTablesNewRequestBody `json:"body,omitempty"`
}

func (o RestoreTablesNewRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreTablesNewRequest struct{}"
	}

	return strings.Join([]string{"RestoreTablesNewRequest", string(data)}, " ")
}
