package ast

import (
	"fmt"
	"strings"
)

type Context struct {
	QueryType        string // select, insert, update, delete
	Variable         map[string]string
	Sqls             map[string]*SqlNode
	DefaultNamespace string // namespace of current mapper
	Config           *Config
}

func NewContext(config *Config) *Context {
	return &Context{
		Variable: map[string]string{},
		Sqls:     map[string]*SqlNode{},
		Config:   config,
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
	// 这里的key某些xml在带有全限定名的时候需要去掉, 避免找不到SQL
	k = c.CutKeyDefaultNameSpace(k)
	sql, ok := c.Sqls[k]
	if ok {
		return sql, true
	}
	// 当存在跨namespace引用时，需要通过namespace区分引用的SQL id
	sql, ok = c.Sqls[fmt.Sprintf("%v.%v", c.DefaultNamespace, k)]
	return sql, ok
}

type Config struct {
	SkipErrorQuery   bool
	RestoreOriginSql bool
}

func (c *Context) CutKeyDefaultNameSpace(key string) string {
	afterKey, _ := strings.CutPrefix(key, c.DefaultNamespace+".")
	return afterKey
}

type ConfigFn func() func(*Config)
