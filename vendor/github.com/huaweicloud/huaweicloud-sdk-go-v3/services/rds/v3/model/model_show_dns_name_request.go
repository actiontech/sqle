package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDnsNameRequest Request Object
type ShowDnsNameRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例域名类型，当前只支持\"private\"。
	DnsType string `json:"dns_type"`
}

func (o ShowDnsNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDnsNameRequest struct{}"
	}

	return strings.Join([]string{"ShowDnsNameRequest", string(data)}, " ")
}
