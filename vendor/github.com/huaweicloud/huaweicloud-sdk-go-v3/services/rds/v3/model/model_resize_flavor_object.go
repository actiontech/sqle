package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type ResizeFlavorObject struct {

	// 资源规格编码。例如：rds.mysql.m1.xlarge。其中，rds代表RDS产品，mysql代表数据库引擎，m1.xlarge代表性能规格，为高内存类型。带\"rr\"的表示只读实例规格，反之表示单实例和HA实例规格。
	SpecCode string `json:"spec_code"`

	// 变更包周期实例的规格时可指定，表示是否自动从客户的账户中支付。 - true，为自动支付。 - false，为手动支付，默认该方式。
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`
}

func (o ResizeFlavorObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResizeFlavorObject struct{}"
	}

	return strings.Join([]string{"ResizeFlavorObject", string(data)}, " ")
}
