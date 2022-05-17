//go:build enterprise
// +build enterprise

package model

const (
	// additional operation code list for ee
	// SqlQuery: SQL查询 reserved 40000-49999
	OP_SQL_QUERY_QUERY = 40100
)

func getConfigurableOperationCodeListForEE() []uint{
	return []uint{
		// Sql Query: SQL查询
		OP_SQL_QUERY_QUERY,
	}
}

func additionalOperationForEE(opCode uint) string {
	switch opCode {
	case OP_SQL_QUERY_QUERY:
		return "SQL查询"
	}
	return "未知动作"
}