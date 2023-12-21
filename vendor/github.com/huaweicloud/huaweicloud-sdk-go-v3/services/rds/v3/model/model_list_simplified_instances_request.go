package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListSimplifiedInstancesRequest Request Object
type ListSimplifiedInstancesRequest struct {

	// 语言
	XLanguage *ListSimplifiedInstancesRequestXLanguage `json:"X-Language,omitempty"`

	Body *SimplifiedInstancesRequest `json:"body,omitempty"`
}

func (o ListSimplifiedInstancesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListSimplifiedInstancesRequest struct{}"
	}

	return strings.Join([]string{"ListSimplifiedInstancesRequest", string(data)}, " ")
}

type ListSimplifiedInstancesRequestXLanguage struct {
	value string
}

type ListSimplifiedInstancesRequestXLanguageEnum struct {
	ZH_CN ListSimplifiedInstancesRequestXLanguage
	EN_US ListSimplifiedInstancesRequestXLanguage
}

func GetListSimplifiedInstancesRequestXLanguageEnum() ListSimplifiedInstancesRequestXLanguageEnum {
	return ListSimplifiedInstancesRequestXLanguageEnum{
		ZH_CN: ListSimplifiedInstancesRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListSimplifiedInstancesRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListSimplifiedInstancesRequestXLanguage) Value() string {
	return c.value
}

func (c ListSimplifiedInstancesRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListSimplifiedInstancesRequestXLanguage) UnmarshalJSON(b []byte) error {
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
