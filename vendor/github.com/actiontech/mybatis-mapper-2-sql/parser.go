package parser

import (
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/actiontech/mybatis-mapper-2-sql/ast"
)

// ParseXML is a parser for parse all query in XML to string.
func ParseXML(data string) (string, error) {
	r := strings.NewReader(data)
	d := xml.NewDecoder(r)
	n, err := parse(d, nil)
	if err != nil {
		return "", err
	}
	if n == nil {
		return "", nil
	}
	stmt, err := n.GetStmt(ast.NewContext())
	if err != nil {
		return "", err
	}
	return stmt, nil
}

// ParseXMLQuery is a parser for parse all query in XML to []string one by one;
// you can set `skipErrorQuery` true to ignore invalid query.
func ParseXMLQuery(data string, skipErrorQuery bool) ([]string, error) {
	r := strings.NewReader(data)
	d := xml.NewDecoder(r)
	n, err := parse(d, nil)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, nil
	}
	m, ok := n.(*ast.Mapper)
	if !ok {
		return nil, fmt.Errorf("the mapper is not found")
	}
	stmts, err := m.GetStmts(ast.NewContext(), skipErrorQuery)
	if err != nil {
		return nil, err
	}
	return stmts, nil
}

func parse(d *xml.Decoder, start *xml.StartElement) (node ast.Node, err error) {
	if start != nil {
		node, err = scan(start)
		if err != nil {
			return nil, err
		}
	}

	for {
		t, err := d.Token()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return nil, err
		}

		switch tt := t.(type) {
		case xml.StartElement:
			child, err := parse(d, &tt)
			if err != nil {
				return nil, err
			}
			if child == nil {
				continue
			}
			if node == nil {
				node = child
			} else {
				err := node.AddChildren(child)
				if err != nil {
					return nil, err
				}
			}
		case xml.EndElement:
			if start != nil && tt.Name == start.Name {
				return node, nil
			}
		case xml.CharData:
			s := string(tt)
			if strings.TrimSpace(s) == "" {
				continue
			}
			d := ast.NewData(tt)
			d.ScanData()
			if node != nil {
				node.AddChildren(d)
			}
		default:
			continue
		}
	}
	return node, nil
}

func scan(start *xml.StartElement) (ast.Node, error) {
	var node ast.Node
	switch start.Name.Local {
	case "mapper":
		node = ast.NewMapper()
	case "sql":
		node = ast.NewSqlNode()
	case "include":
		node = ast.NewIncludeNode()
	case "property":
		node = ast.NewPropertyNode()
	case "select", "update", "delete", "insert":
		node = ast.NewQueryNode()
	case "if":
		node = ast.NewIfNode()
	case "choose":
		node = ast.NewChooseNode()
	case "when":
		node = ast.NewWhenNode()
	case "otherwise":
		node = ast.NewOtherwiseNode()
	case "where", "set", "trim":
		node = ast.NewTrimNode()
	case "foreach":
		node = ast.NewForeachNode()
	default:
		return nil, nil
		//return node, fmt.Errorf("unknow xml <%s>", start.Name.Local)
	}
	node.Scan(start)
	return node, nil
}
