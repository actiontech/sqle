// Package parse implements parsing of Go structs in files and runtime.
//
// This package, despite containing exported types, methods and functions,
// is an internal part of implementation of 'reform' command, also used by generated files,
// and not a part of public stable API.
package parse // import "gopkg.in/reform.v1/parse"

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldInfo represents information about struct field.
type FieldInfo struct {
	Name   string // field name as defined in source file, e.g. Name
	Type   string // field type as defined in source file, e.g. string; always present for primary key, may be absent otherwise
	Column string // SQL database column name from "reform:" struct field tag, e.g. name
}

// fieldInfoInSync returns true if FieldInfo fields that are set by both file and runtime parser are equal.
func fieldInfoInSync(fi1, fi2 *FieldInfo) bool {
	if fi1 == nil {
		panic("fi1 is nil")
	}
	if fi2 == nil {
		panic("fi2 is nil")
	}

	return fi1.Name == fi2.Name &&
		fi1.Type == fi2.Type &&
		fi1.Column == fi2.Column
}

// GoString returns struct field information as Go code string.
func (fi *FieldInfo) GoString() string {
	return fmt.Sprintf("{Name: %q, Type: %q, Column: %q}", fi.Name, fi.Type, fi.Column)
}

// StructInfo represents information about struct.
type StructInfo struct {
	Type         string      // struct type as defined in source file, e.g. User
	SQLSchema    string      // SQL database schema name from magic "reform:" comment, e.g. public
	SQLName      string      // SQL database view or table name from magic "reform:" comment, e.g. users
	Fields       []FieldInfo // fields info
	PKFieldIndex int         // index of primary key field in Fields, -1 if none
}

// structInfoInSync returns true if FieldInfo fields that are set by both file and runtime parser are equal.
func structInfoInSync(si1, si2 *StructInfo) bool {
	if si1 == nil {
		panic("si1 is nil")
	}
	if si2 == nil {
		panic("si2 is nil")
	}

	inSync := si1.Type == si2.Type &&
		si1.SQLSchema == si2.SQLSchema &&
		si1.SQLName == si2.SQLName &&
		si1.PKFieldIndex == si2.PKFieldIndex
	if !inSync {
		return false
	}

	if len(si1.Fields) != len(si2.Fields) {
		return false
	}
	for i := range si1.Fields {
		if inSync = fieldInfoInSync(&si1.Fields[i], &si2.Fields[i]); !inSync {
			return false
		}
	}
	return true
}

// GoString returns struct information as Go code string.
func (s *StructInfo) GoString() string {
	res := "parse.StructInfo{\n"

	res += fmt.Sprintf("\tType: %q,\n", s.Type)
	if s.SQLSchema != "" {
		res += fmt.Sprintf("\tSQLSchema: %q,\n", s.SQLSchema)
	}
	res += fmt.Sprintf("\tSQLName: %q,\n", s.SQLName)

	res += "\tFields: []parse.FieldInfo{\n"
	for _, f := range s.Fields {
		res += fmt.Sprintf("\t\t%s,\n", f.GoString())
	}
	res += "\t},\n"

	res += fmt.Sprintf("\tPKFieldIndex: %d,\n", s.PKFieldIndex)

	res += "}"
	return res
}

// Columns returns a new slice of column names.
func (s *StructInfo) Columns() []string {
	res := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		res[i] = f.Column
	}
	return res
}

// ColumnsGoString returns column names as Go code string.
func (s *StructInfo) ColumnsGoString() string {
	res := make([]string, len(s.Fields))
	for i, f := range s.Fields {
		res[i] = strconv.Quote(f.Column)
	}
	return "[]string{\n\t" + strings.Join(res, ",\n\t") + ",\n}"
}

// IsTable returns true if this object represent information for table, false for view.
func (s *StructInfo) IsTable() bool {
	return s.PKFieldIndex >= 0
}

// PKField returns a primary key field, panics for views.
func (s *StructInfo) PKField() FieldInfo {
	if !s.IsTable() {
		panic("reform: not a table")
	}
	return s.Fields[s.PKFieldIndex]
}

// AssertUpToDate checks that given StructInfo matches given object.
// It is used during program initialization to check that generated files are up-to-date.
func AssertUpToDate(si *StructInfo, obj interface{}) {
	msg := fmt.Sprintf(`reform:
		%s struct information is not up-to-date.
		Typically this means that %s type definition was changed, but 'reform' command / 'go generate' was not run.

		`, si.Type, si.Type)
	si2, err := Object(obj, si.SQLSchema, si.SQLName)
	if err != nil {
		panic(msg + err.Error())
	}
	if !structInfoInSync(si, si2) {
		panic(msg)
	}
}

// parseStructFieldTag is used by both file and runtime parsers
func parseStructFieldTag(tag string) (sqlName string, isPK bool) {
	parts := strings.Split(tag, ",")
	if len(parts) == 0 || len(parts) > 2 {
		return
	}

	if len(parts) == 2 {
		switch parts[1] {
		case "pk":
			isPK = true
		default:
			return
		}
	}

	sqlName = parts[0]
	return
}

// checkFields is used by both file and runtime parsers
func checkFields(res *StructInfo) error {
	if len(res.Fields) == 0 {
		return fmt.Errorf(`reform: %s has no fields with "reform:" tag, it is not allowed`, res.Type)
	}

	dupes := make(map[string]string)
	for _, f := range res.Fields {
		if f2, ok := dupes[f.Column]; ok {
			return fmt.Errorf(`reform: %s has field %s with "reform:" tag with duplicate column name %s (used by %s), it is not allowed`,
				res.Type, f.Name, f.Column, f2)
		}
		dupes[f.Column] = f.Name
	}

	return nil
}
