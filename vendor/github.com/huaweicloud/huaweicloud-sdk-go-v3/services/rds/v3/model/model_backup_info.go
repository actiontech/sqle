package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BackupInfo 备份信息。
type BackupInfo struct {

	// 备份ID。
	Id string `json:"id"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 备份名称。
	Name string `json:"name"`

	// 备份描述。
	Description *string `json:"description,omitempty"`

	// 只支持Microsoft SQL Server，局部备份的用户自建数据库名列表，当有此参数时以局部备份为准。
	Databases *[]BackupDatabase `json:"databases,omitempty"`

	// 备份开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	BeginTime string `json:"begin_time"`

	// 备份状态，取值：  - BUILDING: 备份中。 - COMPLETED: 备份完成。 - FAILED：备份失败。 - DELETING：备份删除中。
	Status BackupInfoStatus `json:"status"`

	// 备份类型，取值：  - “auto”: 自动全量备份 - “manual”: 手动全量备份 - “fragment”: 差异全量备份 - “incremental”: 自动增量备份
	Type BackupInfoType `json:"type"`
}

func (o BackupInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupInfo struct{}"
	}

	return strings.Join([]string{"BackupInfo", string(data)}, " ")
}

type BackupInfoStatus struct {
	value string
}

type BackupInfoStatusEnum struct {
	BUILDING  BackupInfoStatus
	COMPLETED BackupInfoStatus
	FAILED    BackupInfoStatus
	DELETING  BackupInfoStatus
}

func GetBackupInfoStatusEnum() BackupInfoStatusEnum {
	return BackupInfoStatusEnum{
		BUILDING: BackupInfoStatus{
			value: "BUILDING",
		},
		COMPLETED: BackupInfoStatus{
			value: "COMPLETED",
		},
		FAILED: BackupInfoStatus{
			value: "FAILED",
		},
		DELETING: BackupInfoStatus{
			value: "DELETING",
		},
	}
}

func (c BackupInfoStatus) Value() string {
	return c.value
}

func (c BackupInfoStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupInfoStatus) UnmarshalJSON(b []byte) error {
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

type BackupInfoType struct {
	value string
}

type BackupInfoTypeEnum struct {
	AUTO        BackupInfoType
	MANUAL      BackupInfoType
	FRAGMENT    BackupInfoType
	INCREMENTAL BackupInfoType
}

func GetBackupInfoTypeEnum() BackupInfoTypeEnum {
	return BackupInfoTypeEnum{
		AUTO: BackupInfoType{
			value: "auto",
		},
		MANUAL: BackupInfoType{
			value: "manual",
		},
		FRAGMENT: BackupInfoType{
			value: "fragment",
		},
		INCREMENTAL: BackupInfoType{
			value: "incremental",
		},
	}
}

func (c BackupInfoType) Value() string {
	return c.value
}

func (c BackupInfoType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupInfoType) UnmarshalJSON(b []byte) error {
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
