package ast

import (
	"encoding/xml"
	"fmt"
)

type IncludeNode struct {
	RefId      DataNode
	Properties map[string]*PropertyNode
}

func NewIncludeNode() *IncludeNode {
	return &IncludeNode{
		Properties: map[string]*PropertyNode{},
	}
}

func (i *IncludeNode) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "refid" {
			data := NewMyBatisData([]byte(attr.Value))
			err := data.ScanData()
			if err != nil {
				return err
			}

			if len(data.Nodes) != 1 {
				return fmt.Errorf("include refid must be a variable or a string")
			}
			switch data.Nodes[0].(type) {
			case Value, *Variable:
			default:
				return fmt.Errorf("include refid must be a variable or a string")
			}
			i.RefId = data.Nodes[0]
		}
	}
	return nil
}

func (i *IncludeNode) AddChildren(ns ...Node) error {
	for _, n := range ns {
		switch nt := n.(type) {
		case *PropertyNode:
			if _, ok := i.Properties[nt.Name]; ok {
				return fmt.Errorf("property name %s is repeat", nt.Name)
			}
			i.Properties[nt.Name] = nt
		default:
		}
	}
	return nil
}

func (i *IncludeNode) GetStmt(ctx *Context) (string, error) {
	var refId string
	for _, p := range i.Properties {
		ctx.SetVariable(p.Name, p.Value)
	}
	switch it := i.RefId.(type) {
	case Value:
		refId = string(it)
	case *Variable:
		variable, ok := ctx.GetVariable(it.Name)
		if !ok {
			return "", fmt.Errorf("variable %s undifine", it.Name)
		}
		refId = variable
	}
	sql, ok := ctx.GetSql(refId)
	if !ok {
		return "", fmt.Errorf("sql %s is not exist", refId)
	}
	data, err := sql.GetStmt(ctx)
	if err != nil {
		return "", err
	}
	return data, nil
}

type PropertyNode struct {
	Name  string
	Value string
}

func NewPropertyNode() *PropertyNode {
	return &PropertyNode{}
}

func (p *PropertyNode) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "name" {
			p.Name = attr.Value
		}
		if attr.Name.Local == "value" {
			p.Value = attr.Value
		}
	}
	return nil
}

func (p *PropertyNode) AddChildren(ns ...Node) error {
	return nil
}

func (p *PropertyNode) GetStmt(ctx *Context) (string, error) {
	return "", nil
}
