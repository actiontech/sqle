package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OffsiteBackupInstance 跨区域备份实例信息。
type OffsiteBackupInstance struct {

	// 实例ID。
	Id string `json:"id"`

	// 创建的实例名称。
	Name *string `json:"name,omitempty"`

	// 源区域。
	SourceRegion *string `json:"source_region,omitempty"`

	// 租户在源区域下的project ID。
	SourceProjectId *string `json:"source_project_id,omitempty"`

	Datastore *ParaGroupDatastore `json:"datastore,omitempty"`

	// 跨区域备份所在区域。
	DestinationRegion *string `json:"destination_region,omitempty"`

	// 租户在目标区域下的project ID。
	DestinationProjectId *string `json:"destination_project_id,omitempty"`

	// 跨区域备份保留天数。
	KeepDays *int64 `json:"keep_days,omitempty"`
}

func (o OffsiteBackupInstance) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OffsiteBackupInstance struct{}"
	}

	return strings.Join([]string{"OffsiteBackupInstance", string(data)}, " ")
}
