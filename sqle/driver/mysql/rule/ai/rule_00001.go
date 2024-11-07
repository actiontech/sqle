package ai

import (
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/opcode"
	parserdriver "github.com/pingcap/tidb/types/parser_driver"
)

const (
	SQLE00001 = "SQLE00001"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00001,
			Desc:       "在 MySQL 中, 禁止SQL语句不带WHERE条件或者WHERE条件为永真",
			Annotation: "使用有效的WHERE条件能够避免全表扫描，提高SQL执行效率；而恒为TRUE的WHERE条件，如where 1=1、where true=true等，在执行时会进行全表扫描产生额外开销。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 禁止SQL语句不带WHERE条件或者WHERE条件为永真",
		AllowOffline: false,
		Func:         RuleSQLE00001,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00001): "在 MySQL 中，禁止SQL语句不带WHERE条件或者WHERE条件为永真."
您应遵循以下逻辑：
1. 对于DML语句的SELECT子句（包括SELECT、INSERT、UPDATE、DELETE、UNION语句中的子SELECT语句），如果满足以下任一条件，则报告违反规则：
   1. SQL语句未包含WHERE条件。
   2. 使用辅助函数IsExprConstTrue检查WHERE条件是否为永真表达式。
   3. WHERE条件的最外层使用OR，并且OR条件中包含恒真表达式。

2. 对于"WITH.."语句，执行与上述相同的检查。

3. 对于DML语句的SELECT子句（包括SELECT、INSERT、UPDATE、DELETE、UNION语句中的子SELECT语句），如果满足以下任一条件，则报告违反规则：
   1. WHERE条件为column IS NOT NULL，且column是不可为空的字段。
   2. WHERE条件的最外层使用OR，并且OR条件中包含上述条件。

4. 对于"WITH.."语句，执行与上述相同的检查。
==== Prompt end ====
*/

// ==== Rule code start ====

