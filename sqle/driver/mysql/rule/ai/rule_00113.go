package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
)

const (
	SQLE00113 = "SQLE00113"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00113,
			Desc:       "对于MySQL的DML, 不建议对条件字段使用负向查询",
			Annotation: "SQL查询条件中存在NOT IN、NOT LIKE、NOT EXISTS、不等于等负向查询条件，将导致全表扫描，出现慢SQL",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeIndexInvalidation,
		},
		Message: "对于MySQL的DML, 不建议对条件字段使用负向查询",
		AllowOffline: true,
		Func:    RuleSQLE00113,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00113): "For dml, using negative queries against WHERE conditional fields is prohibited".
You should follow the following logic:
1. For "SELECT..." Statement, checks for the presence of a negative (NOT(except NOT NULL), NOT IN, NOT LIKE, NOT EXISTS, NOT BETWEEN, not equal to) in the WHERE condition in the sentence, and if so, reports a rule violation.
2. For "DELETE..." Statement, perform the same checks as above for each SELECT clause in the statement. The DELETE statement itself is checked for the presence of a negative query in its WHERE condition, and if so, a rule violation is reported.
3. For "INSERT..." Statement, perform the same checks as above for each SELECT clause in the statement.
4. For "UPDATE..." Statement, perform the same checks as above for each SELECT clause in the statement. The UPDATE statement itself is checked for the presence of a negative query in its WHERE condition, and if so, a rule violation is reported.
5. For "UNION..." Statement, perform the same checks as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00113(input *rulepkg.RuleHandlerInput) error {

	// get where condition from DML statement
	whereList := util.GetWhereExprFromDMLStmt(input.Node)

	negative := false
	util.ScanWhereStmt(func(expr ast.ExprNode) (skip bool) {
		switch x := expr.(type) {
		case *ast.UnaryOperationExpr:
			// NOT
			if x.Op == opcode.Not {
				negative = true
				return true
			}
		case *ast.BinaryOperationExpr:
			//  not equal to, NOT
			if x.Op == opcode.NE || x.Op == opcode.NullEQ {
				negative = true
				return true
			}
			if x.Op == opcode.Not {
				// except not null
				if v, ok := x.R.(*ast.ValuesExpr); ok && v.Type.Tp == mysql.TypeNull {
					return true
				}
				negative = true
				return true
			}
		case *ast.PatternInExpr:
			// NOT IN
			if x.Not {
				negative = true
				return true
			}
		case *ast.PatternLikeExpr:
			// NOT LIKE
			if x.Not {
				negative = true
				return true
			}
		case *ast.ExistsSubqueryExpr:
			// NOT EXISTS
			if v, ok := x.Sel.(*ast.SubqueryExpr); ok && x.Not && v.Exists {
				negative = true
				return true
			}
		case *ast.BetweenExpr:
			// NOT BETWEEN
			if x.Not {
				negative = true
				return true
			}
		}
		return false
	}, whereList...)

	if negative {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00113)
	}
	return nil
}

// ==== Rule code end ====