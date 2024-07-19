package state

type PostgresRelationBloat struct {
	SchemaName   string
	RelationName string
	TotalBytes   int64
	BloatBytes   int64
}

type PostgresIndexBloat struct {
	SchemaName string
	IndexName  string
	TotalBytes int64
	BloatBytes int64
}

type PostgresBloatStats struct {
	DatabaseName string
	Relations    []PostgresRelationBloat
	Indices      []PostgresIndexBloat
}
