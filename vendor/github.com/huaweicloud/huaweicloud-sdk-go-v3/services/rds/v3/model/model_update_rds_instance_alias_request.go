package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateRdsInstanceAliasRequest struct {

	// 长度可在0~64个字符之间，由字母、数字、汉字、英文句号、下划线、中划线组成。此参数为空时可以清空原有备注。
	Alias *string `json:"alias,omitempty"`
}

func (o UpdateRdsInstanceAliasRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateRdsInstanceAliasRequest struct{}"
	}

	return strings.Join([]string{"UpdateRdsInstanceAliasRequest", string(data)}, " ")
}
