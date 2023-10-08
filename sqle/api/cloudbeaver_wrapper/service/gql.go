package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	gqlClient "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/client"
)

var QueryGQL GetQueryGQL = CloudBeaverV2321{}

var (
	Version2215 = CBVersion{
		version: []int{22, 1, 5},
	}

	Version2221 = CBVersion{
		version: []int{22, 2, 1},
	}

	Version2223 = CBVersion{
		version: []int{22, 2, 3},
	}
	Version2321 = CBVersion{
		version: []int{23, 2, 1},
	}
)

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

	version, err := NewCBVersion(resp.ServerConfig.Version)
	if err != nil {
		return err
	}
	// QueryGQL 默认值是 CloudBeaverV2223{}
	if version.LessThan(Version2321) {
		QueryGQL = CloudBeaverV2223{}
	}
	if version.LessThan(Version2223) {
		QueryGQL = CloudBeaverV2221{}
	}
	if version.LessThan(Version2221) {
		QueryGQL = CloudBeaverV2215{}
	}
	if version.LessThan(Version2215) {
		return fmt.Errorf("CloudBeaver version less than 22.1.5 are not supported temporarily, your version is %v", resp.ServerConfig.Version)
	}

	return nil
}

// CloudBeaver 版本号格式一般为 X.X.X.X 格式,例如 '22.3.1.202212261505' , 其中前三位为版本号
type CBVersion struct {
	version []int // version是版本号使用'.'进行分割后的数组
}

func NewCBVersion(versionStr string) (CBVersion, error) {
	versions := strings.Split(versionStr, ".")
	if len(versions) < 3 {
		return CBVersion{}, fmt.Errorf("CloudBeaver version number that cannot be resolved")
	}
	cb := CBVersion{
		version: []int{},
	}
	for _, version := range versions {
		v, err := strconv.Atoi(version)
		if err != nil {
			return CBVersion{}, fmt.Errorf("CloudBeaver version number that cannot be resolved")
		}
		cb.version = append(cb.version, v)
	}
	return cb, nil
}

// 只比较前三位, 因为只有前三位与版本有关
func (v CBVersion) LessThan(version CBVersion) bool {
	if v.version[0] < version.version[0] {
		return true
	}
	if v.version[0] > version.version[0] {
		return false
	}
	if v.version[1] < version.version[1] {
		return true
	}
	if v.version[1] > version.version[1] {
		return false
	}
	return v.version[2] < version.version[2]

}

// 不同版本的CloudBeaver间存在不兼容查询语句
// 说明: 接口传参时不要删除旧版查询语句的查询参数,多余的参数对新接口没有影响但是可以兼容旧版本
type GetQueryGQL interface {
	CreateConnectionQuery() string
	UpdateConnectionQuery() string
	GetUserConnectionsQuery() string
	SetUserConnectionsQuery() string
	IsUserExistQuery(userId string) (string, map[string]interface{})
	UpdatePasswordQuery() string
	CreateUserQuery() string
	GrantUserRoleQuery() string
	LoginQuery() string
	GetActiveUserQuery() string
}

// TODO 暂时无法确定这套查询语句是兼容到22.1.5版本还是22.1.4版本, 因为虽然找到了22.1.4版本的镜像, 但没找到22.1.4版本的代码
type CloudBeaverV2215 struct{}

func (CloudBeaverV2215) CreateConnectionQuery() string {
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

func (CloudBeaverV2215) UpdateConnectionQuery() string {
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

func (CloudBeaverV2215) GetUserConnectionsQuery() string {
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

func (CloudBeaverV2215) SetUserConnectionsQuery() string {
	return `
query setConnections($userId: ID!, $connections: [ID!]!) {
  grantedConnections: setSubjectConnectionAccess(
    subjectId: $userId
    connections: $connections
  )
}
`
}

func (CloudBeaverV2215) IsUserExistQuery(userId string) (string, map[string]interface{}) {
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
}`, map[string]interface{}{"userId": userId}
}

func (CloudBeaverV2215) UpdatePasswordQuery() string {
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

func (CloudBeaverV2215) CreateUserQuery() string {
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

func (CloudBeaverV2215) GrantUserRoleQuery() string {
	return `
query grantUserRole($userId: ID!, $roleId: ID!) {
  grantUserRole(userId: $userId, roleId: $roleId)
}`
}

func (CloudBeaverV2215) LoginQuery() string {
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

func (CloudBeaverV2215) GetActiveUserQuery() string {
	return `
	query getActiveUser {
  		user: activeUser {
    		userId
  		}
	}
`
}

type CloudBeaverV2221 struct {
	CloudBeaverV2215
}

func (CloudBeaverV2221) CreateUserQuery() string {
	return `
query createUser(
  $userId: ID!
) {
  user: createUser(userId: $userId, enabled: true) {
    ...adminUserInfo
  }
}
fragment adminUserInfo on AdminUserInfo {
	userId
}
`
}

type CloudBeaverV2223 struct {
	CloudBeaverV2221
}

func (CloudBeaverV2223) CreateUserQuery() string {
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

func (CloudBeaverV2223) GrantUserRoleQuery() string {
	return `
query grantUserTeam($userId: ID!, $teamId: ID!) {
  grantUserTeam(userId: $userId, teamId: $teamId)
}`
}

type CloudBeaverV2321 struct {
	CloudBeaverV2221
}

func (CloudBeaverV2321) IsUserExistQuery(userId string) (string, map[string]interface{}) {

	return `
query getUserList(
	$page: PageInput!
	$filter: AdminUserFilterInput!
){
	listUsers(page: $page, filter: $filter) {
		...adminUserInfo
	}
}
fragment adminUserInfo on AdminUserInfo {
	userId
}
`, map[string]interface{}{
			"page":   map[string]interface{}{"offset": 0, "limit": 100},
			"filter": map[string]interface{}{"userIdMask": userId, "enabledState": true},
		}
}

func (CloudBeaverV2321) GetActiveUserQuery() string {
	return `
	query getActiveUser {
  		user: activeUser {
    		userId
  		}
	}
`
}
