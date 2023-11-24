package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DeleteBackupResult struct {

	// 备份id
	BackupId *string `json:"backup_id,omitempty"`

	// 提交任务异常时返回的错误编码。0表示成功。
	ErrorCode *string `json:"error_code,omitempty"`

	// 提交任务异常时返回的错误描述信息。
	ErrorMsg *string `json:"error_msg,omitempty"`
}

func (o DeleteBackupResult) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteBackupResult struct{}"
	}

	return strings.Join([]string{"DeleteBackupResult", string(data)}, " ")
}
