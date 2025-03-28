package ai

import (
	"fmt"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/model"
	"github.com/pingcap/parser/mysql"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00097 = "SQLE00097"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00097,
			Desc:       plocale.Rule00097Desc,
			Annotation: plocale.Rule00097Annotation,
			Category:   plocale.RuleTypeDMLConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagBusiness.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDML.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagPerformance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOnline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level: driverV2.RuleLevelError,
			Params: []*rulepkg.SourceParam{{
				Key:   rulepkg.DefaultSingleParamKeyName,
				Value: "100",
				Desc:  plocale.Rule00097Params1,
				Type:  params.ParamTypeString,
				Enums: nil,
			}},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: false,
			Version:      2,
		},
		Message: plocale.Rule00097Message,
		Func:    RuleSQLE00097,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00097): "在 MySQL 中，禁止对长字段排序.默认参数描述: 排序字段的最大长度, 默认参数值: 100"
您应遵循以下逻辑：
1. 对于 "SELECT ..." 语句：
   1. 检查是否包含 ORDER BY、DISTINCT、GROUP BY 子句。
   2. 提取这些子句中涉及的字段名，并记录对应的表名。
   3. 连接数据库，通过语法节点匹配查询每个字段的类型：
      - 使用辅助函数GetCreateTableStmt获取表的所有字段信息。
      - 使用辅助函数GetColumnWidth获取字段的长度。
      - 如果字段类型为 "TINYTEXT", "TEXT", "MEDIUMTEXT", "LONGTEXT", "TINYBLOB", "BLOB", "MEDIUMBLOB", "LONGBLOB"，则报告违反规则。
   	  - 如果字段类型为VARCHAR且字段长度超过设定阈值，则报告违反规则。

2. 对于 "UPDATE ..." 语句，执行与 SELECT 语句相同的检查流程。

3. 对于 "DELETE ..." 语句，执行与 SELECT 语句相同的检查流程。

4. 对于 "INSERT ... SELECT ..." 语句，执行与 SELECT 语句相同的检查流程。

5. 对于 UNION... 语句，逐一检查所有 SELECT 子句，按照 SELECT 语句的检查流程执行。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00097(input *rulepkg.RuleHandlerInput) error {
	// 获取数值类型的阈值参数
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	threshold := param.Int()
	if threshold <= 0 {
		return fmt.Errorf("param value should be greater than 0")
	}

	createMap := make(map[*ast.TableName]*ast.CreateTableStmt) // 缓存CreateTableStmt

	getCreateTableStmt := func(table *ast.TableName) (*ast.CreateTableStmt, error) {
		createTableStmt, ok := createMap[table]
		if !ok {
			createTableStmt, err := util.GetCreateTableStmt(input.Ctx, table)
			if err != nil {
				return nil, err
			}
			createMap[table] = createTableStmt
			return createTableStmt, nil
		}
		return createTableStmt, nil
	}

	// 定义内部辅助函数获取列定义
	getColumnDef := func(createTableStmt *ast.CreateTableStmt, columnName string) *ast.ColumnDef {
		for _, col := range createTableStmt.Cols {
			if col.Name.Name.L == columnName {
				return col
			}
		}
		return nil
	}

	checkViolate := func(table *ast.TableName, col string) (bool, error) {
		createTableStmt, err := getCreateTableStmt(table)
		if err != nil {
			return false, fmt.Errorf("Failed to get CREATE TABLE statement for table %s: %v", table.Name.String(), err)
		}
		columnDef := getColumnDef(createTableStmt, col)
		if columnDef == nil {
			return false, nil
		}

		// 获取列类型
		colType := columnDef.Tp.Tp

		// 检查是否为 TEXT 或 BLOB 类型
		if colType == mysql.TypeLongBlob || colType == mysql.TypeBlob ||
			colType == mysql.TypeTinyBlob || colType == mysql.TypeMediumBlob {
			return true, nil
		}

		// 检查是否为 VARCHAR 类型
		if colType == mysql.TypeVarchar {
			width := util.GetColumnWidth(columnDef)
			if width > threshold {
				return true, nil
			}
		}
		return false, nil
	}

	extractFieldsFromExpr := func(expr ast.ExprNode, aliasMap []*util.TableAliasInfo) bool {
		if expr == nil {
			return false
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

			table := &ast.TableName{Name: model.NewCIStr(tableName), Schema: model.NewCIStr(schemaName)}
			col := field.Name.Name.String()
			isViolate, err := checkViolate(table, col)
			if err != nil {
				log.NewEntry().Errorf("checkViolate err: %s", err)
				continue
			}
			if isViolate {
				return true
			}
		}
		return false
	}

	gatherColFromOrderByClause := func(orderBy *ast.OrderByClause, aliasInfo []*util.TableAliasInfo) bool {
		if orderBy != nil {
			for _, item := range orderBy.Items {
				if extractFieldsFromExpr(item.Expr, aliasInfo) {
					return true
				}
			}
		}
		return false
	}

	gatherColFromSelectStmt := func(stmt *ast.SelectStmt, aliasInfo []*util.TableAliasInfo) bool {
		if gatherColFromOrderByClause(stmt.OrderBy, aliasInfo) {
			return true
		}
		if stmt.GroupBy != nil {
			for _, item := range stmt.GroupBy.Items {
				if extractFieldsFromExpr(item.Expr, aliasInfo) {
					return true
				}
			}
		}
		if stmt.Distinct {
			if stmt.Fields != nil {
				for _, field := range stmt.Fields.Fields {
					if extractFieldsFromExpr(field.Expr, aliasInfo) {
						return true
					}
				}
			}
		}
		return false
	}

	// 定义一个函数来处理单个 SELECT 语句
	processSelectStmt := func(selectStmt *ast.SelectStmt) bool {

		if selectStmt.From == nil || selectStmt.From.TableRefs == nil {
			// 跳过
			return false
		}
		// 获取表的别名信息
		aliasInfo := util.GetTableAliasInfoFromJoin(selectStmt.From.TableRefs)
		if gatherColFromSelectStmt(selectStmt, aliasInfo) {
			return true
		}
		return false
	}

	// 提取所有的 SELECT 语句，包括子查询和 UNION
	selectStmts := util.GetSelectStmt(input.Node)
	for _, sq := range selectStmts {
		if processSelectStmt(sq) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00097)
			return nil
		}
	}

	// DML
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UnionStmt, *ast.InsertStmt:
		// 上面 util.GetSelectStmt 已经实现了
	case *ast.UpdateStmt:
		aliasInfos := util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		if gatherColFromOrderByClause(stmt.Order, aliasInfos) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00097)
			return nil
		}
	case *ast.DeleteStmt:
		aliasInfos := util.GetTableAliasInfoFromJoin(stmt.TableRefs.TableRefs)
		if gatherColFromOrderByClause(stmt.Order, aliasInfos) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00097)
			return nil
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
