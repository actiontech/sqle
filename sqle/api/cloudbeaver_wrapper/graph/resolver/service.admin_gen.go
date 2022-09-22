package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// ListUsers is the resolver for the listUsers field.
func (r *queryResolver) ListUsers(ctx context.Context, userID *string) ([]*model.AdminUserInfo, error) {
	panic(fmt.Errorf("not implemented: ListUsers - listUsers"))
}

// ListRoles is the resolver for the listRoles field.
func (r *queryResolver) ListRoles(ctx context.Context, roleID *string) ([]*model.AdminRoleInfo, error) {
	panic(fmt.Errorf("not implemented: ListRoles - listRoles"))
}

// ListPermissions is the resolver for the listPermissions field.
func (r *queryResolver) ListPermissions(ctx context.Context) ([]*model.AdminPermissionInfo, error) {
	panic(fmt.Errorf("not implemented: ListPermissions - listPermissions"))
}

// CreateUser is the resolver for the createUser field.
func (r *queryResolver) CreateUser(ctx context.Context, userID string) (*model.AdminUserInfo, error) {
	panic(fmt.Errorf("not implemented: CreateUser - createUser"))
}

// DeleteUser is the resolver for the deleteUser field.
func (r *queryResolver) DeleteUser(ctx context.Context, userID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteUser - deleteUser"))
}

// CreateRole is the resolver for the createRole field.
func (r *queryResolver) CreateRole(ctx context.Context, roleID string, roleName *string, description *string) (*model.AdminRoleInfo, error) {
	panic(fmt.Errorf("not implemented: CreateRole - createRole"))
}

// UpdateRole is the resolver for the updateRole field.
func (r *queryResolver) UpdateRole(ctx context.Context, roleID string, roleName *string, description *string) (*model.AdminRoleInfo, error) {
	panic(fmt.Errorf("not implemented: UpdateRole - updateRole"))
}

// DeleteRole is the resolver for the deleteRole field.
func (r *queryResolver) DeleteRole(ctx context.Context, roleID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteRole - deleteRole"))
}

// GrantUserRole is the resolver for the grantUserRole field.
func (r *queryResolver) GrantUserRole(ctx context.Context, userID string, roleID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: GrantUserRole - grantUserRole"))
}

// RevokeUserRole is the resolver for the revokeUserRole field.
func (r *queryResolver) RevokeUserRole(ctx context.Context, userID string, roleID string) (*bool, error) {
	panic(fmt.Errorf("not implemented: RevokeUserRole - revokeUserRole"))
}

// SetSubjectPermissions is the resolver for the setSubjectPermissions field.
func (r *queryResolver) SetSubjectPermissions(ctx context.Context, roleID string, permissions []string) ([]*model.AdminPermissionInfo, error) {
	panic(fmt.Errorf("not implemented: SetSubjectPermissions - setSubjectPermissions"))
}

// SetUserCredentials is the resolver for the setUserCredentials field.
func (r *queryResolver) SetUserCredentials(ctx context.Context, userID string, providerID string, credentials interface{}) (*bool, error) {
	panic(fmt.Errorf("not implemented: SetUserCredentials - setUserCredentials"))
}

// EnableUser is the resolver for the enableUser field.
func (r *queryResolver) EnableUser(ctx context.Context, userID string, enabled bool) (*bool, error) {
	panic(fmt.Errorf("not implemented: EnableUser - enableUser"))
}

// AllConnections is the resolver for the allConnections field.
func (r *queryResolver) AllConnections(ctx context.Context, id *string) ([]*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: AllConnections - allConnections"))
}

// SearchConnections is the resolver for the searchConnections field.
func (r *queryResolver) SearchConnections(ctx context.Context, hostNames []string) ([]*model.AdminConnectionSearchInfo, error) {
	panic(fmt.Errorf("not implemented: SearchConnections - searchConnections"))
}

// CreateConnectionConfiguration is the resolver for the createConnectionConfiguration field.
func (r *queryResolver) CreateConnectionConfiguration(ctx context.Context, config model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: CreateConnectionConfiguration - createConnectionConfiguration"))
}

// CopyConnectionConfiguration is the resolver for the copyConnectionConfiguration field.
func (r *queryResolver) CopyConnectionConfiguration(ctx context.Context, nodePath string, config *model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: CopyConnectionConfiguration - copyConnectionConfiguration"))
}

// UpdateConnectionConfiguration is the resolver for the updateConnectionConfiguration field.
func (r *queryResolver) UpdateConnectionConfiguration(ctx context.Context, id string, config model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: UpdateConnectionConfiguration - updateConnectionConfiguration"))
}

// DeleteConnectionConfiguration is the resolver for the deleteConnectionConfiguration field.
func (r *queryResolver) DeleteConnectionConfiguration(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: DeleteConnectionConfiguration - deleteConnectionConfiguration"))
}

