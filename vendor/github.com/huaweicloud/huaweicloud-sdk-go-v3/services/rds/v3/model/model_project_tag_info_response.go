package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ProjectTagInfoResponse 项目标签信息。
type ProjectTagInfoResponse struct {

	// 标签键。
	Key string `json:"key"`

	// 标签值列表。
	Values []string `json:"values"`
}

func (o ProjectTagInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProjectTagInfoResponse struct{}"
	}

	return strings.Join([]string{"ProjectTagInfoResponse", string(data)}, " ")
}
