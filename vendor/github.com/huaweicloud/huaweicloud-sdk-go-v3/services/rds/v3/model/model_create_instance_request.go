package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateInstanceRequest Request Object
type CreateInstanceRequest struct {

	// 语言
	XLanguage *CreateInstanceRequestXLanguage `json:"X-Language,omitempty"`

	// 保证客户端请求幂等性的标识。 该标识为32位UUID格式，由客户端生成，且需确保72小时内不同请求之间该标识具有唯一性。
	XClientToken *string `json:"X-Client-Token,omitempty"`

	Body *InstanceRequest `json:"body,omitempty"`
}

func (o CreateInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceRequest struct{}"
	}

	return strings.Join([]string{"CreateInstanceRequest", string(data)}, " ")
}

type CreateInstanceRequestXLanguage struct {
	value string
}

type CreateInstanceRequestXLanguageEnum struct {
	ZH_CN CreateInstanceRequestXLanguage
	EN_US CreateInstanceRequestXLanguage
}

func GetCreateInstanceRequestXLanguageEnum() CreateInstanceRequestXLanguageEnum {
	return CreateInstanceRequestXLanguageEnum{
		ZH_CN: CreateInstanceRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: CreateInstanceRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c CreateInstanceRequestXLanguage) Value() string {
	return c.value
}

func (c CreateInstanceRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateInstanceRequestXLanguage) UnmarshalJSON(b []byte) error {
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
