package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListInstancesInfoDiagnosisResponse Response Object
type ListInstancesInfoDiagnosisResponse struct {

	// 诊断项
	Diagnosis *ListInstancesInfoDiagnosisResponseDiagnosis `json:"diagnosis,omitempty"`

	// 实例数量
	TotalCount *int32 `json:"total_count,omitempty"`

	// 实例信息
	Instances      *[]DiagnosisInstancesInfoResult `json:"instances,omitempty"`
	HttpStatusCode int                             `json:"-"`
}

func (o ListInstancesInfoDiagnosisResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListInstancesInfoDiagnosisResponse struct{}"
	}

	return strings.Join([]string{"ListInstancesInfoDiagnosisResponse", string(data)}, " ")
}

type ListInstancesInfoDiagnosisResponseDiagnosis struct {
	value string
}

type ListInstancesInfoDiagnosisResponseDiagnosisEnum struct {
	HIGH_PRESSURE         ListInstancesInfoDiagnosisResponseDiagnosis
	LOCK_WAIT             ListInstancesInfoDiagnosisResponseDiagnosis
	INSUFFICIENT_CAPACITY ListInstancesInfoDiagnosisResponseDiagnosis
	SLOW_SQL_FREQUENCY    ListInstancesInfoDiagnosisResponseDiagnosis
	DISK_PERFORMANCE_CAP  ListInstancesInfoDiagnosisResponseDiagnosis
	MEM_OVERRUN           ListInstancesInfoDiagnosisResponseDiagnosis
	AGE_EXCEED            ListInstancesInfoDiagnosisResponseDiagnosis
	CONNECTIONS_EXCEED    ListInstancesInfoDiagnosisResponseDiagnosis
}

func GetListInstancesInfoDiagnosisResponseDiagnosisEnum() ListInstancesInfoDiagnosisResponseDiagnosisEnum {
	return ListInstancesInfoDiagnosisResponseDiagnosisEnum{
		HIGH_PRESSURE: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "high_pressure",
		},
		LOCK_WAIT: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "lock_wait",
		},
		INSUFFICIENT_CAPACITY: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "insufficient_capacity",
		},
		SLOW_SQL_FREQUENCY: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "slow_sql_frequency",
		},
		DISK_PERFORMANCE_CAP: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "disk_performance_cap",
		},
		MEM_OVERRUN: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "mem_overrun",
		},
		AGE_EXCEED: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "age_exceed",
		},
		CONNECTIONS_EXCEED: ListInstancesInfoDiagnosisResponseDiagnosis{
			value: "connections_exceed",
		},
	}
}

func (c ListInstancesInfoDiagnosisResponseDiagnosis) Value() string {
	return c.value
}

func (c ListInstancesInfoDiagnosisResponseDiagnosis) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListInstancesInfoDiagnosisResponseDiagnosis) UnmarshalJSON(b []byte) error {
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
