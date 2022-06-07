package driver

import (
	"context"

	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/sirupsen/logrus"
)

// AnalysisDriver is a driver for SQL analysis and getting table metadata
type AnalysisDriver interface {
	ListTablesInSchema(ctx context.Context, conf *ListTablesInSchemaConf) (*ListTablesInSchemaResult, error)
	GetTableMetaByTableName(ctx context.Context, conf *GetTableMetaByTableNameConf) (*GetTableMetaByTableNameResult, error)
	GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error)
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

type GetTableMetaByTableNameConf struct {
	Schema string
	Table  string
}

type GetTableMetaBySQLConf struct {
	Sql string
}

type ColumnsInfo struct {
	AnalysisInfoInTableFormat
}

type IndexesInfo struct {
	AnalysisInfoInTableFormat
}

type TableMetaItem struct {
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
}

type GetTableMetaByTableNameResult struct {
	TableMeta TableMetaItem
}

type GetTableMetaBySQLResult struct {
	TableMetas []TableMetaItem
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
