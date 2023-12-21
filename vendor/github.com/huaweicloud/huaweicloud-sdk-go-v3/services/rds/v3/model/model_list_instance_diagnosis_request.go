package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstanceDiagnosisRequest Request Object
type ListInstanceDiagnosisRequest struct {

	// 引擎类型
	Engine ListInstanceDiagnosisRequestEngine `json:"engine"`
}

func (o ListInstanceDiagnosisRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstanceDiagnosisRequest struct{}"
	}

	return strings.Join([]string{"ListInstanceDiagnosisRequest", string(data)}, " ")
}

type ListInstanceDiagnosisRequestEngine struct {
	value string
}

type ListInstanceDiagnosisRequestEngineEnum struct {
	MYSQL      ListInstanceDiagnosisRequestEngine
	POSTGRESQL ListInstanceDiagnosisRequestEngine
	SQLSERVER  ListInstanceDiagnosisRequestEngine
}

func GetListInstanceDiagnosisRequestEngineEnum() ListInstanceDiagnosisRequestEngineEnum {
	return ListInstanceDiagnosisRequestEngineEnum{
		MYSQL: ListInstanceDiagnosisRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: ListInstanceDiagnosisRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: ListInstanceDiagnosisRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c ListInstanceDiagnosisRequestEngine) Value() string {
	return c.value
}

func (c ListInstanceDiagnosisRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstanceDiagnosisRequestEngine) UnmarshalJSON(b []byte) error {
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
