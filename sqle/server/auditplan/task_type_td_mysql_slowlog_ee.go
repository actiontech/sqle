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

func (at *TDMySQLSlowLogTaskV2) Params() func(instanceId ...string) params.Params {
	return func(instanceId ...string) params.Params {
		return []*params.Param{
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		}
	}
}
