package sqltext

import (
	"context"
	"fmt"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/mybatis"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/sirupsen/logrus"
)

type SqlText struct {
	l *logrus.Entry
	c *scanner.Client

	needTrigger bool
	sqls        []scanners.SQL

	allSQL []driver.Node
	getAll chan struct{}

	apName string
	sqlDir string
	sql    string
	audit  string
	synctype string
}

type Params struct {
	SQL    string
	SQLDir string
	APName string
	AUDIT  string
	SYNCTYPE string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SqlText, error) {
	return &SqlText{
		sql:    params.SQL,
		sqlDir: params.SQLDir,
		apName: params.APName,
		audit:  params.AUDIT,
		synctype:  params.SYNCTYPE,
		l:      l,
		c:      c,
		getAll: make(chan struct{}),
	}, nil
}

func (mb *SqlText) Run(ctx context.Context) error {
	sqls, err := GetSQLFromCli(mb.sql)

	if sqls == nil {
		sqls, err = GetSQLFromPath(mb.sqlDir)
	}

	if err != nil {
		return err
	}

	mb.allSQL = sqls
	close(mb.getAll)

	<-ctx.Done()
	return nil
}

func (mb *SqlText) SQLs() <-chan scanners.SQL {
	// todo: channel size configurable
	sqlCh := make(chan scanners.SQL, 10240)

	go func() {
		<-mb.getAll
		for _, sql := range mb.allSQL {
			sqlCh <- scanners.SQL{
				Fingerprint: sql.Fingerprint,
				RawText:     sql.Text,
			}
		}
		mb.needTrigger = true
		close(sqlCh)
	}()
	return sqlCh
}

func (mb *SqlText) Upload(ctx context.Context, sqls []scanners.SQL) error {
	mb.sqls = append(mb.sqls, sqls...)

	if !mb.needTrigger {
		return nil
	}

	// key=fingerPrint val=count
	counterMap := make(map[string]uint, len(mb.sqls))

	nodeList := make([]scanners.SQL, 0, len(mb.sqls))
	for _, node := range mb.sqls {
		counterMap[node.Fingerprint]++
		if counterMap[node.Fingerprint] <= 1 {
			nodeList = append(nodeList, node)
		}
	}

	reqBody := make([]scanner.AuditPlanSQLReq, 0, len(nodeList))
	now := time.Now().Format(time.RFC3339)
	for _, sql := range nodeList {
		reqBody = append(reqBody, scanner.AuditPlanSQLReq{
			Fingerprint:          sql.Fingerprint,
			Counter:              fmt.Sprintf("%v", counterMap[sql.Fingerprint]),
			LastReceiveText:      sql.RawText,
			LastReceiveTimestamp: now,
		})
	}

	reqUrl :=  scanner.FullUpload
	if  mb.synctype == "2" {
		reqUrl = scanner.PartialUpload
	}

	err := mb.c.UploadReq(reqUrl, mb.apName, reqBody)
	if err != nil {
		return err
	}

	if mb.audit == "1" || mb.audit == "true" { //trigger audit immediately 立即触发审核
		reportID, err := mb.c.TriggerAuditReq(mb.apName)
		if err != nil {
			return err
		}
		return mb.c.GetAuditReportReq(mb.apName, reportID)
	}

	return nil
}

func GetSQLFromPath(pathName string) (allSQL []driver.Node, err error) {
	if !path.IsAbs(pathName) {
		pwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		pathName = path.Join(pwd, pathName)
	}

	fileInfos, err := ioutil.ReadDir(pathName)
	if err != nil {
		return nil, err
	}
	for _, fi := range fileInfos {
		var sqlList []driver.Node
		if fi.IsDir() {
			sqlList, err = GetSQLFromPath(path.Join(pathName, fi.Name()))
		} else if strings.HasSuffix(fi.Name(), "sql.txt") {
			sqlList, err = GetSQLFromFile(path.Join(pathName, fi.Name()))
		}

		if err != nil {
			return nil, err
		}
		allSQL = append(allSQL, sqlList...)
	}
	return allSQL, err
}

func GetSQLFromCli(sql string) (r []driver.Node, err error) {
	if sql == "" {
		return nil, err
	}
	return mybatis.Parse(context.TODO(), sql)
}

func GetSQLFromFile(file string) (r []driver.Node, err error) {
	sql, err := ReadFileContent(file)
	if err != nil {
		return nil, err
	}
	return mybatis.Parse(context.TODO(), sql)
}

func ReadFileContent(file string) (content string, err error) {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", err
	}
	return string(data), err
}
