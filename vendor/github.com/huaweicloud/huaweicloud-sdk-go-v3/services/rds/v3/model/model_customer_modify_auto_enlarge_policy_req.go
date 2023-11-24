package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type CustomerModifyAutoEnlargePolicyReq struct {

	// 是否开启自动扩容,true为开启,false为关闭
	SwitchOption bool `json:"switch_option"`

	// 扩容上限，单位GB, 取值范围40~4000，需要大于等于实例当前存储空间总大小，switch_option为true必填
	LimitSize *int32 `json:"limit_size,omitempty"`

	// 可用存储空间百分比，小于等于此值或者10GB时触发扩容，switch_option为true时必填
	TriggerThreshold *CustomerModifyAutoEnlargePolicyReqTriggerThreshold `json:"trigger_threshold,omitempty"`
}

func (o CustomerModifyAutoEnlargePolicyReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CustomerModifyAutoEnlargePolicyReq struct{}"
	}

	return strings.Join([]string{"CustomerModifyAutoEnlargePolicyReq", string(data)}, " ")
}

type CustomerModifyAutoEnlargePolicyReqTriggerThreshold struct {
	value int32
}

type CustomerModifyAutoEnlargePolicyReqTriggerThresholdEnum struct {
	E_10 CustomerModifyAutoEnlargePolicyReqTriggerThreshold
	E_15 CustomerModifyAutoEnlargePolicyReqTriggerThreshold
	E_20 CustomerModifyAutoEnlargePolicyReqTriggerThreshold
}

func GetCustomerModifyAutoEnlargePolicyReqTriggerThresholdEnum() CustomerModifyAutoEnlargePolicyReqTriggerThresholdEnum {
	return CustomerModifyAutoEnlargePolicyReqTriggerThresholdEnum{
		E_10: CustomerModifyAutoEnlargePolicyReqTriggerThreshold{
			value: 10,
		}, E_15: CustomerModifyAutoEnlargePolicyReqTriggerThreshold{
			value: 15,
		}, E_20: CustomerModifyAutoEnlargePolicyReqTriggerThreshold{
			value: 20,
		},
	}
}

func (c CustomerModifyAutoEnlargePolicyReqTriggerThreshold) Value() int32 {
	return c.value
}

func (c CustomerModifyAutoEnlargePolicyReqTriggerThreshold) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CustomerModifyAutoEnlargePolicyReqTriggerThreshold) UnmarshalJSON(b []byte) error {
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
