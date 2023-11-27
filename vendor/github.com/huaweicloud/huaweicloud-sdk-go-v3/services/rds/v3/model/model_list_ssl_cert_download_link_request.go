package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListSslCertDownloadLinkRequest Request Object
type ListSslCertDownloadLinkRequest struct {
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ListSslCertDownloadLinkRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSslCertDownloadLinkRequest struct{}"
	}

	return strings.Join([]string{"ListSslCertDownloadLinkRequest", string(data)}, " ")
}
