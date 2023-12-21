package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// EnlargeVolumeObject 扩容实例磁盘时必填。
type EnlargeVolumeObject struct {

	// 每次扩容最小容量为10GB，实例所选容量大小必须为10的整数倍，取值范围：40GB~4000GB。 - MySQL部分用户支持11GB~10000GB，如果您想开通该功能，请联系客服。 - PostgreSQL部分用户支持40GB~15000GB，如果您想开通该功能，请联系客服。
	Size int32 `json:"size"`

	// 变更包周期实例的规格时可指定，表示是否自动从客户的账户中支付。 - true，为自动支付。 - false，为手动支付，默认该方式。
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`
}

func (o EnlargeVolumeObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EnlargeVolumeObject struct{}"
	}

	return strings.Join([]string{"EnlargeVolumeObject", string(data)}, " ")
}
