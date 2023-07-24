package parse

import (
	"context"

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

	ns := make([]driverV2.Node, len(nodes))
	for i := range nodes {
		n := driverV2.Node{}
		fingerprint, err := util.Fingerprint(nodes[i].Text(), true)
		if err != nil {
			return nil, err
		}
		n.Fingerprint = fingerprint
		n.Text = nodes[i].Text()
		switch nodes[i].(type) {
		case ast.DMLNode:
			n.Type = driverV2.SQLTypeDML
		default:
			n.Type = driverV2.SQLTypeDDL
		}

		ns[i] = n
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
