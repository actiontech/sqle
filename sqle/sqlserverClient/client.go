package sqlserverClient

import (
	"context"
	"fmt"
	"actiontech.cloud/universe/sqle/v3/sqle/errors"
	"actiontech.cloud/universe/sqle/v3/sqle/log"
	"actiontech.cloud/universe/sqle/v3/sqle/model"
	"actiontech.cloud/universe/sqle/v3/sqle/sqlserver/SqlserverProto"
	"time"

	"github.com/pingcap/tidb/ast"
	"google.golang.org/grpc"
)

var GrpcClient *Client

func GetClient() *Client {
	return GrpcClient
}

func GetSqlserverMeta(user, password, host, port, dbName, schemaName string) *SqlserverProto.SqlserverMeta {
	return &SqlserverProto.SqlserverMeta{
		User:            user,
		Password:        password,
		Host:            host,
		Port:            port,
		CurrentDatabase: dbName,
		CurrentSchema:   schemaName,
	}
}

type Client struct {
	version string
	conn    *grpc.ClientConn
	client  SqlserverProto.SqlserverServiceClient
}

func InitClient(ip, port string) error {
	log.Logger().Infof("connecting to SQLServer parser server %s:%s", ip, port)
	c := &Client{}
	err := c.Conn(ip, port)
	if err != nil {
		log.Logger().Warnf("connect to SQLServer parser server failed, error: %v", err)
		return err
	}
	log.Logger().Info("connected to SQLServer parser server")
	GrpcClient = c
	return nil
}

func (c *Client) Conn(ip, port string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%s", ip, port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.New(errors.CONNECT_SQLSERVER_RPC_ERROR, err)
	}
	c.conn = conn
	c.client = SqlserverProto.NewSqlserverServiceClient(conn)
	return nil
}

func (c *Client) ParseSql(sql string) ([]ast.Node, error) {
	out, err := c.client.GetSplitSqls(context.Background(), &SqlserverProto.SplitSqlsInput{
		Sqls:    sql,
		Version: c.version,
	})
	sqls := out.GetSplitSqls()
	stmts := make([]ast.Node, 0, len(sqls))
	for _, s := range sqls {
		stmts = append(stmts, NewSqlServerStmt(s.Sql, s.IsDDL, s.IsDML, s.IsProcedure, s.IsFunction))
	}
	return stmts, errors.New(errors.CONNECT_SQLSERVER_RPC_ERROR, err)
}

func (c *Client) Advise(sqls []string, ruleNames []string, meta *SqlserverProto.SqlserverMeta, ddlContextSqls []*SqlserverProto.DDLContext) (*SqlserverProto.AdviseOutput, error) {
	out, err := c.client.Advise(context.Background(), &SqlserverProto.AdviseInput{
		Version:        c.version,
		Sqls:           sqls,
		RuleNames:      ruleNames,
		SqlserverMeta:  meta,
		DDLContextSqls: ddlContextSqls,
	})
	if err != nil {
		return nil, errors.New(errors.CONNECT_SQLSERVER_RPC_ERROR, err)
	}

	if len(out.GetResults()) != len(sqls) {
		return nil, errors.New(errors.CONNECT_REMOTE_DB_ERROR, fmt.Errorf("don't match sql advise result"))
	}

	return out, nil
}

func (c *Client) GenerateAllRollbackSql(commitSqls []*model.CommitSql, config *SqlserverProto.Config, meta *SqlserverProto.SqlserverMeta) ([]*SqlserverProto.Sql, error) {
	sqls := []string{}
	for _, commitSql := range commitSqls {
		sqls = append(sqls, commitSql.Content)
	}
	out, err := c.client.GetRollbackSqls(context.Background(), &SqlserverProto.GetRollbackSqlsInput{
		Version:        c.version,
		Sqls:           sqls,
		SqlserverMeta:  meta,
		RollbackConfig: config,
	})
	return out.GetRollbackSqls(), err
}

func (c *Client) GetProcedureFunctionBackupSql(sql string, meta *SqlserverProto.SqlserverMeta) ([]string, error) {
	out, err := c.client.GetProcedureFunctionBackupSql(context.Background(), &SqlserverProto.GetProcedureFunctionBackupSqlInput{
		Sql:           sql,
		SqlserverMeta: meta,
	})
	return out.GetBackupSqls(), err
}
