package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
)

const (
	SQLE00012 = "SQLE00012"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00012,
			Desc:       "对于MySQL的DDL, 建议用BIGINT类型代替DECIMAL",
			Annotation: "因为CPU不支持对DECIMAL的直接运算，只是MySQL自身实现了DECIMAL的高精度计算，但是计算代价高，并且存储同样范围值的时候，空间占用也更多；使用BIGINT代替DECIMAL，可根据小数的位数乘以相应的倍数，即可达到精确的浮点存储计算，避免DECIMAL计算代价高的问题",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDDLConvention,
		},
		Message: "对于MySQL的DDL, 建议用BIGINT类型代替DECIMAL. 不符合规定的字段: %v",
		AllowOffline: true,
		Func:    RuleSQLE00012,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00012): "In table definition, BIGINT type should be used instead of DECIMAL".
You should follow the following logic:
1. For "create table ..." statement, check every column, if its type is not Decimal-type, otherwise, add the column name to violation-list
2. For "alter table ... add column ..." statement, check the column, if its type is not Decimal-type, otherwise, add the column name to violation-list
3. For "alter table ... modify column ..." statement, check the modified column definition, if its type is not Decimal-type, otherwise, add the column name to violation-list
4. For "alter table ... change column ..." statement, check the new column's definition, if its type is not Decimal-type, otherwise, add the column name to violation-list
5. Generate a violation message as the checking result, including column names which violate the rule, if there is any violations
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00012(input *rulepkg.RuleHandlerInput) error {
	violateColumns := []*ast.ColumnDef{}
	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		for _, col := range stmt.Cols {
			if util.IsColumnTypeEqual(col, mysql.TypeNewDecimal) {
				violateColumns = append(violateColumns, col)
			}
		}
	case *ast.AlterTableStmt:
		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn) {
			for _, col := range spec.NewColumns {
				if util.IsColumnTypeEqual(col, mysql.TypeNewDecimal) {
					violateColumns = append(violateColumns, col)
				}
			}
		}
	default:
		return nil
	}

	if len(violateColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00012, util.JoinColumnNames(violateColumns))
	}
	return nil
}

// ==== Rule code end ====