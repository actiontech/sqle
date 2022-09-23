package resolver

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
)

// OpenSession is the resolver for the openSession field.
func (r *mutationResolver) OpenSession(ctx context.Context, defaultLocale *string) (*model.SessionInfo, error) {
	panic(fmt.Errorf("not implemented: OpenSession - openSession"))
}

// CloseSession is the resolver for the closeSession field.
func (r *mutationResolver) CloseSession(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented: CloseSession - closeSession"))
}

// TouchSession is the resolver for the touchSession field.
func (r *mutationResolver) TouchSession(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented: TouchSession - touchSession"))
}

// RefreshSessionConnections is the resolver for the refreshSessionConnections field.
func (r *mutationResolver) RefreshSessionConnections(ctx context.Context) (*bool, error) {
	panic(fmt.Errorf("not implemented: RefreshSessionConnections - refreshSessionConnections"))
}

// ChangeSessionLanguage is the resolver for the changeSessionLanguage field.
func (r *mutationResolver) ChangeSessionLanguage(ctx context.Context, locale *string) (*bool, error) {
	panic(fmt.Errorf("not implemented: ChangeSessionLanguage - changeSessionLanguage"))
}

// CreateConnection is the resolver for the createConnection field.
func (r *mutationResolver) CreateConnection(ctx context.Context, config model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: CreateConnection - createConnection"))
}

// UpdateConnection is the resolver for the updateConnection field.
func (r *mutationResolver) UpdateConnection(ctx context.Context, config model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: UpdateConnection - updateConnection"))
}

// DeleteConnection is the resolver for the deleteConnection field.
func (r *mutationResolver) DeleteConnection(ctx context.Context, id string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteConnection - deleteConnection"))
}

// CreateConnectionFromTemplate is the resolver for the createConnectionFromTemplate field.
func (r *mutationResolver) CreateConnectionFromTemplate(ctx context.Context, templateID string, connectionName *string) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: CreateConnectionFromTemplate - createConnectionFromTemplate"))
}

// CreateConnectionFolder is the resolver for the createConnectionFolder field.
func (r *mutationResolver) CreateConnectionFolder(ctx context.Context, parentFolderPath *string, folderName string) (*model.ConnectionFolderInfo, error) {
	panic(fmt.Errorf("not implemented: CreateConnectionFolder - createConnectionFolder"))
}

// DeleteConnectionFolder is the resolver for the deleteConnectionFolder field.
func (r *mutationResolver) DeleteConnectionFolder(ctx context.Context, folderPath string) (bool, error) {
	panic(fmt.Errorf("not implemented: DeleteConnectionFolder - deleteConnectionFolder"))
}

// CopyConnectionFromNode is the resolver for the copyConnectionFromNode field.
func (r *mutationResolver) CopyConnectionFromNode(ctx context.Context, nodePath string, config *model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: CopyConnectionFromNode - copyConnectionFromNode"))
}

// TestConnection is the resolver for the testConnection field.
func (r *mutationResolver) TestConnection(ctx context.Context, config model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: TestConnection - testConnection"))
}

// TestNetworkHandler is the resolver for the testNetworkHandler field.
func (r *mutationResolver) TestNetworkHandler(ctx context.Context, config model.NetworkHandlerConfigInput) (*model.NetworkEndpointInfo, error) {
	panic(fmt.Errorf("not implemented: TestNetworkHandler - testNetworkHandler"))
}

// InitConnection is the resolver for the initConnection field.
func (r *mutationResolver) InitConnection(ctx context.Context, id string, credentials interface{}, networkCredentials []*model.NetworkHandlerConfigInput, saveCredentials *bool) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: InitConnection - initConnection"))
}

// CloseConnection is the resolver for the closeConnection field.
func (r *mutationResolver) CloseConnection(ctx context.Context, id string) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: CloseConnection - closeConnection"))
}

// SetConnectionNavigatorSettings is the resolver for the setConnectionNavigatorSettings field.
func (r *mutationResolver) SetConnectionNavigatorSettings(ctx context.Context, id string, settings model.NavigatorSettingsInput) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: SetConnectionNavigatorSettings - setConnectionNavigatorSettings"))
}

