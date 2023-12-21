package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListOffSiteInstancesRequest Request Object
type ListOffSiteInstancesRequest struct {
	ContentType *string `json:"Content-Type,omitempty"`

	// 语言
	XLanguage *ListOffSiteInstancesRequestXLanguage `json:"X-Language,omitempty"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *interface{} `json:"offset,omitempty"`

	// 查询记录数。默认为100，不能为负数，最小值为1，最大值为100。
	Limit *interface{} `json:"limit,omitempty"`
}

func (o ListOffSiteInstancesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListOffSiteInstancesRequest struct{}"
	}

	return strings.Join([]string{"ListOffSiteInstancesRequest", string(data)}, " ")
}

type ListOffSiteInstancesRequestXLanguage struct {
	value string
}

type ListOffSiteInstancesRequestXLanguageEnum struct {
	ZH_CN ListOffSiteInstancesRequestXLanguage
	EN_US ListOffSiteInstancesRequestXLanguage
}

func GetListOffSiteInstancesRequestXLanguageEnum() ListOffSiteInstancesRequestXLanguageEnum {
	return ListOffSiteInstancesRequestXLanguageEnum{
		ZH_CN: ListOffSiteInstancesRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListOffSiteInstancesRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListOffSiteInstancesRequestXLanguage) Value() string {
	return c.value
}

func (c ListOffSiteInstancesRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListOffSiteInstancesRequestXLanguage) UnmarshalJSON(b []byte) error {
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
