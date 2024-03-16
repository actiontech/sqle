package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ResizeFlavorRequest 变更实例规格时必填。
type ResizeFlavorRequest struct {
	ResizeFlavor *ResizeFlavorObject `json:"resize_flavor"`
}

func (o ResizeFlavorRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeFlavorRequest struct{}"
	}

	return strings.Join([]string{"ResizeFlavorRequest", string(data)}, " ")
}
