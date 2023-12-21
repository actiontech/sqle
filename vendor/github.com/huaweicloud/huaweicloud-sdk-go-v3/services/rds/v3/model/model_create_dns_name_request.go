package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateDnsNameRequest Request Object
type CreateDnsNameRequest struct {

	// 语言
	XLanguage *CreateDnsNameRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateDnsNameRequestBody `json:"body,omitempty"`
}

func (o CreateDnsNameRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateDnsNameRequest struct{}"
	}

	return strings.Join([]string{"CreateDnsNameRequest", string(data)}, " ")
}

type CreateDnsNameRequestXLanguage struct {
	value string
}

type CreateDnsNameRequestXLanguageEnum struct {
	ZH_CN CreateDnsNameRequestXLanguage
	EN_US CreateDnsNameRequestXLanguage
}

func GetCreateDnsNameRequestXLanguageEnum() CreateDnsNameRequestXLanguageEnum {
	return CreateDnsNameRequestXLanguageEnum{
		ZH_CN: CreateDnsNameRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: CreateDnsNameRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c CreateDnsNameRequestXLanguage) Value() string {
	return c.value
}

func (c CreateDnsNameRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateDnsNameRequestXLanguage) UnmarshalJSON(b []byte) error {
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
