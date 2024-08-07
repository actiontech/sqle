package parser

import (
	"github.com/actiontech/mybatis-mapper-2-sql/ast"
)

func SkipErrorQuery() func(*ast.Config) {
	return func(c *ast.Config) {
		c.SkipErrorQuery = true
	}
}

func RestoreOriginSql() func(*ast.Config) {
	return func(c *ast.Config) {
		c.RestoreOriginSql = true
	}
}
