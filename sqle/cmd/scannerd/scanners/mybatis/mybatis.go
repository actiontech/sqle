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
	"github.com/actiontech/sqle/sqle/cmd/scannerd/config"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/utils/parser"
	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
	"github.com/sirupsen/logrus"
)

type MyBatis struct {
	l *logrus.Entry
	c *scanner.Client

	needTrigger bool
	sqls        []scanners.SQL

	apName string
	xmlDir string
}

type Params struct {
	XMLDir string
	APName string
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*MyBatis, error) {
	return &MyBatis{
		xmlDir: params.XMLDir,
		apName: params.APName,
		l:      l,
	}, nil
}

func (mb *MyBatis) Run(ctx context.Context) error {
	return nil
}

func (mb *MyBatis) SQLs() <-chan scanners.SQL {
	// todo: channel size configurable
	sqlCh := make(chan scanners.SQL, 10240)

	sqls, err := GetSQLFromPath(mb.xmlDir)
	if err != nil {
		mb.l.Errorf("parse sql error:%v", err)
		close(sqlCh)
		return sqlCh
	}

	go func() {
		for _, sql := range sqls {
			sqlCh <- scanners.SQL{
				Fingerprint: sql.Fingerprint,
				RawText:     sql.Text,
			}
		}
		mb.needTrigger = true
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

func MybatisScanner(cfg *config.Config) error {
	sqlList, err := PrepareSQL(cfg)
	if err != nil {
		return err
	}

	client := scanner.NewSQLEClient(time.Second, cfg).WithToken(cfg.Token)
	err = client.UploadReq(scanner.FullUpload, cfg.AuditPlanName, sqlList)
	if err != nil {
		return err
	}

	reportID, err := client.TriggerAuditReq(cfg.AuditPlanName)
	if err != nil {
		return err
	}

	err = client.GetAuditReportReq(cfg.AuditPlanName, reportID)
	if err != nil {
		return err
	}
	return nil
}

func PrepareSQL(cfg *config.Config) ([]scanner.AuditPlanSQLReq, error) {
	allSQL, err := GetSQLFromPath(cfg.Dir)
	if err != nil {
		return nil, err
	}

	// key=fingerPrint val=count
	counterMap := make(map[string]uint, len(allSQL))

	nodeList := make([]driver.Node, 0, len(allSQL))
	for _, node := range allSQL {
		counterMap[node.Fingerprint]++
		if counterMap[node.Fingerprint] <= 1 {
			nodeList = append(nodeList, node)
		}
	}

	reqBody := make([]scanner.AuditPlanSQLReq, 0, len(nodeList))
	for _, sql := range nodeList {
		reqBody = append(reqBody, scanner.AuditPlanSQLReq{
			Fingerprint:          sql.Fingerprint,
			Counter:              fmt.Sprintf("%v", counterMap[sql.Fingerprint]),
			LastReceiveText:      sql.Text,
			LastReceiveTimestamp: time.Now().Format(time.RFC3339),
		})
	}

	return reqBody, nil
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
		} else if strings.HasSuffix(fi.Name(), "xml") {
			sqlList, err = GetSQLFromFile(path.Join(pathName, fi.Name()))
		}

		if err != nil {
			return nil, err
		}
		allSQL = append(allSQL, sqlList...)
	}
	return allSQL, err
}

func GetSQLFromFile(file string) (r []driver.Node, err error) {
	content, err := ReadFileContent(file)
	if err != nil {
		return nil, err
	}
	sql, err := mybatisParser.ParseXML(content)
	if err != nil {
		return nil, err
	}
	return parser.Parse(context.TODO(), sql)
}

func ReadFileContent(file string) (content string, err error) {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", err
	}
	return string(data), err
}
