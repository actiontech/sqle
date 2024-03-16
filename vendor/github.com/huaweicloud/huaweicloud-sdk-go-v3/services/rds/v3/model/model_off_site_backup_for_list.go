package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// OffSiteBackupForList 跨区域备份信息。
type OffSiteBackupForList struct {

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
	Status OffSiteBackupForListStatus `json:"status"`

	// 备份类型，取值：  - “auto”: 自动全量备份 - “incremental”: 自动增量备份
	Type OffSiteBackupForListType `json:"type"`

	// 备份大小，单位为KB。
	Size int64 `json:"size"`

	Datastore *ParaGroupDatastore `json:"datastore"`

	// 是否已被DDM实例关联。
	AssociatedWithDdm *bool `json:"associated_with_ddm,omitempty"`
}

func (o OffSiteBackupForList) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "OffSiteBackupForList struct{}"
	}

	return strings.Join([]string{"OffSiteBackupForList", string(data)}, " ")
}

type OffSiteBackupForListStatus struct {
	value string
}

type OffSiteBackupForListStatusEnum struct {
	BUILDING  OffSiteBackupForListStatus
	COMPLETED OffSiteBackupForListStatus
	FAILED    OffSiteBackupForListStatus
	DELETING  OffSiteBackupForListStatus
}

func GetOffSiteBackupForListStatusEnum() OffSiteBackupForListStatusEnum {
	return OffSiteBackupForListStatusEnum{
		BUILDING: OffSiteBackupForListStatus{
			value: "BUILDING",
		},
		COMPLETED: OffSiteBackupForListStatus{
			value: "COMPLETED",
		},
		FAILED: OffSiteBackupForListStatus{
			value: "FAILED",
		},
		DELETING: OffSiteBackupForListStatus{
			value: "DELETING",
		},
	}
}

func (c OffSiteBackupForListStatus) Value() string {
	return c.value
}

func (c OffSiteBackupForListStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OffSiteBackupForListStatus) UnmarshalJSON(b []byte) error {
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

type OffSiteBackupForListType struct {
	value string
}

type OffSiteBackupForListTypeEnum struct {
	AUTO        OffSiteBackupForListType
	INCREMENTAL OffSiteBackupForListType
}

func GetOffSiteBackupForListTypeEnum() OffSiteBackupForListTypeEnum {
	return OffSiteBackupForListTypeEnum{
		AUTO: OffSiteBackupForListType{
			value: "auto",
		},
		INCREMENTAL: OffSiteBackupForListType{
			value: "incremental",
		},
	}
}

func (c OffSiteBackupForListType) Value() string {
	return c.value
}

func (c OffSiteBackupForListType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *OffSiteBackupForListType) UnmarshalJSON(b []byte) error {
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
