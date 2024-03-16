package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainNameRequest Request Object
type ShowDomainNameRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 实例域名类型，当前只支持private。
	DnsType string `json:"dns_type"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`
}

func (o ShowDomainNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainNameRequest struct{}"
	}

	return strings.Join([]string{"ShowDomainNameRequest", string(data)}, " ")
}
