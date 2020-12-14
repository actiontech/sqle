package inspector

import (
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

func (i *Inspect) Fingerprint(oneSql string) (fingerprint string, err error) {
	node, err := parseOneSql(i.Task.Instance.DbType, oneSql)
	node.Accept(&FingerprintVisitor{})
	q, err := restoreToSqlWithFlag(format.RestoreKeyWordUppercase|format.RestoreNameBackQuotes, node)
	if err != nil {
		return "", err
	}
	return q, nil
}

type FingerprintVisitor struct{}

func (f *FingerprintVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if v, ok := n.(*driver.ValueExpr); ok {
		v.Type.Charset = ""
		v.SetValue([]byte("?"))
	}
	return n, false
}

func (f *FingerprintVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}

