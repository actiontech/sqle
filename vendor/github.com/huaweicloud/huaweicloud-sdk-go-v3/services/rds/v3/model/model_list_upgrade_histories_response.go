package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListUpgradeHistoriesResponse Response Object
type ListUpgradeHistoriesResponse struct {

	// 总记录数。
	TotalCount *int32 `json:"total_count,omitempty"`

	// 升级报告信息。
	UpgradeReports *[]UpgradeReports `json:"upgrade_reports,omitempty"`
	HttpStatusCode int               `json:"-"`
}

func (o ListUpgradeHistoriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListUpgradeHistoriesResponse struct{}"
	}

	return strings.Join([]string{"ListUpgradeHistoriesResponse", string(data)}, " ")
}
