package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlExtensionRequest Request Object
type ListPostgresqlExtensionRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 数据库名称。
	DatabaseName string `json:"database_name"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *int32 `json:"offset,omitempty"`

	// 查询记录数。默认为100，不能为负数，最小值为1，最大值为100。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListPostgresqlExtensionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlExtensionRequest struct{}"
	}

	return strings.Join([]string{"ListPostgresqlExtensionRequest", string(data)}, " ")
}
