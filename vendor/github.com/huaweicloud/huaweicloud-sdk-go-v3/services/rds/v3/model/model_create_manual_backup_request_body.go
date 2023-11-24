package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type CreateManualBackupRequestBody struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 备份名称，4~64个字符，必须以英文字母开头，区分大小写，可以包含英文字母、数字、中划线或者下划线，不能包含其他特殊字符。
	Name string `json:"name"`

	// 备份描述，不能包含>!<\"&'=特殊字符，不大于256个字符。
	Description *string `json:"description,omitempty"`

	// 只支持Microsoft SQL Server，局部备份的用户自建数据库名列表，当有此参数时以局部备份为准。
	Databases *[]BackupDatabase `json:"databases,omitempty"`
}

func (o CreateManualBackupRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateManualBackupRequestBody struct{}"
	}

	return strings.Join([]string{"CreateManualBackupRequestBody", string(data)}, " ")
}
