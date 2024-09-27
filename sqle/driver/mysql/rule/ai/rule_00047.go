package ai

import (
	"fmt"
	"strings"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser/ast"
)

const (
	SQLE00047 = "SQLE00047"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00047,
			Desc:       "在 MySQL 中, 数据库对象名称的字符个数不建议超过阈值",
			Annotation: "通过配置该规则可以规范指定业务的对象命名长度，具体长度可以自定义设置。",
			Level:      driverV2.RuleLevelWarn,
			Category:   rulepkg.RuleTypeNamingConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "64",
					Desc:  "字符个数",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message:      "在 MySQL 中, 数据库对象名称的字符个数不建议超过阈值:%v",
		AllowOffline: true,
		Func:         RuleSQLE00047,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00047): "在 MySQL 中，数据库对象名称的字符个数不建议超过阈值.字符个数:64"
您应遵循以下逻辑：
1、检查当前句子是ALTER还是CREATE类型，如果是ALTER句子，则进入检查步骤4；否则，进入检查步骤2。
2、检查CREATE句子中的目标对象名的长度是否超过阈值，如果是，报告违反规则。
3、提供触发规则的SQL列表，并退出检查流程。
4、检查ALTER句子中是否存在ADD语法节点，如果存在，则进入下一步检查。
5、检查ADD目标对象名的长度是否超过阈值，如果是，报告违反规则。
6、检查当前句子是ALTER还是RENAME类型，如果是ALTER句子，则进入检查步骤7；否则，进入检查步骤8。
7、检查ALTER句子中是否存在RENAME语法节点，如果存在，则进入下一步检查。
8、检查RENAME目标对象名的长度是否超过阈值，如果是，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00047(input *rulepkg.RuleHandlerInput) error {
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}
	maxLength := param.Int()

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
	case *ast.CreateUserStmt:
		for _, spec := range stmt.Specs {
			objectNames = append(objectNames, spec.User.Username)
		}
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
				objectNames = append(objectNames, match1[1])
				break
			}
		} else if strings.HasPrefix(stmtPreifx, "ALTER") {
			// alter event ...
			match2 := alterEventReg.FindStringSubmatch(input.Node.Text())
			if len(match2) > 2 {
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

	// 数据库对象名称的字符个数不建议超过阈值:64
	for _, name := range objectNames {
		if name == "" {
			continue
		}
		//fmt.Printf("....name: '%s'", name)
		if len(name) > maxLength {
			rulepkg.AddResult(input.Res, input.Rule, SQLE00047, maxLength)
			break
		}
	}
	return nil
}

// ==== Rule code end ====
