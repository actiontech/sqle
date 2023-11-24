package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ChargeInfo 计费类型信息，支持包年包月和按需，默认为按需。
type ChargeInfo struct {

	// 计费模式。  取值范围：  - prePaid：预付费，即包年/包月。 - postPaid：后付费，即按需付费。
	ChargeMode ChargeInfoChargeMode `json:"charge_mode"`

	// 订购周期类型。  取值范围：  - month：包月。 - year：包年。
	PeriodType *ChargeInfoPeriodType `json:"period_type,omitempty"`

	// “charge_mode”为“prePaid”时生效，且为必选值，指定订购的时间。  取值范围：  当“period_type”为“month”时，取值为1~9。 当“period_type”为“year”时，取值为1~3。
	PeriodNum *int32 `json:"period_num,omitempty"`

	// 创建包周期实例时可指定，表示是否自动续订，续订的周期和原周期相同，且续订时会自动支付。  - true，为自动续订。 - false，为不自动续订，默认该方式。
	IsAutoRenew *bool `json:"is_auto_renew,omitempty"`

	// 创建包周期时可指定，表示是否自动从客户的账户中支付，此字段不影响自动续订的支付方式。  - true，为自动支付。 - false，为手动支付，默认该方式。
	IsAutoPay *bool `json:"is_auto_pay,omitempty"`
}

func (o ChargeInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ChargeInfo struct{}"
	}

	return strings.Join([]string{"ChargeInfo", string(data)}, " ")
}

type ChargeInfoChargeMode struct {
	value string
}

type ChargeInfoChargeModeEnum struct {
	PRE_PAID  ChargeInfoChargeMode
	POST_PAID ChargeInfoChargeMode
}

func GetChargeInfoChargeModeEnum() ChargeInfoChargeModeEnum {
	return ChargeInfoChargeModeEnum{
		PRE_PAID: ChargeInfoChargeMode{
			value: "prePaid",
		},
		POST_PAID: ChargeInfoChargeMode{
			value: "postPaid",
		},
	}
}

func (c ChargeInfoChargeMode) Value() string {
	return c.value
}

func (c ChargeInfoChargeMode) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ChargeInfoChargeMode) UnmarshalJSON(b []byte) error {
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

type ChargeInfoPeriodType struct {
	value string
}

type ChargeInfoPeriodTypeEnum struct {
	MONTH ChargeInfoPeriodType
	YEAR  ChargeInfoPeriodType
}

func GetChargeInfoPeriodTypeEnum() ChargeInfoPeriodTypeEnum {
	return ChargeInfoPeriodTypeEnum{
		MONTH: ChargeInfoPeriodType{
			value: "month",
		},
		YEAR: ChargeInfoPeriodType{
			value: "year",
		},
	}
}

func (c ChargeInfoPeriodType) Value() string {
	return c.value
}

func (c ChargeInfoPeriodType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ChargeInfoPeriodType) UnmarshalJSON(b []byte) error {
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
