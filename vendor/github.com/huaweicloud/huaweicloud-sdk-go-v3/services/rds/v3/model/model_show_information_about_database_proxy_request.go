package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowInformationAboutDatabaseProxyRequest Request Object
type ShowInformationAboutDatabaseProxyRequest struct {

	// 语言
	XLanguage *ShowInformationAboutDatabaseProxyRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`
}

func (o ShowInformationAboutDatabaseProxyRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowInformationAboutDatabaseProxyRequest struct{}"
	}

	return strings.Join([]string{"ShowInformationAboutDatabaseProxyRequest", string(data)}, " ")
}

type ShowInformationAboutDatabaseProxyRequestXLanguage struct {
	value string
}

type ShowInformationAboutDatabaseProxyRequestXLanguageEnum struct {
	ZH_CN ShowInformationAboutDatabaseProxyRequestXLanguage
	EN_US ShowInformationAboutDatabaseProxyRequestXLanguage
}

func GetShowInformationAboutDatabaseProxyRequestXLanguageEnum() ShowInformationAboutDatabaseProxyRequestXLanguageEnum {
	return ShowInformationAboutDatabaseProxyRequestXLanguageEnum{
		ZH_CN: ShowInformationAboutDatabaseProxyRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ShowInformationAboutDatabaseProxyRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ShowInformationAboutDatabaseProxyRequestXLanguage) Value() string {
	return c.value
}

func (c ShowInformationAboutDatabaseProxyRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowInformationAboutDatabaseProxyRequestXLanguage) UnmarshalJSON(b []byte) error {
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
