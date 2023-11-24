package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type GetBackupDownloadLinkFiles struct {

	// 文件名。
	Name string `json:"name"`

	// 文件大小，单位为KB。
	Size int64 `json:"size"`

	// 文件下载链接。
	DownloadLink string `json:"download_link"`

	// 下载链接过期时间，格式为“yyyy-mm-ddThh:mm:ssZ”。其中，T指某个时间的开始，Z指时区偏移量，例如北京时间偏移显示为+0800。
	LinkExpiredTime string `json:"link_expired_time"`

	// 数据库名。若文件不是数据库备份，则返回空
	DatabaseName string `json:"database_name"`
}

func (o GetBackupDownloadLinkFiles) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetBackupDownloadLinkFiles struct{}"
	}

	return strings.Join([]string{"GetBackupDownloadLinkFiles", string(data)}, " ")
}
