package driver

import (
	"context"
	"fmt"
	"sync"
)

// AnalysisDriver is a driver for SQL analysis and getting table metadata
type AnalysisDriver interface {
	ListTablesInSchema(ctx context.Context, conf *ListTablesInSchemaConf) (*ListTablesInSchemaResult, error)
	GetTableMetaByTableName(ctx context.Context, conf *GetTableMetaByTableNameConf) (*GetTableMetaByTableNameResult, error)
	GetTableMetaBySQL(ctx context.Context, conf *GetTableMetaBySQLConf) (*GetTableMetaBySQLResult, error)
	Explain(ctx context.Context, conf *ExplainConf) (*ExplainResult, error)
}

type ListTablesInSchemaConf struct {
	Schema string
}

type Table struct {
	Name string
}

type ListTablesInSchemaResult struct {
	Tables []Table
}

type AnalysisInfoHead struct {
	Name string
	Desc string
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
	Column []AnalysisInfoHead
	Rows   [][]string
}

type GetTableMetaByTableNameConf struct {
	Schema string
	Table  string
}

type GetTableMetaBySQLConf struct {
	// this SQL should be a single SQL
	Sql string
}

type ColumnsInfo struct {
	AnalysisInfoInTableFormat
}

type IndexesInfo struct {
	AnalysisInfoInTableFormat
}

type TableMetaItem struct {
	Name           string
	Schema         string
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
}

type GetTableMetaByTableNameResult struct {
	TableMeta TableMetaItem
}

type TableMetaItemBySQL struct {
	Name           string
	Schema         string
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
	Message        string
}

type GetTableMetaBySQLResult struct {
	TableMetas []TableMetaItemBySQL
}

type ExplainConf struct {
	// this SQL should be a single SQL
	Sql string
}

type ExplainClassicResult struct {
	AnalysisInfoInTableFormat
}

type ExplainResult struct {
	ClassicResult ExplainClassicResult
}

var analysisDriverMu = &sync.RWMutex{}
var analysisDrivers = make(map[string]struct{})

// RegisterAnalysisDriver like sql.RegisterAuditDriver.
//
// RegisterAnalysisDriver makes a database driver available by the provided driver name.
// RegisterAnalysisDriver's initialize handler registered by RegisterAnalysisDriver.
func RegisterAnalysisDriver(name string) {
	analysisDriverMu.RLock()
	_, exist := analysisDrivers[name]
	analysisDriverMu.RUnlock()
	if exist {
		panic(fmt.Sprintf("duplicated driver name %v", name))
	}

	analysisDriverMu.Lock()
	analysisDrivers[name] = struct{}{}
	analysisDriverMu.Unlock()
}
