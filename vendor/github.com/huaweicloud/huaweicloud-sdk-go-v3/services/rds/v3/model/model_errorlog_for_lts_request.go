package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ErrorlogForLtsRequest 查询实例的错误日志对象
type ErrorlogForLtsRequest struct {

	// 开始日期，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartTime string `json:"start_time"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。只能查询当前时间前一个月内的慢日志。
	EndTime string `json:"end_time"`

	// 日志级别，默认为ALL。
	Level *ErrorlogForLtsRequestLevel `json:"level,omitempty"`

	// 日志单行序列号，第一次查询时不需要此参数，后续分页查询时需要使用，可从上次查询的返回信息中获取。line_num应在start_time和end_time之间。
	LineNum *string `json:"line_num,omitempty"`

	// 每页多少条记录（查询结果），取值范围是1~100，不填时默认为10。
	Limit *int32 `json:"limit,omitempty"`

	// 搜索方式。默认forwards。配合line_num使用，以line_num为起点，向前搜索或向后搜索。
	SearchType *string `json:"search_type,omitempty"`
}

func (o ErrorlogForLtsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ErrorlogForLtsRequest struct{}"
	}

	return strings.Join([]string{"ErrorlogForLtsRequest", string(data)}, " ")
}

type ErrorlogForLtsRequestLevel struct {
	value string
}

type ErrorlogForLtsRequestLevelEnum struct {
	ALL     ErrorlogForLtsRequestLevel
	INFO    ErrorlogForLtsRequestLevel
	LOG     ErrorlogForLtsRequestLevel
	WARNING ErrorlogForLtsRequestLevel
	ERROR   ErrorlogForLtsRequestLevel
	FATAL   ErrorlogForLtsRequestLevel
	PANIC   ErrorlogForLtsRequestLevel
	NOTE    ErrorlogForLtsRequestLevel
}

func GetErrorlogForLtsRequestLevelEnum() ErrorlogForLtsRequestLevelEnum {
	return ErrorlogForLtsRequestLevelEnum{
		ALL: ErrorlogForLtsRequestLevel{
			value: "ALL",
		},
		INFO: ErrorlogForLtsRequestLevel{
			value: "INFO",
		},
		LOG: ErrorlogForLtsRequestLevel{
			value: "LOG",
		},
		WARNING: ErrorlogForLtsRequestLevel{
			value: "WARNING",
		},
		ERROR: ErrorlogForLtsRequestLevel{
			value: "ERROR",
		},
		FATAL: ErrorlogForLtsRequestLevel{
			value: "FATAL",
		},
		PANIC: ErrorlogForLtsRequestLevel{
			value: "PANIC",
		},
		NOTE: ErrorlogForLtsRequestLevel{
			value: "NOTE",
		},
	}
}

func (c ErrorlogForLtsRequestLevel) Value() string {
	return c.value
}

func (c ErrorlogForLtsRequestLevel) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ErrorlogForLtsRequestLevel) UnmarshalJSON(b []byte) error {
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
