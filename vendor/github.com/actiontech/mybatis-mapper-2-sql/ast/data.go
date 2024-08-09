package ast

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"
)

type DataNode interface {
	String() string
}

type Value string

func (v Value) String() string {
	return string(v)
}

type Param struct {
	Name string
}

func (p *Param) String() string {
	return "?"
}

type Variable struct {
	Name string
}

func (p *Variable) String() string {
	return "$"
}

type Data struct {
	tmp    bytes.Buffer
	reader *bytes.Reader
	Nodes  []DataNode
}

func newData(data []byte) *Data {
	return &Data{
		tmp:    bytes.Buffer{},
		reader: bytes.NewReader(data),
		Nodes:  []DataNode{},
	}
}

func (d *Data) read() (rune, error) {
	r, _, err := d.reader.ReadRune()
	if err != nil {
		return r, err
	}
	_, err = d.tmp.WriteRune(r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func (d *Data) unRead() {
	d.reader.Seek(-1, io.SeekCurrent)
}

func (d *Data) clean() {
	d.tmp.Reset()
}

func (d *Data) Scan(start *xml.StartElement) error {
	return nil
}

func (d *Data) AddChildren(ns ...Node) error {
	return nil
}

func (d *Data) GetStmt(ctx *Context) (string, error) {
	buff := bytes.Buffer{}
	for _, child := range d.Nodes {
		switch dt := child.(type) {
		case Value:
			buff.WriteString(dt.String())
		case *Param:
			buff.WriteString(dt.String())
		case *Variable:
			variable, ok := ctx.GetVariable(dt.Name)
			if !ok {
				buff.WriteString("?")
				//return "", fmt.Errorf("variable %s undifine", dt.Name)
			} else {
				buff.WriteString(variable)
			}
		}
	}
	return buff.String(), nil
}

func (d *Data) String() string {
	buff := bytes.Buffer{}
	for _, child := range d.Nodes {
		buff.WriteString(child.String())
	}
	return buff.String()
}

type MyBatisData struct {
	*Data
}

func NewMyBatisData(data []byte) *MyBatisData {
	d := &MyBatisData{}
	d.Data = newData(data)
	return d
}

func (d *MyBatisData) ScanData() error {
	for {
		var err error
		r, err := d.read()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return err
		}

		switch r {
		case '#':
			s, err := d.read()
			if err == io.EOF { // found end of element
				break
			}
			if s == '{' {
				err := d.scanParam()
				if err != nil {
					return err
				}
			}
		case '$':
			s, err := d.read()
			if err == io.EOF { // found end of element
				break
			}
			if s == '{' {
				err := d.scanVariable()
				if err != nil {
					return err
				}
			}
		default:
			err := d.scanValue()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *MyBatisData) scanParam() error {
	d.clean()
	for {
		r, err := d.read()
		if err == io.EOF {
			return fmt.Errorf("data is invalid, not found \"}\" for param")
		}
		if err != nil {
			return err
		}
		if r == '}' {
			break
		}
	}
	data := strings.TrimSuffix(d.tmp.String(), "}")
	d.Nodes = append(d.Nodes, &Param{Name: data})
	d.clean()
	return nil
}

func (d *MyBatisData) scanVariable() error {
	d.clean()
	for {
		r, err := d.read()
		if err == io.EOF {
			return fmt.Errorf("data is invalid, not found \"}\" for vaiable")
		}
		if err != nil {
			return err
		}
		if r == '}' {
			break
		}
	}
	data := strings.TrimSuffix(d.tmp.String(), "}")
	d.Nodes = append(d.Nodes, &Variable{Name: data})
	d.clean()
	return nil
}

func (d *MyBatisData) scanValue() error {
	var first rune
	var second rune
	for {
		r, err := d.read()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return err
		}
		if r == '#' || r == '$' {
			first = r
			s, err := d.read()
			if err == io.EOF { // found end of element
				break
			}
			second = s
			if s == '{' {
				d.unRead()
				d.unRead()
				break
			}
		}
	}
	data := strings.TrimSuffix(d.tmp.String(), string([]rune{first, second}))
	d.Nodes = append(d.Nodes, Value(data))
	d.clean()
	return nil
}

type IBatisData struct {
	*Data
}

func NewIBatisData(data []byte) *IBatisData {
	d := &IBatisData{}
	d.Data = newData(data)
	return d
}

func (d *IBatisData) ScanData() error {
	for {
		var err error
		r, err := d.read()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return err
		}
		switch r {
		case '#':
			err := d.scanParam()
			if err != nil {
				return err
			}
		case '$':
			err := d.scanVariable()
			if err != nil {
				return err
			}
		default:
			err := d.scanValue()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (d *IBatisData) scanParam() error {
	d.clean()
	for {
		r, err := d.read()
		if err == io.EOF {
			return fmt.Errorf("data is invalid, not found \"#\" for param")
		}
		if err != nil {
			return err
		}
		if r == '#' {
			break
		}
	}
	data := strings.TrimSuffix(d.tmp.String(), "#")
	d.Nodes = append(d.Nodes, &Param{Name: data})
	d.clean()
	return nil
}

func (d *IBatisData) scanVariable() error {
	d.clean()
	for {
		r, err := d.read()
		if err == io.EOF {
			return fmt.Errorf("data is invalid, not found \"$\" for vaiable")
		}
		if err != nil {
			return err
		}
		if r == '$' || r == '}' {
			break
		}
	}
	data := strings.TrimSuffix(d.tmp.String(), "$")
	d.Nodes = append(d.Nodes, &Variable{Name: data})
	d.clean()
	return nil
}

func (d *IBatisData) scanValue() error {
	var end rune
	for {
		r, err := d.read()
		if err == io.EOF { // found end of element
			break
		}
		if err != nil {
			return err
		}
		if r == '#' || r == '$' {
			end = r
			d.unRead()
			break
		}
	}
	data := strings.TrimSuffix(d.tmp.String(), string([]rune{end}))
	d.Nodes = append(d.Nodes, Value(data))
	d.clean()
	return nil
}
