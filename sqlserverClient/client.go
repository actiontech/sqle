package sqlserverClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"sqle/errors"
	"sqle/model"
	"sqle/sqlserver/SqlserverProto"
)

var GrpcClient *Client

func GetClient() *Client {
	return GrpcClient
}

type Client struct {
	version string
	conn    *grpc.ClientConn
	client  SqlserverProto.SqlserverServiceClient
}

func InitClient(ip, port string) error {
	c := &Client{}
	err := c.Conn(ip, port)
	if err != nil {
		return err
	}
	GrpcClient = c
	return nil
}

func (c *Client) Conn(ip, port string) error {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", ip, port), grpc.WithInsecure())
	if err != nil {
		return errors.New(errors.CONNECT_SQLSERVER_RPC_ERROR, err)
	}
	c.conn = conn
	c.client = SqlserverProto.NewSqlserverServiceClient(conn)
	return nil
}

func (c *Client) SplitSql(sql string) ([]string, error) {
	out, err := c.client.GetSplitSqls(context.Background(), &SqlserverProto.SplitSqlsInput{
		Sqls:    sql,
		Version: c.version,
	})
	return out.GetSqls(), errors.New(errors.CONNECT_SQLSERVER_RPC_ERROR, err)
}

func (c *Client) Advise(commitSqls []*model.CommitSql, rules []model.Rule) error {
	sqls := []string{}
	ruleNames := []string{}
	for _, commitSql := range commitSqls {
		sqls = append(sqls, commitSql.Content)
	}

	for _, rule := range rules {
		ruleNames = append(ruleNames, rule.Name)
	}
	out, err := c.client.Audit(context.Background(), &SqlserverProto.AuditInput{
		Version:   c.version,
		Sqls:      sqls,
		RuleNames: ruleNames,
	})
	results := out.GetAuditResults()
	if len(results) != len(commitSqls) {
		return errors.New(errors.CONNECT_REMOTE_DB_ERROR, fmt.Errorf("don't match sql advise result"))
	}

	for n, result := range results {
		commitSql := commitSqls[n]
		commitSql.InspectLevel = result.AuditLevel
		commitSql.InspectResult = result.AuditResultMessage
		commitSql.InspectStatus = model.TASK_ACTION_DONE
	}
	return errors.New(errors.CONNECT_SQLSERVER_RPC_ERROR, err)
}
