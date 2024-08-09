package ast

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type IfNode struct {
	*ChildrenNode
	Expression string
}

func NewIfNode() *IfNode {
	n := &IfNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *IfNode) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "test" {
			n.Expression = attr.Value
		}
	}
	return nil
}

type ChooseNode struct {
	When      []*WhenNode
	Otherwise *OtherwiseNode
}

func NewChooseNode() *ChooseNode {
	return &ChooseNode{
		When: []*WhenNode{},
	}
}

func (n *ChooseNode) Scan(start *xml.StartElement) error {
	return nil
}

func (n *ChooseNode) AddChildren(ns ...Node) error {
	for _, node := range ns {
		switch nt := node.(type) {
		case *WhenNode:
			n.When = append(n.When, nt)
		case *OtherwiseNode:
			if n.Otherwise != nil {
				return fmt.Errorf("otherwise is repeat in <choose>")
			}
			n.Otherwise = nt
		default:
			return fmt.Errorf("data is invalid in <choose>")
		}
	}
	return nil
}

func (n *ChooseNode) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}

	for i, a := range n.When {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		// https://github.com/actiontech/sqle/issues/302
		// In some cases, users like to write XML like:
		/*
			<select id="selectUserByState" resultType="com.bz.model.entity.User">
			    SELECT * FROM user
			    <choose>
			        <when test="state == 1">
			            where name = #{name1}
			        </when>
			        <otherwise>
			            where name = #{name2}
			        </otherwise>
			    </choose>
			</select>
		*/
		// parer it as "where name = ? and name = ?".
		//strings.
		if i > 0 {
			data = replaceWhere(data)
		}
		buff.WriteString(data)
	}
	// https://github.com/actiontech/sqle/issues/639
	// otherwise can be not defined. so ChooseNode -> Otherwise may be nil.
	/*
		<choose>
			<when test="state == 1">
				where name = #{name1}
			</when>
		</choose>
	*/
	if n.Otherwise != nil {
		data, err := n.Otherwise.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff.WriteString(replaceWhere(data))
	}
	return buff.String(), nil
}

type WhenNode struct {
	*ChildrenNode
	Expression string
}

func NewWhenNode() *WhenNode {
	n := &WhenNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *WhenNode) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "test" {
			n.Expression = attr.Value
		}
	}
	return nil
}

type OtherwiseNode struct {
	*ChildrenNode
}

func NewOtherwiseNode() *OtherwiseNode {
	n := &OtherwiseNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *OtherwiseNode) Scan(start *xml.StartElement) error {
	return nil
}

type TrimNode struct {
	*ChildrenNode
	Name            string
	Prefix          string
	Suffix          string
	PrefixOverrides []string
	SuffixOverrides []string
}

func NewTrimNode() *TrimNode {
	n := &TrimNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *TrimNode) Scan(start *xml.StartElement) error {
	n.Name = start.Name.Local
	switch start.Name.Local {
	case "where":
		n.Prefix = "WHERE"
		n.PrefixOverrides = []string{"and ", "or ", "AND ", "OR "}
	case "set":
		n.Prefix = "SET"
		n.SuffixOverrides = []string{","}
	case "trim":
		for _, attr := range start.Attr {
			if attr.Name.Local == "prefix" {
				n.Prefix = attr.Value
			}
			if attr.Name.Local == "suffix" {
				n.Suffix = attr.Value
			}
			if attr.Name.Local == "prefixOverrides" {
				n.PrefixOverrides = strings.Split(attr.Value, "|")
			}
			if attr.Name.Local == "suffixOverrides" {
				n.SuffixOverrides = strings.Split(attr.Value, "|")
			}
		}
	}
	return nil
}

func (n *TrimNode) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	for _, a := range n.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff.WriteString(data)
	}
	body := strings.TrimSpace(buff.String())
	for _, po := range n.PrefixOverrides {
		body = strings.TrimPrefix(body, po)
	}
	for _, so := range n.SuffixOverrides {
		body = strings.TrimSuffix(body, so)
	}
	buff.Reset()
	buff.WriteString(n.Prefix)
	buff.WriteString(" ")
	buff.WriteString(body)
	buff.WriteString(n.Suffix)
	return buff.String(), nil
}

type ForeachNode struct {
	*ChildrenNode
	Open      string
	Close     string
	Separator string
}

func NewForeachNode() *ForeachNode {
	n := &ForeachNode{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *ForeachNode) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "open" {
			n.Open = attr.Value
		}
		if attr.Name.Local == "close" {
			n.Close = attr.Value
		}
		if attr.Name.Local == "separator" {
			n.Separator = attr.Value
		}
	}
	return nil
}

func (n *ForeachNode) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}

	body := make([]string, 0, len(n.Children))
	for _, a := range n.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		body = append(body, data)
	}
	if len(body) == 0 {
		return "", nil
	}
	if n.Open != "" {
		buff.WriteString(n.Open)
	}
	if len(body) == 1 {
		buff.WriteString(body[0])
		if n.Separator != "" {
			buff.WriteString(n.Separator)
		}
		buff.WriteString(body[0])
	} else {
		buff.WriteString(strings.Join(body, fmt.Sprintf(" %s ", n.Separator)))
	}
	if n.Close != "" {
		buff.WriteString(n.Close)
	}
	return buff.String(), nil
}

