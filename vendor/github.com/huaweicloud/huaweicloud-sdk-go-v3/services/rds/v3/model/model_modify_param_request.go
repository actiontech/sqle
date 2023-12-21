package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyParamRequest struct {

	// 参数值。
	Value string `json:"value"`
}

func (o ModifyParamRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyParamRequest struct{}"
	}

	return strings.Join([]string{"ModifyParamRequest", string(data)}, " ")
}
