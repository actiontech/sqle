package ast

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type Mapper struct {
	version        int
	NameSpace      string
	SqlNodes       map[string]*SqlNode
	QueryNodeIndex map[string]*QueryNode
	QueryNodes     []*QueryNode
	FilePath       string
}

func NewMapper() *Mapper {
	return &Mapper{
		SqlNodes:       map[string]*SqlNode{},
		QueryNodeIndex: map[string]*QueryNode{},
		QueryNodes:     []*QueryNode{},
	}
}

func (m *Mapper) AddChildren(ns ...Node) error {
	for _, n := range ns {
		switch nt := n.(type) {
		case *SqlNode:
			if _, ok := m.SqlNodes[nt.Id]; ok {
				return fmt.Errorf("sql id %s is repeat", nt.Id)
			}
			m.SqlNodes[nt.Id] = nt
		case *QueryNode:
			if _, ok := m.QueryNodeIndex[nt.Id]; ok {
				return fmt.Errorf("%s id %s is repeat", nt.Type, nt.Id)
			}
			m.QueryNodeIndex[nt.Id] = nt
			m.QueryNodes = append(m.QueryNodes, nt)
		}
	}
	return nil
}

func (m *Mapper) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "namespace" {
			m.NameSpace = attr.Value
		}
	}
	return nil
}

func (m *Mapper) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	ctx.Sqls = m.SqlNodes
	for _, a := range m.QueryNodes {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff.WriteString(data)
		if !strings.HasSuffix(strings.TrimSpace(data), ";") {
			buff.WriteString(";")
		}
		buff.WriteString("\n")
	}
	return strings.TrimSuffix(buff.String(), "\n"), nil
}

func (m *Mapper) GetStmts(ctx *Context, skipErrorQuery bool) ([]StmtInfo, error) {
	var stmts []StmtInfo
	if len(ctx.Sqls) == 0 {
		ctx.Sqls = m.SqlNodes
	}
	ctx.DefaultNamespace = m.NameSpace
	for _, a := range m.QueryNodes {
		data, err := a.GetStmt(ctx)
		if err == nil {
			stmts = append(stmts, StmtInfo{SQL: data, StartLine: a.StartLine})
			continue
		}
		if skipErrorQuery {
			continue
		}
		return nil, err
	}
	return stmts, nil
}
