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

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00040 = "SQLE00040"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00040,
			Desc:       plocale.Rule00040Desc,
			Annotation: plocale.Rule00040Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			Level:      driverV2.RuleLevelWarn,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "idx_",
				Desc:  plocale.Rule00040Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
		},
		Message: plocale.Rule00040Message,
		Func:    RuleSQLE00040,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00040): "在 MySQL 中，普通索引必须使用固定前缀.默认参数描述: 固定前缀, 默认参数值: idx_"
您应遵循以下逻辑：
1、检查当前语句类型：如果是ALTER语句，进入步骤2；如果是CREATE语句，进入步骤3。
2、对于ALTER语句：
   a. 检查是否包含ADD语法节点：如果包含，进入步骤4。
   b. 检查是否包含RENAME语法节点：如果包含，进入步骤5。
3、对于CREATE语句：
   a. 检查是否包含UNIQUE、FULLTEXT、SPATIAL等关键词，仅定义了INDEX标识：如果是，进入步骤4。
4、检查目标索引名是否包含固定前缀：如果不包含，报告违反规则。
5、从表中获取当前索引信息：
   a. 检查RENAME的目标对象类型是否为普通索引（非唯一索引）：如果是，进入步骤4。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00040(input *rulepkg.RuleHandlerInput) error {
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
			// alter table ... ADD index
			if util.IsAlterTableCommand(spec, ast.AlterTableAddConstraint) && spec.Constraint.Tp == ast.ConstraintIndex {
				if !strings.HasPrefix(spec.Constraint.Name, requiredPrefix) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00040)
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
			constraints := util.GetTableConstraints(createTableStmt.Constraints, ast.ConstraintIndex)
			if len(constraints) > 0 {
				for _, constraint := range constraints { // origin index name
					if newIndexName, ok := renameto_map[constraint.Name]; ok {
						if !strings.HasPrefix(newIndexName, requiredPrefix) {
							rulepkg.AddResult(input.Res, input.Rule, SQLE00040)
							return nil
						}
					}
				}
			}
		}
	case *ast.CreateTableStmt:
		// index
		for _, constraint := range stmt.Constraints {
			if constraint.Tp == ast.ConstraintIndex {
				if !strings.HasPrefix(constraint.Name, requiredPrefix) {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00040)
					return nil
				}
			}
		}
	case *ast.CreateIndexStmt:
		if stmt.KeyType == ast.IndexKeyTypeNone {
			if !strings.HasPrefix(stmt.IndexName, requiredPrefix) {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00040)
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
