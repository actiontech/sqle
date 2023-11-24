package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UserForList 数据库用户信息。
type UserForList struct {

	// 数据库用户名称。
	Name string `json:"name"`
}

func (o UserForList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UserForList struct{}"
	}

	return strings.Join([]string{"UserForList", string(data)}, " ")
}