// 规则函数实现开始
func RuleSQLE00001(input *rulepkg.RuleHandlerInput) error {

	// 获取列的信息
	extractFieldsFromExpr := func(expr ast.ExprNode, aliasMap []*util.TableAliasInfo) (*ast.TableName, string) {
		if expr == nil {
			return nil, ""
		}
		fields := util.GetColumnNameInExpr(expr)
		for _, field := range fields {
			tableName := field.Name.Table.String()
			schemaName := field.Name.Schema.String()
			if tableName != "" {
				// 如果字段有表前缀，通过别名映射获取真实表名
				for _, alias := range aliasMap {
					if alias.TableAliasName == tableName {
						tableName = alias.TableName
						schemaName = alias.SchemaName
						break
					}
				}
			} else {
				// 如果字段没有表前缀，尝试将其映射到第一个表
				for _, alias := range aliasMap {
					tableName = alias.TableName
					schemaName = alias.SchemaName
					break
				}
			}
			if tableName == "" {
				return nil, ""
			}
			return &ast.TableName{Name: model.NewCIStr(tableName), Schema: model.NewCIStr(schemaName)}, field.Name.Name.String()
		}
		return nil, ""
	}

	// 在线，当有 is not null时，则当前列是否有not null约束，如果有则为恒真
	checkIsNotNull := func(table *ast.TableName, col string) (bool, error) {
		createTableStmt, err := util.GetCreateTableStmt(input.Ctx, table)
		if err != nil {
			log.NewEntry().Errorf("GetCreateTableStmt failed, error: %v", err)
			return false, err
		}
		for _, columnDef := range createTableStmt.Cols {
			if strings.EqualFold(columnDef.Name.OrigColName(), col) && util.IsColumnHasOption(columnDef, ast.ColumnOptionNotNull) {
				return true, nil
			}
		}
		return false, nil
	}

	// 恒真判断
	var IsExprConstTrue2 func(node ast.ExprNode, aliasInfo []*util.TableAliasInfo) (bool, error)
	IsExprConstTrue2 = func(node ast.ExprNode, aliasInfo []*util.TableAliasInfo) (bool, error) {
		switch expr := node.(type) {
		case *ast.BinaryOperationExpr:
			if expr != nil && expr.Op == opcode.LogicOr { // 处理or场景
				left, err := IsExprConstTrue2(expr.L, aliasInfo)
				right := false
				if !left {
					right, err = IsExprConstTrue2(expr.R, aliasInfo)
				}
				return (left || right), err
			} else {
				return util.IsExprConstTrue(node), nil
			}
		case *ast.ParenthesesExpr: // 含有括号()
			return IsExprConstTrue2(expr.Expr, aliasInfo)
		case *ast.IsNullExpr:
			if expr.Not { // 需要[在线]获取列是否有 not null约束
				table, col := extractFieldsFromExpr(expr.Expr, aliasInfo)
				if table != nil && col != "" {
					yes, err := checkIsNotNull(table, col)
					if err != nil {
						return false, err
					}
					return yes, nil
				}
			}
			return false, nil
		case *ast.ExistsSubqueryExpr:
			if !expr.Not { // 处理 exists
				if len(util.GetTableNames(expr)) == 0 { // 不存在表时，则为恒真
					return true, nil
				}
			}
			return false, nil
		case *ast.PatternInExpr:
			if left_value, ok := expr.Expr.(*parserdriver.ValueExpr); ok { // 当左边为确定值
				// 1 in (1,2,3)
				if expr.Sel == nil {
					for _, expr := range expr.List {
						if value, okk := expr.(*parserdriver.ValueExpr); okk {
							result, err := util.EqualValueExpr(left_value, value)
							if err != nil {
								log.NewEntry().Errorf("EqualValueExpr failed, error: %v", err)
								return false, err
							}
							if result {
								return true, nil
							}
						}
					}
				} else {
					// 1 in (select 1 from dual union select 2 from dual)
					if len(util.GetTableNames(expr)) == 0 { // 不存在表时
						sList := util.GetSelectStmt(expr)
						for _, ss := range sList {
							if ss.From == nil {
								if field_value, okk := (ss.Fields.Fields[0]).Expr.(*parserdriver.ValueExpr); okk {
									result, err := util.EqualValueExpr(left_value, field_value)
									if err != nil {
										log.NewEntry().Errorf("EqualValueExpr failed, error: %v", err)
										return false, err
									}
									if result {
										return true, nil
									}
								}
							}
						}
					}
				}
			}
			return false, nil
		default:
			return util.IsExprConstTrue(node), nil
		}
	}

	// 检索出所有的select语句中的 where
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt, *ast.DeleteStmt, *ast.UpdateStmt:
		selectList := util.GetSelectStmt(stmt)
		for _, sel := range selectList {
			if sel.From != nil {
				aliasInfo := util.GetTableAliasInfoFromJoin(sel.From.TableRefs)
				if sel.Where != nil {
					isConst, err := IsExprConstTrue2(sel.Where, aliasInfo)
					if err != nil {
						return err
					}
					if isConst {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
						return nil
					}
				} else {
					rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
					return nil
				}
			}
		}
	}

	// 单独处理 delete、update的 where
	switch stmt2 := input.Node.(type) {
	case *ast.DeleteStmt:
		if stmt2.Where == nil {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
			return nil
		} else {
			aliasInfos := util.GetTableAliasInfoFromJoin(stmt2.TableRefs.TableRefs)
			isConst, err := IsExprConstTrue2(stmt2.Where, aliasInfos)
			if err != nil {
				return err
			}
			if isConst {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
				return nil
			}
		}
	case *ast.UpdateStmt:
		if stmt2.Where == nil {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
			return nil
		} else {
			aliasInfos := util.GetTableAliasInfoFromJoin(stmt2.TableRefs.TableRefs)
			isConst, err := IsExprConstTrue2(stmt2.Where, aliasInfos)
			if err != nil {
				return err
			}
			if isConst {
				rulepkg.AddResult(input.Res, input.Rule, SQLE00001)
				return nil
			}
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
