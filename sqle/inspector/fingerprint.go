package inspector

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	driver "github.com/pingcap/tidb/types/parser_driver"
)

func Fingerprint(oneSql string, isCaseSensitive bool) (fingerprint string, err error) {
	stmts, _, err := parser.New().PerfectParse(oneSql, "", "")
	if err != nil {
		return "", err
	}
	if len(stmts) != 1 {
		return "", parser.ErrSyntax
	}

	stmts[0].Accept(&FingerprintVisitor{})
	if !isCaseSensitive {
		stmts[0].Accept(&CapitalizeProcessor{
			capitalizeTableName:      true,
			capitalizeTableAliasName: true,
			capitalizeDatabaseName:   true,
		})
	}
	fingerprint, err = restoreToSqlWithFlag(format.RestoreKeyWordUppercase|format.RestoreNameBackQuotes, stmts[0])
	if err != nil {
		return "", err
	}
	return
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
