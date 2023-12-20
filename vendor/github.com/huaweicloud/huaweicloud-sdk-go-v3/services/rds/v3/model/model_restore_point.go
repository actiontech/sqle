package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"strings"
)

type RestorePoint struct {

	// 源实例ID。
	InstanceId string `json:"instance_id"`

	// 表示恢复方式，枚举值：  - “backup”，表示使用备份文件恢复，按照此方式恢复时，“type”字段为非必选，“backup_id”必选。 - “timestamp”，表示按时间点恢复，按照此方式恢复时，“type”字段必选，“restore_time”必选。
	Type string `json:"type"`

	// 用于恢复的备份ID。当使用备份文件恢复时需要指定该参数。
	BackupId *string `json:"backup_id,omitempty"`

	// 恢复数据的时间点，格式为UNIX时间戳，单位是毫秒，时区为UTC。
	RestoreTime *int64 `json:"restore_time,omitempty"`

	// 仅适用于SQL Server引擎，当有此参数时表示支持局部恢复和重命名恢复，恢复数据以局部恢复为主。  - 新数据库名称不可与源实例数据库名称重名，新数据库名称为空，默认按照原数据库名进行恢复。   注意：   不填写该字段时，默认恢复全部数据库。    示例：”database_name”:{“原库名”:”新库名”}  - 新数据库名不能包含rdsadmin、master、msdb、tempdb、model或resource字段（不区分大小写）。 - 数据库名称长度在1~64个字符之间，包含字母、数字、下划线或中划线，不能包含其他特殊字符。
	DatabaseName map[string]string `json:"database_name,omitempty"`
}

func (o RestorePoint) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestorePoint struct{}"
	}

	return strings.Join([]string{"RestorePoint", string(data)}, " ")
}
