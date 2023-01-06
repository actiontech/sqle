package service

var QueryGQL GetQueryGQL = CloudBeaverV2220{}

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

type CloudBeaverV2231 struct{}

func (CloudBeaverV2231) CreateConnectionQuery() string {
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

func (CloudBeaverV2231) UpdateConnectionQuery() string {
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

func (CloudBeaverV2231) GetUserConnectionsQuery() string {
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

func (CloudBeaverV2231) SetUserConnectionsQuery() string {
	return `
query setConnections($userId: ID!, $connections: [ID!]!) {
  grantedConnections: setSubjectConnectionAccess(
    subjectId: $userId
    connections: $connections
  )
}
`
}

func (CloudBeaverV2231) IsUserExistQuery() string {
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

func (CloudBeaverV2231) UpdatePasswordQuery() string {
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

func (CloudBeaverV2231) LoginQuery() string {
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

func (CloudBeaverV2231) GetActiveUserQuery() string {
	return `
	query getActiveUser {
  		user: activeUser {
    		userId
  		}
	}
`
}
