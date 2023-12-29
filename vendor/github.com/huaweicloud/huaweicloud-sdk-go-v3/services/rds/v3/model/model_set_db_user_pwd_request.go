package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SetDbUserPwdRequest Request Object
type SetDbUserPwdRequest struct {

	// 语言
	XLanguage *SetDbUserPwdRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *DbUserPwdRequest `json:"body,omitempty"`
}

func (o SetDbUserPwdRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetDbUserPwdRequest struct{}"
	}

	return strings.Join([]string{"SetDbUserPwdRequest", string(data)}, " ")
}

type SetDbUserPwdRequestXLanguage struct {
	value string
}

type SetDbUserPwdRequestXLanguageEnum struct {
	ZH_CN SetDbUserPwdRequestXLanguage
	EN_US SetDbUserPwdRequestXLanguage
}

func GetSetDbUserPwdRequestXLanguageEnum() SetDbUserPwdRequestXLanguageEnum {
	return SetDbUserPwdRequestXLanguageEnum{
		ZH_CN: SetDbUserPwdRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: SetDbUserPwdRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c SetDbUserPwdRequestXLanguage) Value() string {
	return c.value
}

func (c SetDbUserPwdRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SetDbUserPwdRequestXLanguage) UnmarshalJSON(b []byte) error {
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
