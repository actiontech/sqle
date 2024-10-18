package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	util "github.com/actiontech/sqle/sqle/driver/mysql/rule/ai/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00122 = "SQLE00122"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00122,
			Desc:       "在 MySQL 中, 避免对值全为NULL的列直接使用 SUM或COUNT函数",
			Annotation: "当某一列的值全是NULL时，COUNT(COL)的返回结果为0，但SUM(COL)的返回结果为NULL，因此使用SUM()时需注意NPE问题（指数据返回NULL）；如业务需避免NPE问题，建议开启此规则",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeDMLConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 避免对值全为NULL的列直接使用 SUM或COUNT函数. 违反规则的列名: %s",
		AllowOffline: false,
		Func:         RuleSQLE00122,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00122): "在 MySQL 中，避免对值全为NULL的列直接使用 SUM或COUNT函数."
您应遵循以下逻辑：
1. 对于所有 DML 语句，检查语句中是否包含 SUM 或 COUNT 语法节点，如果存在，则进一步检查。
2. 登录数据库。
3. 使用辅助函数 IsColumnAllNull 检查语句中 SUM 或 COUNT 所涉及的字段（如 age）的值是否全为 NULL，如果是，报告违反规则。
4. 验证某个字段是否全为 NULL 的方法，执行以下语句：如果结果为 0，表示该字段 age 的值全为 NULL。
   ```sql
   SELECT (SELECT COUNT(*) FROM CUSTOMERS)-(SELECT COUNT(*) from customers WHERE AGE IS NULL) RESULT;
   ```
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00122(input *rulepkg.RuleHandlerInput) error {
	var alias []*util.TableAliasInfo
	var defaultTable string
	getTableName := func(col *ast.ColumnName) string {
		if col.Table.L != "" {
			for _, a := range alias {
				if a.TableAliasName == col.Table.String() {
					return a.TableName
				}
			}
			return col.Table.L
		}
		return defaultTable
	}

	// 检查DML语句类型
	switch stmt := input.Node.(type) {
	case *ast.SelectStmt, *ast.UpdateStmt, *ast.DeleteStmt, *ast.InsertStmt, *ast.UnionStmt:
		for _, selectStmt := range util.GetSelectStmt(stmt) {

			// get default table name
			if t := util.GetDefaultTable(selectStmt); t != nil {
				defaultTable = t.Name.O
			}

			// get table alias info
			if selectStmt.From != nil && selectStmt.From.TableRefs != nil {
				alias = util.GetTableAliasInfoFromJoin(selectStmt.From.TableRefs)
			}

			// 获取所有SUM或COUNT表达式中的列名和表名
			var table2colNames = map[string] /*table name*/ map[string]bool /*col names*/ {}
			funcInfos := util.GetAllFunc(selectStmt)
			for _, funcInfo := range funcInfos {
				if strings.ToUpper(funcInfo.FuncName) == "SUM" || strings.ToUpper(funcInfo.FuncName) == "COUNT" {
					for _, column := range funcInfo.Columns {
						tableName := getTableName(column)
						if tableName == "" {
							continue
						}
						// 对列去重
						if table2colNames[tableName] == nil {
							table2colNames[tableName] = make(map[string]bool)
						}
						table2colNames[tableName][column.Name.L] = true
					}
				}
			}

			// 使用辅助函数检查该列的值是否全为NULL
			for tableName, colMap := range table2colNames {
				for colName := range colMap {
					// 列已去重
					isAllNull, err := util.IsColumnAllNull(input.Ctx, tableName, colName)
					if err != nil {
						return fmt.Errorf("error checking if column '%s' in table '%s' is all NULL: %v", colName, tableName, err)
					}
					// 如果该列的值全为NULL，报告规则违规
					if isAllNull {
						rulepkg.AddResult(input.Res, input.Rule, SQLE00122, colName)
						return nil
					}
				}
			}

		}

	default:
		// 如果节点不是DML语句，则不执行任何检查
		return nil
	}

	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
