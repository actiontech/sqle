package state

// PostgresExtension - an installed extension on a database
type PostgresExtension struct {
	DatabaseOid   Oid
	ExtensionName string
	Version       string
	SchemaName    string
}
