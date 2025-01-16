package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00177 = "SQLE00177"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00177,
			Desc:       "建议Order By字段个数不超过指定阈值",
			Annotation: "使用过多的Order By字段会增加排序操作的复杂性，并可能导致性能下降。排序时，MySQL需要对结果集中的每一行进行多字段比较，这可能会耗费更多的CPU和内存资源。如果排序数据集大小超过了可用内存，则可能会导致创建临时表并在磁盘上进行排序，从而增加I/O开销。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "3",
					Desc:  "order by字段个数最大值",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "建议Order By字段个数不超过指定阈值. 阈值: %v",
		AllowOffline: true,
		Func:    RuleSQLE00177,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00177): "For dml, the number of Order By fields should not exceed the specified threshold".the threshold should be a parameter whose default value is 3.
You should follow the following logic:
1. For "select..." The statement,
* Check if the ORDER BY keyword is present in the SQL statement, if it is, proceed to the next step
* Count the number of fields in the ORDER BY clause and report a rule violation if the number of fields exceeds a threshold
2. For INSERT... Statement to perform the same check as above on the SELECT clause in the INSERT statement.
3. For UNION... Statement, does the same check as above for each SELECT clause in the statement.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00177(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxOrderByCount := param.Int()

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.InsertStmt, *ast.UnionStmt:
		// "select", "insert", "union"
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			if selectStmt.OrderBy == nil {
				continue
			}

			if len(selectStmt.OrderBy.Items) > maxOrderByCount {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00177, maxOrderByCount)
			}
		}
	}
	return nil
}

// ==== Rule code end ====
