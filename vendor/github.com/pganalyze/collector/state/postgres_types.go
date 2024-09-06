package state

import "github.com/guregu/null"

type Oid uint64
type Xid uint32  // 32-bit transaction ID
type Xid8 uint64 // 64-bit transaction ID

// PostgresType - User-defined custom data types
type PostgresType struct {
	Oid               Oid
	ArrayOid          Oid
	DatabaseOid       Oid
	SchemaName        string
	Name              string
	Type              string
	DomainType        null.String
	DomainNotNull     bool
	DomainDefault     null.String
	DomainConstraints []string
	EnumValues        []string
	CompositeAttrs    [][2]string
}
