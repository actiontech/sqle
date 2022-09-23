package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// SQLContextCreate is the resolver for the sqlContextCreate field.
func (r *mutationResolver) SQLContextCreate(ctx context.Context, connectionID string, defaultCatalog *string, defaultSchema *string) (*model.SQLContextInfo, error) {
	panic(fmt.Errorf("not implemented: SQLContextCreate - sqlContextCreate"))
}

// SQLContextSetDefaults is the resolver for the sqlContextSetDefaults field.
func (r *mutationResolver) SQLContextSetDefaults(ctx context.Context, connectionID string, contextID string, defaultCatalog *string, defaultSchema *string) (bool, error) {
	panic(fmt.Errorf("not implemented: SQLContextSetDefaults - sqlContextSetDefaults"))
}

// SQLContextDestroy is the resolver for the sqlContextDestroy field.
func (r *mutationResolver) SQLContextDestroy(ctx context.Context, connectionID string, contextID string) (bool, error) {
	panic(fmt.Errorf("not implemented: SQLContextDestroy - sqlContextDestroy"))
}

// AsyncSQLExecuteQuery is the resolver for the asyncSqlExecuteQuery field.
func (r *mutationResolver) AsyncSQLExecuteQuery(ctx context.Context, connectionID string, contextID string, sql string, resultID *string, filter *model.SQLDataFilter, dataFormat *model.ResultDataFormat) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: AsyncSQLExecuteQuery - asyncSqlExecuteQuery"))
}

// AsyncReadDataFromContainer is the resolver for the asyncReadDataFromContainer field.
func (r *mutationResolver) AsyncReadDataFromContainer(ctx context.Context, connectionID string, contextID string, containerNodePath string, resultID *string, filter *model.SQLDataFilter, dataFormat *model.ResultDataFormat) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: AsyncReadDataFromContainer - asyncReadDataFromContainer"))
}

// SQLResultClose is the resolver for the sqlResultClose field.
func (r *mutationResolver) SQLResultClose(ctx context.Context, connectionID string, contextID string, resultID string) (bool, error) {
	panic(fmt.Errorf("not implemented: SQLResultClose - sqlResultClose"))
}

// UpdateResultsDataBatch is the resolver for the updateResultsDataBatch field.
func (r *mutationResolver) UpdateResultsDataBatch(ctx context.Context, connectionID string, contextID string, resultsID string, updatedRows []*model.SQLResultRow, deletedRows []*model.SQLResultRow, addedRows []*model.SQLResultRow) (*model.SQLExecuteInfo, error) {
	panic(fmt.Errorf("not implemented: UpdateResultsDataBatch - updateResultsDataBatch"))
}

// UpdateResultsDataBatchScript is the resolver for the updateResultsDataBatchScript field.
func (r *mutationResolver) UpdateResultsDataBatchScript(ctx context.Context, connectionID string, contextID string, resultsID string, updatedRows []*model.SQLResultRow, deletedRows []*model.SQLResultRow, addedRows []*model.SQLResultRow) (string, error) {
	panic(fmt.Errorf("not implemented: UpdateResultsDataBatchScript - updateResultsDataBatchScript"))
}

// ReadLobValue is the resolver for the readLobValue field.
func (r *mutationResolver) ReadLobValue(ctx context.Context, connectionID string, contextID string, resultsID string, lobColumnIndex int, row []*model.SQLResultRow) (string, error) {
	panic(fmt.Errorf("not implemented: ReadLobValue - readLobValue"))
}

// AsyncSQLExecuteResults is the resolver for the asyncSqlExecuteResults field.
func (r *mutationResolver) AsyncSQLExecuteResults(ctx context.Context, taskID string) (*model.SQLExecuteInfo, error) {
	panic(fmt.Errorf("not implemented: AsyncSQLExecuteResults - asyncSqlExecuteResults"))
}

// AsyncSQLExplainExecutionPlan is the resolver for the asyncSqlExplainExecutionPlan field.
func (r *mutationResolver) AsyncSQLExplainExecutionPlan(ctx context.Context, connectionID string, contextID string, query string, configuration interface{}) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: AsyncSQLExplainExecutionPlan - asyncSqlExplainExecutionPlan"))
}

// AsyncSQLExplainExecutionPlanResult is the resolver for the asyncSqlExplainExecutionPlanResult field.
func (r *mutationResolver) AsyncSQLExplainExecutionPlanResult(ctx context.Context, taskID string) (*model.SQLExecutionPlan, error) {
	panic(fmt.Errorf("not implemented: AsyncSQLExplainExecutionPlanResult - asyncSqlExplainExecutionPlanResult"))
}

// SQLDialectInfo is the resolver for the sqlDialectInfo field.
func (r *queryResolver) SQLDialectInfo(ctx context.Context, connectionID string) (*model.SQLDialectInfo, error) {
	panic(fmt.Errorf("not implemented: SQLDialectInfo - sqlDialectInfo"))
}

// SQLListContexts is the resolver for the sqlListContexts field.
func (r *queryResolver) SQLListContexts(ctx context.Context, connectionID *string, contextID *string) ([]*model.SQLContextInfo, error) {
	panic(fmt.Errorf("not implemented: SQLListContexts - sqlListContexts"))
}

// SQLCompletionProposals is the resolver for the sqlCompletionProposals field.
func (r *queryResolver) SQLCompletionProposals(ctx context.Context, connectionID string, contextID string, query string, position int, maxResults *int, simpleMode *bool) ([]*model.SQLCompletionProposal, error) {
	panic(fmt.Errorf("not implemented: SQLCompletionProposals - sqlCompletionProposals"))
}

// SQLFormatQuery is the resolver for the sqlFormatQuery field.
func (r *queryResolver) SQLFormatQuery(ctx context.Context, connectionID string, contextID string, query string) (string, error) {
	panic(fmt.Errorf("not implemented: SQLFormatQuery - sqlFormatQuery"))
}

// SQLSupportedOperations is the resolver for the sqlSupportedOperations field.
func (r *queryResolver) SQLSupportedOperations(ctx context.Context, connectionID string, contextID string, resultsID string, attributeIndex int) ([]*model.DataTypeLogicalOperation, error) {
	panic(fmt.Errorf("not implemented: SQLSupportedOperations - sqlSupportedOperations"))
}

// SQLEntityQueryGenerators is the resolver for the sqlEntityQueryGenerators field.
func (r *queryResolver) SQLEntityQueryGenerators(ctx context.Context, nodePathList []string) ([]*model.SQLQueryGenerator, error) {
	panic(fmt.Errorf("not implemented: SQLEntityQueryGenerators - sqlEntityQueryGenerators"))
}

// SQLGenerateEntityQuery is the resolver for the sqlGenerateEntityQuery field.
func (r *queryResolver) SQLGenerateEntityQuery(ctx context.Context, generatorID string, options interface{}, nodePathList []string) (string, error) {
	panic(fmt.Errorf("not implemented: SQLGenerateEntityQuery - sqlGenerateEntityQuery"))
}

// SQLParseScript is the resolver for the sqlParseScript field.
func (r *queryResolver) SQLParseScript(ctx context.Context, connectionID string, script string) (*model.SQLScriptInfo, error) {
	panic(fmt.Errorf("not implemented: SQLParseScript - sqlParseScript"))
}

// SQLParseQuery is the resolver for the sqlParseQuery field.
func (r *queryResolver) SQLParseQuery(ctx context.Context, connectionID string, script string, position int) (*model.SQLScriptQuery, error) {
	panic(fmt.Errorf("not implemented: SQLParseQuery - sqlParseQuery"))
}
