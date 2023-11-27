package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// StartInstanceSingleToHaActionRequest Request Object
type StartInstanceSingleToHaActionRequest struct {

	// 语言
	XLanguage *StartInstanceSingleToHaActionRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *Single2Ha `json:"body,omitempty"`
}

func (o StartInstanceSingleToHaActionRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "StartInstanceSingleToHaActionRequest struct{}"
	}

	return strings.Join([]string{"StartInstanceSingleToHaActionRequest", string(data)}, " ")
}

type StartInstanceSingleToHaActionRequestXLanguage struct {
	value string
}

type StartInstanceSingleToHaActionRequestXLanguageEnum struct {
	ZH_CN StartInstanceSingleToHaActionRequestXLanguage
	EN_US StartInstanceSingleToHaActionRequestXLanguage
}

func GetStartInstanceSingleToHaActionRequestXLanguageEnum() StartInstanceSingleToHaActionRequestXLanguageEnum {
	return StartInstanceSingleToHaActionRequestXLanguageEnum{
		ZH_CN: StartInstanceSingleToHaActionRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: StartInstanceSingleToHaActionRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c StartInstanceSingleToHaActionRequestXLanguage) Value() string {
	return c.value
}

func (c StartInstanceSingleToHaActionRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *StartInstanceSingleToHaActionRequestXLanguage) UnmarshalJSON(b []byte) error {
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
