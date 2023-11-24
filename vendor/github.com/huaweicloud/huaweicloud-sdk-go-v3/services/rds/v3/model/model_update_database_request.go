package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpdateDatabaseRequest Request Object
type UpdateDatabaseRequest struct {

	// 语言
	XLanguage *UpdateDatabaseRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateDatabaseReq `json:"body,omitempty"`
}

func (o UpdateDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateDatabaseRequest struct{}"
	}

	return strings.Join([]string{"UpdateDatabaseRequest", string(data)}, " ")
}

type UpdateDatabaseRequestXLanguage struct {
	value string
}

type UpdateDatabaseRequestXLanguageEnum struct {
	ZH_CN UpdateDatabaseRequestXLanguage
	EN_US UpdateDatabaseRequestXLanguage
}

func GetUpdateDatabaseRequestXLanguageEnum() UpdateDatabaseRequestXLanguageEnum {
	return UpdateDatabaseRequestXLanguageEnum{
		ZH_CN: UpdateDatabaseRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: UpdateDatabaseRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c UpdateDatabaseRequestXLanguage) Value() string {
	return c.value
}

func (c UpdateDatabaseRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateDatabaseRequestXLanguage) UnmarshalJSON(b []byte) error {
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
