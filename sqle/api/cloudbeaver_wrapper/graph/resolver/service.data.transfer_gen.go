package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// DataTransferAvailableStreamProcessors is the resolver for the dataTransferAvailableStreamProcessors field.
func (r *queryResolver) DataTransferAvailableStreamProcessors(ctx context.Context) ([]*model.DataTransferProcessorInfo, error) {
	panic(fmt.Errorf("not implemented: DataTransferAvailableStreamProcessors - dataTransferAvailableStreamProcessors"))
}

// DataTransferExportDataFromContainer is the resolver for the dataTransferExportDataFromContainer field.
func (r *queryResolver) DataTransferExportDataFromContainer(ctx context.Context, connectionID string, containerNodePath string, parameters model.DataTransferParameters) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: DataTransferExportDataFromContainer - dataTransferExportDataFromContainer"))
}

// DataTransferExportDataFromResults is the resolver for the dataTransferExportDataFromResults field.
func (r *queryResolver) DataTransferExportDataFromResults(ctx context.Context, connectionID string, contextID string, resultsID string, parameters model.DataTransferParameters) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: DataTransferExportDataFromResults - dataTransferExportDataFromResults"))
}

// DataTransferRemoveDataFile is the resolver for the dataTransferRemoveDataFile field.
func (r *queryResolver) DataTransferRemoveDataFile(ctx context.Context, dataFileID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DataTransferRemoveDataFile - dataTransferRemoveDataFile"))
}
