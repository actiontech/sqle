package sqlfile

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/parse"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/sirupsen/logrus"
)

type SQLFile struct {
	l *logrus.Entry
	c *scanner.Client

	sqls []scanners.SQL

	allSQL []driverV2.Node
	getAll chan struct{}

	apName           string
	sqlDir           string
	skipErrorQuery   bool
	skipErrorSqlFile bool
	skipAudit        bool
}

type Params struct {
	SQLDir           string
	APName           string
	SkipErrorQuery   bool
	SkipErrorSqlFile bool
	SkipAudit        bool
}

func New(params *Params, l *logrus.Entry, c *scanner.Client) (*SQLFile, error) {
	return &SQLFile{
		sqlDir:           params.SQLDir,
		apName:           params.APName,
		skipErrorQuery:   params.SkipErrorQuery,
		skipErrorSqlFile: params.SkipErrorSqlFile,
		skipAudit:        params.SkipAudit,
		l:                l,
		c:                c,
		getAll:           make(chan struct{}),
	}, nil
}

func (sf *SQLFile) Run(ctx context.Context) error {
	sqls, err := GetSQLFromPath(sf.sqlDir, sf.skipErrorQuery, sf.skipErrorSqlFile)
	if err != nil {
		return err
	}

	sf.allSQL = sqls
	close(sf.getAll)

	<-ctx.Done()
	return nil
}

func (sf *SQLFile) SQLs() <-chan scanners.SQL {
	// todo: channel size configurable
	sqlCh := make(chan scanners.SQL, 10240)

	go func() {
		<-sf.getAll
		for _, sql := range sf.allSQL {
			sqlCh <- scanners.SQL{
				Fingerprint: sql.Fingerprint,
				RawText:     sql.Text,
			}
		}
		close(sqlCh)
	}()
	return sqlCh
}

func (sf *SQLFile) Upload(ctx context.Context, sqls []scanners.SQL) error {
	sf.sqls = append(sf.sqls, sqls...)

	// key=fingerPrint val=count
	counterMap := make(map[string]uint, len(sf.sqls))

	nodeList := make([]scanners.SQL, 0, len(sf.sqls))
	for _, node := range sf.sqls {
		counterMap[node.Fingerprint]++
		if counterMap[node.Fingerprint] <= 1 {
			nodeList = append(nodeList, node)
		}
	}

	reqBody := make([]*scanner.AuditPlanSQLReq, 0, len(nodeList))
	now := time.Now().Format(time.RFC3339)
	for _, sql := range nodeList {
		reqBody = append(reqBody, &scanner.AuditPlanSQLReq{
			Fingerprint:          sql.Fingerprint,
			Counter:              fmt.Sprintf("%v", counterMap[sql.Fingerprint]),
			LastReceiveText:      sql.RawText,
			LastReceiveTimestamp: now,
		})
	}

	err := sf.c.UploadReq(scanner.FullUpload, sf.apName, reqBody)
	if err != nil {
		return err
	}

	if sf.skipAudit {
		return nil
	}

	reportID, err := sf.c.TriggerAuditReq(sf.apName)
	if err != nil {
		return err
	}
	return sf.c.GetAuditReportReq(sf.apName, reportID)
}


func GetSQLFromPath(pathName string, skipErrorQuery, skipErrorSqlFile bool) (allSQL []driverV2.Node, err error) {
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
		var sqlList []driverV2.Node
		pathJoin := path.Join(pathName, fi.Name())

		if fi.IsDir() {
			sqlList, err = GetSQLFromPath(pathJoin, skipErrorQuery, skipErrorSqlFile)
		} else if strings.HasSuffix(fi.Name(), "sql") {
			sqlList, err = GetSQLFromFile(pathJoin)
		}

		if err != nil {
			if skipErrorSqlFile {
				fmt.Printf("[parse sql file error] parse file %s error: %v\n", pathJoin, err)
			} else {
				return nil, fmt.Errorf("parse file %s error: %v", pathJoin, err)
			}
		}
		allSQL = append(allSQL, sqlList...)
	}
	return allSQL, err
}

func GetSQLFromFile(file string) (r []driverV2.Node, err error) {
	content, err := ReadFileContent(file)
	if err != nil {
		return nil, err
	}
	n, err := parse.Parse(context.TODO(), content)
	if err != nil {
		return nil, err
	}
	r = append(r, n...)
	return r, nil
}

func ReadFileContent(file string) (content string, err error) {
	data, err := ioutil.ReadFile(filepath.Clean(file))
	if err != nil {
		return "", err
	}
	return string(data), err
}
