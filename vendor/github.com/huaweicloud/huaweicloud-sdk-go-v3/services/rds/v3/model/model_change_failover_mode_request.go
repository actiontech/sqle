package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ChangeFailoverModeRequest Request Object
type ChangeFailoverModeRequest struct {

	// 语言
	XLanguage *ChangeFailoverModeRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *FailoverModeRequest `json:"body,omitempty"`
}

func (o ChangeFailoverModeRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChangeFailoverModeRequest struct{}"
	}

	return strings.Join([]string{"ChangeFailoverModeRequest", string(data)}, " ")
}

type ChangeFailoverModeRequestXLanguage struct {
	value string
}

type ChangeFailoverModeRequestXLanguageEnum struct {
	ZH_CN ChangeFailoverModeRequestXLanguage
	EN_US ChangeFailoverModeRequestXLanguage
}

func GetChangeFailoverModeRequestXLanguageEnum() ChangeFailoverModeRequestXLanguageEnum {
	return ChangeFailoverModeRequestXLanguageEnum{
		ZH_CN: ChangeFailoverModeRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ChangeFailoverModeRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ChangeFailoverModeRequestXLanguage) Value() string {
	return c.value
}

func (c ChangeFailoverModeRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ChangeFailoverModeRequestXLanguage) UnmarshalJSON(b []byte) error {
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
