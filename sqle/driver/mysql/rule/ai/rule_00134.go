package ai

import (
	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00134 = "SQLE00134"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00134,
			Desc:       "在 MySQL 中, 避免对主键值进行修改",
			Annotation: "主键在大多数数据库系统中用于定义数据的唯一性，并且常常与数据的物理存储结构密切相关。更新主键会导致底层存储结构（如聚簇索引）的重大重新组织，引发性能下降。此外，主键的更改可能影响数据一致性，尤其在涉及复杂事务处理和高并发操作的场景中。",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免对主键值进行修改",
		AllowOffline: false,
		Func:         RuleSQLE00134,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00134): "在 MySQL 中，避免对主键值进行修改."
您应遵循以下逻辑：
1. 对于 "UPDATE..." 语句，执行以下步骤：
2. 提取语句中的 SET 列，存入集合。
3. 使用辅助函数GetCreateTableStmt获取表的主键字段。
4. 比较集合中的字段与主键字段：
   - 如果集合中包含主键字段，则报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00134(input *rulepkg.RuleHandlerInput) error {
	// 检查输入的 AST 节点是否为 UPDATE 语句
	updateStmt, ok := input.Node.(*ast.UpdateStmt)
	if !ok {
		// 不是 UPDATE 语句，跳过检查
		return nil
	}

	// 获取要更新的表名
	tableNames := util.GetTableNames(updateStmt.TableRefs.TableRefs)
	if len(tableNames) == 0 {
		// 没有表名，无法进行检查
		return nil
	}

	// 使用辅助函数获取表的 CREATE TABLE 语句
	createTableStmt, err := util.GetCreateTableStmt(input.Ctx, tableNames[0])
	if err != nil {
		return err
	}

	// 内部辅助函数提取主键字段
	extractPrimaryKeys := func(createTableStmt *ast.CreateTableStmt) []string {
		var primaryKeys []string
		// 检查列定义中的主键
		for _, col := range createTableStmt.Cols {
			if util.IsColumnPrimaryKey(col) {
				primaryKeys = append(primaryKeys, util.GetColumnName(col))
			}
		}
		// 检查表级约束中的主键
		constraint := util.GetTableConstraint(createTableStmt.Constraints, ast.ConstraintPrimaryKey)
		if constraint != nil {
			for _, key := range constraint.Keys {
				primaryKeys = append(primaryKeys, util.GetIndexColName(key))
			}
		}
		return primaryKeys
	}

	// 提取主键字段
	primaryKeys := extractPrimaryKeys(createTableStmt)

	// 如果没有主键定义，则不进行进一步检查
	if len(primaryKeys) == 0 {
		return nil
	}

	// 内部辅助函数提取 SET 子句中的所有列名
	extractSetColumns := func(updateStmt *ast.UpdateStmt) []string {
		var setColumns []string
		for _, item := range updateStmt.List {
			// 提取被赋值的列名
			columnName := item.Column.Name.String()
			if columnName != "" {
				setColumns = append(setColumns, columnName)
			}
		}
		return setColumns
	}

	// 提取 SET 子句中的所有列名
	setColumns := extractSetColumns(updateStmt)

	// 比较 SET 列与主键列，找出重复的列
	var violatingColumns []string
	for _, pk := range primaryKeys {
		if util.IsStrInSlice(pk, setColumns) {
			violatingColumns = append(violatingColumns, pk)
		}
	}

	// 如果存在违反规则的列，报告违规
	if len(violatingColumns) > 0 {
		rulepkg.AddResult(input.Res, input.Rule, SQLE00134)
	}

	return nil
}

// 规则函数实现结束

// ==== Rule code end ====
