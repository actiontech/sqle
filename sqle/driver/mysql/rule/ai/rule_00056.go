package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00056 = "SQLE00056"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00056,
			Desc:       "在 MySQL 中, 表建议使用指定的字符集",
			Annotation: "数据库内使用非标准的字符集，可能导致字符无法编码或者编码不全引起的乱码，最终出现应用写入数据失败或者查询结果显示乱码，影响数据库服务可用性。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "UTF8MB4",
					Desc:  "标准字符集",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 表建议使用指定的字符集: %v",
		AllowOffline: true,
		Func:         RuleSQLE00056,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00056): "在 MySQL 中，表建议使用指定的字符集.默认参数描述: 标准字符集, 默认参数值: UTF8MB4"
您应遵循以下逻辑：
1. 对于 "CREATE TABLE ..."语句，
   1. 存在 charset或CHARACTER 语法节点，
   2. 且charset或CHARACTER 的值与阈值不一致时，则报告违反规则。
2. 对于"ALTER TABLE..." 语句，执行上述相同的检查。
3. 对于 "CREATE DATABASE ..."语句，执行上述相同的检查。
4. 对于 "ALTER DATABASE ..."语句，执行上述相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00056(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	// 获取规则指定的标准字符集
	expectedCharset := param.String()

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 检查 CREATE TABLE 语句中的字符集
		if option := util.GetTableOption(stmt.Options, ast.TableOptionCharset); option != nil {
			if !strings.EqualFold(option.StrValue, expectedCharset) {
				rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name, expectedCharset)
			}
		}
	case *ast.AlterTableStmt:
		// 检查 ALTER TABLE 语句中的字符集
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableOption) {
			for _, option := range spec.Options {
				if option.Tp == ast.TableOptionCharset {
					if !strings.EqualFold(option.StrValue, expectedCharset) {
						rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name, expectedCharset)
						break
					}
				}
			}
		}
	case *ast.CreateDatabaseStmt:
		// 检查 CREATE DATABASE 语句中的字符集
		if option := util.GetDatabaseOption(stmt.Options, ast.DatabaseOptionCharset); option != nil {
			if !strings.EqualFold(option.Value, expectedCharset) {
				rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name, expectedCharset)
			}
		}
	case *ast.AlterDatabaseStmt:
		// 检查 ALTER DATABASE 语句中的字符集
		if option := util.GetDatabaseOption(stmt.Options, ast.DatabaseOptionCharset); option != nil {
			if !strings.EqualFold(option.Value, expectedCharset) {
				rulepkg.AddResult(input.Res, input.Rule, input.Rule.Name, expectedCharset)
			}
		}
	default:
		return nil
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====