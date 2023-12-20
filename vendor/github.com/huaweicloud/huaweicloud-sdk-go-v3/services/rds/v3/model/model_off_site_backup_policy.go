package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

// OffSiteBackupPolicy 备份策略对象，包括备份类型、备份保留天数、目标区域ID和目标project ID。
type OffSiteBackupPolicy struct {

	// 指定备份的类型。  SQL Server仅支持设置为“all”。  取值如下：  - auto：自动全量备份。 - incremental：自动增量备份。 - all：同时设置所有备份类型。   - MySQL：同时设置自动全量和自动增量备份。   - SQL Server：同时设置自动全量、自动增量备份和手动备份。
	BackupType string `json:"backup_type"`

	// 备份文件可以保存的天数。
	KeepDays int32 `json:"keep_days"`

	// 设置跨区域备份策略的目标区域ID。
	DestinationRegion string `json:"destination_region"`

	// 设置跨区域备份策略的目标project ID。
	DestinationProjectId string `json:"destination_project_id"`
}

func (o OffSiteBackupPolicy) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OffSiteBackupPolicy struct{}"
	}

	return strings.Join([]string{"OffSiteBackupPolicy", string(data)}, " ")
}
