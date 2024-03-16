package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstancesRequest Request Object
type ListInstancesRequest struct {
	ContentType *string `json:"Content-Type,omitempty"`

	// 语言
	XLanguage *ListInstancesRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。  “\\*”为系统保留字符，如果id是以“\\*”起始，表示按照\\*后面的值模糊匹配，否则，按照id精确匹配查询。不能只传入“\\*”。
	Id *string `json:"id,omitempty"`

	// 实例名称。  “\\*”为系统保留字符，如果name是以“\\*”起始，表示按照\\*后面的值模糊匹配，否则，按照name精确匹配查询。不能只传入“\\*”。
	Name *string `json:"name,omitempty"`

	// 按照实例类型查询。取值Single、Ha、Replica，分别对应于单实例、主备实例和只读实例。
	Type *ListInstancesRequestType `json:"type,omitempty"`

	// 数据库类型，区分大小写。 - MySQL - PostgreSQL - SQLServer - MariaDB
	DatastoreType *ListInstancesRequestDatastoreType `json:"datastore_type,omitempty"`

	// 虚拟私有云ID。
	VpcId *string `json:"vpc_id,omitempty"`

	// 子网ID。
	SubnetId *string `json:"subnet_id,omitempty"`

	// 索引位置，偏移量。从第一条数据偏移offset条数据后开始查询，默认为0（偏移0条数据，表示从第一条数据开始查询），必须为数字，不能为负数。
	Offset *int32 `json:"offset,omitempty"`

	// 查询记录数。默认为100，不能为负数，最小值为1，最大值为100。
	Limit *int32 `json:"limit,omitempty"`

	// 根据实例标签键值对进行查询。 {key}表示标签键，不可以为空或重复。最大长度127个unicode字符。key不能为空或者空字符串，不能为空格，使用之前先trim前后半角空格。不能包含+/?#&=,%特殊字符。 {value}表示标签值，可以为空。最大长度255个unicode字符，使用之前先trim 前后半角空格。不能包含+/?#&=,%特殊字符。如果value为空，则表示any_value（查询任意value）。 如果同时使用多个标签键值对进行查询，中间使用逗号分隔开，最多包含10组。
	Tags *string `json:"tags,omitempty"`
}

func (o ListInstancesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesRequest struct{}"
	}

	return strings.Join([]string{"ListInstancesRequest", string(data)}, " ")
}

type ListInstancesRequestXLanguage struct {
	value string
}

type ListInstancesRequestXLanguageEnum struct {
	ZH_CN ListInstancesRequestXLanguage
	EN_US ListInstancesRequestXLanguage
}

func GetListInstancesRequestXLanguageEnum() ListInstancesRequestXLanguageEnum {
	return ListInstancesRequestXLanguageEnum{
		ZH_CN: ListInstancesRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListInstancesRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListInstancesRequestXLanguage) Value() string {
	return c.value
}

func (c ListInstancesRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesRequestXLanguage) UnmarshalJSON(b []byte) error {
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

type ListInstancesRequestType struct {
	value string
}

type ListInstancesRequestTypeEnum struct {
	SINGLE  ListInstancesRequestType
	HA      ListInstancesRequestType
	REPLICA ListInstancesRequestType
}

func GetListInstancesRequestTypeEnum() ListInstancesRequestTypeEnum {
	return ListInstancesRequestTypeEnum{
		SINGLE: ListInstancesRequestType{
			value: "Single",
		},
		HA: ListInstancesRequestType{
			value: "Ha",
		},
		REPLICA: ListInstancesRequestType{
			value: "Replica",
		},
	}
}

func (c ListInstancesRequestType) Value() string {
	return c.value
}

func (c ListInstancesRequestType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesRequestType) UnmarshalJSON(b []byte) error {
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

type ListInstancesRequestDatastoreType struct {
	value string
}

type ListInstancesRequestDatastoreTypeEnum struct {
	MY_SQL      ListInstancesRequestDatastoreType
	POSTGRE_SQL ListInstancesRequestDatastoreType
	SQL_SERVER  ListInstancesRequestDatastoreType
	MARIA_DB    ListInstancesRequestDatastoreType
}

func GetListInstancesRequestDatastoreTypeEnum() ListInstancesRequestDatastoreTypeEnum {
	return ListInstancesRequestDatastoreTypeEnum{
		MY_SQL: ListInstancesRequestDatastoreType{
			value: "MySQL",
		},
		POSTGRE_SQL: ListInstancesRequestDatastoreType{
			value: "PostgreSQL",
		},
		SQL_SERVER: ListInstancesRequestDatastoreType{
			value: "SQLServer",
		},
		MARIA_DB: ListInstancesRequestDatastoreType{
			value: "MariaDB",
		},
	}
}

func (c ListInstancesRequestDatastoreType) Value() string {
	return c.value
}

func (c ListInstancesRequestDatastoreType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesRequestDatastoreType) UnmarshalJSON(b []byte) error {
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
