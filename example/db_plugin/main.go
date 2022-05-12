package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	adaptor "github.com/actiontech/sqle/sqle/pkg/driver"
	"github.com/actiontech/sqle/sqle/pkg/params"
	parser "github.com/pganalyze/pg_query_go/v2"
)

func main() {
	plugin := adaptor.NewAdaptor(&adaptor.PostgresDialector{})

	rule1 := &driver.Rule{
		Name:     "pg_rule_1",           // 规则ID，该值会与插件类型一起作为这条规则在 SQLE 的唯一标识
		Desc:     "避免查询所有的列",            // 规则描述
		Category: "DQL规范",               // 规则分类，用于分组，相同类型的规则会在 SQLE 的页面上展示在一起
		Level:    driver.RuleLevelError, // 规则等级，表示该规则的严重程度
	}
	rule1Handler := func(ctx context.Context, rule *driver.Rule, sql string) (string, error) {
		if strings.Contains(sql, "select *") {
			return rule.Desc, nil
		}
		return "", nil
	}

	// 定义第二条规则
	rule2 := &driver.Rule{
		Name:     "pg_rule_2",
		Desc:     "表字段不建议过多",
		Level:    driver.RuleLevelWarn,
		Category: "DDL规范",
		Params: []*params.Param{ // 自定义参数列表
			&params.Param{
				Key:   "max_column_count",  // 自定义参数的ID
				Value: "50",                // 自定义参数的默认值
				Desc:  "最大字段个数",            // 自定义参数在页面上的描述
				Type:  params.ParamTypeInt, // 自定义参数的值类型
			},
		},
	}

	// 这时处理函数的参数是 interface{} 类型，需要将其断言成 AST 语法树。
	rule2Handler := func(ctx context.Context, rule *driver.Rule, ast interface{}) (string, error) {
		node, ok := ast.(*parser.RawStmt)
		if !ok {
			return "", nil
		}
		switch stmt := node.GetStmt().GetNode().(type) {
		case *parser.Node_CreateStmt:
			columnCounter := 0
			for _, elt := range stmt.CreateStmt.TableElts {
				switch elt.GetNode().(type) {
				case *parser.Node_ColumnDef:
					columnCounter++
				}
			}
			// 读取 SQLE 传递过来的该参数配置的值
			count := rule.Params.GetParam("max_column_count").Int()
			if count > 0 && columnCounter > count {
				return fmt.Sprintf("表字段不建议超过%d个，目前有%d个", count, columnCounter), nil
			}
		}
		return "", nil
	}

	plugin.AddRule(rule1, rule1Handler)
	plugin.AddRuleWithSQLParser(rule2, rule2Handler)

	// 需要将 SQL 解析的方法注册到插件中。
	plugin.Serve(adaptor.WithSQLParser(func(sql string) (ast interface{}, err error) {
		// parser.Parse 使用 PostgreSQL 的解析器，将 sql 解析成 AST 语法树。
		result, err := parser.Parse(sql)
		if err != nil {
			return nil, fmt.Errorf("parse sql error")
		}
		if len(result.Stmts) != 1 {
			return nil, fmt.Errorf("unexpected statement count: %d", len(result.Stmts))
		}
		// 将 SQL 的语法树返回。
		return result.Stmts[0], nil
	}))

	plugin.Serve()
}
