package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00046 = "SQLE00046"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00046,
			Desc:       "在 MySQL 中, 数据库对象命名不建议大小写字母混合",
			Annotation: "数据库对象命名规范，不推荐采用大小写混用的形式建议词语之间使用下划线连接，提高代码可读性",
			Level:      driverV2.RuleLevelNotice,
			Category:   rulepkg.RuleTypeNamingConvention,
			Params:     params.Params{},
		},
		Message:      "在 MySQL 中, 数据库对象命名不建议大小写字母混合",
		AllowOffline: true,
		Func:         RuleSQLE00046,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00046): "在 MySQL 中，数据库对象命名不建议大小写字母混合."
您应遵循以下逻辑：
1、检查当前句子是ALTER还是CREATE类型，如果是ALTER句子，则进入检查步骤4；否则，进入检查步骤2。
2、检查CREATE句子中的目标对象名的语法节点是否存在大小写字母混合，如果存在，报告违反规则。
3、提供触发规则的SQL列表，并退出检查流程。
4、检查ALTER句子中是否存在ADD语法节点，如果存在，则进入下一步检查。
5、检查ADD目标对象名的语法节点是否存在大小写字母混合，如果存在，报告违反规则。
1、检查当前句子是ALTER还是RENAME类型，如果是ALTER句子，则进入检查步骤2；否则，进入检查步骤3。
2、检查ALTER句子中是否存在RENAME或CHANGE语法节点，如果存在，则进入下一步检查。
3、检查RENAME或CHANGE目标的新对象名的语法节点是否存在大小写字母混合，如果存在，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00046(input *rulepkg.RuleHandlerInput) error {
	objectNames := []string{}

	switch stmt := input.Node.(type) {
	case *ast.CreateDatabaseStmt:
		// schema
		objectNames = append(objectNames, stmt.Name)
	case *ast.CreateTableStmt:
		// table
		objectNames = append(objectNames, stmt.Table.Name.String())
		// column
		for _, col := range stmt.Cols {
			objectNames = append(objectNames, col.Name.Name.String())
		}
		// index
		for _, constraint := range stmt.Constraints {
			switch constraint.Tp {
			case ast.ConstraintUniqKey, ast.ConstraintKey, ast.ConstraintUniqIndex, ast.ConstraintIndex:
				objectNames = append(objectNames, constraint.Name)
			}
		}
	case *ast.CreateViewStmt:
		objectNames = append(objectNames, stmt.ViewName.Name.String())
	case *ast.CreateIndexStmt:
		objectNames = append(objectNames, stmt.IndexName)
	case *ast.AlterTableStmt:
		for _, spec := range stmt.Specs {
			switch spec.Tp {
			case ast.AlterTableRenameTable:
				// rename table
				objectNames = append(objectNames, spec.NewTable.Name.String())
			case ast.AlterTableAddColumns:
				// new column
				for _, col := range spec.NewColumns {
					objectNames = append(objectNames, col.Name.Name.String())
				}
			case ast.AlterTableChangeColumn:
				// rename column
				for _, col := range spec.NewColumns {
					objectNames = append(objectNames, col.Name.Name.String())
				}
			case ast.AlterTableAddConstraint:
				objectNames = append(objectNames, spec.Constraint.Name)
			case ast.AlterTableRenameIndex:
				objectNames = append(objectNames, spec.ToKey.String())
			}
		}
	case *ast.RenameTableStmt:
		for _, totable := range stmt.TableToTables {
			objectNames = append(objectNames, totable.NewTable.Name.String())
		}
	case *ast.UnparsedStmt:
		stmtPreifx := strings.ToUpper(input.Node.Text())
		if strings.HasPrefix(stmtPreifx, "CREATE") {
			// create event ...
			match1 := createEventReg.FindStringSubmatch(input.Node.Text())
			if len(match1) > 1 {
				fmt.Printf("sql '%s': %s\n", input.Node.Text(), match1[1])
				objectNames = append(objectNames, match1[1])
				break
			}
		} else if strings.HasPrefix(stmtPreifx, "ALTER") {
			// alter event ...
			match2 := alterEventReg.FindStringSubmatch(input.Node.Text())
			if len(match2) > 2 {
				fmt.Printf("sql '%s': %s\n", input.Node.Text(), match2[2])
				objectNames = append(objectNames, match2[2])
				break
			}
		} else if strings.HasPrefix(stmtPreifx, "RENAME") {
			// https://dev.mysql.com/doc/refman/8.4/en/rename-user.html
			// rename user old_user to new_user [,old_user to new_user] ...
			match3 := renameUserReg1.MatchString(input.Node.Text())
			if match3 {
				matches4 := renameUserReg2.FindAllStringSubmatch(input.Node.Text(), -1)
				for _, match := range matches4 {
					if len(match) > 1 {
						objectNames = append(objectNames, match[1])
					}
				}
				break
			}
		}
	}

	// 数据库对象命名不建议大小写字母混合
	for _, name := range objectNames {
		if name == "" {
			continue
		}
		// fmt.Printf("....name: '%s'", name)
		if utils.IsUpperAndLowerLetterMixed(name) {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00046)
			break
		}
	}
	return nil
}

// 规则函数实现结束
// ==== Rule code end ====
