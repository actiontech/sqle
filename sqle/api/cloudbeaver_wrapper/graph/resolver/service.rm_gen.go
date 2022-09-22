package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// RmCreateResource is the resolver for the rmCreateResource field.
func (r *mutationResolver) RmCreateResource(ctx context.Context, projectID string, resourcePath string, isFolder bool) (string, error) {
	panic(fmt.Errorf("not implemented: RmCreateResource - rmCreateResource"))
}

// RmMoveResource is the resolver for the rmMoveResource field.
func (r *mutationResolver) RmMoveResource(ctx context.Context, projectID string, oldResourcePath string, newResourcePath *string) (string, error) {
	panic(fmt.Errorf("not implemented: RmMoveResource - rmMoveResource"))
}

// RmDeleteResource is the resolver for the rmDeleteResource field.
func (r *mutationResolver) RmDeleteResource(ctx context.Context, projectID string, resourcePath string, recursive bool) (*bool, error) {
	panic(fmt.Errorf("not implemented: RmDeleteResource - rmDeleteResource"))
}

// RmWriteResourceStringContent is the resolver for the rmWriteResourceStringContent field.
func (r *mutationResolver) RmWriteResourceStringContent(ctx context.Context, projectID string, resourcePath string, data string) (string, error) {
	panic(fmt.Errorf("not implemented: RmWriteResourceStringContent - rmWriteResourceStringContent"))
}

// RmListProjects is the resolver for the rmListProjects field.
func (r *queryResolver) RmListProjects(ctx context.Context) ([]*model.RMProject, error) {
	panic(fmt.Errorf("not implemented: RmListProjects - rmListProjects"))
}

// RmListResources is the resolver for the rmListResources field.
func (r *queryResolver) RmListResources(ctx context.Context, projectID string, folder *string, nameMask *string, readProperties *bool, readHistory *bool) ([]*model.RMResource, error) {
	panic(fmt.Errorf("not implemented: RmListResources - rmListResources"))
}

// RmReadResourceAsString is the resolver for the rmReadResourceAsString field.
func (r *queryResolver) RmReadResourceAsString(ctx context.Context, projectID string, resourcePath string) (string, error) {
	panic(fmt.Errorf("not implemented: RmReadResourceAsString - rmReadResourceAsString"))
}
