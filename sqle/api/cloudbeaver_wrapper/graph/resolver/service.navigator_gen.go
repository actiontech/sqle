package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// NavRenameNode is the resolver for the navRenameNode field.
func (r *mutationResolver) NavRenameNode(ctx context.Context, nodePath string, newName string) (*string, error) {
	panic(fmt.Errorf("not implemented: NavRenameNode - navRenameNode"))
}

// NavDeleteNodes is the resolver for the navDeleteNodes field.
func (r *mutationResolver) NavDeleteNodes(ctx context.Context, nodePaths []string) (*int, error) {
	panic(fmt.Errorf("not implemented: NavDeleteNodes - navDeleteNodes"))
}

// NavMoveNodesToFolder is the resolver for the navMoveNodesToFolder field.
func (r *mutationResolver) NavMoveNodesToFolder(ctx context.Context, nodePaths []string, folderPath string) (*bool, error) {
	panic(fmt.Errorf("not implemented: NavMoveNodesToFolder - navMoveNodesToFolder"))
}

// NavNodeChildren is the resolver for the navNodeChildren field.
func (r *queryResolver) NavNodeChildren(ctx context.Context, parentPath string, offset *int, limit *int, onlyFolders *bool) ([]*model.NavigatorNodeInfo, error) {
	panic(fmt.Errorf("not implemented: NavNodeChildren - navNodeChildren"))
}

// NavNodeParents is the resolver for the navNodeParents field.
func (r *queryResolver) NavNodeParents(ctx context.Context, nodePath string) ([]*model.NavigatorNodeInfo, error) {
	panic(fmt.Errorf("not implemented: NavNodeParents - navNodeParents"))
}

// NavNodeInfo is the resolver for the navNodeInfo field.
func (r *queryResolver) NavNodeInfo(ctx context.Context, nodePath string) (*model.NavigatorNodeInfo, error) {
	panic(fmt.Errorf("not implemented: NavNodeInfo - navNodeInfo"))
}

// NavRefreshNode is the resolver for the navRefreshNode field.
func (r *queryResolver) NavRefreshNode(ctx context.Context, nodePath string) (*bool, error) {
	panic(fmt.Errorf("not implemented: NavRefreshNode - navRefreshNode"))
}

// NavGetStructContainers is the resolver for the navGetStructContainers field.
func (r *queryResolver) NavGetStructContainers(ctx context.Context, connectionID string, contextID *string, catalog *string) (*model.DatabaseStructContainers, error) {
	panic(fmt.Errorf("not implemented: NavGetStructContainers - navGetStructContainers"))
}
