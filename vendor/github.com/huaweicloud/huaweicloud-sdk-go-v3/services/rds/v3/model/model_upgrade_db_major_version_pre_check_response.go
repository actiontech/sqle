package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// UpgradeDbMajorVersionPreCheckResponse Response Object
type UpgradeDbMajorVersionPreCheckResponse struct {

	// 检查报告ID。
	ReportId       *string `json:"report_id,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o UpgradeDbMajorVersionPreCheckResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbMajorVersionPreCheckResponse struct{}"
	}

	return strings.Join([]string{"UpgradeDbMajorVersionPreCheckResponse", string(data)}, " ")
}
