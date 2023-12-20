package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// BackupForList 备份信息。
type BackupForList struct {

	// 备份ID。
	Id string `json:"id"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 备份名称。
	Name string `json:"name"`

	// 备份的数据库。
	Databases *[]BackupDatabase `json:"databases,omitempty"`

	// 备份开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	BeginTime string `json:"begin_time"`

	// 备份结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	EndTime string `json:"end_time"`

	// 备份状态，取值：  - BUILDING: 备份中。 - COMPLETED: 备份完成。 - FAILED：备份失败。 - DELETING：备份删除中。
	Status BackupForListStatus `json:"status"`

	// 备份类型，取值：  - “auto”: 自动全量备份 - “manual”: 手动全量备份 - “fragment”: 差异全量备份 - “incremental”: 自动增量备份
	Type BackupForListType `json:"type"`

	// 备份大小，单位为KB。
	Size int64 `json:"size"`

	Datastore *BackupDatastore `json:"datastore"`

	// 是否已被DDM实例关联。
	AssociatedWithDdm *bool `json:"associated_with_ddm,omitempty"`
}

func (o BackupForList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BackupForList struct{}"
	}

	return strings.Join([]string{"BackupForList", string(data)}, " ")
}

type BackupForListStatus struct {
	value string
}

type BackupForListStatusEnum struct {
	BUILDING  BackupForListStatus
	COMPLETED BackupForListStatus
	FAILED    BackupForListStatus
	DELETING  BackupForListStatus
}

func GetBackupForListStatusEnum() BackupForListStatusEnum {
	return BackupForListStatusEnum{
		BUILDING: BackupForListStatus{
			value: "BUILDING",
		},
		COMPLETED: BackupForListStatus{
			value: "COMPLETED",
		},
		FAILED: BackupForListStatus{
			value: "FAILED",
		},
		DELETING: BackupForListStatus{
			value: "DELETING",
		},
	}
}

func (c BackupForListStatus) Value() string {
	return c.value
}

func (c BackupForListStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupForListStatus) UnmarshalJSON(b []byte) error {
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

type BackupForListType struct {
	value string
}

type BackupForListTypeEnum struct {
	AUTO        BackupForListType
	MANUAL      BackupForListType
	FRAGMENT    BackupForListType
	INCREMENTAL BackupForListType
}

func GetBackupForListTypeEnum() BackupForListTypeEnum {
	return BackupForListTypeEnum{
		AUTO: BackupForListType{
			value: "auto",
		},
		MANUAL: BackupForListType{
			value: "manual",
		},
		FRAGMENT: BackupForListType{
			value: "fragment",
		},
		INCREMENTAL: BackupForListType{
			value: "incremental",
		},
	}
}

func (c BackupForListType) Value() string {
	return c.value
}

func (c BackupForListType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BackupForListType) UnmarshalJSON(b []byte) error {
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
