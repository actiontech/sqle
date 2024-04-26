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
	SQLE00039 = "SQLE00039"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00039,
			Desc:       "对于MySQL的索引, 建议索引字段的区分度大于阈值",
			Annotation: "选择区分度高的字段作为索引，可快速定位数据；区分度太低，无法有效利用索引，甚至可能需要扫描大量数据页，拖慢SQL；具体规则阈值可以根据业务需求调整，默认值：70",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeIndexOptimization,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "70",
					Desc:  "区分度（百分比）",
					Type:  params.ParamTypeInt,
				},
			},
		},
		Message: "索引列 %v 未超过区分度阈值 百分之%v, 不建议选为索引",
		Func:    RuleSQLE00039,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00039): "For index, It is suggested that the discrimination of the index field is greater than the percentage threshold. the threshold should be a parameter whose default value is 70.".
You should follow the following logic:
1. For the  "CREATE INDEX ..." statements, check if discrimination of the index which on the table option or on the column less than the threshold. If it does, report a violation. The calculation of index discrimination needs to be obtained from the online database.
2. For the  "ALTER TABLE ... ADD INDEX ..." statements, perform the same check as above.
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00039(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	threshold := param.Int()
	if threshold <= 0 || threshold > 100 {
		return fmt.Errorf("param %s should be in range (0, 100]", rulepkg.DefaultSingleParamKeyName)
	}

	var tableName *ast.TableName
	indexColumns := make([]string, 0)
	switch stmt := input.Node.(type) {
	case *ast.CreateIndexStmt:
		// "create index..."

		tableName = stmt.Table
		for _, col := range stmt.IndexPartSpecifications {
			//"create index... column..."
			indexColumns = append(indexColumns, util.GetIndexColName(col))
		}
	case *ast.AlterTableStmt:
		// "alter table"
		tableName = stmt.Table

		for _, spec := range util.GetAlterTableCommandsByTypes(stmt, ast.AlterTableAddConstraint) {
			// "alter table... add index..."
			constraints := util.GetTableConstraints([]*ast.Constraint{spec.Constraint}, util.GetIndexConstraintTypes()...)

			for _, constraint := range constraints {
				for _, col := range constraint.Keys {
					indexColumns = append(indexColumns, util.GetIndexColName(col))
				}
			}

		}
	}
	if len(indexColumns) == 0 {
		// the table has no index
		return nil
	}

	discrimination, err := util.CalculateIndexDiscrimination(input.Ctx, tableName, indexColumns)
	if err != nil {
		log.NewEntry().Errorf("get index discrimination failed, sqle: %v, error: %v", input.Node.Text(), err)
		return nil
	}

	// the table has index, check the discrimination
	for col, percentage := range discrimination {
		if percentage < float64(threshold) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00039, col, threshold)
		}
	}
	return nil
}

// ==== Rule code end ====
