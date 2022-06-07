package driver

import (
	"context"

	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// AnalysisDriver is a driver for SQL analysis and getting table metadata
type AnalysisDriver interface {
	ListTablesInSchema(ctx context.Context, conf *ListTablesInSchemaConf) (*ListTablesInSchemaResult, error)
	GetTableMetas(ctx context.Context, conf *GetTableMetasConf) (*GetTableMetasResult, error)
	Explain(ctx context.Context, conf *ExplainConf) (*ExplainResult, error)
}

type ListTablesInSchemaConf struct {
	schema string
}

type Table struct {
	Name string
}

type ListTablesInSchemaResult struct {
	Tables []Table
}

// AnalysisInfoInTableFormat
// the field Column represents the column name of a table
// the field Rows represents the data of the table
// their relationship is as follows
/*
	| Column[0]  | Column[1]  | Column[2]  |
	| Rows[0][0] | Rows[0][1] | Rows[0][2] |
	| Rows[1][0] | Rows[1][1] | Rows[1][2] |
*/
type AnalysisInfoInTableFormat struct {
	Column params.Params
	Rows   [][]string
}

type GetTableMetasConf struct {
	SchemaToTables map[string] /* schema name */ []string /* tables' names */
}

type ColumnsInfo struct {
	AnalysisInfoInTableFormat
}

type IndexesInfo struct {
	AnalysisInfoInTableFormat
}

type GetTableMetaItem struct {
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
}

type SchemaTableMetas struct {
	Schema     string
	TableMetas map[string] /* table name */ GetTableMetaItem
}

type GetTableMetasResult struct {
	SchemaTableMetas []SchemaTableMetas
}

type ExplainConf struct {
	sql string
}

type ExplainClassicResult struct {
	AnalysisInfoInTableFormat
}

type ExplainResult struct {
	ClassicResult ExplainClassicResult
}

func NewAnalysisDriver(log *logrus.Entry, dbType string, cfg *DSN) (AnalysisDriver, error) {
	return nil, nil
}
