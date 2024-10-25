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
	SQLE00015 = "SQLE00015"
)

func init() {
	rh := rulepkg.RuleHandler{
		Rule: driverV2.Rule{
			Name:       SQLE00015,
			Desc:       "对于MySQL的DDL, 建议使用规定的数据库排序规则",
			Annotation: "通过该规则约束全局的数据库排序规则，避免创建非预期的数据库排序规则，防止业务侧出现排序结果非预期等问题。建议项目内库表使用统一的字符集和字符集排序，部分连表查询的情况下字段的字符集或排序规则不一致可能会导致索引失效且不易发现",
			Level:      driverV2.RuleLevelError,
			Category:   rulepkg.RuleTypeDDLConvention,
			Params: params.Params{
				&params.Param{
					Key:   rulepkg.DefaultSingleParamKeyName,
					Value: "utf8mb4_0900_ai_ci",
					Desc:  "数据库排序规则",
					Type:  params.ParamTypeString,
				},
			},
		},
		Message: "对于MySQL的DDL, 建议使用规定的数据库排序规则为%s",
		AllowOffline: false,
		Func:    RuleSQLE00015,
	}
	rulepkg.RuleHandlers = append(rulepkg.RuleHandlers, rh)
	rulepkg.RuleHandlerMap[rh.Rule.Name] = rh
}

/*
==== Prompt start ====
In MySQL, you should check if the SQL violate the rule(SQLE00015): "In DDL, the specified collation should be used, the specified collation should be a parameter whose default value is utf8mb4_0900_ai_ci.".
You should follow the following logic:
1. For "create table ... collate= ..." statement, check collate which on the table option or on the column should be the same as the specified collation, otherwise, report a violation
2. For "alter table ... collate= ..." statement, check collate which on the table option or on the column should be the same as the specified collation, otherwise, report a violation
2. For "create database ... collate= ..." statement, check collate should be the same as the specified collation, otherwise, report a violation
2. For "alter database ... collate= ..." statement, check collate should be the same as the specified collation, otherwise, report a violation
==== Prompt end ====
*/

// ==== Rule code start ====
func RuleSQLE00015(input *rulepkg.RuleHandlerInput) error {
	// get expected param value
	param := input.Rule.Params.GetParam(rulepkg.DefaultSingleParamKeyName)
	if param == nil {
		return fmt.Errorf("param %s not found", rulepkg.DefaultSingleParamKeyName)
	}

	expectedCollation := param.Value
	var columnCollations []string

	switch stmt := input.Node.(type) {
	case *ast.CreateTableStmt:
		// "create table ..."

		// get collate in column definition
		for _, col := range stmt.Cols {
			if option := util.GetColumnOption(col, ast.ColumnOptionCollate); nil != option {
				columnCollations = append(columnCollations, option.StrValue)
			}
		}

		// get collate in table option
		if option := util.GetTableOption(stmt.Options, ast.TableOptionCollate); nil != option {
			columnCollations = append(columnCollations, option.StrValue)
		}

		// if create table not define collate, using default
		if len(columnCollations) == 0 {
			c, err := input.Ctx.GetCollationDatabase(stmt.Table, "")
			if err != nil {
				return err
			}
			columnCollations = append(columnCollations, c)
		}

	case *ast.AlterTableStmt:
		// "alter table"

		for _, spec := range stmt.Specs {

			// get collate in column definition
			for _, col := range spec.NewColumns {
				if option := util.GetColumnOption(col, ast.ColumnOptionCollate); nil != option {
					columnCollations = append(columnCollations, option.StrValue)
				}
			}

			// get collate in table option
			if option := util.GetTableOption(spec.Options, ast.TableOptionCollate); nil != option {
				columnCollations = append(columnCollations, option.StrValue)
			}
		}

	case *ast.CreateDatabaseStmt:
		// "create database ..."

		if option := util.GetDatabaseOption(stmt.Options, ast.DatabaseOptionCollate); nil != option {
			columnCollations = append(columnCollations, option.Value)
		}

		// if create database not define collate, using default
		if len(columnCollations) == 0 {
			c, err := input.Ctx.GetCollationDatabase(nil, stmt.Name)
			if err != nil {
				return err
			}
			columnCollations = append(columnCollations, c)
		}

	case *ast.AlterDatabaseStmt:
		// "alter database ..."

		if option := util.GetDatabaseOption(stmt.Options, ast.DatabaseOptionCollate); nil != option {
			columnCollations = append(columnCollations, option.Value)
		}
	default:
		return nil
	}

	for _, cs := range columnCollations {
		if !strings.EqualFold(cs, expectedCollation) {
			// the collate is not the same as param
			rulepkg.AddResult(input.Res, input.Rule, SQLE00015, expectedCollation)
		}
	}
	return nil
}

// ==== Rule code end ====
