package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// StartupInstanceRequest Request Object
type StartupInstanceRequest struct {

	// 语言
	XLanguage *StartupInstanceRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o StartupInstanceRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartupInstanceRequest struct{}"
	}

	return strings.Join([]string{"StartupInstanceRequest", string(data)}, " ")
}

type StartupInstanceRequestXLanguage struct {
	value string
}

type StartupInstanceRequestXLanguageEnum struct {
	ZH_CN StartupInstanceRequestXLanguage
	EN_US StartupInstanceRequestXLanguage
}

func GetStartupInstanceRequestXLanguageEnum() StartupInstanceRequestXLanguageEnum {
	return StartupInstanceRequestXLanguageEnum{
		ZH_CN: StartupInstanceRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: StartupInstanceRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c StartupInstanceRequestXLanguage) Value() string {
	return c.value
}

func (c StartupInstanceRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *StartupInstanceRequestXLanguage) UnmarshalJSON(b []byte) error {
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
