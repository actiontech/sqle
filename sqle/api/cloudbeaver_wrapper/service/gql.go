package service

import (
	"context"
	"fmt"
	"strings"

	gqlClient "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/client"
)

var QueryGQL GetQueryGQL = CloudBeaverV2220{}

const getServerConfigQuery = `
query serverConfig {
  serverConfig {
    version
  }
}`

func InitGQLVersion() error {
	client := gqlClient.NewClient(GetGqlServerURI())
	req := gqlClient.NewRequest(getServerConfigQuery, map[string]interface{}{})
	resp := struct {
		ServerConfig struct {
			Version string `json:"version"`
		} `json:"serverConfig"`
	}{}
	err := client.Run(context.Background(), req, &resp)
	if err != nil {
		return err
	}
	versions := strings.Split(resp.ServerConfig.Version, ".")
	if len(versions) < 3 {
		return fmt.Errorf("CloudBeaver version number that cannot be resolved")
	}

	switch fmt.Sprintf("%s%s%s", versions[0], versions[1], versions[2]) {
	case Version2231:
		QueryGQL = CloudBeaverV2231{}
	}
	return nil
}

// 不同版本的CloudBeaver间存在不兼容查询语句
type GetQueryGQL interface {
	CreateConnectionQuery() string
	UpdateConnectionQuery() string
	GetUserConnectionsQuery() string
	SetUserConnectionsQuery() string
	IsUserExistQuery() string
	UpdatePasswordQuery() string
	CreateUserQuery() string
	GrantUserRoleQuery() string
	LoginQuery() string
	GetActiveUserQuery() string
}

type CloudBeaverV2220 struct{}

func (CloudBeaverV2220) CreateConnectionQuery() string {
	return `
mutation createConnection(
  $projectId: ID!
  $config: ConnectionConfig!
) {
  connection: createConnection(projectId: $projectId, config: $config) {
    ...DatabaseConnection
  }
}

fragment DatabaseConnection on ConnectionInfo {
  id
}
`
}

func (CloudBeaverV2220) UpdateConnectionQuery() string {
	return `
mutation updateConnection(
  $projectId: ID!
  $config: ConnectionConfig!
) {
  connection: updateConnection(projectId: $projectId, config: $config) {
    ...DatabaseConnection
  }
}

fragment DatabaseConnection on ConnectionInfo {
  id
}
`
}

func (CloudBeaverV2220) GetUserConnectionsQuery() string {
	return `
query getUserConnections (
  $projectId: ID
  $connectionId: ID
){
  connections: userConnections(projectId: $projectId, id: $connectionId) {
    ...DatabaseConnection
  }
}

fragment DatabaseConnection on ConnectionInfo {
  id
}
`
}

func (CloudBeaverV2220) SetUserConnectionsQuery() string {
	return `
query setConnections($userId: ID!, $connections: [ID!]!) {
  grantedConnections: setSubjectConnectionAccess(
    subjectId: $userId
    connections: $connections
  )
}
`
}

func (CloudBeaverV2220) IsUserExistQuery() string {
	return `
query getUserList(
	$userId: ID
){
	listUsers(userId: $userId) {
		...adminUserInfo
	}
}

fragment adminUserInfo on AdminUserInfo {
	userId
}
`
}

func (CloudBeaverV2220) UpdatePasswordQuery() string {
	return `
query setUserCredentials(
  $userId: ID!
  $credentials: Object!
) {
  setUserCredentials(
    userId: $userId
    providerId: "local"
    credentials: $credentials
  )
}
`
}

func (CloudBeaverV2220) CreateUserQuery() string {
	return `
query createUser(
  $userId: ID!
) {
  user: createUser(userId: $userId) {
    ...adminUserInfo
  }
}

fragment adminUserInfo on AdminUserInfo {
	userId
}
`
}

func (CloudBeaverV2220) GrantUserRoleQuery() string {
	return `
query grantUserRole($userId: ID!, $roleId: ID!) {
  grantUserRole(userId: $userId, roleId: $roleId)
}`
}

func (CloudBeaverV2220) LoginQuery() string {
	return `
query authLogin(
  $credentials: Object
) {
	authInfo: authLogin(
    	provider: "local"
    	configuration: null
    	credentials: $credentials
    	linkUser: false
  ){
    authId
  }
}
`
}

func (CloudBeaverV2220) GetActiveUserQuery() string {
	return `
	query getActiveUser {
  		user: activeUser {
    		userId
  		}
	}
`
}

type CloudBeaverV2231 struct {
	CloudBeaverV2220
}

func (CloudBeaverV2231) CreateUserQuery() string {
	return `
query createUser(
  $userId: ID!
) {
  user: createUser(userId: $userId, enabled: true, authRole: null) {
    ...adminUserInfo
  }
}

fragment adminUserInfo on AdminUserInfo {
	userId
}
`
}

func (CloudBeaverV2231) GrantUserRoleQuery() string {
	return `
query grantUserTeam($userId: ID!, $teamId: ID!) {
  grantUserTeam(userId: $userId, teamId: $teamId)
}`
}
