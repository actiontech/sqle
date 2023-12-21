package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInstancesSupportFastRestoreRequest Request Object
type ListInstancesSupportFastRestoreRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	Body *ListInstancesSupportFastRestoreRequestBody `json:"body,omitempty"`
}

func (o ListInstancesSupportFastRestoreRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesSupportFastRestoreRequest struct{}"
	}

	return strings.Join([]string{"ListInstancesSupportFastRestoreRequest", string(data)}, " ")
}
