package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListErrorLogsNewRequest Request Object
type ListErrorLogsNewRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartDate string `json:"start_date"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。 只能查询当前时间前一个月内的错误日志。
	EndDate string `json:"end_date"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *int64 `json:"offset,omitempty"`

	// 每页多少条记录，取值范围是1~100，不填时默认为10。
	Limit *int64 `json:"limit,omitempty"`

	// 日志级别，默认为ALL。
	Level *ListErrorLogsNewRequestLevel `json:"level,omitempty"`
}

func (o ListErrorLogsNewRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListErrorLogsNewRequest struct{}"
	}

	return strings.Join([]string{"ListErrorLogsNewRequest", string(data)}, " ")
}

type ListErrorLogsNewRequestLevel struct {
	value string
}

type ListErrorLogsNewRequestLevelEnum struct {
	ALL     ListErrorLogsNewRequestLevel
	INFO    ListErrorLogsNewRequestLevel
	LOG     ListErrorLogsNewRequestLevel
	WARNING ListErrorLogsNewRequestLevel
	ERROR   ListErrorLogsNewRequestLevel
	FATAL   ListErrorLogsNewRequestLevel
	PANIC   ListErrorLogsNewRequestLevel
	NOTE    ListErrorLogsNewRequestLevel
}

func GetListErrorLogsNewRequestLevelEnum() ListErrorLogsNewRequestLevelEnum {
	return ListErrorLogsNewRequestLevelEnum{
		ALL: ListErrorLogsNewRequestLevel{
			value: "ALL",
		},
		INFO: ListErrorLogsNewRequestLevel{
			value: "INFO",
		},
		LOG: ListErrorLogsNewRequestLevel{
			value: "LOG",
		},
		WARNING: ListErrorLogsNewRequestLevel{
			value: "WARNING",
		},
		ERROR: ListErrorLogsNewRequestLevel{
			value: "ERROR",
		},
		FATAL: ListErrorLogsNewRequestLevel{
			value: "FATAL",
		},
		PANIC: ListErrorLogsNewRequestLevel{
			value: "PANIC",
		},
		NOTE: ListErrorLogsNewRequestLevel{
			value: "NOTE",
		},
	}
}

func (c ListErrorLogsNewRequestLevel) Value() string {
	return c.value
}

func (c ListErrorLogsNewRequestLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListErrorLogsNewRequestLevel) UnmarshalJSON(b []byte) error {
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
