package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpdatePostgresqlDatabaseRequest Request Object
type UpdatePostgresqlDatabaseRequest struct {

	// 语言
	XLanguage *UpdatePostgresqlDatabaseRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *UpdateDatabaseReq `json:"body,omitempty"`
}

func (o UpdatePostgresqlDatabaseRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdatePostgresqlDatabaseRequest struct{}"
	}

	return strings.Join([]string{"UpdatePostgresqlDatabaseRequest", string(data)}, " ")
}

type UpdatePostgresqlDatabaseRequestXLanguage struct {
	value string
}

type UpdatePostgresqlDatabaseRequestXLanguageEnum struct {
	ZH_CN UpdatePostgresqlDatabaseRequestXLanguage
	EN_US UpdatePostgresqlDatabaseRequestXLanguage
}

func GetUpdatePostgresqlDatabaseRequestXLanguageEnum() UpdatePostgresqlDatabaseRequestXLanguageEnum {
	return UpdatePostgresqlDatabaseRequestXLanguageEnum{
		ZH_CN: UpdatePostgresqlDatabaseRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: UpdatePostgresqlDatabaseRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c UpdatePostgresqlDatabaseRequestXLanguage) Value() string {
	return c.value
}

func (c UpdatePostgresqlDatabaseRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdatePostgresqlDatabaseRequestXLanguage) UnmarshalJSON(b []byte) error {
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