type ConditionStmt struct {
	*ChildrenNode
	typ     string
	prepEnd string
}

func NewConditionStmt() *ConditionStmt {
	n := &ConditionStmt{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *ConditionStmt) Scan(start *xml.StartElement) error {
	n.typ = start.Name.Local
	for _, attr := range start.Attr {
		if attr.Name.Local == "prepend" {
			n.prepEnd = attr.Value
		}
	}
	return nil
}

func (n *ConditionStmt) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	data, err := n.ChildrenNode.GetStmt(ctx)
	if err != nil {
		return "", err
	}
	if n.prepEnd != "" {
		buff.WriteString(n.prepEnd)
		buff.WriteString(" ")
	}
	buff.WriteString(data)
	return buff.String(), nil
}

type IterateStmt struct {
	*ForeachNode
	prepEnd string
}

func NewIterateStmt() *IterateStmt {
	n := &IterateStmt{}
	n.ForeachNode = NewForeachNode()
	return n
}

func (n *IterateStmt) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "open" {
			n.Open = attr.Value
		}
		if attr.Name.Local == "close" {
			n.Close = attr.Value
		}
		if attr.Name.Local == "conjunction" {
			n.Separator = attr.Value
		}
		if attr.Name.Local == "prepend" {
			n.prepEnd = attr.Value
		}
	}
	return nil
}

func (n *IterateStmt) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	data, err := n.ForeachNode.GetStmt(ctx)
	if err != nil {
		return "", err
	}
	if data != "" && n.prepEnd != "" {
		buff.WriteString(n.prepEnd)
	}
	buff.WriteString(data)
	return buff.String(), nil
}

type StrStmt struct {
	*ChildrenNode
}

func NewStrStmt() *StrStmt {
	n := &StrStmt{}
	n.ChildrenNode = NewNode()
	return n
}

func (s *StrStmt) Scan(start *xml.StartElement) error {
	return nil
}

func (s *StrStmt) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	for _, a := range s.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		// 如果解析出来的语句仅仅只有一个变量,不处理,直接跳过,因为没有意义
		// <str type="Str"><![CDATA[${setID}]]></str>
		if data == "?" {
			continue
		}

		buff.WriteString(" ")
		buff.WriteString(data)
		buff.WriteString(" ")
	}
	return buff.String(), nil
}

type DynamicStmt struct {
	*ChildrenNode
	prepEnd string
}

func NewDynamicStmt() *DynamicStmt {
	n := &DynamicStmt{}
	n.ChildrenNode = NewNode()
	return n
}

func (n *DynamicStmt) Scan(start *xml.StartElement) error {
	for _, attr := range start.Attr {
		if attr.Name.Local == "prepend" {
			n.prepEnd = attr.Value
		}
	}
	return nil
}

func (n *DynamicStmt) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	for i, a := range n.Children {
		// the first condition
		if i == 0 {
			switch c := a.(type) {
			case *ConditionStmt:
				c.prepEnd = n.prepEnd
			case *IterateStmt:
				c.prepEnd = n.prepEnd
			default:
				buff.WriteString(n.prepEnd)
			}
		}
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		if i != 0 {
			d := strings.ToLower(strings.TrimSpace(data))
			if !strings.HasPrefix(d, "and") && !strings.HasPrefix(d, "or") &&
				!strings.HasPrefix(d, strings.ToLower(n.prepEnd)) {
				buff.WriteString("AND ")
			}
		}
		buff.WriteString(data)
	}
	return buff.String(), nil
}

type AndNodeStmt struct {
	*ChildrenNode
}

func NewAndNode() *AndNodeStmt {
	n := &AndNodeStmt{}
	n.ChildrenNode = NewNode()
	return n
}

func (a *AndNodeStmt) Scan(start *xml.StartElement) error {
	return nil
}

func (a *AndNodeStmt) AddChildren(ns ...Node) error {
	a.Children = append(a.Children, ns...)
	return nil
}

func (a *AndNodeStmt) GetStmt(ctx *Context) (string, error) {
	andList := make([]string, 0, len(a.Children))
	for _, a := range a.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff := bytes.Buffer{}
		buff.WriteString("(")
		buff.WriteString(data)
		buff.WriteString(")")
		andList = append(andList, buff.String())
	}
	return strings.Join(andList, "AND"), nil
}

type OrNodeStmt struct {
	*ChildrenNode
}

func NewOrNode() *OrNodeStmt {
	n := &OrNodeStmt{}
	n.ChildrenNode = NewNode()
	return n
}

func (a *OrNodeStmt) Scan(start *xml.StartElement) error {
	return nil
}

func (a *OrNodeStmt) AddChildren(ns ...Node) error {
	a.Children = append(a.Children, ns...)
	return nil
}

func (a *OrNodeStmt) GetStmt(ctx *Context) (string, error) {
	orList := make([]string, 0, len(a.Children))
	for _, a := range a.Children {
		data, err := a.GetStmt(ctx)
		if err != nil {
			return "", err
		}
		buff := bytes.Buffer{}
		buff.WriteString("(")
		buff.WriteString(data)
		buff.WriteString(")")
		orList = append(orList, buff.String())
	}
	return strings.Join(orList, "OR"), nil
}
