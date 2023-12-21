package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPredefinedTagResponse Response Object
type ListPredefinedTagResponse struct {

	// 标签集合
	Tags           *[]TagResp `json:"tags,omitempty"`
	HttpStatusCode int        `json:"-"`
}

func (o ListPredefinedTagResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPredefinedTagResponse struct{}"
	}

	return strings.Join([]string{"ListPredefinedTagResponse", string(data)}, " ")
}
