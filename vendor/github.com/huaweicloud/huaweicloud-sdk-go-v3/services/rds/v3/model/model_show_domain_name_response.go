package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowDomainNameResponse Response Object
type ShowDomainNameResponse struct {

	// 实例ID。
	InstanceId *string `json:"instance_id,omitempty"`

	// 实例域名。
	DnsName *string `json:"dns_name,omitempty"`

	// 实例域名类型，当前只支持private。
	DnsType *string `json:"dns_type,omitempty"`

	// 实例内网IPv4地址。
	Ipv4Address *string `json:"ipv4_address,omitempty"`

	// 域名状态
	Status         *string `json:"status,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowDomainNameResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowDomainNameResponse struct{}"
	}

	return strings.Join([]string{"ShowDomainNameResponse", string(data)}, " ")
}
