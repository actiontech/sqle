package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type BatchDeleteBackupRequestBody struct {

	// 需要删除的手动备份ID列表。列表大小范围：[1-50]
	BackupIds []string `json:"backup_ids"`
}

func (o BatchDeleteBackupRequestBody) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchDeleteBackupRequestBody struct{}"
	}

	return strings.Join([]string{"BatchDeleteBackupRequestBody", string(data)}, " ")
}
