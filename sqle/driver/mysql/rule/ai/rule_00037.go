package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00037 = "SQLE00037"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00037,
			Desc:       "避免一张表内二级索引的个数过多",
			Annotation: "在表上建立的每个索引都会增加存储开销，索引对于插入、删除、更新操作也会增加维护索引处理上的开销（TPS），且太多与不充分、不正确的索引对性能都毫无益处。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "二级索引个数",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "避免一张表内二级索引的个数过多",
		AllowOffline: false,
		Func:         RuleSQLE00037,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00037): "在 MySQL 中，避免一张表内二级索引的个数过多.默认参数描述: 二级索引个数, 默认参数值: 5"
您应遵循以下逻辑：
1. 对于"CREATE TABLE..." 语句，
    1. 初始化一个索引计数变量。
    2. 解析语法树，统计语句中包含的索引定义节点的数量，并更新到索引计数变量。
    3. 检查索引计数变量的值，如果超过允许的最大索引数，则报告违反规则。

2. 对于"ALTER TABLE..." 语句，
    1. 初始化一个索引计数变量。
    2. 解析语法树，统计语句中新增索引定义节点的数量，并更新到索引计数变量。
    3. 使用辅助函数GetCreateTableStmt获取目标表上现有的二级索引数，并加到索引计数变量。
    4. 检查索引计数变量的值，如果超过允许的最大索引数，则报告违反规则。

3. 对于"CREATE INDEX ..." 语句，
    1. 初始化一个索引计数变量。
    2. 解析语法树，统计语句中新增的索引数量，并更新到索引计数变量。
    3. 使用辅助函数GetCreateTableStmt获取目标表上现有的二级索引数，并加到索引计数变量。
    4. 检查索引计数变量的值，如果超过允许的最大索引数，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00037(input *rulepkg.RuleHandlerInput) error {
	// 获取数值类型的规则参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	maxIndexCount := param.Int()
	if maxIndexCount <= 0 {
		return fmt.Errorf("param value should be greater than 0")
	}

	indexCount := 0
	secondaryIndexes := []ast.ConstraintType{
		ast.ConstraintKey,
		ast.ConstraintIndex,
		ast.ConstraintUniq,
		ast.ConstraintUniqKey,
		ast.ConstraintUniqIndex,
		ast.ConstraintFulltext,
		ast.ConstraintSpatial,
	}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// 处理 "CREATE TABLE..." 语句
		// 统计 CREATE TABLE 语句中定义的二级索引数量
		constraints := util.GetTableConstraints(stmt.Constraints, secondaryIndexes...)
		indexCount = len(constraints)

		// 检查索引计数是否超过最大允许值
		if indexCount > maxIndexCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00037)

		}

	case *ast.AlterTableStmt:
		// 处理 "ALTER TABLE..." 语句
		// 统计 ALTER TABLE ... ADD/MODIFY/CHANGE 语句中新增的二级索引数量
		var alterConstraints []*ast.Constraint
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint, ast.AlterTableModifyColumn, ast.AlterTableChangeColumn) {
			if spec.Constraint != nil {
				alterConstraints = append(alterConstraints, spec.Constraint)
			}
		}
		alterSecondaryConstraints := util.GetTableConstraints(alterConstraints, secondaryIndexes...)

		// 获取现有的二级索引数
		createTable, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if err != nil {
			log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
			return err
		}
		constraints := util.GetTableConstraints(createTable.Constraints, secondaryIndexes...)
		indexCount = len(alterSecondaryConstraints) + len(constraints)
		// 检查索引计数是否超过最大允许值
		if indexCount > maxIndexCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00037)
			return nil
		}

	case *ast.CreateIndexStmt:
		// 处理 "CREATE INDEX..." 语句
		// 每个 CREATE INDEX 语句只创建一个二级索引
		// 获取现有的二级索引数
		createTable, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
		if err != nil {
			log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
			return err
		}
		constraints := util.GetTableConstraints(createTable.Constraints, secondaryIndexes...)
		indexCount = 1 + len(constraints)
		// 检查索引计数是否超过最大允许值
		if indexCount > maxIndexCount {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00037)
			return nil
		}

	default:
		// 其他类型的语句不处理
		return nil
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
