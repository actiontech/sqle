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
)

const (
	SQLE00110 = "SQLE00110"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00110,
			Desc:       "在 MySQL 中, 建议为SQL查询条件建立索引",
			Annotation: "为SQL查询条件建立索引可以显著提高查询性能，减少I/O操作，并提高查询效率。特别是在处理大数据量的表时，索引可以大幅度缩短查询时间，优化数据库性能。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 建议为SQL查询条件建立索引. 不符合条件的字段有: %v",
		AllowOffline: false,
		Func:         RuleSQLE00110,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00110): "在 MySQL 中，建议为SQL查询条件建立索引."
您应遵循以下逻辑：
1. 对于 DML 语句中的 SELECT 子句（包括 SELECT、INSERT、UPDATE、DELETE、UNION 中的子 SELECT）：
   1. 提取表名、WHERE 子句、GROUP BY 和 ORDER BY 节点中的字段，存入集合。
   2. 使用辅助函数 GetTableIndexes 检查集合中字段在对应表中是否有索引。
   3. 若无索引，报告规则违规。

2. 对于 "UPDATE ..." 语句：
   1. 提取表名和 WHERE 子句节点中的字段，存入集合。
   2. 使用辅助函数 GetTableIndexes 检查集合中字段在对应表中是否有索引。
   3. 若无索引，报告规则违规。

3. 对于 "DELETE ..." 语句：
   1. 提取表名和 WHERE 子句节点中的字段，存入集合。
   2. 使用辅助函数 GetTableIndexes 检查集合中字段在对应表中是否有索引。
   3. 若无索引，报告规则违规。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00110(input *rulepkg.RuleHandlerInput) error {
	// 初始化一个映射，用于存储每个表需要检查的字段
	fieldsToCheck := make(map[*ast.TableName]map[string]struct{})

	// 从表达式中提取字段并更新到 fieldsToCheck
	extractFieldsFromExpr := func(expr ast.ExprNode, aliasMap []*util.TableAliasInfo, fieldsToCheck map[*ast.TableName]map[string]struct{}) {
		if expr == nil {
			return
		}
		fields := util.GetColumnNameInExpr(expr)
	NEXT_FIELD:
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

			for table := range fieldsToCheck {
				if table.Name.String() == tableName && table.Schema.String() == schemaName {
					if _, exists := fieldsToCheck[table][field.Name.Name.String()]; !exists {
						fieldsToCheck[table][field.Name.Name.String()] = struct{}{}
						continue NEXT_FIELD
					}
				}
			}
			fieldsToCheck[&ast.TableName{Name: model.NewCIStr(tableName), Schema: model.NewCIStr(schemaName)}] = map[string]struct{}{field.Name.Name.String(): {}}
		}
		return
	}

	// 提取所有的 SELECT 语句，包括子查询和 UNION
	selectStmts := util.GetSelectStmt(input.Node)
	subqueries := util.GetSubquery(input.Node)
	for _, sq := range subqueries {
		selectStmts = append(selectStmts, util.GetSelectStmt(sq.Query)...)
	}

	for _, selectStmt := range selectStmts {
		// 如果 FROM 子句为空，跳过该 SELECT 语句
		if selectStmt.From == nil || selectStmt.From.TableRefs == nil {
			continue
		}

		// 获取表的别名信息
		aliasInfo := util.GetTableAliasInfoFromJoin(selectStmt.From.TableRefs)

		// 从 WHERE 子句中提取字段
		extractFieldsFromExpr(selectStmt.Where, aliasInfo, fieldsToCheck)

		// 从 GROUP BY 子句中提取字段
		if selectStmt.GroupBy != nil {
			for _, item := range selectStmt.GroupBy.Items {
				extractFieldsFromExpr(item.Expr, aliasInfo, fieldsToCheck)
			}
		}

		// 从 ORDER BY 子句中提取字段
		if selectStmt.OrderBy != nil {
			for _, item := range selectStmt.OrderBy.Items {
				extractFieldsFromExpr(item.Expr, aliasInfo, fieldsToCheck)
			}
		}
	}

	// 处理 UPDATE 语句及其子查询
	if updateStmt, ok := input.Node.(*ast.UpdateStmt); ok {
		if updateStmt.TableRefs != nil {
			aliasInfos := util.GetTableAliasInfoFromJoin(updateStmt.TableRefs.TableRefs)

			// 从 WHERE 子句中提取字段
			extractFieldsFromExpr(updateStmt.Where, aliasInfos, fieldsToCheck)
		}
	}

	// 处理 DELETE 语句及其子查询
	if deleteStmt, ok := input.Node.(*ast.DeleteStmt); ok {
		if deleteStmt.TableRefs != nil {
			aliasInfos := util.GetTableAliasInfoFromJoin(deleteStmt.TableRefs.TableRefs)

			// 从 WHERE 子句中提取字段
			extractFieldsFromExpr(deleteStmt.Where, aliasInfos, fieldsToCheck)
		}
	}

	// 遍历收集到的每个表和对应的字段
	for table, fields := range fieldsToCheck {

		// 使用辅助函数获取表的索引信息
		indexes, err := util.GetTableIndexes(input.Ctx, table.Name.String(), table.Schema.String())
		if err != nil {
			// 记录错误日志并继续检查下一个表
			log.NewEntry().Errorf("获取表 %s 的索引失败: %v", table, err)
			continue
		}

		// 创建一个映射，用于存储所有索引列的名称（忽略大小写）
		indexColumns := make(map[string]struct{})
		for _, index := range indexes {
			for _, col := range index {
				indexColumns[strings.ToLower(col)] = struct{}{}
			}
		}

		// 收集未被索引的字段
		var missingFields []string
		for field := range fields {
			if _, exists := indexColumns[strings.ToLower(field)]; !exists {
				missingFields = append(missingFields, field)
			}
		}

		// 如果存在未被索引的字段，报告规则违规
		if len(missingFields) > 0 {
			// 使用 strings.Join 将字段名称合并为一个字符串，使用逗号分隔
			violateFields := strings.Join(missingFields, ", ")
			rulepkg.AddResult(input.Res, input.Rule, SQLE00110, violateFields)
		}
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
