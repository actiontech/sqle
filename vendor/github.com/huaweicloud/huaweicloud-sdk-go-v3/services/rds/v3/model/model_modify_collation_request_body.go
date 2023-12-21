package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ModifyCollationRequestBody struct {

	// 字符集。 取值范围：根据查询SQL Server可用字符集查询可设置的字符集。
	Collation string `json:"collation"`
}

func (o ModifyCollationRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifyCollationRequestBody struct{}"
	}

	return strings.Join([]string{"ModifyCollationRequestBody", string(data)}, " ")
}
