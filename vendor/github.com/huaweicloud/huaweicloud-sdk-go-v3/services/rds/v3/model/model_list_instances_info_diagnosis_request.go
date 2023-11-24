package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstancesInfoDiagnosisRequest Request Object
type ListInstancesInfoDiagnosisRequest struct {

	// 引擎类型
	Engine ListInstancesInfoDiagnosisRequestEngine `json:"engine"`

	// 诊断项
	Diagnosis ListInstancesInfoDiagnosisRequestDiagnosis `json:"diagnosis"`

	// offset
	Offset *int32 `json:"offset,omitempty"`

	// limit
	Limit *int32 `json:"limit,omitempty"`
}

func (o ListInstancesInfoDiagnosisRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesInfoDiagnosisRequest struct{}"
	}

	return strings.Join([]string{"ListInstancesInfoDiagnosisRequest", string(data)}, " ")
}

type ListInstancesInfoDiagnosisRequestEngine struct {
	value string
}

type ListInstancesInfoDiagnosisRequestEngineEnum struct {
	MYSQL      ListInstancesInfoDiagnosisRequestEngine
	POSTGRESQL ListInstancesInfoDiagnosisRequestEngine
	SQLSERVER  ListInstancesInfoDiagnosisRequestEngine
}

func GetListInstancesInfoDiagnosisRequestEngineEnum() ListInstancesInfoDiagnosisRequestEngineEnum {
	return ListInstancesInfoDiagnosisRequestEngineEnum{
		MYSQL: ListInstancesInfoDiagnosisRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: ListInstancesInfoDiagnosisRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: ListInstancesInfoDiagnosisRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c ListInstancesInfoDiagnosisRequestEngine) Value() string {
	return c.value
}

func (c ListInstancesInfoDiagnosisRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesInfoDiagnosisRequestEngine) UnmarshalJSON(b []byte) error {
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

type ListInstancesInfoDiagnosisRequestDiagnosis struct {
	value string
}

type ListInstancesInfoDiagnosisRequestDiagnosisEnum struct {
	HIGH_PRESSURE         ListInstancesInfoDiagnosisRequestDiagnosis
	LOCK_WAIT             ListInstancesInfoDiagnosisRequestDiagnosis
	INSUFFICIENT_CAPACITY ListInstancesInfoDiagnosisRequestDiagnosis
	SLOW_SQL_FREQUENCY    ListInstancesInfoDiagnosisRequestDiagnosis
	DISK_PERFORMANCE_CAP  ListInstancesInfoDiagnosisRequestDiagnosis
	MEM_OVERRUN           ListInstancesInfoDiagnosisRequestDiagnosis
	AGE_EXCEED            ListInstancesInfoDiagnosisRequestDiagnosis
	CONNECTIONS_EXCEED    ListInstancesInfoDiagnosisRequestDiagnosis
}

func GetListInstancesInfoDiagnosisRequestDiagnosisEnum() ListInstancesInfoDiagnosisRequestDiagnosisEnum {
	return ListInstancesInfoDiagnosisRequestDiagnosisEnum{
		HIGH_PRESSURE: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "high_pressure",
		},
		LOCK_WAIT: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "lock_wait",
		},
		INSUFFICIENT_CAPACITY: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "insufficient_capacity",
		},
		SLOW_SQL_FREQUENCY: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "slow_sql_frequency",
		},
		DISK_PERFORMANCE_CAP: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "disk_performance_cap",
		},
		MEM_OVERRUN: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "mem_overrun",
		},
		AGE_EXCEED: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "age_exceed",
		},
		CONNECTIONS_EXCEED: ListInstancesInfoDiagnosisRequestDiagnosis{
			value: "connections_exceed",
		},
	}
}

func (c ListInstancesInfoDiagnosisRequestDiagnosis) Value() string {
	return c.value
}

func (c ListInstancesInfoDiagnosisRequestDiagnosis) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesInfoDiagnosisRequestDiagnosis) UnmarshalJSON(b []byte) error {
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
