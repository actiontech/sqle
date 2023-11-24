package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateXelLogDownloadRequest Request Object
type CreateXelLogDownloadRequest struct {

	// 语言
	XLanguage *CreateXelLogDownloadRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *CreateXelLogDownloadRequestBody `json:"body,omitempty"`
}

func (o CreateXelLogDownloadRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateXelLogDownloadRequest struct{}"
	}

	return strings.Join([]string{"CreateXelLogDownloadRequest", string(data)}, " ")
}

type CreateXelLogDownloadRequestXLanguage struct {
	value string
}

type CreateXelLogDownloadRequestXLanguageEnum struct {
	ZH_CN CreateXelLogDownloadRequestXLanguage
	EN_US CreateXelLogDownloadRequestXLanguage
}

func GetCreateXelLogDownloadRequestXLanguageEnum() CreateXelLogDownloadRequestXLanguageEnum {
	return CreateXelLogDownloadRequestXLanguageEnum{
		ZH_CN: CreateXelLogDownloadRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: CreateXelLogDownloadRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c CreateXelLogDownloadRequestXLanguage) Value() string {
	return c.value
}

func (c CreateXelLogDownloadRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateXelLogDownloadRequestXLanguage) UnmarshalJSON(b []byte) error {
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
