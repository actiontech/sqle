package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/sdktime"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// ListPostgresqlHbaInfoHistoryRequest Request Object
type ListPostgresqlHbaInfoHistoryRequest struct {

	// 实例id
	InstanceId string `json:"instance_id"`

	// 开始时间,不传默认当天0点（UTC时区）
	StartTime *sdktime.SdkTime `json:"start_time,omitempty"`

	// 结束时间,不传默认当前时间（UTC时区）
	EndTime *sdktime.SdkTime `json:"end_time,omitempty"`
}

func (o ListPostgresqlHbaInfoHistoryRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPostgresqlHbaInfoHistoryRequest struct{}"
	}

	return strings.Join([]string{"ListPostgresqlHbaInfoHistoryRequest", string(data)}, " ")
}
