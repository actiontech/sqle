package ai

import (
	"fmt"
	"strconv"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00139 = "SQLE00139"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00139,
			Desc:       "不建议使用全表扫描",
			Annotation: "全表扫描是数据库执行查询时读取表中每一行来查找匹配记录的过程。对于大型表来说，出现全表扫描的SQL会导致显著的性能下降和资源消耗，影响业务稳定运行。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "5",
					Desc:  "表大小(GB)",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "不建议使用全表扫描. 表大小阈值: %v GB",
		AllowOffline: false,
		Func:    RuleSQLE00139,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00139): "For dml, Full table scan are prohibited when the table size is larger than a threshold". The threshold should be a parameter whose default value is "5"
You should follow the following logic:
1. For "select..." The statement,
  * Look at the execution plan of the SQL statement. If the type column is ALL, go to the next step. The execution plan needs to be obtained online,
  * See the size of the table with the type column ALL, The size of the table can be determined by 'select round((index_length+data_length)/1024/1024/1024) 'size_GB' from information_schema.tables where table_name=' table name ' 'to get, the table size needs to be obtained online, if the size of the table is greater than or equal to the threshold, then report the rule violation.
2. For "union..." Statement, performs the same check as above.
3. For "update..." Statement, performs the same check as above.
4. For "delete... Statement, performs the same check as above.
5. For "insert..." Statement, performs the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00139(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxSize, err := strconv.Atoi(param.Value)
	if err != nil {
		return fmt.Errorf("param %s should be a number", param.Value)
	}

	var tableNames []string
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.UpdateStmt, *ast.DeleteStmt:
		// "select..." "union..." "insert..." "update..." "delete..."
		explain, err := util.GetExecutionPlan(input.Ctx, stmt.Text())
		if err != nil {
			log.NewEntry().Errorf("get execution plan failed, sqle: %v, error: %v", stmt.Text(), err)
			return nil
		}
		for _, record := range explain.Plan {
			if record.Type == executor.ExplainRecordAccessTypeAll {
				// full table scan
				tableNames = append(tableNames, record.Table)
			}
		}

	}

	for _, table := range tableNames {
		// get the size of table
		size, err := util.GetTableSizeMB(input.Ctx, table)
		if err != nil {
			log.NewEntry().Errorf("get table size failed, sqle: %v, error: %v", input.Node.Text(), err)
			return nil
		}
		if size >= int64(maxSize*1024) {
			// report rule violation
			rulepkg.AddResult(input.Res, input.Rule, SQLE00139, maxSize)
		}
	}

	return nil
}

// ==== Rule code end ====
