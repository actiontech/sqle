package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00174 = "SQLE00174"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00174,
			Desc:       "在 MySQL 中, 禁止GRANT 授予过高权限",
			Annotation: "授予过高权限，可能会带来严重的安全风险。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "ALL,SUPER,WITH GRANT OPTION",
					Desc:  "高权限范围",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 禁止GRANT 授予过高权限",
		AllowOffline: true,
		Func:         RuleSQLE00174,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00174): "在 MySQL 中，禁止GRANT 授予过高权限.默认参数描述: 高权限范围, 默认参数值: ALL,SUPER,WITH GRANT OPTION"
您应遵循以下逻辑：
1. 对于 "GRANT ..." 语句：
   1. 提取语句中的权限（如 SELECT、INSERT、UPDATE、DELETE、CREATE、ALTER、DROP、SUPER、WITH GRANT OPTION 等），并将这些权限存储在一个集合中。
   2. 将集合中的权限与规则中定义的高权限列表进行对比，如果存在交集，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====

// 规则函数实现开始
func RuleSQLE00174(input *rulepkg.RuleHandlerInput) error {
	// 检查输入的SQL语句是否为GRANT语句
	grantStmt, ok := input.Node.(*ast.GrantStmt)
	if !ok {
		// 不是GRANT语句，忽略
		return nil
	}

	if grantStmt.WithGrant { // WITH GRANT OPTION
		rulepkg.AddResult(input.Res, input.Rule, SQLE00174)
		return nil
	}
	// 提取GRANT语句中的权限
	for _, priv := range grantStmt.Privs {
		if priv.Priv == mysql.AllPriv || priv.Priv == mysql.SuperPriv { // ALL、SUPER
			rulepkg.AddResult(input.Res, input.Rule, SQLE00174)
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
