//go:build enterprise
// +build enterprise

package tidb_audit_log

import (
	"github.com/actiontech/sqle/sqle/driver/mysql/util"

	"github.com/pingcap/parser/ast"
)

func GetMissingSchema(sql string, possibleSchema []string) (string, error) {
	stmts, err := util.ParseOneSql(sql)
	if err != nil {
		return "", err
	}

	visitor := &GetMissingSchemaVisitor{
		possibleSchema: possibleSchema,
		existSchema:    map[string]struct{}{},
	}
	stmts.Accept(visitor)
	return visitor.GetMissingSchema(), nil
}

// GetMissingSchemaVisitor implements ast.Visitor interface.
type GetMissingSchemaVisitor struct {
	possibleSchema []string
	existSchema    map[string]struct{}
}

// 因为一个SQL只能在一个库执行, 所以最多存在一个缺省的库名, 也可能没有缺失的库名
func (g *GetMissingSchemaVisitor) GetMissingSchema() string {
	for _, s := range g.possibleSchema {
		if _, ok := g.existSchema[s]; !ok {
			return s
		}
	}
	return ""
}

func (g *GetMissingSchemaVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if stmt, ok := n.(*ast.TableName); ok {
		if stmt.Schema.L != "" {
			g.existSchema[stmt.Schema.L] = struct{}{}
		}
	}
	return n, false
}

func (g *GetMissingSchemaVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}
