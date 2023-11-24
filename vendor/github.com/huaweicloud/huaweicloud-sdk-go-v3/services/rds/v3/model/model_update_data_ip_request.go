package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpdateDataIpRequest Request Object
type UpdateDataIpRequest struct {

	// 语言
	XLanguage *UpdateDataIpRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *DataIpRequest `json:"body,omitempty"`
}

func (o UpdateDataIpRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDataIpRequest struct{}"
	}

	return strings.Join([]string{"UpdateDataIpRequest", string(data)}, " ")
}

type UpdateDataIpRequestXLanguage struct {
	value string
}

type UpdateDataIpRequestXLanguageEnum struct {
	ZH_CN UpdateDataIpRequestXLanguage
	EN_US UpdateDataIpRequestXLanguage
}

func GetUpdateDataIpRequestXLanguageEnum() UpdateDataIpRequestXLanguageEnum {
	return UpdateDataIpRequestXLanguageEnum{
		ZH_CN: UpdateDataIpRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: UpdateDataIpRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c UpdateDataIpRequestXLanguage) Value() string {
	return c.value
}

func (c UpdateDataIpRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateDataIpRequestXLanguage) UnmarshalJSON(b []byte) error {
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
