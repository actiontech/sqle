package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSlowlogStatisticsRequest Request Object
type ListSlowlogStatisticsRequest struct {

	// 语言
	XLanguage *ListSlowlogStatisticsRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 当前页号
	CurPage int32 `json:"cur_page"`

	// 每页多少条记录，取值范围0~100
	PerPage int32 `json:"per_page"`

	// 开始时间，格式为“yyyy-mm-ddThh:mm:ssZ”。其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	StartDate string `json:"start_date"`

	// 结束时间，格式为“yyyy-mm-ddThh:mm:ssZ”。其中，T指某个时间的开始；Z指时区偏移量，例如北京时间偏移显示为+0800。
	EndDate string `json:"end_date"`

	// 语句类型，ALL表示查询所有语句类型，也可指定日志类型 - INSERT - UPDATE - SELECT - DELETE - CREATE - ALL
	Type ListSlowlogStatisticsRequestType `json:"type"`

	// 取值范围：\"executeTime\",表示按执行时间降序排序，不传或者传其他表示按执行次数降序排序
	Sort *string `json:"sort,omitempty"`
}

func (o ListSlowlogStatisticsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSlowlogStatisticsRequest struct{}"
	}

	return strings.Join([]string{"ListSlowlogStatisticsRequest", string(data)}, " ")
}

type ListSlowlogStatisticsRequestXLanguage struct {
	value string
}

type ListSlowlogStatisticsRequestXLanguageEnum struct {
	ZH_CN ListSlowlogStatisticsRequestXLanguage
	EN_US ListSlowlogStatisticsRequestXLanguage
}

func GetListSlowlogStatisticsRequestXLanguageEnum() ListSlowlogStatisticsRequestXLanguageEnum {
	return ListSlowlogStatisticsRequestXLanguageEnum{
		ZH_CN: ListSlowlogStatisticsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListSlowlogStatisticsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListSlowlogStatisticsRequestXLanguage) Value() string {
	return c.value
}

func (c ListSlowlogStatisticsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSlowlogStatisticsRequestXLanguage) UnmarshalJSON(b []byte) error {
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

type ListSlowlogStatisticsRequestType struct {
	value string
}

type ListSlowlogStatisticsRequestTypeEnum struct {
	INSERT ListSlowlogStatisticsRequestType
	UPDATE ListSlowlogStatisticsRequestType
	SELECT ListSlowlogStatisticsRequestType
	DELETE ListSlowlogStatisticsRequestType
	CREATE ListSlowlogStatisticsRequestType
	ALL    ListSlowlogStatisticsRequestType
}

func GetListSlowlogStatisticsRequestTypeEnum() ListSlowlogStatisticsRequestTypeEnum {
	return ListSlowlogStatisticsRequestTypeEnum{
		INSERT: ListSlowlogStatisticsRequestType{
			value: "INSERT",
		},
		UPDATE: ListSlowlogStatisticsRequestType{
			value: "UPDATE",
		},
		SELECT: ListSlowlogStatisticsRequestType{
			value: "SELECT",
		},
		DELETE: ListSlowlogStatisticsRequestType{
			value: "DELETE",
		},
		CREATE: ListSlowlogStatisticsRequestType{
			value: "CREATE",
		},
		ALL: ListSlowlogStatisticsRequestType{
			value: "ALL",
		},
	}
}

func (c ListSlowlogStatisticsRequestType) Value() string {
	return c.value
}

func (c ListSlowlogStatisticsRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSlowlogStatisticsRequestType) UnmarshalJSON(b []byte) error {
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
