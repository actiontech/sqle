package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ModifiyInstanceNameRequest 修改实例名称必填。
type ModifiyInstanceNameRequest struct {

	// 实例名称。  用于表示实例的名称，同一租户下，同类型的实例名可重名。取值范围如下：  - MySQL数据库支持的字符长度是4~64个字符，必须以字母开头，区分大小写，可以包含字母、数字、中文字符、中划线或者下划线，不能包含其他的特殊字符。 - PostgreSQL和SQL Server数据库支持的字符长度是4~64个字符，必须以字母开头，区分大小写，可以包含字母、数字、中划线或者下划线，不能包含其他的特殊字符。
	Name string `json:"name"`
}

func (o ModifiyInstanceNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ModifiyInstanceNameRequest struct{}"
	}

	return strings.Join([]string{"ModifiyInstanceNameRequest", string(data)}, " ")
}
