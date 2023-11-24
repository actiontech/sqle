package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ShowSecondLevelMonitoringResponse Response Object
type ShowSecondLevelMonitoringResponse struct {

	// 秒级监控开关
	SwitchOption *bool `json:"switch_option,omitempty"`

	// 监控间隔, 支持1秒和5秒
	Interval       *ShowSecondLevelMonitoringResponseInterval `json:"interval,omitempty"`
	HttpStatusCode int                                        `json:"-"`
}

func (o ShowSecondLevelMonitoringResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowSecondLevelMonitoringResponse struct{}"
	}

	return strings.Join([]string{"ShowSecondLevelMonitoringResponse", string(data)}, " ")
}

type ShowSecondLevelMonitoringResponseInterval struct {
	value int32
}

type ShowSecondLevelMonitoringResponseIntervalEnum struct {
	E_1 ShowSecondLevelMonitoringResponseInterval
	E_5 ShowSecondLevelMonitoringResponseInterval
}

func GetShowSecondLevelMonitoringResponseIntervalEnum() ShowSecondLevelMonitoringResponseIntervalEnum {
	return ShowSecondLevelMonitoringResponseIntervalEnum{
		E_1: ShowSecondLevelMonitoringResponseInterval{
			value: 1,
		}, E_5: ShowSecondLevelMonitoringResponseInterval{
			value: 5,
		},
	}
}

func (c ShowSecondLevelMonitoringResponseInterval) Value() int32 {
	return c.value
}

func (c ShowSecondLevelMonitoringResponseInterval) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowSecondLevelMonitoringResponseInterval) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: int32")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(int32); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to int32 error")
	}
}
