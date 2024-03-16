package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowBackupDownloadLinkResponse Response Object
type ShowBackupDownloadLinkResponse struct {

	// 备份文件信息。
	Files *[]GetBackupDownloadLinkFiles `json:"files,omitempty"`

	// OBS桶名。
	Bucket         *string `json:"bucket,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowBackupDownloadLinkResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowBackupDownloadLinkResponse struct{}"
	}

	return strings.Join([]string{"ShowBackupDownloadLinkResponse", string(data)}, " ")
}
