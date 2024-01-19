package parser

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/actiontech/mybatis-mapper-2-sql/ast"
)

// ParseXML is a parser for parse all query in XML to string.
func ParseXML(data string) (string, error) {
	r := strings.NewReader(data)
	d := xml.NewDecoder(r)
	n, err := parse(d)
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

type XmlFile struct {
	FilePath string
	Content  string
}

// ParseXMLs is a parser for parse all query in several XML files to []ast.StmtInfo one by one;
// you can set `skipErrorQuery` true to ignore invalid query.
func ParseXMLs(data []XmlFile, skipErrorQuery bool) ([]ast.StmtInfo, error) {
	ms := ast.NewMappers()
	for _, data := range data {
		r := strings.NewReader(data.Content)
		d := xml.NewDecoder(r)
		n, err := parse(d)
		if err != nil {
			if skipErrorQuery {
				continue
			} else {
				return nil, err
			}
		}

		if n == nil {
			continue
		}

		m, ok := n.(*ast.Mapper)
		if !ok {
			if skipErrorQuery {
				continue
			} else {
				return nil, errors.New("the mapper is not found")
			}
		}
		m.FilePath = data.FilePath
		err = ms.AddMapper(m)
		if err != nil && !skipErrorQuery {
			return nil, fmt.Errorf("add mapper failed: %v", err)
		}
	}
	stmts, err := ms.GetStmts(skipErrorQuery)
	if err != nil {
		return nil, err
	}

	return stmts, nil
}

// ParseXMLQuery is a parser for parse all query in XML to []string one by one;
// you can set `skipErrorQuery` true to ignore invalid query.
func ParseXMLQuery(data string, skipErrorQuery bool) ([]string, error) {
	r := strings.NewReader(data)
	d := xml.NewDecoder(r)
	n, err := parse(d)
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
	sqls := []string{}
	for _, stmt := range stmts {
		sqls = append(sqls, stmt.SQL)
	}
	return sqls, nil
}

func parse(d *xml.Decoder) (node ast.Node, err error) {
	for {
		t, err := d.Token()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return nil, err
		}
		// get first start element
		if st, ok := t.(xml.StartElement); ok {
			switch st.Name.Local {
			case "mapper":
				return parseMyBatis(d, &st)
			case "sqlMap":
				return parseIBatis(d, &st)
			}
		}
	}
	return nil, nil
}

func parseMyBatis(d *xml.Decoder, start *xml.StartElement) (node ast.Node, err error) {
	node, err = scanMyBatis(d, start)
	if err != nil {
		return nil, err
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
			child, err := parseMyBatis(d, &tt)
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
			if tt.Name == start.Name {
				return node, nil
			}
		case xml.CharData:
			s := string(tt)
			if strings.TrimSpace(s) == "" {
				continue
			}
			d := ast.NewMyBatisData(tt)
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

func scanMyBatis(d *xml.Decoder, start *xml.StartElement) (ast.Node, error) {
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
		startLine, _ := d.InputPos()
		node = ast.NewQueryNode(uint64(startLine))
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

// ref: https://ibatis.apache.org/docs/java/pdf/iBATIS-SqlMaps-2_cn.pdf
func parseIBatis(d *xml.Decoder, start *xml.StartElement) (node ast.Node, err error) {
	node, err = scanIBatis(d, start)
	if err != nil {
		return nil, err
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
			child, err := parseIBatis(d, &tt)
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
			if tt.Name == start.Name {
				return node, nil
			}
		case xml.CharData:
			s := string(tt)
			if strings.TrimSpace(s) == "" {
				continue
			}
			d := ast.NewIBatisData(tt)
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

func scanIBatis(d *xml.Decoder, start *xml.StartElement) (ast.Node, error) {
	var node ast.Node
	switch start.Name.Local {
	case "sqlMap":
		node = ast.NewMapper()
	case "sql":
		node = ast.NewSqlNode()
	case "include":
		node = ast.NewIncludeNode()
	case "select", "update", "delete", "insert", "statement":
		startLine, _ := d.InputPos()
		node = ast.NewQueryNode(uint64(startLine))
	case "isEqual", "isNotEqual", "isGreaterThan", "isGreaterEqual", "isLessEqual",
		"isPropertyAvailable", "isNotPropertyAvailable", "isNull", "isNotNull", "isEmpty", "isNotEmpty":
		node = ast.NewConditionStmt()
	case "dynamic":
		node = ast.NewDynamicStmt()
	case "iterate":
		node = ast.NewIterateStmt()
	default:
		return nil, nil
		//return node, fmt.Errorf("unknow xml <%s>", start.Name.Local)
	}
	node.Scan(start)
	return node, nil
}
