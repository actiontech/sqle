package differ

import (
	"fmt"
	"regexp"
	"strings"
)

// ObjectType defines a class of object in a relational database system.
type ObjectType string

// Constants enumerating valid object types.
// Currently we do not define separate types for sub-types such as columns,
// indexes, foreign keys, etc as these are handled within the table logic.
const (
	ObjectTypeNil      ObjectType = ""
	ObjectTypeDatabase ObjectType = "database"
	ObjectTypeTable    ObjectType = "table"
	ObjectTypeProc     ObjectType = "procedure"
	ObjectTypeFunc     ObjectType = "function"
)

// Caps returns the object type as an uppercase string.
func (ot ObjectType) Caps() string {
	return strings.ToUpper(string(ot))
}

// ObjectKey is useful as a map key for indexing database objects within a
// single schema.
type ObjectKey struct {
	Type ObjectType
	Name string
}

func (key ObjectKey) String() string {
	return fmt.Sprintf("%s %s", key.Type, EscapeIdentifier(key.Name))
}

// ObjectKey inception as a syntactic sugar hack: this allows keys to be
// passed directly for any arg expecting an ObjectKeyer interface.
func (key ObjectKey) ObjectKey() ObjectKey {
	return key
}

// ObjectKeyer is an interface implemented by each type of database object,
// providing a generic way of obtaining an object's type and name.
type ObjectKeyer interface {
	ObjectKey() ObjectKey
}

// DefKeyer is an interface that extends ObjectKeyer with an additional Def
// method, for returning a CREATE statement corresponding to the object. No
// guarantees are made as to whether this corresponds to a normalized value
// obtained from SHOW CREATE, an imputed value based on a particular Flavor,
// or an arbitrarily-formatted CREATE obtained from some other source. This
// flexibility allows DefKeyer to be used for purposes beyond just representing
// live database objects.
type DefKeyer interface {
	ObjectKeyer
	Def() string
}

// ObjectPattern is a regular expression matched against an object name, but
// only for a specific object type.
type ObjectPattern struct {
	Type    ObjectType
	Pattern *regexp.Regexp
}

// Match returns true if p's Type equals obj's ObjectKey.Type and p's Pattern
// matches obj's ObjectKey.Name.
func (p *ObjectPattern) Match(obj ObjectKeyer) bool {
	if p == nil || p.Pattern == nil {
		return false
	}
	key := obj.ObjectKey()
	return p.Type == key.Type && key.Name != "" && p.Pattern.MatchString(key.Name)
}

func (p *ObjectPattern) String() string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%s %s", p.Type, p.Pattern)
}
