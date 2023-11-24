package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListDatastoresRequest Request Object
type ListDatastoresRequest struct {

	// 语言
	XLanguage *string `json:"X-Language,omitempty"`

	// 数据库引擎。支持的引擎如下，不区分大小写： MySQL PostgreSQL SQLServer
	DatabaseName ListDatastoresRequestDatabaseName `json:"database_name"`
}

func (o ListDatastoresRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListDatastoresRequest struct{}"
	}

	return strings.Join([]string{"ListDatastoresRequest", string(data)}, " ")
}

type ListDatastoresRequestDatabaseName struct {
	value string
}

type ListDatastoresRequestDatabaseNameEnum struct {
	MY_SQL      ListDatastoresRequestDatabaseName
	POSTGRE_SQL ListDatastoresRequestDatabaseName
	SQL_SERVER  ListDatastoresRequestDatabaseName
	MARIA_DB    ListDatastoresRequestDatabaseName
}

func GetListDatastoresRequestDatabaseNameEnum() ListDatastoresRequestDatabaseNameEnum {
	return ListDatastoresRequestDatabaseNameEnum{
		MY_SQL: ListDatastoresRequestDatabaseName{
			value: "MySQL",
		},
		POSTGRE_SQL: ListDatastoresRequestDatabaseName{
			value: "PostgreSQL",
		},
		SQL_SERVER: ListDatastoresRequestDatabaseName{
			value: "SQLServer",
		},
		MARIA_DB: ListDatastoresRequestDatabaseName{
			value: "MariaDB",
		},
	}
}

func (c ListDatastoresRequestDatabaseName) Value() string {
	return c.value
}

func (c ListDatastoresRequestDatabaseName) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListDatastoresRequestDatabaseName) UnmarshalJSON(b []byte) error {
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
