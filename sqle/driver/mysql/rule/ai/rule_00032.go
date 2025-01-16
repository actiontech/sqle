package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00032 = "SQLE00032"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00032,
			Desc:       "数据库名称必须使用固定后缀结尾",
			Annotation: "通过配置该规则可以规范指定业务的数据库命名规则，具体命名规范可以自定义设置。",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeNamingConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "_DB",
					Desc:  "固定后缀",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "数据库名称必须使用固定后缀结尾:%s",
		AllowOffline: true,
		Func:         RuleSQLE00032,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00032): "在 MySQL 中，数据库名称必须使用固定后缀结尾.固定后缀:_DB"
您应遵循以下逻辑：
1、检查CREATE句子中是否存在DATABASE语法节点，如果存在，则进入下一步检查。
2、检查数据库对象名是否遵从固定后缀要求，如果不遵从，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00032(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	// 判断是否为创建数据库的语句
	switch stmt := input.Node.(type) {
	case *ast.CreateDatabaseStmt:
		suffix := param.String()
		// 检查数据库对象名是否遵从固定后缀要求
		if !strings.HasSuffix(stmt.Name, suffix) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00032, suffix)
		}
	}
	return nil
}

// ==== Rule code end ====
