package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type Single2Ha struct {
	SingleToHa *Single2HaObject `json:"single_to_ha"`
}

func (o Single2Ha) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Single2Ha struct{}"
	}

	return strings.Join([]string{"Single2Ha", string(data)}, " ")
}
