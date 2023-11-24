package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSlowLogStatisticsForLtsRequest Request Object
type ListSlowLogStatisticsForLtsRequest struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 语言。默认en-us。
	XLanguage *ListSlowLogStatisticsForLtsRequestXLanguage `json:"X-Language,omitempty"`

	Body *SlowLogStatisticsForLtsRequest `json:"body,omitempty"`
}

func (o ListSlowLogStatisticsForLtsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSlowLogStatisticsForLtsRequest struct{}"
	}

	return strings.Join([]string{"ListSlowLogStatisticsForLtsRequest", string(data)}, " ")
}

type ListSlowLogStatisticsForLtsRequestXLanguage struct {
	value string
}

type ListSlowLogStatisticsForLtsRequestXLanguageEnum struct {
	ZH_CN ListSlowLogStatisticsForLtsRequestXLanguage
	EN_US ListSlowLogStatisticsForLtsRequestXLanguage
}

func GetListSlowLogStatisticsForLtsRequestXLanguageEnum() ListSlowLogStatisticsForLtsRequestXLanguageEnum {
	return ListSlowLogStatisticsForLtsRequestXLanguageEnum{
		ZH_CN: ListSlowLogStatisticsForLtsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListSlowLogStatisticsForLtsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListSlowLogStatisticsForLtsRequestXLanguage) Value() string {
	return c.value
}

func (c ListSlowLogStatisticsForLtsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSlowLogStatisticsForLtsRequestXLanguage) UnmarshalJSON(b []byte) error {
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
