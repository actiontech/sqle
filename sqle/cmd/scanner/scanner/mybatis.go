package scanner

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"actiontech.cloud/sqle/sqle/sqle/cmd/scanner/config"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scanner/sqle"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scanner/utils/parser"
	"actiontech.cloud/sqle/sqle/sqle/driver"
	mybatisParser "github.com/actiontech/mybatis-mapper-2-sql"
)

func MybatisScanner(cfg *config.Config) error {
	sqlList, err := PrepareSQL(cfg)
	if err != nil {
		return err
	}

	client := sqle.NewSQLEClient(time.Second, cfg).WithToken(cfg.Token)
	err = client.UploadReq(sqle.FullUpload, cfg.AuditPlanName, sqlList)
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

func PrepareSQL(cfg *config.Config) ([]sqle.AuditPlanSQLReq, error) {
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

	reqBody := make([]sqle.AuditPlanSQLReq, 0, len(nodeList))
	for _, sql := range nodeList {
		reqBody = append(reqBody, sqle.AuditPlanSQLReq{
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
