package ai

import (
	"fmt"
	"strconv"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00098 = "SQLE00098"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00098,
			Desc:       "避免在单个SQL语句中对同一张表进行多次连接或查询",
			Annotation: "在设计SQL语句时，应避免对同一张表进行多次连接或查询。这种做法可能导致查询性能显著下降，因为它会增加数据库的I/O操作，CPU处理以及内存使用，从而影响整体查询效率",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "3",
					Desc:  "max_table_join_count",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message:      "避免在单个SQL语句中对同一张表进行多次连接或查询. 违反规则的表名: %s",
		AllowOffline: true,
		Func:         RuleSQLE00098,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00098): "For dml, Avoid joining or querying the same table multiple times in a single SQL statement".the threshold should be a parameter whose default value is 3.
You should follow the following logic:
1. For the SELECT clause in all DML statements, check the FROM clause in the SELECT to get the tables involved, record how many times the same table appears, and report the rule violation if the times is greater than or equal to the threshold
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00098(input *rulepkg.RuleHandlerInput) error {

	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxTableCount, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s should be an integer", param.Value)
	}

	// get the table count
	if stmt, ok := input.Node.(ast.DMLNode); ok {
		// dml statement
		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// select ...

			// define a map to record the table count in the select statement
			tableCount := map[string] /*table name*/ int{}
			for _, tableName := range util.GetTableNames(selectStmt) {
				// get table name
				key := fmt.Sprintf("%s.%s", util.GetSchemaName(input.Ctx, tableName.Schema.L), tableName.Name.L)
				tableCount[key]++
			}

			// check if the table count is greater than the threshold
			for tableName, count := range tableCount {
				if count >= maxTableCount {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00098, tableName)
				}
			}
		}
	}

	return nil
}

// ==== Rule code end ====
