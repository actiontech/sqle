package ai

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"

	rulepkg "github.com/actiontech/sqle/sqle/driver/mysql/rule"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/pingcap/parser/ast"

	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
)

const (
	SQLE00048 = "SQLE00048"
)

func init() {
	rh := rulepkg.SourceHandler{
		Rule: rulepkg.SourceRule{
			Name:       SQLE00048,
			Desc:       plocale.Rule00048Desc,
			Annotation: plocale.Rule00048Annotation,
			Category:   plocale.RuleTypeNamingConvention,
			CategoryTags: map[string][]string{
				plocale.RuleCategoryOperand.ID:              {plocale.RuleTagDatabase.ID, plocale.RuleTagTable.ID, plocale.RuleTagColumn.ID, plocale.RuleTagIndex.ID, plocale.RuleTagView.ID, plocale.RuleTagProcedure.ID, plocale.RuleTagFunction.ID, plocale.RuleTagTrigger.ID, plocale.RuleTagEvent.ID},
				plocale.RuleCategorySQL.ID:                  {plocale.RuleTagDDL.ID},
				plocale.RuleCategoryAuditPurpose.ID:         {plocale.RuleTagMaintenance.ID},
				plocale.RuleCategoryAuditAccuracy.ID:        {plocale.RuleTagOffline.ID},
				plocale.RuleCategoryAuditPerformanceCost.ID: {},
			},
			Level:        driverV2.RuleLevelWarn,
			Params:       []*rulepkg.SourceParam{},
			Knowledge:    driverV2.RuleKnowledge{},
			AllowOffline: true,
		},
		Message: plocale.Rule00048Message,
		Func:    RuleSQLE00048,
	}
	sourceRuleHandlers = append(sourceRuleHandlers, &rh)
}

/*
==== Prompt start ====
在 MySQL 中，您应该检查 SQL 是否违反了规则(SQLE00048): "在 MySQL 中，数据库对象命名只能使用英文、下划线或数字，首字母必须是英文."
您应遵循以下逻辑：
1、检查当前句子是ALTER还是CREATE类型，如果是ALTER句子，则进入检查步骤5；否则，进入检查步骤2。
2、检查CREATE句子中的目标对象名的首个字符是否英文字母，如果不是，报告违反规则。
3、检查CREATE句子中的目标对象名是否存在除了英文字母、下划线、数字外的其他字符，如果是，报告违反规则。
4、提供触发规则的SQL列表，并退出检查流程。
5、检查ALTER句子中是否存在ADD语法节点，如果存在，则进入下一步检查。
6、检查ADD目标对象名的首个字符是否英文字母，如果不是，报告违反规则。
7、检查ADD目标对象名是否存在除了英文字母、下划线、数字外的其他字符，如果是，报告违反规则。
1、检查当前句子是ALTER还是RENAME类型，如果是ALTER句子，则进入检查步骤2；否则，进入检查步骤3。
2、检查ALTER句子中是否存在RENAME或CHANGE语法节点，如果存在，则进入下一步检查。
3、检查RENAME或CHANGE目标对象名的首个字符是否英文字母，如果不是，报告违反规则。
4、检查RENAME或CHANGE目标的新对象名是否存在除了英文字母、下划线、数字外的其他字符，如果是，报告违反规则。
==== Prompt end ====
*/

// ==== Rule code start ====
// 规则函数实现开始
func RuleSQLE00048(input *rulepkg.RuleHandlerInput) error {
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

	// 数据库对象命名只能使用英文、下划线或数字，首字母必须是英文
	for _, name := range objectNames {
		if name == "" {
			continue
		}
		// fmt.Printf("....name: '%s'", name)
		if !unicode.Is(unicode.Latin, rune(name[0])) ||
			bytes.IndexFunc([]byte(name), func(r rune) bool {
				return !(unicode.Is(unicode.Latin, r) || string(r) == "_" || unicode.IsDigit(r))
			}) != -1 {

			rulepkg.AddResult(input.Res, input.Rule, SQLE00048)
			break
		}
	}
	return nil
}

// 规则函数实现结束

var createEventReg = regexp.MustCompile(`(?i)\bcreate\s+(?:definer\s*=\s*'(?:[^'\\]|\\.)*'@\s*'(?:[^'\\]|\\.)*'\s+)?event\s+(?:if\s+not\s+exists\s+)?(\w+)\s+ON SCHEDULE`)
var alterEventReg = regexp.MustCompile(`(?i)\balter\s+(?:definer\s*=\s*\S+\s+)?event\s+(\w+)\s+(?:on\s+scheduler\s+\S+\s+)?(?:on\s+completion\s+(?:not\s+)?preserve\s+)?rename\s+to\s+(\w+)\s?`)

var renameUserReg1 = regexp.MustCompile(`(?i)\bRENAME\s+USER\s+.+`)
var renameUserReg2 = regexp.MustCompile(`(?i)TO ['"]?(\w+)['"]?@?['"]?[\w.]*['"]?`)

// ==== Rule code end ====
