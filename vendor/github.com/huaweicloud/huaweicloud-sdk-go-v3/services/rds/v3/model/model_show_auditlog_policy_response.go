package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ShowAuditlogPolicyResponse Response Object
type ShowAuditlogPolicyResponse struct {

	// 审计日志保存天数，取值范围0~732。0表示关闭审计日志策略。
	KeepDays       *int32 `json:"keep_days,omitempty"`
	HttpStatusCode int    `json:"-"`
}

func (o ShowAuditlogPolicyResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowAuditlogPolicyResponse struct{}"
	}

	return strings.Join([]string{"ShowAuditlogPolicyResponse", string(data)}, " ")
}
