package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SlowLogStatisticsForLtsRequest 查询实例的慢日志对象
type SlowLogStatisticsForLtsRequest struct {

	// 开始日期，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartTime string `json:"start_time"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。 其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。只能查询当前时间前一个月内的慢日志。
	EndTime string `json:"end_time"`

	// 索引位置，偏移量。默认为0，表示从第一条数据开始查询。
	Offset *int32 `json:"offset,omitempty"`

	// 每页多少条记录（查询结果），取值范围是1~100，不填时默认为10。
	Limit *int32 `json:"limit,omitempty"`

	// 语句类型，取空值，表示查询所有语句类型。
	Type *SlowLogStatisticsForLtsRequestType `json:"type,omitempty"`

	// 数据库名称。
	Database *string `json:"database,omitempty"`

	// 指定排序字段。\"executeTime\"，表示按照执行时间降序排序。字段为空或传入其他值，表示按照执行次数降序排序。
	Sort *string `json:"sort,omitempty"`

	// 排序顺序。默认desc。
	Order *SlowLogStatisticsForLtsRequestOrder `json:"order,omitempty"`
}

func (o SlowLogStatisticsForLtsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SlowLogStatisticsForLtsRequest struct{}"
	}

	return strings.Join([]string{"SlowLogStatisticsForLtsRequest", string(data)}, " ")
}

type SlowLogStatisticsForLtsRequestType struct {
	value string
}

type SlowLogStatisticsForLtsRequestTypeEnum struct {
	INSERT SlowLogStatisticsForLtsRequestType
	UPDATE SlowLogStatisticsForLtsRequestType
	SELECT SlowLogStatisticsForLtsRequestType
	DELETE SlowLogStatisticsForLtsRequestType
	CREATE SlowLogStatisticsForLtsRequestType
	ALL    SlowLogStatisticsForLtsRequestType
}

func GetSlowLogStatisticsForLtsRequestTypeEnum() SlowLogStatisticsForLtsRequestTypeEnum {
	return SlowLogStatisticsForLtsRequestTypeEnum{
		INSERT: SlowLogStatisticsForLtsRequestType{
			value: "INSERT",
		},
		UPDATE: SlowLogStatisticsForLtsRequestType{
			value: "UPDATE",
		},
		SELECT: SlowLogStatisticsForLtsRequestType{
			value: "SELECT",
		},
		DELETE: SlowLogStatisticsForLtsRequestType{
			value: "DELETE",
		},
		CREATE: SlowLogStatisticsForLtsRequestType{
			value: "CREATE",
		},
		ALL: SlowLogStatisticsForLtsRequestType{
			value: "ALL",
		},
	}
}

func (c SlowLogStatisticsForLtsRequestType) Value() string {
	return c.value
}

func (c SlowLogStatisticsForLtsRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SlowLogStatisticsForLtsRequestType) UnmarshalJSON(b []byte) error {
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

type SlowLogStatisticsForLtsRequestOrder struct {
	value string
}

type SlowLogStatisticsForLtsRequestOrderEnum struct {
	DESC SlowLogStatisticsForLtsRequestOrder
	ASC  SlowLogStatisticsForLtsRequestOrder
}

func GetSlowLogStatisticsForLtsRequestOrderEnum() SlowLogStatisticsForLtsRequestOrderEnum {
	return SlowLogStatisticsForLtsRequestOrderEnum{
		DESC: SlowLogStatisticsForLtsRequestOrder{
			value: "desc",
		},
		ASC: SlowLogStatisticsForLtsRequestOrder{
			value: "asc",
		},
	}
}

func (c SlowLogStatisticsForLtsRequestOrder) Value() string {
	return c.value
}

func (c SlowLogStatisticsForLtsRequestOrder) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SlowLogStatisticsForLtsRequestOrder) UnmarshalJSON(b []byte) error {
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
