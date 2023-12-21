package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateRestoreInstanceRequest Request Object
type CreateRestoreInstanceRequest struct {

	// 语言
	XLanguage *CreateRestoreInstanceRequestXLanguage `json:"X-Language,omitempty"`

	Body *CreateRestoreInstanceRequestBody `json:"body,omitempty"`
}

func (o CreateRestoreInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRestoreInstanceRequest struct{}"
	}

	return strings.Join([]string{"CreateRestoreInstanceRequest", string(data)}, " ")
}

type CreateRestoreInstanceRequestXLanguage struct {
	value string
}

type CreateRestoreInstanceRequestXLanguageEnum struct {
	ZH_CN CreateRestoreInstanceRequestXLanguage
	EN_US CreateRestoreInstanceRequestXLanguage
}

func GetCreateRestoreInstanceRequestXLanguageEnum() CreateRestoreInstanceRequestXLanguageEnum {
	return CreateRestoreInstanceRequestXLanguageEnum{
		ZH_CN: CreateRestoreInstanceRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: CreateRestoreInstanceRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c CreateRestoreInstanceRequestXLanguage) Value() string {
	return c.value
}

func (c CreateRestoreInstanceRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateRestoreInstanceRequestXLanguage) UnmarshalJSON(b []byte) error {
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
