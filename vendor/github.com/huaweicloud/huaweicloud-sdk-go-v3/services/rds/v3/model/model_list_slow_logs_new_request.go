package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSlowLogsNewRequest Request Object
type ListSlowLogsNewRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 开始日期，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartDate string `json:"start_date"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。只能查询当前时间前一个月内的慢日志。
	EndDate string `json:"end_date"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *int64 `json:"offset,omitempty"`

	// 每页多少条记录，取值范围是1~100，不填时默认为10。
	Limit *int64 `json:"limit,omitempty"`

	// 语句类型，取空值，表示查询所有语句类型。
	Type *ListSlowLogsNewRequestType `json:"type,omitempty"`
}

func (o ListSlowLogsNewRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSlowLogsNewRequest struct{}"
	}

	return strings.Join([]string{"ListSlowLogsNewRequest", string(data)}, " ")
}

type ListSlowLogsNewRequestType struct {
	value string
}

type ListSlowLogsNewRequestTypeEnum struct {
	INSERT ListSlowLogsNewRequestType
	UPDATE ListSlowLogsNewRequestType
	SELECT ListSlowLogsNewRequestType
	DELETE ListSlowLogsNewRequestType
	CREATE ListSlowLogsNewRequestType
}

func GetListSlowLogsNewRequestTypeEnum() ListSlowLogsNewRequestTypeEnum {
	return ListSlowLogsNewRequestTypeEnum{
		INSERT: ListSlowLogsNewRequestType{
			value: "INSERT",
		},
		UPDATE: ListSlowLogsNewRequestType{
			value: "UPDATE",
		},
		SELECT: ListSlowLogsNewRequestType{
			value: "SELECT",
		},
		DELETE: ListSlowLogsNewRequestType{
			value: "DELETE",
		},
		CREATE: ListSlowLogsNewRequestType{
			value: "CREATE",
		},
	}
}

func (c ListSlowLogsNewRequestType) Value() string {
	return c.value
}

func (c ListSlowLogsNewRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSlowLogsNewRequestType) UnmarshalJSON(b []byte) error {
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
