package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListStorageTypesRequest Request Object
type ListStorageTypesRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 数据库引擎。支持的引擎如下，不区分大小写： MySQL PostgreSQL SQLServer
	DatabaseName ListStorageTypesRequestDatabaseName `json:"database_name"`

	// 数据库版本号。
	VersionName string `json:"version_name"`

	// 主备模式： single：单机模式。 ha：主备模式。 replica：只读模式。
	HaMode *ListStorageTypesRequestHaMode `json:"ha_mode,omitempty"`
}

func (o ListStorageTypesRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListStorageTypesRequest struct{}"
	}

	return strings.Join([]string{"ListStorageTypesRequest", string(data)}, " ")
}

type ListStorageTypesRequestDatabaseName struct {
	value string
}

type ListStorageTypesRequestDatabaseNameEnum struct {
	MY_SQL      ListStorageTypesRequestDatabaseName
	POSTGRE_SQL ListStorageTypesRequestDatabaseName
	SQL_SERVER  ListStorageTypesRequestDatabaseName
	MARIA_DB    ListStorageTypesRequestDatabaseName
}

func GetListStorageTypesRequestDatabaseNameEnum() ListStorageTypesRequestDatabaseNameEnum {
	return ListStorageTypesRequestDatabaseNameEnum{
		MY_SQL: ListStorageTypesRequestDatabaseName{
			value: "MySQL",
		},
		POSTGRE_SQL: ListStorageTypesRequestDatabaseName{
			value: "PostgreSQL",
		},
		SQL_SERVER: ListStorageTypesRequestDatabaseName{
			value: "SQLServer",
		},
		MARIA_DB: ListStorageTypesRequestDatabaseName{
			value: "MariaDB",
		},
	}
}

func (c ListStorageTypesRequestDatabaseName) Value() string {
	return c.value
}

func (c ListStorageTypesRequestDatabaseName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListStorageTypesRequestDatabaseName) UnmarshalJSON(b []byte) error {
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

type ListStorageTypesRequestHaMode struct {
	value string
}

type ListStorageTypesRequestHaModeEnum struct {
	HA      ListStorageTypesRequestHaMode
	SINGLE  ListStorageTypesRequestHaMode
	REPLICA ListStorageTypesRequestHaMode
}

func GetListStorageTypesRequestHaModeEnum() ListStorageTypesRequestHaModeEnum {
	return ListStorageTypesRequestHaModeEnum{
		HA: ListStorageTypesRequestHaMode{
			value: "ha",
		},
		SINGLE: ListStorageTypesRequestHaMode{
			value: "single",
		},
		REPLICA: ListStorageTypesRequestHaMode{
			value: "replica",
		},
	}
}

func (c ListStorageTypesRequestHaMode) Value() string {
	return c.value
}

func (c ListStorageTypesRequestHaMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListStorageTypesRequestHaMode) UnmarshalJSON(b []byte) error {
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
