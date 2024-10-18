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
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/opcode"
	"github.com/pingcap/tidb/types"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
)

const (
	SQLE00112 = "SQLE00112"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00112,
			Desc:       "在 MySQL 中, 禁止WHERE子句中条件字段与值的数据类型不一致",
			Annotation: "WHERE子句中条件字段与值数据类型不一致会引发隐式数据类型转换，导致优化器选择错误的执行计划，在高并发、大数据量的情况下，不走索引会使得数据库的查询性能严重下降",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止WHERE子句中条件字段与值的数据类型不一致",
		AllowOffline: false,
		Func:         RuleSQLE00112,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00112): "在 MySQL 中，禁止WHERE子句中条件字段与值的数据类型不一致."
您应遵循以下逻辑：
1. 对所有 DML 语句，解析 SQL 语句，获取所有的 WHERE 和 ON 条件字段。
   1. 若条件两侧均为列字段：
      - 使用辅助函数GetCreateTableStmt获取列字段类型。
      - 检查两列字段类型是否一致，不一致则报告违规。

   2. 若左侧为列字段，右侧为常量：
      - 使用辅助函数GetCreateTableStmt获取列字段类型。
      - 检查常量类型是否与列字段类型一致，不一致则报告违规。

   3. 若左侧为常量，右侧为列字段：
      - 使用辅助函数GetCreateTableStmt获取列字段类型。
      - 检查常量类型是否与列字段类型一致，不一致则报告违规。

   4. 若条件两侧均为常量：
      - 检查两常量类型是否一致，不一致则报告违规。

2. 对所有 DML 语句，解析 SQL 语句，获取 USING 条件字段。
   1. 使用辅助函数GetCreateTableStmt获取USING字段在关联表中的数据类型。
   2. 若 USING 字段在两表中的数据类型不一致，则报告违规。

