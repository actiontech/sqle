package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// SetUserConfigurationParameter is the resolver for the setUserConfigurationParameter field.
func (r *mutationResolver) SetUserConfigurationParameter(ctx context.Context, name string, value interface{}) (bool, error) {
	panic(fmt.Errorf("not implemented: SetUserConfigurationParameter - setUserConfigurationParameter"))
}

// AuthLogin is the resolver for the authLogin field.
func (r *queryResolver) AuthLogin(ctx context.Context, provider string, configuration *string, credentials interface{}, linkUser *bool) (*model.AuthInfo, error) {
	panic(fmt.Errorf("not implemented: AuthLogin - authLogin"))
}

// AuthUpdateStatus is the resolver for the authUpdateStatus field.
func (r *queryResolver) AuthUpdateStatus(ctx context.Context, authID string, linkUser *bool) (*model.AuthInfo, error) {
	panic(fmt.Errorf("not implemented: AuthUpdateStatus - authUpdateStatus"))
}

// AuthLogout is the resolver for the authLogout field.
func (r *queryResolver) AuthLogout(ctx context.Context, provider *string, configuration *string) (*bool, error) {
	panic(fmt.Errorf("not implemented: AuthLogout - authLogout"))
}

// ActiveUser is the resolver for the activeUser field.
func (r *queryResolver) ActiveUser(ctx context.Context) (*model.UserInfo, error) {
	panic(fmt.Errorf("not implemented: ActiveUser - activeUser"))
}

// AuthProviders is the resolver for the authProviders field.
func (r *queryResolver) AuthProviders(ctx context.Context) ([]*model.AuthProviderInfo, error) {
	panic(fmt.Errorf("not implemented: AuthProviders - authProviders"))
}

// AuthChangeLocalPassword is the resolver for the authChangeLocalPassword field.
func (r *queryResolver) AuthChangeLocalPassword(ctx context.Context, oldPassword string, newPassword string) (bool, error) {
	panic(fmt.Errorf("not implemented: AuthChangeLocalPassword - authChangeLocalPassword"))
}

// ListUserProfileProperties is the resolver for the listUserProfileProperties field.
func (r *queryResolver) ListUserProfileProperties(ctx context.Context) ([]*model.ObjectPropertyInfo, error) {
	panic(fmt.Errorf("not implemented: ListUserProfileProperties - listUserProfileProperties"))
}
