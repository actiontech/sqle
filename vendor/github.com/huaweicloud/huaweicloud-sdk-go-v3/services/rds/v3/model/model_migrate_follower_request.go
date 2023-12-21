package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// MigrateFollowerRequest Request Object
type MigrateFollowerRequest struct {

	// 语言
	XLanguage *MigrateFollowerRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *FollowerMigrateRequest `json:"body,omitempty"`
}

func (o MigrateFollowerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MigrateFollowerRequest struct{}"
	}

	return strings.Join([]string{"MigrateFollowerRequest", string(data)}, " ")
}

type MigrateFollowerRequestXLanguage struct {
	value string
}

type MigrateFollowerRequestXLanguageEnum struct {
	ZH_CN MigrateFollowerRequestXLanguage
	EN_US MigrateFollowerRequestXLanguage
}

func GetMigrateFollowerRequestXLanguageEnum() MigrateFollowerRequestXLanguageEnum {
	return MigrateFollowerRequestXLanguageEnum{
		ZH_CN: MigrateFollowerRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: MigrateFollowerRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c MigrateFollowerRequestXLanguage) Value() string {
	return c.value
}

func (c MigrateFollowerRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MigrateFollowerRequestXLanguage) UnmarshalJSON(b []byte) error {
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
