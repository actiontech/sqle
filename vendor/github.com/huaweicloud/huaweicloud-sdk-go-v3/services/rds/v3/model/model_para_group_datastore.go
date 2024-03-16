package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type ParaGroupDatastore struct {

	// 数据库引擎，不区分大小写： - MySQL - PostgreSQL - SQLServer - MariaDB
	Type ParaGroupDatastoreType `json:"type"`

	// 数据库版本。  - MySQL引擎支持5.6、5.7、8.0版本。取值示例：5.7。具有相应权限的用户才可使用8.0，您可联系华为云客服人员申请。 - PostgreSQL引擎支持9.5、9.6、10、11版本。取值示例：9.6。 - Microsoft SQL Server：仅支持2017 企业版、2017 标准版、2017 web版、2014 标准版、2014 企业版、2016 标准版、2016 企业版、2012 企业版、2012 标准版、2012 web版、2008 R2 企业版、2008 R2 web版、2014 web版、2016 web版。取值示例2014_SE。 例如：2017标准版可填写2017_SE，2017企业版可填写2017_EE，2017web版可以填写2017_WEB
	Version string `json:"version"`
}

func (o ParaGroupDatastore) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ParaGroupDatastore struct{}"
	}

	return strings.Join([]string{"ParaGroupDatastore", string(data)}, " ")
}

type ParaGroupDatastoreType struct {
	value string
}

type ParaGroupDatastoreTypeEnum struct {
	MY_SQL      ParaGroupDatastoreType
	POSTGRE_SQL ParaGroupDatastoreType
	SQL_SERVER  ParaGroupDatastoreType
	MARIA_DB    ParaGroupDatastoreType
}

func GetParaGroupDatastoreTypeEnum() ParaGroupDatastoreTypeEnum {
	return ParaGroupDatastoreTypeEnum{
		MY_SQL: ParaGroupDatastoreType{
			value: "MySQL",
		},
		POSTGRE_SQL: ParaGroupDatastoreType{
			value: "PostgreSQL",
		},
		SQL_SERVER: ParaGroupDatastoreType{
			value: "SQLServer",
		},
		MARIA_DB: ParaGroupDatastoreType{
			value: "MariaDB",
		},
	}
}

func (c ParaGroupDatastoreType) Value() string {
	return c.value
}

func (c ParaGroupDatastoreType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ParaGroupDatastoreType) UnmarshalJSON(b []byte) error {
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
