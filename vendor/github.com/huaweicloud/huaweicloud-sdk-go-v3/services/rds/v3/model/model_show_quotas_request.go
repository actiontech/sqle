package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowQuotasRequest Request Object
type ShowQuotasRequest struct {

	// 语言
	XLanguage *ShowQuotasRequestXLanguage `json:"X-Language,omitempty"`
}

func (o ShowQuotasRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowQuotasRequest struct{}"
	}

	return strings.Join([]string{"ShowQuotasRequest", string(data)}, " ")
}

type ShowQuotasRequestXLanguage struct {
	value string
}

type ShowQuotasRequestXLanguageEnum struct {
	ZH_CN ShowQuotasRequestXLanguage
	EN_US ShowQuotasRequestXLanguage
}

func GetShowQuotasRequestXLanguageEnum() ShowQuotasRequestXLanguageEnum {
	return ShowQuotasRequestXLanguageEnum{
		ZH_CN: ShowQuotasRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ShowQuotasRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ShowQuotasRequestXLanguage) Value() string {
	return c.value
}

func (c ShowQuotasRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowQuotasRequestXLanguage) UnmarshalJSON(b []byte) error {
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
