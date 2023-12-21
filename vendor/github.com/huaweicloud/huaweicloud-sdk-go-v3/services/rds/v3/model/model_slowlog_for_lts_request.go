package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SlowlogForLtsRequest 查询实例的慢日志对象
type SlowlogForLtsRequest struct {

	// 开始日期，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartTime string `json:"start_time"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。只能查询当前时间前一个月内的慢日志。
	EndTime string `json:"end_time"`

	// 语句类型，取空值，表示查询所有语句类型。
	Type *SlowlogForLtsRequestType `json:"type,omitempty"`

	// 日志单行序列号，第一次查询时不需要此参数，后续分页查询时需要使用，可从上次查询的返回信息中获取。line_num应在start_time和end_time之间。
	LineNum *string `json:"line_num,omitempty"`

	// 每页多少条记录（查询结果），取值范围是1~100，不填时默认为10。
	Limit *int32 `json:"limit,omitempty"`

	// 搜索方式。默认forwards。配合line_num使用，以line_num为起点，向前搜索或向后搜索。
	SearchType *SlowlogForLtsRequestSearchType `json:"search_type,omitempty"`

	// 数据库名称。
	Database *string `json:"database,omitempty"`
}

func (o SlowlogForLtsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SlowlogForLtsRequest struct{}"
	}

	return strings.Join([]string{"SlowlogForLtsRequest", string(data)}, " ")
}

type SlowlogForLtsRequestType struct {
	value string
}

type SlowlogForLtsRequestTypeEnum struct {
	INSERT SlowlogForLtsRequestType
	UPDATE SlowlogForLtsRequestType
	SELECT SlowlogForLtsRequestType
	DELETE SlowlogForLtsRequestType
	CREATE SlowlogForLtsRequestType
}

func GetSlowlogForLtsRequestTypeEnum() SlowlogForLtsRequestTypeEnum {
	return SlowlogForLtsRequestTypeEnum{
		INSERT: SlowlogForLtsRequestType{
			value: "INSERT",
		},
		UPDATE: SlowlogForLtsRequestType{
			value: "UPDATE",
		},
		SELECT: SlowlogForLtsRequestType{
			value: "SELECT",
		},
		DELETE: SlowlogForLtsRequestType{
			value: "DELETE",
		},
		CREATE: SlowlogForLtsRequestType{
			value: "CREATE",
		},
	}
}

func (c SlowlogForLtsRequestType) Value() string {
	return c.value
}

func (c SlowlogForLtsRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SlowlogForLtsRequestType) UnmarshalJSON(b []byte) error {
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

type SlowlogForLtsRequestSearchType struct {
	value string
}

type SlowlogForLtsRequestSearchTypeEnum struct {
	FORWARDS  SlowlogForLtsRequestSearchType
	BACKWARDS SlowlogForLtsRequestSearchType
}

func GetSlowlogForLtsRequestSearchTypeEnum() SlowlogForLtsRequestSearchTypeEnum {
	return SlowlogForLtsRequestSearchTypeEnum{
		FORWARDS: SlowlogForLtsRequestSearchType{
			value: "forwards",
		},
		BACKWARDS: SlowlogForLtsRequestSearchType{
			value: "backwards",
		},
	}
}

func (c SlowlogForLtsRequestSearchType) Value() string {
	return c.value
}

func (c SlowlogForLtsRequestSearchType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SlowlogForLtsRequestSearchType) UnmarshalJSON(b []byte) error {
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
