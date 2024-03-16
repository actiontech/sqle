package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListPredefinedTagRequest Request Object
type ListPredefinedTagRequest struct {

	// 语言
	XLanguage *ListPredefinedTagRequestXLanguage `json:"X-Language,omitempty"`
}

func (o ListPredefinedTagRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListPredefinedTagRequest struct{}"
	}

	return strings.Join([]string{"ListPredefinedTagRequest", string(data)}, " ")
}

type ListPredefinedTagRequestXLanguage struct {
	value string
}

type ListPredefinedTagRequestXLanguageEnum struct {
	ZH_CN ListPredefinedTagRequestXLanguage
	EN_US ListPredefinedTagRequestXLanguage
}

func GetListPredefinedTagRequestXLanguageEnum() ListPredefinedTagRequestXLanguageEnum {
	return ListPredefinedTagRequestXLanguageEnum{
		ZH_CN: ListPredefinedTagRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListPredefinedTagRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListPredefinedTagRequestXLanguage) Value() string {
	return c.value
}

func (c ListPredefinedTagRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListPredefinedTagRequestXLanguage) UnmarshalJSON(b []byte) error {
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
