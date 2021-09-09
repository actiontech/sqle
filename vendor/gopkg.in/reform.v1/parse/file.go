package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"regexp"
	"strings"
)

var magicReformComment = regexp.MustCompile(`reform:([0-9A-Za-z_\.]+)`)

func fileGoType(x ast.Expr) string {
	switch t := x.(type) {
	case *ast.StarExpr:
		return "*" + fileGoType(t.X)
	case *ast.SelectorExpr:
		return fileGoType(t.X) + "." + t.Sel.String()
	case *ast.Ident:
		s := t.String()
		if s == "byte" {
			return "uint8"
		}
		return s
	case *ast.ArrayType:
		return "[" + fileGoType(t.Len) + "]" + fileGoType(t.Elt)
	case *ast.BasicLit:
		return t.Value
	case nil:
		return ""
	default:
		panic(fmt.Sprintf("reform: fileGoType: unhandled '%s' (%#v). Please report this bug.", x, x))
	}
}

func commentText(g *ast.CommentGroup) string {
	// this code used to just call g.Text(), but the behavior of this method changed in Go 1.15:
	// https://go-review.googlesource.com/c/go/+/224737
	res := make([]string, len(g.List))
	for i, c := range g.List {
		res[i] = c.Text
	}
	return strings.Join(res, " ")
}

func parseStructTypeSpec(ts *ast.TypeSpec, str *ast.StructType) (*StructInfo, error) {
	res := &StructInfo{
		Type:         ts.Name.Name,
		PKFieldIndex: -1,
	}

	var n int
	for _, f := range str.Fields.List {
		// consider only fields with "reform:" tag
		if f.Tag == nil {
			continue
		}
		tag := f.Tag.Value
		if len(tag) < 3 {
			continue
		}
		tag = reflect.StructTag(tag[1 : len(tag)-1]).Get("reform") // strip quotes
		if tag == "" || tag == "-" {
			continue
		}

		// check for anonymous fields
		if len(f.Names) == 0 {
			return nil, fmt.Errorf(`reform: %s has anonymous field %s with "reform:" tag, it is not allowed`, res.Type, f.Type)
		}
		if len(f.Names) != 1 {
			panic(fmt.Sprintf("reform: %d names: %#v. Please report this bug.", len(f.Names), f.Names))
		}

		// check for exported name
		name := f.Names[0]
		if !name.IsExported() {
			return nil, fmt.Errorf(`reform: %s has non-exported field %s with "reform:" tag, it is not allowed`, res.Type, name.Name)
		}

		// parse tag and type
		column, isPK := parseStructFieldTag(tag)
		if column == "" {
			return nil, fmt.Errorf(`reform: %s has field %s with invalid "reform:" tag value, it is not allowed`, res.Type, name.Name)
		}
		typ := fileGoType(f.Type)
		if isPK {
			if strings.HasPrefix(typ, "*") {
				return nil, fmt.Errorf(`reform: %s has pointer field %s with with "pk" label in "reform:" tag, it is not allowed`, res.Type, name.Name)
			}
			if strings.HasPrefix(typ, "[") {
				return nil, fmt.Errorf(`reform: %s has slice field %s with with "pk" label in "reform:" tag, it is not allowed`, res.Type, name.Name)
			}
			if res.PKFieldIndex >= 0 {
				return nil, fmt.Errorf(`reform: %s has field %s with with duplicate "pk" label in "reform:" tag (first used by %s), it is not allowed`, res.Type, name.Name, res.Fields[res.PKFieldIndex].Name)
			}
		}

		res.Fields = append(res.Fields, FieldInfo{
			Name:   name.Name,
			Type:   typ,
			Column: column,
		})
		if isPK {
			res.PKFieldIndex = n
		}
		n++
	}

	if len(res.Fields) == 0 {
		return nil, fmt.Errorf(`reform: %s has no fields with "reform:" tag, it is not allowed`, res.Type)
	}

	if err := checkFields(res); err != nil {
		return nil, err
	}

	return res, nil
}

// File parses given file and returns found structs information.
func File(path string) ([]StructInfo, error) {
	// parse file
	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// consider only top-level struct type declarations with magic comment
	var res []StructInfo
	for _, decl := range fileNode.Decls {
		// ast.Print(fset, decl)

		gd, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			// ast.Print(fset, ts)

			// magic comment may be attached to "type Foo struct" (TypeSpec)
			// or to "type (" (GenDecl)
			doc := ts.Doc
			if doc == nil && len(gd.Specs) == 1 {
				doc = gd.Doc
			}
			if doc == nil {
				continue
			}
			// ast.Print(fset, doc)

			sm := magicReformComment.FindStringSubmatch(commentText(doc))
			if len(sm) < 2 {
				continue
			}
			parts := strings.SplitN(sm[1], ".", 2)
			var schema string
			if len(parts) == 2 {
				schema = parts[0]
			}
			table := parts[len(parts)-1]

			str, ok := ts.Type.(*ast.StructType)
			if !ok {
				continue
			}
			if str.Incomplete {
				continue
			}
			// ast.Print(fset, str)

			s, err := parseStructTypeSpec(ts, str)
			if err != nil {
				return nil, err
			}
			s.SQLSchema = schema
			s.SQLName = table
			res = append(res, *s)
		}
	}

	return res, nil
}
