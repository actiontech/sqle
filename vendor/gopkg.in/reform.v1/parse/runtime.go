package parse

import (
	"fmt"
	"reflect"
	"strings"
)

func objectGoType(t reflect.Type, structT reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		return "*" + objectGoType(t.Elem(), structT)
	}

	s := t.String()

	// drop package name from qualified identifier if type is defined in the same package
	if strings.Contains(s, ".") && t.PkgPath() == structT.PkgPath() {
		s = strings.Join(strings.Split(s, ".")[1:], ".")
	}

	return s
}

// Object extracts struct information from given object.
func Object(obj interface{}, schema, table string) (res *StructInfo, err error) {
	// convert any panic to error
	defer func() {
		p := recover()
		switch p := p.(type) {
		case error:
			err = p
		case nil:
			// nothing
		default:
			err = fmt.Errorf("%s", p)
		}
	}()

	t := reflect.ValueOf(obj).Elem().Type()
	res = &StructInfo{
		Type:         t.Name(),
		SQLSchema:    schema,
		SQLName:      table,
		PKFieldIndex: -1,
	}

	var n int
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		tag := f.Tag.Get("reform")
		if tag == "" || tag == "-" {
			continue
		}

		// check for anonymous fields
		if f.Anonymous {
			return nil, fmt.Errorf(`reform: %s has anonymous field %s with "reform:" tag, it is not allowed`, res.Type, f.Name)
		}

		// check for exported name
		if f.PkgPath != "" {
			return nil, fmt.Errorf(`reform: %s has non-exported field %s with "reform:" tag, it is not allowed`, res.Type, f.Name)
		}

		// parse tag and type
		column, isPK := parseStructFieldTag(tag)
		if column == "" {
			return nil, fmt.Errorf(`reform: %s has field %s with invalid "reform:" tag value, it is not allowed`, res.Type, f.Name)
		}
		typ := objectGoType(f.Type, t)
		if isPK {
			if strings.HasPrefix(typ, "*") {
				return nil, fmt.Errorf(`reform: %s has pointer field %s with with "pk" label in "reform:" tag, it is not allowed`, res.Type, f.Name)
			}
			if strings.HasPrefix(typ, "[") {
				return nil, fmt.Errorf(`reform: %s has slice field %s with with "pk" label in "reform:" tag, it is not allowed`, res.Type, f.Name)
			}
			if res.PKFieldIndex >= 0 {
				return nil, fmt.Errorf(`reform: %s has field %s with with duplicate "pk" label in "reform:" tag (first used by %s), it is not allowed`, res.Type, f.Name, res.Fields[res.PKFieldIndex].Name)
			}
		}

		res.Fields = append(res.Fields, FieldInfo{
			Name:   f.Name,
			Type:   typ,
			Column: column,
		})
		if isPK {
			res.PKFieldIndex = n
		}
		n++
	}

	if err = checkFields(res); err != nil {
		return nil, err
	}

	return
}