// GetConnectionSubjectAccess is the resolver for the getConnectionSubjectAccess field.
func (r *queryResolver) GetConnectionSubjectAccess(ctx context.Context, connectionID *string) ([]*model.AdminConnectionGrantInfo, error) {
	panic(fmt.Errorf("not implemented: GetConnectionSubjectAccess - getConnectionSubjectAccess"))
}

// SetConnectionSubjectAccess is the resolver for the setConnectionSubjectAccess field.
func (r *queryResolver) SetConnectionSubjectAccess(ctx context.Context, connectionID string, subjects []string) (*bool, error) {
	panic(fmt.Errorf("not implemented: SetConnectionSubjectAccess - setConnectionSubjectAccess"))
}

// GetSubjectConnectionAccess is the resolver for the getSubjectConnectionAccess field.
func (r *queryResolver) GetSubjectConnectionAccess(ctx context.Context, subjectID *string) ([]*model.AdminConnectionGrantInfo, error) {
	panic(fmt.Errorf("not implemented: GetSubjectConnectionAccess - getSubjectConnectionAccess"))
}

// SetSubjectConnectionAccess is the resolver for the setSubjectConnectionAccess field.
func (r *queryResolver) SetSubjectConnectionAccess(ctx context.Context, subjectID string, connections []string) (*bool, error) {
	panic(fmt.Errorf("not implemented: SetSubjectConnectionAccess - setSubjectConnectionAccess"))
}

// ListFeatureSets is the resolver for the listFeatureSets field.
func (r *queryResolver) ListFeatureSets(ctx context.Context) ([]*model.WebFeatureSet, error) {
	panic(fmt.Errorf("not implemented: ListFeatureSets - listFeatureSets"))
}

// ListAuthProviderConfigurationParameters is the resolver for the listAuthProviderConfigurationParameters field.
func (r *queryResolver) ListAuthProviderConfigurationParameters(ctx context.Context, providerID string) ([]*model.ObjectPropertyInfo, error) {
	panic(fmt.Errorf("not implemented: ListAuthProviderConfigurationParameters - listAuthProviderConfigurationParameters"))
}

// ListAuthProviderConfigurations is the resolver for the listAuthProviderConfigurations field.
func (r *queryResolver) ListAuthProviderConfigurations(ctx context.Context, providerID *string) ([]*model.AdminAuthProviderConfiguration, error) {
	panic(fmt.Errorf("not implemented: ListAuthProviderConfigurations - listAuthProviderConfigurations"))
}

// SaveAuthProviderConfiguration is the resolver for the saveAuthProviderConfiguration field.
func (r *queryResolver) SaveAuthProviderConfiguration(ctx context.Context, providerID string, id string, displayName *string, disabled *bool, iconURL *string, description *string, parameters interface{}) (*model.AdminAuthProviderConfiguration, error) {
	panic(fmt.Errorf("not implemented: SaveAuthProviderConfiguration - saveAuthProviderConfiguration"))
}

// DeleteAuthProviderConfiguration is the resolver for the deleteAuthProviderConfiguration field.
func (r *queryResolver) DeleteAuthProviderConfiguration(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteAuthProviderConfiguration - deleteAuthProviderConfiguration"))
}

// SaveUserMetaParameter is the resolver for the saveUserMetaParameter field.
func (r *queryResolver) SaveUserMetaParameter(ctx context.Context, id string, displayName string, description *string, required bool) (*model.ObjectPropertyInfo, error) {
	panic(fmt.Errorf("not implemented: SaveUserMetaParameter - saveUserMetaParameter"))
}

// DeleteUserMetaParameter is the resolver for the deleteUserMetaParameter field.
func (r *queryResolver) DeleteUserMetaParameter(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteUserMetaParameter - deleteUserMetaParameter"))
}

// SetUserMetaParameterValues is the resolver for the setUserMetaParameterValues field.
func (r *queryResolver) SetUserMetaParameterValues(ctx context.Context, userID string, parameters interface{}) (bool, error) {
	panic(fmt.Errorf("not implemented: SetUserMetaParameterValues - setUserMetaParameterValues"))
}

// ConfigureServer is the resolver for the configureServer field.
func (r *queryResolver) ConfigureServer(ctx context.Context, configuration model.ServerConfigInput) (bool, error) {
	panic(fmt.Errorf("not implemented: ConfigureServer - configureServer"))
}

// SetDefaultNavigatorSettings is the resolver for the setDefaultNavigatorSettings field.
func (r *queryResolver) SetDefaultNavigatorSettings(ctx context.Context, settings model.NavigatorSettingsInput) (bool, error) {
	panic(fmt.Errorf("not implemented: SetDefaultNavigatorSettings - setDefaultNavigatorSettings"))
}
