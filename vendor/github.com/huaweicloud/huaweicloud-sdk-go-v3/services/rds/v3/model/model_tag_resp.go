package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type TagResp struct {

	// 标签的key
	Key string `json:"key"`

	// 标签value的集合
	Values []string `json:"values"`
}

func (o TagResp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TagResp struct{}"
	}

	return strings.Join([]string{"TagResp", string(data)}, " ")
}
