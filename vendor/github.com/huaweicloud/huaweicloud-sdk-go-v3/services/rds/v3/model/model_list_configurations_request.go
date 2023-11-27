package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListConfigurationsRequest Request Object
type ListConfigurationsRequest struct {

	// 语言
	XLanguage *ListConfigurationsRequestXLanguage `json:"X-Language,omitempty"`
}

func (o ListConfigurationsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListConfigurationsRequest struct{}"
	}

	return strings.Join([]string{"ListConfigurationsRequest", string(data)}, " ")
}

type ListConfigurationsRequestXLanguage struct {
	value string
}

type ListConfigurationsRequestXLanguageEnum struct {
	ZH_CN ListConfigurationsRequestXLanguage
	EN_US ListConfigurationsRequestXLanguage
}

func GetListConfigurationsRequestXLanguageEnum() ListConfigurationsRequestXLanguageEnum {
	return ListConfigurationsRequestXLanguageEnum{
		ZH_CN: ListConfigurationsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListConfigurationsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListConfigurationsRequestXLanguage) Value() string {
	return c.value
}

func (c ListConfigurationsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListConfigurationsRequestXLanguage) UnmarshalJSON(b []byte) error {
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
