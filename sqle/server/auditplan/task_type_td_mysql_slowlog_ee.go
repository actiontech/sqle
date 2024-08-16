//go:build enterprise
// +build enterprise

package auditplan

import "github.com/actiontech/sqle/sqle/pkg/params"

type TDMySQLSlowLogTaskV2 struct {
	SlowLogTaskV2
}

func NewTDMySQLSlowLogTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &TDMySQLSlowLogTaskV2{}
	}
}

func (at *TDMySQLSlowLogTaskV2) InstanceType() string {
	return InstanceTypeTDSQL
}

func (at *TDMySQLSlowLogTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{}
}

func (at *TDMySQLSlowLogTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{}
}
