package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBackupDownloadLinkRequest Request Object
type ShowBackupDownloadLinkRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 备份ID。
	BackupId string `json:"backup_id"`
}

func (o ShowBackupDownloadLinkRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBackupDownloadLinkRequest struct{}"
	}

	return strings.Join([]string{"ShowBackupDownloadLinkRequest", string(data)}, " ")
}
