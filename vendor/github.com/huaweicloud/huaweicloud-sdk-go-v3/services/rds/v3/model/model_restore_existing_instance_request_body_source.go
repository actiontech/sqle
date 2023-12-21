package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// RestoreExistingInstanceRequestBodySource 恢复数据源对象。
type RestoreExistingInstanceRequestBodySource struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 表示恢复方式，枚举值：  - “backup”，表示使用备份文件恢复，按照此方式恢复时，“type”字段为非必选，“backup_id”必选。 - “timestamp”，表示按时间点恢复，按照此方式恢复时，“type”字段必选，“restore_time”必选。
	Type *RestoreExistingInstanceRequestBodySourceType `json:"type,omitempty"`

	// 用于恢复的备份ID。当使用备份文件恢复时需要指定该参数。
	BackupId *string `json:"backup_id,omitempty"`

	// 恢复数据的时间点，格式为UNIX时间戳，单位是毫秒，时区为UTC。
	RestoreTime *int64 `json:"restore_time,omitempty"`

	// 仅适用于SQL Server引擎，当有此参数时表示支持局部恢复和重命名恢复，恢复数据以局部恢复为主。不填写该字段时，默认恢复全部数据库。 - 新数据库名称不可与源实例或目标实例数据库名称重名，新数据库名称为空，默认按照原数据库名进行恢复。 - 新数据库名不能包含rdsadmin、master、msdb、tempdb、model或resource字段（不区分大小写）。 - 数据库名称长度在1~64个字符之间，包含字母、数字、下划线或中划线，不能包含其他特殊字符。
	DatabaseName map[string]string `json:"database_name,omitempty"`

	// 该字段仅适用于SQL Server引擎。是否恢复所有数据库，不填写该字段默认为false，不会恢复所有数据库到目标实例。 - 须知： 如果您想恢复所有数据库到已有实例，必须设置restore_all_database为true。
	RestoreAllDatabase *bool `json:"restore_all_database,omitempty"`
}

func (o RestoreExistingInstanceRequestBodySource) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RestoreExistingInstanceRequestBodySource struct{}"
	}

	return strings.Join([]string{"RestoreExistingInstanceRequestBodySource", string(data)}, " ")
}

type RestoreExistingInstanceRequestBodySourceType struct {
	value string
}

type RestoreExistingInstanceRequestBodySourceTypeEnum struct {
	BACKUP    RestoreExistingInstanceRequestBodySourceType
	TIMESTAMP RestoreExistingInstanceRequestBodySourceType
}

func GetRestoreExistingInstanceRequestBodySourceTypeEnum() RestoreExistingInstanceRequestBodySourceTypeEnum {
	return RestoreExistingInstanceRequestBodySourceTypeEnum{
		BACKUP: RestoreExistingInstanceRequestBodySourceType{
			value: "backup",
		},
		TIMESTAMP: RestoreExistingInstanceRequestBodySourceType{
			value: "timestamp",
		},
	}
}

func (c RestoreExistingInstanceRequestBodySourceType) Value() string {
	return c.value
}

func (c RestoreExistingInstanceRequestBodySourceType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RestoreExistingInstanceRequestBodySourceType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
