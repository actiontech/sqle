package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type AvailableZone struct {

	// 可用区编码。
	Code string `json:"code"`

	// 可用区描述。
	Description string `json:"description"`
}

func (o AvailableZone) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AvailableZone struct{}"
	}

	return strings.Join([]string{"AvailableZone", string(data)}, " ")
}
