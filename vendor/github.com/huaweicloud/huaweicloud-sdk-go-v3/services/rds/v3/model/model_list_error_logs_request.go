package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListErrorLogsRequest Request Object
type ListErrorLogsRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartDate string `json:"start_date"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。  其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。  只能查询当前时间前一个月内的错误日志。
	EndDate string `json:"end_date"`

	// 页数偏移量，如1、2、3、4等，不填时默认为1。
	Offset *int32 `json:"offset,omitempty"`

	// 每页多少条记录，取值范围是1~100，不填时默认为10。
	Limit *int32 `json:"limit,omitempty"`

	// 日志级别，默认为ALL。
	Level *ListErrorLogsRequestLevel `json:"level,omitempty"`
}

func (o ListErrorLogsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListErrorLogsRequest struct{}"
	}

	return strings.Join([]string{"ListErrorLogsRequest", string(data)}, " ")
}

type ListErrorLogsRequestLevel struct {
	value string
}

type ListErrorLogsRequestLevelEnum struct {
	ALL     ListErrorLogsRequestLevel
	INFO    ListErrorLogsRequestLevel
	LOG     ListErrorLogsRequestLevel
	WARNING ListErrorLogsRequestLevel
	ERROR   ListErrorLogsRequestLevel
	FATAL   ListErrorLogsRequestLevel
	PANIC   ListErrorLogsRequestLevel
	NOTE    ListErrorLogsRequestLevel
}

func GetListErrorLogsRequestLevelEnum() ListErrorLogsRequestLevelEnum {
	return ListErrorLogsRequestLevelEnum{
		ALL: ListErrorLogsRequestLevel{
			value: "ALL",
		},
		INFO: ListErrorLogsRequestLevel{
			value: "INFO",
		},
		LOG: ListErrorLogsRequestLevel{
			value: "LOG",
		},
		WARNING: ListErrorLogsRequestLevel{
			value: "WARNING",
		},
		ERROR: ListErrorLogsRequestLevel{
			value: "ERROR",
		},
		FATAL: ListErrorLogsRequestLevel{
			value: "FATAL",
		},
		PANIC: ListErrorLogsRequestLevel{
			value: "PANIC",
		},
		NOTE: ListErrorLogsRequestLevel{
			value: "NOTE",
		},
	}
}

func (c ListErrorLogsRequestLevel) Value() string {
	return c.value
}

func (c ListErrorLogsRequestLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListErrorLogsRequestLevel) UnmarshalJSON(b []byte) error {
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
