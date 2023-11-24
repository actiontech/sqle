package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ResetPwdRequest Request Object
type ResetPwdRequest struct {

	// 语言
	XLanguage *ResetPwdRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *PwdResetRequest `json:"body,omitempty"`
}

func (o ResetPwdRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResetPwdRequest struct{}"
	}

	return strings.Join([]string{"ResetPwdRequest", string(data)}, " ")
}

type ResetPwdRequestXLanguage struct {
	value string
}

type ResetPwdRequestXLanguageEnum struct {
	ZH_CN ResetPwdRequestXLanguage
	EN_US ResetPwdRequestXLanguage
}

func GetResetPwdRequestXLanguageEnum() ResetPwdRequestXLanguageEnum {
	return ResetPwdRequestXLanguageEnum{
		ZH_CN: ResetPwdRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ResetPwdRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ResetPwdRequestXLanguage) Value() string {
	return c.value
}

func (c ResetPwdRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResetPwdRequestXLanguage) UnmarshalJSON(b []byte) error {
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
