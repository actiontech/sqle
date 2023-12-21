package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSslCertDownloadLinkResponse Response Object
type ListSslCertDownloadLinkResponse struct {
	CertInfoList   *[]DownloadInfoRsp `json:"cert_info_list,omitempty"`
	HttpStatusCode int                `json:"-"`
}

func (o ListSslCertDownloadLinkResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSslCertDownloadLinkResponse struct{}"
	}

	return strings.Join([]string{"ListSslCertDownloadLinkResponse", string(data)}, " ")
}
