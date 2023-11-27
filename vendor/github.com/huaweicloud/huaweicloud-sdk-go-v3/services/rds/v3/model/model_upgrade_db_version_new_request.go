package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpgradeDbVersionNewRequest Request Object
type UpgradeDbVersionNewRequest struct {

	// 语言
	XLanguage *UpgradeDbVersionNewRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CustomerUpgradeDatabaseVersionReqNew `json:"body,omitempty"`
}

func (o UpgradeDbVersionNewRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpgradeDbVersionNewRequest struct{}"
	}

	return strings.Join([]string{"UpgradeDbVersionNewRequest", string(data)}, " ")
}

type UpgradeDbVersionNewRequestXLanguage struct {
	value string
}

type UpgradeDbVersionNewRequestXLanguageEnum struct {
	ZH_CN UpgradeDbVersionNewRequestXLanguage
	EN_US UpgradeDbVersionNewRequestXLanguage
}

func GetUpgradeDbVersionNewRequestXLanguageEnum() UpgradeDbVersionNewRequestXLanguageEnum {
	return UpgradeDbVersionNewRequestXLanguageEnum{
		ZH_CN: UpgradeDbVersionNewRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: UpgradeDbVersionNewRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c UpgradeDbVersionNewRequestXLanguage) Value() string {
	return c.value
}

func (c UpgradeDbVersionNewRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpgradeDbVersionNewRequestXLanguage) UnmarshalJSON(b []byte) error {
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
