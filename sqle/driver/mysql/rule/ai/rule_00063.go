package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00063 = "SQLE00063"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00063,
			Desc:       "在 MySQL 中, 唯一索引名必须遵循指定格式",
			Annotation: "通过配置该规则可以规范指定业务的唯一索引命名规则，如索引字段存在多个，则可以拼接字段名，不要超过索引名长度即可。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "固定字符+索引类型+表名+字段名，如IDX_UK_TABLENAME_COLNAME",
					Desc:  "索引命名格式",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 唯一索引名必须遵循指定格式",
		AllowOffline: false,
		Func:         RuleSQLE00063,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00063): "在 MySQL 中，唯一索引名必须遵循指定格式.默认参数描述: 索引命名格式, 默认参数值: 固定字符+索引类型+表名+字段名，如IDX_UK_TABLENAME_COLNAME"
您应遵循以下逻辑：
1、检查当前句子是ALTER还是CREATE类型。
   - 如果是ALTER句子，进入步骤2。
   - 如果是CREATE句子，进入步骤3。

2、对于ALTER句子：
   - 检查是否存在ADD操作节点。
     - 如果存在，进入步骤4。
   - 检查是否存在RENAME操作节点。
     - 如果存在，使用辅助函数GetTableIndexes获取当前索引，检查RENAME的目标对象类型是否为唯一索引。
     - 如果是唯一索引，进入步骤4。

3、对于CREATE句子：
   - 检查句子中是否存在UNIQUE INDEX节点。
     - 如果存在，进入步骤4。

4、检查目标索引名是否遵从固定格式。
   - 如果不遵从，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00063(input *rulepkg.RuleHandlerInput) error {
	// 内部匿名的辅助函数
	isIndexNameViolate := func(indexName string, tableName string, cols []string) bool {
		if !strings.EqualFold(indexName, fmt.Sprintf("IDX_UK_%v_%v", tableName, strings.Join(cols, "_"))) {
			return true
		}
		return false
	}

	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		tableName := stmt.Table.Name.String()
		var createTable *ast.CreateTableStmt
		var err error
		for _, spec := range stmt.Specs {
			// 检查ADD操作节点
			if spec.Tp == ast.AlterTableAddConstraint {
				if spec.Constraint != nil && spec.Constraint.Tp == ast.ConstraintUniq {
					indexName := spec.Constraint.Name
					var indexedCols []string
					for _, key := range spec.Constraint.Keys {
						indexedCols = append(indexedCols, key.Column.Name.String())
					}
					if isIndexNameViolate(indexName, tableName, indexedCols) {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00063)
						return nil
					}
				}
			} else if spec.Tp == ast.AlterTableRenameIndex { // 检查RENAME操作节点 （在线）
				createTable, err = util.GetCreateTableStmt(input.Ctx, stmt.Table)
				if err != nil {
					log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
					return nil
				}
				constraintUniqs := util.GetTableConstraints(createTable.Constraints, ast.ConstraintUniq)
				if len(constraintUniqs) == 0 {
					return nil
				}
				oldIdxname := spec.FromKey.String()
				newIdxname := spec.ToKey.String()
				for _, constraint := range constraintUniqs {
					if strings.EqualFold(oldIdxname, oldIdxname) {
						var indexedCols []string
						for _, key := range constraint.Keys {
							indexedCols = append(indexedCols, key.Column.Name.String())
						}
						if isIndexNameViolate(newIdxname, tableName, indexedCols) {
							rulepkg.AddResult(input.Res, input.Rule, SQLE00063)
							return nil
						}
					}

				}
			}
		}
	case *ast.CreateIndexStmt:
		// 检查CREATE INDEX语句是否为UNIQUE
		if stmt.KeyType == ast.IndexKeyTypeUnique {
			tableName := stmt.Table.Name.String()
			indexName := stmt.IndexName
			var indexedCols []string
			for _, key := range stmt.IndexPartSpecifications {
				indexedCols = append(indexedCols, key.Column.Name.String())
			}
			if isIndexNameViolate(indexName, tableName, indexedCols) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00063)
				return nil
			}
		}

	case *ast.CreateTableStmt:
		// 检查CREATE TABLE语句中的UNIQUE约束
		tableName := stmt.Table.Name.String()
		constraints := util.GetTableConstraints(stmt.Constraints, ast.ConstraintUniq)
		for _, constraint := range constraints {
			indexName := constraint.Name
			var indexedCols []string
			for _, key := range constraint.Keys {
				indexedCols = append(indexedCols, key.Column.Name.String())
			}
			if isIndexNameViolate(indexName, tableName, indexedCols) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00063)
				return nil
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
