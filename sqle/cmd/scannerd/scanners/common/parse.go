package common

import (
	"context"
	"regexp"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/util"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

func Parse(_ context.Context, sqlText string) ([]driverV2.Node, error) {
	nodes, err := ParseSql(sqlText)
	if err != nil {
		return nil, err
	}

	ns := make([]driverV2.Node, 0, len(nodes))
	for i := range nodes {
		n := driverV2.Node{}
		n.Text = nodes[i].Text()
		switch nodes[i].(type) {
		case *ast.UnparsedStmt:
			nodes[i].SetText(clearComments(nodes[i].Text()))
			// TODO https://github.com/actiontech/sqle-ee/issues/1075 未解析节点的类型是未知的，是否应该新增一个unknown类型，作为未解析的SQL的类型
			n.Type = driverV2.SQLTypeDDL
		case ast.DMLNode:
			n.Type = driverV2.SQLTypeDML
		default:
			n.Type = driverV2.SQLTypeDDL
		}

		n.Fingerprint, err = util.Fingerprint(nodes[i].Text(), true)
		if err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, nil
}

// TODO https://github.com/actiontech/sqle-ee/issues/1075 这是临时方案的函数，暂时不放到公共工具中
func clearComments(sqlText string) string {
	// 将注释替换为一个空格，防止语句粘连
	sqlText = regexp.MustCompile(`(?s)/\*.*?\*/`).ReplaceAllString(sqlText, " ")
	// 去除结尾分号后的内容
	idx := strings.Index(sqlText, ";")
	if idx >= 0 {
		sqlText = sqlText[:idx]
	}
	// 去除开头结尾的空格，并且替换中间连续的空格为一个空格
	sqlText = strings.Join(strings.Fields(sqlText), " ")
	return sqlText
}

func ParseSql(sql string) ([]ast.Node, error) {
	stmts, err := parseSql(sql)
	if err != nil {
		return nil, err
	}
	nodes := make([]ast.Node, 0, len(stmts))
	for _, stmt := range stmts {
		// node can only be ast.Node
		//nolint:forcetypeassert
		node := stmt.(ast.Node)
		nodes = append(nodes, node)
	}
	return nodes, nil
}

func parseSql(sql string) ([]ast.StmtNode, error) {
	p := parser.New()
	stmts, _, err := p.PerfectParse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmts, nil
}
