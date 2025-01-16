package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00079 = "SQLE00079"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00079,
			Desc:       "别名不建议与表或列的名字相同",
			Annotation: "表或列的别名与其真实名称相同, 这样的别名会使得查询更难去分辨",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "别名不建议与表或列的名字相同",
		AllowOffline: true,
		Func:         RuleSQLE00079,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00079): "在 MySQL 中，别名不建议与表或列的名字相同."
您应遵循以下逻辑：
1. 对于所有DML、CTE语中句含有SELECT语法节点，则：
  1. 创建2个集合，集合A收集sql中表、列名，集合B收集sql中所有的别名
  2. 判断集合B中所有别名是否有重复，如果有，则报告违反规则
  3. 判断集合A和集合B是否有相同名称，如果有，则报告违反规则
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00079(input *rulepkg.RuleHandlerInput) error {
	var originNameList []string // 记录表名和列名
	var asNamelist []string     // 记录别名

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UpdateStmt:
		// 对于 SELECT 和 UNION 语句，获取所有的 SELECT 子句
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil {
				tableSources := util.GetTableSourcesFromJoin(selectStmt.From.TableRefs)
				for _, tableSource := range tableSources {
					if tableName, ok := tableSource.Source.(*ast.TableName); ok {
						originNameList = append(originNameList, tableName.Name.L)
						if tableSource.AsName.L != "" {
							asNamelist = append(asNamelist, tableSource.AsName.L)
						}
					}
				}
			}

			// 获取列别名和表别名
			if selectStmt.Fields != nil {
				for _, field := range selectStmt.Fields.Fields {
					if selectColumn, ok := field.Expr.(*ast.ColumnNameExpr); ok && selectColumn.Name.Name.L != "" {
						originNameList = append(originNameList, selectColumn.Name.Name.L)
						if field.AsName.L != "" {
							asNamelist = append(asNamelist, field.AsName.L)
						}
					}
				}
			}
		}
		// 对于"WITH..."语句
		// TODO 待实现
	}

	// 1、别名之间是否有重复的
	if util.HasDuplicateInStrings(asNamelist) {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00079)
		return nil

	}
	// 2、别名与 表/列名是否有重复的
	if util.HasDuplicateIn2Strings(originNameList, asNamelist) {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00079)
		return nil
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
