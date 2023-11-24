package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// GetOffSiteBackupPolicy 备份策略对象，包括备份类型、备份保留天数、目标区域ID和目标project ID。
type GetOffSiteBackupPolicy struct {

	// 指定备份的类型。取值如下：  - auto：自动全量备份。 - incremental：自动增量备份。 - manual：手动备份，仅SQL Server返回该备份类型 。
	BackupType *string `json:"backup_type,omitempty"`

	// 备份文件可以保存的天数。
	KeepDays *int32 `json:"keep_days,omitempty"`

	// 设置跨区域备份策略的目标区域ID。
	DestinationRegion *string `json:"destination_region,omitempty"`

	// 设置跨区域备份策略的目标project ID。
	DestinationProjectId *string `json:"destination_project_id,omitempty"`
}

func (o GetOffSiteBackupPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "GetOffSiteBackupPolicy struct{}"
	}

	return strings.Join([]string{"GetOffSiteBackupPolicy", string(data)}, " ")
}
