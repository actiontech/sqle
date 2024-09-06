package state

type PostgresSequenceReport struct {
	DatabaseName string

	Sequences            PostgresSequenceInformationMap
	SerialColumns        []PostgresSerialColumn
	ForeignSerialColumns []PostgresForeignSerialColumn
}

type PostgresSequenceInformationMap map[Oid]PostgresSequenceInformation

type PostgresSequenceInformation struct {
	SchemaName   string
	SequenceName string

	LastValue   int64
	StartValue  int64
	IncrementBy int64
	MaxValue    int64
	MinValue    int64
	CacheValue  int64
	IsCycled    bool
}

type PostgresSerialColumn struct {
	RelationOid  Oid
	SchemaName   string
	RelationName string
	ColumnName   string
	DataType     string
	MaximumValue uint64
	SequenceOid  Oid

	ForeignColumns []PostgresForeignSerialColumn
}

type PostgresForeignSerialColumn struct {
	RelationOid  Oid
	SchemaName   string
	RelationName string
	ColumnName   string

	DataType     string
	MaximumValue uint64

	Inferred bool
}
