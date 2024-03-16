package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type DownloadInfoRsp struct {

	// 证书下载地址
	DownloadLink string `json:"download_link"`

	// 证书类型
	Category string `json:"category"`
}

func (o DownloadInfoRsp) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadInfoRsp struct{}"
	}

	return strings.Join([]string{"DownloadInfoRsp", string(data)}, " ")
}
