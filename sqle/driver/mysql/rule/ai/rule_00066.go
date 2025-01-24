package ai

import (
	"regexp"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00066 = "SQLE00066"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:         SQLE00066,
			Desc:         plocale.Rule00066Desc,
			Annotation:   plocale.Rule00066Annotation,
			Category:     plocale.RuleTypeDDLConvention,
			Level:        driverV2.RuleLevelError,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00066Message,
		Func:    RuleSQLE00066,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00066): "在 MySQL 中，禁止除索引外的DROP 操作."
您应遵循以下逻辑：
1. 对于 "ALTER TABLE ..."语句，如果存在以下任何一项，则报告违反规则：
  1. 语法树中存在DROP操作且操作对象不是索引。
2. 对于 "Drop..."语句，如果操作对象不是索引，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00066(input *rulepkg.RuleHandlerInput) error {
	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		// 检查 ALTER TABLE 语句中的 DROP 操作
		for _, spec := range stmt.Specs {
			// 如果 DROP 的对象不是索引，则报告违规
			switch spec.Tp {
			case ast.AlterTableDropColumn:
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
				break
			case ast.AlterTableDropPrimaryKey:
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
				break
			case ast.AlterTableModifyColumn:
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
				break
			case ast.AlterTableDropForeignKey:
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
				break
			case ast.AlterTableDropCheck:
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
				break
			case ast.AlterTableDropPartition:
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
				break
			case ast.AlterTableAlterColumn:
				// ref code: github.com/pingcap/parser/ast/ddl.go
				// In: func (n *AlterTableSpec) Restore(ctx *format.RestoreCtx) error {
				// case AlterTableAlterColumn: 中的处理方式
				if spec.NewColumns != nil {
					if len(spec.NewColumns[0].Options) == 0 {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
					}
				}
			}
		}
	case *ast.DropDatabaseStmt:
		rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
	case *ast.DropTableStmt: // 包含了 view
		rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
	case *ast.DropSequenceStmt:
		rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
	case *ast.UnparsedStmt:
		// alter table t1 drop constraint ...
		match1 := alter_table1.MatchString(input.Node.Text())
		if match1 {
			matches2 := drop_constraint2.MatchString(input.Node.Text())
			if matches2 {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
			}
		} else {
			// drop EVENT/PROCEDURE/FUNCTION/TRIGGER/VIEW xxx.....
			if dropReg1.MatchString(input.Node.Text()) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00066)
			}
		}

	}
	return nil
}

// 规则函数实现结束
var dropReg1 = regexp.MustCompile(`(?i)\bDROP\s+(EVENT|PROCEDURE|FUNCTION|TRIGGER|VIEW)\s+`)

var alter_table1 = regexp.MustCompile(`(?i)\balter\s+table\s+`)
var drop_constraint2 = regexp.MustCompile(`(?i)drop\s+constraint\s+`)

// ==== Rule code end ====
