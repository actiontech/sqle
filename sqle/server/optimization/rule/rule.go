package optimization

import (
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
)

var OptimizationRuleMap = map[string][]OptimizationRuleHandler{} // ruleCode与plugin重写规则的映射关系

type OptimizationRuleHandler struct {
	Rule     driverV2.Rule
	RuleCode string // SQL优化规则的ruleCode
}

func InitOptimizationRule() {
	initOptimizationRule()
}

// CanOptimizeDbType SQL优化是否支持该数据源类型
func CanOptimizeDbType(dt string) bool {
	_, exist := OptimizationRuleMap[dt]
	return exist
}
