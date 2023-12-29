package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDnsNameResponse Response Object
type ShowDnsNameResponse struct {

	// 实例ID。
	InstanceId *string `json:"instance_id,omitempty"`

	// 实例域名。
	DnsName *string `json:"dns_name,omitempty"`

	// 实例域名类型，当前只支持private。
	DnsType *string `json:"dns_type,omitempty"`

	// 实例内网IPv6地址。
	Ipv6Address *string `json:"ipv6_address,omitempty"`

	// 域名状态。
	Status         *string `json:"status,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowDnsNameResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDnsNameResponse struct{}"
	}

	return strings.Join([]string{"ShowDnsNameResponse", string(data)}, " ")
}
