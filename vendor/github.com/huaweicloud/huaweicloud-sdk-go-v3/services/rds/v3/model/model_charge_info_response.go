package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ChargeInfoResponse 计费类型信息，支持包年包月和按需，默认为按需。
type ChargeInfoResponse struct {

	// 计费模式。  取值范围：  - prePaid：预付费，即包年/包月。 - postPaid：后付费，即按需付费。
	ChargeMode ChargeInfoResponseChargeMode `json:"charge_mode"`
}

func (o ChargeInfoResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChargeInfoResponse struct{}"
	}

	return strings.Join([]string{"ChargeInfoResponse", string(data)}, " ")
}

type ChargeInfoResponseChargeMode struct {
	value string
}

type ChargeInfoResponseChargeModeEnum struct {
	PRE_PAID  ChargeInfoResponseChargeMode
	POST_PAID ChargeInfoResponseChargeMode
}

func GetChargeInfoResponseChargeModeEnum() ChargeInfoResponseChargeModeEnum {
	return ChargeInfoResponseChargeModeEnum{
		PRE_PAID: ChargeInfoResponseChargeMode{
			value: "prePaid",
		},
		POST_PAID: ChargeInfoResponseChargeMode{
			value: "postPaid",
		},
	}
}

func (c ChargeInfoResponseChargeMode) Value() string {
	return c.value
}

func (c ChargeInfoResponseChargeMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ChargeInfoResponseChargeMode) UnmarshalJSON(b []byte) error {
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
