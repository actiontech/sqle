package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type PwdResetRequest struct {

	// 数据库密码
	DbUserPwd string `json:"db_user_pwd"`
}

func (o PwdResetRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "PwdResetRequest struct{}"
	}

	return strings.Join([]string{"PwdResetRequest", string(data)}, " ")
}
