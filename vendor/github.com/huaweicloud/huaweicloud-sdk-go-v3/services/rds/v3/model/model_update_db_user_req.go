package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateDbUserReq struct {

	// 数据库用户备注。
	Comment *string `json:"comment,omitempty"`
}

func (o UpdateDbUserReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDbUserReq struct{}"
	}

	return strings.Join([]string{"UpdateDbUserReq", string(data)}, " ")
}
