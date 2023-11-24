package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListApiVersionNewRequest Request Object
type ListApiVersionNewRequest struct {
}

func (o ListApiVersionNewRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListApiVersionNewRequest struct{}"
	}

	return strings.Join([]string{"ListApiVersionNewRequest", string(data)}, " ")
}
