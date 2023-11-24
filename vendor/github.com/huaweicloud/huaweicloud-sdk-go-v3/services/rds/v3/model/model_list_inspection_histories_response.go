package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListInspectionHistoriesResponse Response Object
type ListInspectionHistoriesResponse struct {

	// 总记录数。
	TotalCount *int32 `json:"total_count,omitempty"`

	// 检查报告信息。
	InspectionReports *[]InspectionReports `json:"inspection_reports,omitempty"`
	HttpStatusCode    int                  `json:"-"`
}

func (o ListInspectionHistoriesResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInspectionHistoriesResponse struct{}"
	}

	return strings.Join([]string{"ListInspectionHistoriesResponse", string(data)}, " ")
}
