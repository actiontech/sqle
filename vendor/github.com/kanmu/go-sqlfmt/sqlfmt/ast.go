package sqlfmt

import (
	"fmt"
	"go/ast"
	"go/token"
	"log"
	"strings"

	"github.com/kanmu/go-sqlfmt/sqlfmt/parser/group"
)

// sqlfmt retrieves all strings from "Query" and "QueryRow" and "Exec" functions in .go file
const (
	QUERY    = "Query"
	QUERYROW = "QueryRow"
	EXEC     = "Exec"
)

// replaceAst replace ast node with formatted SQL statement
func replaceAst(f *ast.File, fset *token.FileSet, options *Options) {
	ast.Inspect(f, func(n ast.Node) bool {
		if x, ok := n.(*ast.CallExpr); ok {
			if fun, ok := x.Fun.(*ast.SelectorExpr); ok {
				funcName := fun.Sel.Name
				if funcName == QUERY || funcName == QUERYROW || funcName == EXEC {
					// not for parsing url.Query
					if len(x.Args) > 0 {
						if arg, ok := x.Args[0].(*ast.BasicLit); ok {
							sqlStmt := arg.Value
							if !strings.HasPrefix(sqlStmt, "`") {
								return true
							}
							src := strings.Trim(sqlStmt, "`")
							res, err := Format(src, options)
							if err != nil {
								log.Println(fmt.Sprintf("Format failed at %s: %v", fset.Position(arg.Pos()), err))
								return true
							}
							// FIXME
							// more elegant
							arg.Value = "`" + res + strings.Repeat(group.WhiteSpace, options.Distance) + "`"
						}
					}
				}
			}
		}
		return true
	})
}
