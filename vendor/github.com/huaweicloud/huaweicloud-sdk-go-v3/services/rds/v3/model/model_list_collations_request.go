package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListCollationsRequest Request Object
type ListCollationsRequest struct {

	// 语言
	XLanguage *ListCollationsRequestXLanguage `json:"X-Language,omitempty"`
}

func (o ListCollationsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListCollationsRequest struct{}"
	}

	return strings.Join([]string{"ListCollationsRequest", string(data)}, " ")
}

type ListCollationsRequestXLanguage struct {
	value string
}

type ListCollationsRequestXLanguageEnum struct {
	ZH_CN ListCollationsRequestXLanguage
	EN_US ListCollationsRequestXLanguage
}

func GetListCollationsRequestXLanguageEnum() ListCollationsRequestXLanguageEnum {
	return ListCollationsRequestXLanguageEnum{
		ZH_CN: ListCollationsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListCollationsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListCollationsRequestXLanguage) Value() string {
	return c.value
}

func (c ListCollationsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListCollationsRequestXLanguage) UnmarshalJSON(b []byte) error {
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
