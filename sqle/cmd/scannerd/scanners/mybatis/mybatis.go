package mybatis

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	mybatisParser "github.com/actiontech/mybatis-mapper-2-sql"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/sirupsen/logrus"
)

type MyBatis struct {
	l *logrus.Entry
	c *scanner.Client

	needTrigger bool
	sqls        []scanners.SQL

	allSQL []driver.Node
	getAll chan struct{}

	apName         string
	xmlDir         string
	skipErrorQuery bool
}

type Params struct {
	XMLDir         string
	APName         string
	SkipErrorQuery bool
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*MyBatis, error) {
	return &MyBatis{
		xmlDir:         params.XMLDir,
		apName:         params.APName,
		skipErrorQuery: params.SkipErrorQuery,
		l:              l,
		c:              c,
		getAll:         make(chan struct{}),
	}, nil
}

func (mb *MyBatis) Run(ctx context.Context) error {
	sqls, err := GetSQLFromPath(mb.xmlDir, mb.skipErrorQuery)
	if err != nil {
		return err
	}

	mb.allSQL = sqls
	close(mb.getAll)

	<-ctx.Done()
	return nil
}

func (mb *MyBatis) SQLs() <-chan scanners.SQL {
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

func (mb *MyBatis) Upload(ctx context.Context, sqls []scanners.SQL) error {
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

	err := mb.c.UploadReq(scanner.FullUpload, mb.apName, reqBody)
	if err != nil {
		return err
	}

	reportID, err := mb.c.TriggerAuditReq(mb.apName)
	if err != nil {
		return err
	}
	return mb.c.GetAuditReportReq(mb.apName, reportID)
}

func GetSQLFromPath(pathName string, skipErrorQuery bool) (allSQL []driver.Node, err error) {
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
			sqlList, err = GetSQLFromPath(path.Join(pathName, fi.Name()), skipErrorQuery)
		} else if strings.HasSuffix(fi.Name(), "xml") {
			sqlList, err = GetSQLFromFile(path.Join(pathName, fi.Name()), skipErrorQuery)
		}

		if err != nil {
			return nil, err
		}
		allSQL = append(allSQL, sqlList...)
	}
	return allSQL, err
}

func GetSQLFromFile(file string, skipErrorQuery bool) (r []driver.Node, err error) {
	content, err := ReadFileContent(file)
	if err != nil {
		return nil, err
	}
	sqls, err := mybatisParser.ParseXMLQuery(content, skipErrorQuery)
	if err != nil {
		return nil, err
	}
	for _, sql := range sqls {
		n, err := Parse(context.TODO(), sql)
		if err != nil {
			return nil, err
		}
		r = append(r, n...)
	}
	return r, nil
}

func ReadFileContent(file string) (content string, err error) {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", err
	}
	return string(data), err
}
