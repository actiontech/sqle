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
	SQLE00041 = "SQLE00041"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00041,
			Desc:       "在 MySQL 中, 唯一索引必须使用固定前缀",
			Annotation: "通过配置该规则可以规范指定业务的唯一索引命名规则，具体命名规范可以自定义设置。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "uniq_",
					Desc:  "固定前缀",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 唯一索引必须使用固定前缀",
		AllowOffline: false,
		Func:         RuleSQLE00041,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00041): "在 MySQL 中，唯一索引必须使用固定前缀.默认参数描述: 固定前缀, 默认参数值: uniq_"
您应遵循以下逻辑：
1、检查当前句子是ALTER还是CREATE类型。
   - 如果是ALTER句子，进入步骤2。
   - 如果是CREATE句子，进入步骤4。

2、检查ALTER句子中是否存在ADD语法节点。
   - 如果存在，进入步骤4。
   - 如果不存在，进入步骤3。

3、检查ALTER句子中是否存在RENAME语法节点。
   - 检查RENAME的目标对象类型是否为唯一索引。
     - 如果是，进入步骤5。

4、检查句子中是否存在UNIQUE INDEX语法节点。
   - 如果存在，进入步骤5。

5、检查目标索引名是否包含固定前缀。
   - 如果不包含，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00041(input *rulepkg.RuleHandlerInput) error {
	// 获取规则参数中的固定前缀
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	requiredPrefix := param.String()

	switch stmt := input.Node.(type) {
	case *ast.AlterTableStmt:
		renameto_map := make(map[string]string) /*oldIndexName, newIndexName*/

		// 检查ALTER语句中的各个语法节点
		for _, spec := range stmt.Specs {
			// alter table ... ADD UNIQUE index
			if util.IsAlterTableCommand(spec, ast.AlterTableAddConstraint) &&
				(spec.Constraint.Tp == ast.ConstraintUniq || spec.Constraint.Tp == ast.ConstraintUniqIndex) {
				if !strings.HasPrefix(spec.Constraint.Name, requiredPrefix) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00041)
					return nil
				}
				// alter table ... RENAME index ....
			} else if util.IsAlterTableCommand(spec, ast.AlterTableRenameIndex) {
				renameto_map[spec.FromKey.String()] = spec.ToKey.String()
			}
		}
		if len(renameto_map) > 0 {
			// 获取获取表的信息
			createTableStmt, err := util.GetCreateTableStmt(input.Ctx, stmt.Table)
			if err != nil {
				log.NewEntry().Errorf("GetCreateTableStmt failed, sqle: %v, error: %v", stmt.Text(), err)
				return err
			}

			// 检查表中的index
			constraints := util.GetTableConstraints(createTableStmt.Constraints, ast.ConstraintUniqIndex, ast.ConstraintUniq)
			if len(constraints) > 0 {
				for _, constraint := range constraints { // origin index name
					if newIndexName, ok := renameto_map[constraint.Name]; ok {
						if !strings.HasPrefix(newIndexName, requiredPrefix) {
							rulepkg.AddResult(input.Res, input.Rule, SQLE00041)
							return nil
						}
					}
				}
			}
		}
	case *ast.CreateTableStmt:
		// UNIQUE index
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintUniq || constraint.Tp == ast.ConstraintUniqIndex {
				if !strings.HasPrefix(constraint.Name, requiredPrefix) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00041)
					return nil
				}
			}
		}
	case *ast.CreateIndexStmt:
		if stmt.KeyType == ast.IndexKeyTypeNone {
			if !strings.HasPrefix(stmt.IndexName, requiredPrefix) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00041)
				return nil
			}
		}
	default:
		return nil
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
