package service

import (
	"fmt"
	"sync"

	gqlClient "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/client"
	"github.com/actiontech/sqle/sqle/config"

	"github.com/labstack/echo/v4"
)

type SQLQueryConfig struct {
	config.SQLQueryConfig
	SqlePort        int
	SqleEnableHttps bool
}

var (
	cfg   = &SQLQueryConfig{}
	cfgMu = &sync.RWMutex{}
)

const (
	CbRootUri  = "/sql_query"
	CbGqlApi   = "/api/gql"
	CBUserRole = "user"
)

// 这个客户端会用当前用户操作, 请求会发给SQLE, 由SQLE转发到CB
func GetSQLEGQLClientWithCurrentUser(ctx echo.Context) (*gqlClient.Client, error) {
	return gqlClient.NewClient(GetSQLEGqlServerURI(), gqlClient.WithCookie(ctx.Cookies())), nil
}

// 这个客户端会用CB管理员操作, 请求会直接发到CB
func GetGQLClientWithRootUser() (*gqlClient.Client, error) {
	cookies, err := LoginToCBServer(GetSQLQueryConfig().CloudBeaverAdminUser, GetSQLQueryConfig().CloudBeaverAdminPassword)
	if err != nil {
		return nil, err
	}
	return gqlClient.NewClient(GetGqlServerURI(), gqlClient.WithCookie(cookies)), nil
}

// 这个客户端会用指定用户操作, 请求会直接发到CB
func GetGQLClient(username, password string) (*gqlClient.Client, error) {
	cookies, err := LoginToCBServer(username, password)
	if err != nil {
		return nil, err
	}
	return gqlClient.NewClient(GetGqlServerURI(), gqlClient.WithCookie(cookies)), nil
}

func IsCloudBeaverConfigured() bool {
	c := GetSQLQueryConfig()
	return c.SqlePort != 0 && c.CloudBeaverHost != "" && c.CloudBeaverPort != "" && c.CloudBeaverAdminUser != "" && c.CloudBeaverAdminPassword != ""
}

func GetGqlServerURI() string {
	protocol := "http"
	if cfg.EnableHttps {
		protocol = "https"
	}

	c := GetSQLQueryConfig()

	return fmt.Sprintf("%v://%v:%v%v%v", protocol, c.CloudBeaverHost, c.CloudBeaverPort, CbRootUri, CbGqlApi)
}

func GetSQLEGqlServerURI() string {
	protocol := "http"
	if cfg.EnableHttps {
		protocol = "https"
	}

	c := GetSQLQueryConfig()

	return fmt.Sprintf("%v://localhost:%v%v%v", protocol, c.SqlePort, CbRootUri, CbGqlApi)
}

func InitSQLQueryConfig(sqlePort int, sqleEnableHttps bool, c config.SQLQueryConfig) {
	cfgMu.Lock()
	cfg.SqlePort = sqlePort
	cfg.SQLQueryConfig = c
	cfg.SqleEnableHttps = sqleEnableHttps
	cfgMu.Unlock()
}

func GetSQLQueryConfig() *SQLQueryConfig {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return cfg
}
