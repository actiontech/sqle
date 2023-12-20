package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DownloadSlowlogRequest Request Object
type DownloadSlowlogRequest struct {

	// 语言
	XLanguage *DownloadSlowlogRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *SlowlogDownloadRequest `json:"body,omitempty"`
}

func (o DownloadSlowlogRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DownloadSlowlogRequest struct{}"
	}

	return strings.Join([]string{"DownloadSlowlogRequest", string(data)}, " ")
}

type DownloadSlowlogRequestXLanguage struct {
	value string
}

type DownloadSlowlogRequestXLanguageEnum struct {
	ZH_CN DownloadSlowlogRequestXLanguage
	EN_US DownloadSlowlogRequestXLanguage
}

func GetDownloadSlowlogRequestXLanguageEnum() DownloadSlowlogRequestXLanguageEnum {
	return DownloadSlowlogRequestXLanguageEnum{
		ZH_CN: DownloadSlowlogRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: DownloadSlowlogRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c DownloadSlowlogRequestXLanguage) Value() string {
	return c.value
}

func (c DownloadSlowlogRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DownloadSlowlogRequestXLanguage) UnmarshalJSON(b []byte) error {
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
