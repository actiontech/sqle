package mybatis

import (
	"context"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/driver/mysql"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
)

func Parse(_ context.Context, sqlText string) ([]driver.Node, error) {
	nodes, err := ParseSql(sqlText)
	if err != nil {
		return nil, err
	}

	var ns []driver.Node
	for i := range nodes {
		n := driver.Node{}
		fingerprint, err := mysql.Fingerprint(nodes[i].Text(), true)
		if err != nil {
			return nil, err
		}
		n.Fingerprint = fingerprint
		n.Text = nodes[i].Text()
		switch nodes[i].(type) {
		case ast.DMLNode:
			n.Type = driver.SQLTypeDML
		default:
			n.Type = driver.SQLTypeDDL
		}

		ns = append(ns, n)
	}
	return ns, nil
}

func ParseSql(sql string) ([]ast.Node, error) {
	stmts, err := parseSql(sql)
	if err != nil {
		return nil, err
	}
	nodes := make([]ast.Node, 0, len(stmts))
	for _, stmt := range stmts {
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
