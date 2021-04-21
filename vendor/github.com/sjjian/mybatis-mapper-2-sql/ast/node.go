package ast

import (
	"bytes"
	"encoding/xml"
)

type Node interface {
	Scan(start *xml.StartElement) error
	AddChildren(ns ...Node) error
	GetStmt(ctx *Context) (string, error)
}

type ChildrenNode struct {
	Children []Node
}

func NewNode() *ChildrenNode {
	return &ChildrenNode{
		Children: []Node{},
	}
}

func (n *ChildrenNode) AddChildren(ns ...Node) error {
	n.Children = append(n.Children, ns...)
	return nil
}

func (n *ChildrenNode) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	for _, a := range n.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff.WriteString(data)
	}
	return buff.String(), nil
}
