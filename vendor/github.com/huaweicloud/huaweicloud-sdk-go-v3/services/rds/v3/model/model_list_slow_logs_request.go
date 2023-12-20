package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSlowLogsRequest Request Object
type ListSlowLogsRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 开始日期，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartDate string `json:"start_date"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。只能查询当前时间前一个月内的慢日志。
	EndDate string `json:"end_date"`

	// 页数偏移量，如1、2、3、4等，不填时默认为1。
	Offset *int32 `json:"offset,omitempty"`

	// 每页多少条记录，取值范围是1~100，不填时默认为10。
	Limit *int32 `json:"limit,omitempty"`

	// 语句类型，取空值，表示查询所有语句类型。
	Type *ListSlowLogsRequestType `json:"type,omitempty"`
}

func (o ListSlowLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSlowLogsRequest struct{}"
	}

	return strings.Join([]string{"ListSlowLogsRequest", string(data)}, " ")
}

type ListSlowLogsRequestType struct {
	value string
}

type ListSlowLogsRequestTypeEnum struct {
	INSERT ListSlowLogsRequestType
	UPDATE ListSlowLogsRequestType
	SELECT ListSlowLogsRequestType
	DELETE ListSlowLogsRequestType
	CREATE ListSlowLogsRequestType
}

func GetListSlowLogsRequestTypeEnum() ListSlowLogsRequestTypeEnum {
	return ListSlowLogsRequestTypeEnum{
		INSERT: ListSlowLogsRequestType{
			value: "INSERT",
		},
		UPDATE: ListSlowLogsRequestType{
			value: "UPDATE",
		},
		SELECT: ListSlowLogsRequestType{
			value: "SELECT",
		},
		DELETE: ListSlowLogsRequestType{
			value: "DELETE",
		},
		CREATE: ListSlowLogsRequestType{
			value: "CREATE",
		},
	}
}

func (c ListSlowLogsRequestType) Value() string {
	return c.value
}

func (c ListSlowLogsRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSlowLogsRequestType) UnmarshalJSON(b []byte) error {
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
