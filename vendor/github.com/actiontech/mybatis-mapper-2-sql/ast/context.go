package ast

import "fmt"

type Context struct {
	QueryType        string // select, insert, update, delete
	Variable         map[string]string
	Sqls             map[string]*SqlNode
	DefaultNamespace string // namespace of current mapper
}

func NewContext() *Context {
	return &Context{
		Variable: map[string]string{},
		Sqls:     map[string]*SqlNode{},
	}
}

func (c *Context) GetVariable(k string) (string, bool) {
	variable, ok := c.Variable[k]
	return variable, ok
}

func (c *Context) SetVariable(k, v string) {
	c.Variable[k] = v
}

func (c *Context) GetSql(k string) (*SqlNode, bool) {
	sql, ok := c.Sqls[k]
	if ok {
		return sql, true
	}
	// 当存在跨namespace引用时，需要通过namespace区分引用的SQL id
	sql, ok = c.Sqls[fmt.Sprintf("%v.%v", c.DefaultNamespace, k)]
	return sql, ok
}
