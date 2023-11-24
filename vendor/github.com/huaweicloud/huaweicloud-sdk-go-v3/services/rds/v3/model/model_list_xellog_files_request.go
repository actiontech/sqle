package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListXellogFilesRequest Request Object
type ListXellogFilesRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 索引位置，偏移量。  从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *int32 `json:"offset,omitempty"`

	// 查询记录数。取值范围[1, 100]。
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListXellogFilesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListXellogFilesRequest struct{}"
	}

	return strings.Join([]string{"ListXellogFilesRequest", string(data)}, " ")
}
