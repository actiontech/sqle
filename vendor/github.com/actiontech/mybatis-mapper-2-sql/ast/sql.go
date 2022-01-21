package ast

import (
	"bytes"
	"encoding/xml"
)

type SqlNode struct {
	*ChildrenNode
	Id string
}

func NewSqlNode() *SqlNode {
	n := &SqlNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (s *SqlNode) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "id" {
			s.Id = attr.Value
		}
	}
	return nil
}

func (s *SqlNode) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	for _, a := range s.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff.WriteString(data)
	}
	return buff.String(), nil
}
