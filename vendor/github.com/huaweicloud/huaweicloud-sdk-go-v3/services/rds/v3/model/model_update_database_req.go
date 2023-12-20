package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type UpdateDatabaseReq struct {

	// 数据库名称。
	Name string `json:"name"`

	// 数据库备注。
	Comment *string `json:"comment,omitempty"`
}

func (o UpdateDatabaseReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDatabaseReq struct{}"
	}

	return strings.Join([]string{"UpdateDatabaseReq", string(data)}, " ")
}
