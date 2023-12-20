package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RevokeRequestBodyUsers struct {

	// 数据库用户名称。
	Name string `json:"name"`
}

func (o RevokeRequestBodyUsers) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RevokeRequestBodyUsers struct{}"
	}

	return strings.Join([]string{"RevokeRequestBodyUsers", string(data)}, " ")
}
