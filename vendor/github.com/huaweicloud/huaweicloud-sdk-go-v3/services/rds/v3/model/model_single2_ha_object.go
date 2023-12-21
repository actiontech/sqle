package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// Single2HaObject 单机转主备时必填。
type Single2HaObject struct {

	// 实例节点可用区码（AZ）。
	AzCodeNewNode string `json:"az_code_new_node"`

	// Dec用户专属存储ID，每个az配置的专属存储不同，实例节点所在专属存储ID，仅支持DEC用户创建时使用。
	DsspoolId *string `json:"dsspool_id,omitempty"`

	// 仅包周期实例进行单机转主备时可指定，表示是否自动从客户的账户中支付。 - true，为自动支付。 - false，为手动支付，默认该方式。
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`

	AdDomainInfo *AdDomainInfo `json:"ad_domain_info,omitempty"`
}

func (o Single2HaObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Single2HaObject struct{}"
	}

	return strings.Join([]string{"Single2HaObject", string(data)}, " ")
}