3. 对所有 DML 语句，解析 SQL 语句中的 WHERE 子句，考虑其在 SELECT、UPDATE、DELETE 语句中的位置。
   - 对于 SELECT 语句中的 WHERE 子句：用于筛选符合条件的记录。
   - 对于 UPDATE 语句中的 WHERE 子句：用于指定更新哪些记录。
   - 对于 DELETE 语句中的 WHERE 子句：用于指定删除哪些记录。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00112(input *rulepkg.RuleHandlerInput) error {

	var defaultTable string
	var alias []*util.TableAliasInfo
	getTableName := func(col *ast.ColumnNameExpr) string {
		if col.Name.Table.L != "" {
			for _, a := range alias {
				if a.TableAliasName == col.Name.Table.String() {
					return a.TableName
				}
			}
			return col.Name.Table.O
		}
		return defaultTable
	}

	// 内部辅助函数：获取表达式的类型
	getExprType := func(expr ast.ExprNode) (byte, error) {
		switch node := expr.(type) {
		case *ast.ColumnNameExpr:
			// 获取列的表名
			tableName := getTableName(node)

			// 获取CREATE TABLE语句
			createTableStmt, err := util.GetCreateTableStmt(input.Ctx, &ast.TableName{Name: model.NewCIStr(tableName)})
			if err != nil {
				return 0, err
			}
			// 获取列定义
			for _, colDef := range createTableStmt.Cols {
				if strings.EqualFold(util.GetColumnName(colDef), node.Name.Name.O) {
					return colDef.Tp.Tp, nil
				}
			}
			return 0, fmt.Errorf("列未找到: %s", node.Name.Name.O)
		case *ast.FuncCallExpr:
			names := util.GetFuncName(node)
			if len(names) == 1 {
				switch strings.ToUpper(names[0]) {
				case "CURRENT_DATE":
					return mysql.TypeDate, nil
				case "CURRENT_TIME", "NOW":
					return mysql.TypeDatetime, nil
				case "CURRENT_TIMESTAMP":
					return mysql.TypeTimestamp, nil
				case "CURDATE":
					return mysql.TypeDate, nil
				}
			}
			return 0, fmt.Errorf("不支持的函数: %s", strings.Join(names, "."))
		case *parserdriver.ValueExpr:
			switch node.Datum.Kind() {
			case types.KindInt64, types.KindUint64:
				return mysql.TypeLong, nil
			case types.KindFloat32, types.KindFloat64:
				return mysql.TypeDouble, nil
			case types.KindString:
				return mysql.TypeVarchar, nil
			case types.KindBytes:
				return mysql.TypeBlob, nil
			case types.KindBinaryLiteral:
				return mysql.TypeBit, nil
			case types.KindMysqlDecimal:
				return mysql.TypeNewDecimal, nil
			case types.KindMysqlDuration:
				return mysql.TypeDuration, nil
			case types.KindMysqlTime:
				return mysql.TypeDatetime, nil
			case types.KindMysqlEnum:
				return mysql.TypeEnum, nil
			case types.KindMysqlSet:
				return mysql.TypeSet, nil
			case types.KindMysqlJSON:
				return mysql.TypeJSON, nil
			default:
				return 0, fmt.Errorf("不支持的常量类型: %d", node.Datum.Kind())
			}
		default:
			return 0, fmt.Errorf("不支持的表达式类型: %T", expr)
		}
	}

	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UnionStmt:
		// 获取DELETE、UPDATE语句的表
		alias = util.GetTableAliasInfoFromJoin(util.GetFirstJoinNodeFromStmt(stmt))

		for _, selectStmt := range util.GetSelectStmt(stmt) {
			// 获取默认表名
			if t := util.GetDefaultTable(selectStmt); t != nil {
				defaultTable = t.Name.O
			}

			// 获取表别名信息
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil {
				// 获取FROM子句中的所有表及其别名
				alias = append(alias, util.GetTableAliasInfoFromJoin(selectStmt.From.TableRefs)...)
			}

			// 获取所有WHERE和ON条件表达式
			whereExprs := util.GetWhereExprFromDMLStmt(selectStmt)
			onExprs := []ast.ExprNode{}
			usingExprs := [][]*ast.ColumnName{}
			for _, join := range util.GetAllJoinsFromNode(selectStmt) {
					if join != nil && join.On != nil {
						onExprs = append(onExprs, join.On.Expr)
					}
					if join != nil && join.Using != nil {
						usingExprs = append(usingExprs, join.Using)
					}
			}

			// Combine WHERE and ON expressions
			allExprs := append(whereExprs, onExprs...)

			// 遍历所有条件表达式
			for _, expr := range allExprs {
				// 遍历表达式中的所有二元操作
				util.ScanWhereStmt(func(subExpr ast.ExprNode) (skip bool) {
					binExpr, ok := subExpr.(*ast.BinaryOperationExpr)
					if !ok {
						return false
					}

					// 仅处理等于操作符
					if binExpr.Op != opcode.EQ {
						return false
					}

					// 获取左侧和右侧的类型
					leftType, err := getExprType(binExpr.L)
					if err != nil {
						// 记录错误但不阻止其他检查
						log.NewEntry().Errorf("获取左侧表达式类型失败: %v", err)
						return false
					}

					rightType, err := getExprType(binExpr.R)
					if err != nil {
						log.NewEntry().Errorf("获取右侧表达式类型失败: %v", err)
						return false
					}

					// 比较类型是否一致
					if leftType != rightType {
						// 报告违规
						rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
					}
					return false
				}, expr)
			}
			for _, using := range usingExprs {
				for _, column := range using {
					colTyp := make(map[byte]struct{})
					for _, aliasInfo := range alias {
						tableName := &ast.TableName{Schema: model.NewCIStr(aliasInfo.SchemaName), Name: model.NewCIStr(aliasInfo.TableName)}
						createTableStmt, err := util.GetCreateTableStmt(input.Ctx, tableName)
						if err != nil {
							log.NewEntry().Errorf("获取表 %s 的CREATE TABLE语句失败: %v", tableName.Name.L, err)
							continue
						}

						for _, colDef := range createTableStmt.Cols {
							if strings.EqualFold(util.GetColumnName(colDef), column.Name.L) {
								colTyp[colDef.Tp.Tp] = struct{}{}
							}
						}
					}

					if len(colTyp) > 1 {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
					}
				}
			}
		}

	}

	switch stmt := input.Node.(type) {
	case *ast.UpdateStmt:
		// 获取默认表名
		if t := util.GetTableNames(stmt); len(t) != 0 {
			defaultTable = t[0].Name.O
		}

		// 获取表别名信息
		if stmt.TableRefs != nil && stmt.TableRefs.TableRefs != nil {
			// 获取FROM子句中的所有表及其别名
			alias = util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		}

		// 获取所有WHERE和ON条件表达式
		whereExprs := util.GetWhereExprFromDMLStmt(stmt)
		onExprs := []ast.ExprNode{}
		joinExpr := util.GetFirstJoinNodeFromStmt(stmt)
		if joinExpr != nil && joinExpr.On != nil {
			onExprs = append(onExprs, joinExpr.On.Expr)
		}
		// Combine WHERE and ON expressions
		allExprs := append(whereExprs, onExprs...)

		// 遍历所有条件表达式
		for _, expr := range allExprs {
			// 遍历表达式中的所有二元操作
			util.ScanWhereStmt(func(subExpr ast.ExprNode) (skip bool) {
				binExpr, ok := subExpr.(*ast.BinaryOperationExpr)
				if !ok {
					return false
				}

				// 仅处理等于操作符
				if binExpr.Op != opcode.EQ {
					return false
				}

				// 获取左侧和右侧的类型
				leftType, err := getExprType(binExpr.L)
				if err != nil {
					// 记录错误但不阻止其他检查
					log.NewEntry().Errorf("获取左侧表达式类型失败: %v", err)
					return false
				}

				rightType, err := getExprType(binExpr.R)
				if err != nil {
					log.NewEntry().Errorf("获取右侧表达式类型失败: %v", err)
					return false
				}

				// 比较类型是否一致
				if leftType != rightType {
					// 报告违规
					rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
				}
				return false
			}, expr)
		}
	case *ast.DeleteStmt:
		// 获取默认表名
		if t := util.GetTableNames(stmt); len(t) != 0 {
			defaultTable = t[0].Name.O
		}

		// 获取表别名信息
		if stmt.TableRefs != nil && stmt.TableRefs.TableRefs != nil {
			// 获取FROM子句中的所有表及其别名
			alias = util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		}

		// 获取所有WHERE和ON条件表达式
		whereExprs := util.GetWhereExprFromDMLStmt(stmt)
		onExprs := []ast.ExprNode{}
		joinExpr := util.GetFirstJoinNodeFromStmt(stmt)
		if joinExpr != nil && joinExpr.On != nil {
			onExprs = append(onExprs, joinExpr.On.Expr)
		}
		// Combine WHERE and ON expressions
		allExprs := append(whereExprs, onExprs...)
		// 遍历所有条件表达式
		for _, expr := range allExprs {
			// 遍历表达式中的所有二元操作
			util.ScanWhereStmt(func(subExpr ast.ExprNode) (skip bool) {
				binExpr, ok := subExpr.(*ast.BinaryOperationExpr)
				if !ok {
					return false
				}

				// 仅处理等于操作符
				if binExpr.Op != opcode.EQ {
					return false
				}

				// 获取左侧和右侧的类型
				leftType, err := getExprType(binExpr.L)
				if err != nil {
					// 记录错误但不阻止其他检查
					log.NewEntry().Errorf("获取左侧表达式类型失败: %v", err)
					return false
				}

				rightType, err := getExprType(binExpr.R)
				if err != nil {
					log.NewEntry().Errorf("获取右侧表达式类型失败: %v", err)
					return false
				}

				// 比较类型是否一致
				if leftType != rightType {
					// 报告违规
					rulepkg.AddResult(input.Res, input.Rule, SQLE00112)
				}
				return false
			}, expr)
		}
	}

	return nil
}

// ==== Rule code end ====