// AsyncTaskCancel is the resolver for the asyncTaskCancel field.
func (r *mutationResolver) AsyncTaskCancel(ctx context.Context, id string) (*bool, error) {
	panic(fmt.Errorf("not implemented: AsyncTaskCancel - asyncTaskCancel"))
}

// AsyncTaskInfo is the resolver for the asyncTaskInfo field.
func (r *mutationResolver) AsyncTaskInfo(ctx context.Context, id string, removeOnFinish bool) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: AsyncTaskInfo - asyncTaskInfo"))
}

// OpenConnection is the resolver for the openConnection field.
func (r *mutationResolver) OpenConnection(ctx context.Context, config model.ConnectionConfig) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: OpenConnection - openConnection"))
}

// AsyncTaskStatus is the resolver for the asyncTaskStatus field.
func (r *mutationResolver) AsyncTaskStatus(ctx context.Context, id string) (*model.AsyncTaskInfo, error) {
	panic(fmt.Errorf("not implemented: AsyncTaskStatus - asyncTaskStatus"))
}

// ServerConfig is the resolver for the serverConfig field.
func (r *queryResolver) ServerConfig(ctx context.Context) (*model.ServerConfig, error) {
	panic(fmt.Errorf("not implemented: ServerConfig - serverConfig"))
}

// SessionState is the resolver for the sessionState field.
func (r *queryResolver) SessionState(ctx context.Context) (*model.SessionInfo, error) {
	panic(fmt.Errorf("not implemented: SessionState - sessionState"))
}

// SessionPermissions is the resolver for the sessionPermissions field.
func (r *queryResolver) SessionPermissions(ctx context.Context) ([]*string, error) {
	panic(fmt.Errorf("not implemented: SessionPermissions - sessionPermissions"))
}

// DriverList is the resolver for the driverList field.
func (r *queryResolver) DriverList(ctx context.Context, id *string) ([]*model.DriverInfo, error) {
	panic(fmt.Errorf("not implemented: DriverList - driverList"))
}

// AuthModels is the resolver for the authModels field.
func (r *queryResolver) AuthModels(ctx context.Context) ([]*model.DatabaseAuthModel, error) {
	panic(fmt.Errorf("not implemented: AuthModels - authModels"))
}

// NetworkHandlers is the resolver for the networkHandlers field.
func (r *queryResolver) NetworkHandlers(ctx context.Context) ([]*model.NetworkHandlerDescriptor, error) {
	panic(fmt.Errorf("not implemented: NetworkHandlers - networkHandlers"))
}

// UserConnections is the resolver for the userConnections field.
func (r *queryResolver) UserConnections(ctx context.Context, id *string) ([]*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: UserConnections - userConnections"))
}

// TemplateConnections is the resolver for the templateConnections field.
func (r *queryResolver) TemplateConnections(ctx context.Context) ([]*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: TemplateConnections - templateConnections"))
}

// ConnectionFolders is the resolver for the connectionFolders field.
func (r *queryResolver) ConnectionFolders(ctx context.Context, path *string) ([]*model.ConnectionFolderInfo, error) {
	panic(fmt.Errorf("not implemented: ConnectionFolders - connectionFolders"))
}

// ConnectionState is the resolver for the connectionState field.
func (r *queryResolver) ConnectionState(ctx context.Context, id string) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: ConnectionState - connectionState"))
}

// ConnectionInfo is the resolver for the connectionInfo field.
func (r *queryResolver) ConnectionInfo(ctx context.Context, id string) (*model.ConnectionInfo, error) {
	panic(fmt.Errorf("not implemented: ConnectionInfo - connectionInfo"))
}

// ReadSessionLog is the resolver for the readSessionLog field.
func (r *queryResolver) ReadSessionLog(ctx context.Context, maxEntries *int, clearEntries *bool) ([]*model.LogEntry, error) {
	panic(fmt.Errorf("not implemented: ReadSessionLog - readSessionLog"))
}
