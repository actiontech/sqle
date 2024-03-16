package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListDrRelationsRequest Request Object
type ListDrRelationsRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ListDrRelationsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDrRelationsRequest struct{}"
	}

	return strings.Join([]string{"ListDrRelationsRequest", string(data)}, " ")
}
