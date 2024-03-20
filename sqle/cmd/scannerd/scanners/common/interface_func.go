package common

import (
	"context"
	"fmt"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/pkg/scanner"
)

func Upload(ctx context.Context, sqls []scanners.SQL, c *scanner.Client, apName string) error {
	// key=fingerPrint val=count
	counterMap := make(map[string]uint, len(sqls))

	nodeList := make([]scanners.SQL, 0, len(sqls))
	for _, node := range sqls {
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

	err := c.UploadReq(scanner.FullUpload, apName, reqBody)
	return err
}

func Audit(c *scanner.Client, apName string) error {
	reportID, err := c.TriggerAuditReq(apName)
	if err != nil {
		return err
	}
	return c.GetAuditReportReq(apName, reportID)
}

func UploadForDmSlowLog(ctx context.Context, sqls []scanners.SQL, c *scanner.Client, apName string) error {
	// key=fingerPrint val=count
	counterMap := make(map[string]uint, len(sqls))
	queryTimeTotalMap := make(map[string]float64)
	queryTimeMaxMap := make(map[string]float64)

	nodeList := make([]scanners.SQL, 0, len(sqls))
	for _, node := range sqls {
		counterMap[node.Fingerprint]++
		if counterMap[node.Fingerprint] <= 1 {
			nodeList = append(nodeList, node)
		}
		queryTimeTotal, ok := queryTimeTotalMap[node.Fingerprint]
		if !ok {
			queryTimeTotalMap[node.Fingerprint] = node.QueryTime
		} else {
			queryTimeTotalMap[node.Fingerprint] = queryTimeTotal + node.QueryTime
		}
		queryTimeMax, ok := queryTimeMaxMap[node.Fingerprint]
		if !ok || node.QueryTime > queryTimeMax {
		    queryTimeMaxMap[node.Fingerprint] = node.QueryTime
		}
	}

	reqBody := make([]*scanner.AuditPlanSQLReq, 0, len(nodeList))
	for _, sql := range nodeList {
		queryTimeTotal, _ := queryTimeTotalMap[sql.Fingerprint]
		counter, _ := counterMap[sql.Fingerprint]
		queryTimeAvg := queryTimeTotal / float64(counter)
		queryTimeMax, _ := queryTimeMaxMap[sql.Fingerprint]
		reqBody = append(reqBody, &scanner.AuditPlanSQLReq{
			Fingerprint:          sql.Fingerprint,
			Counter:              fmt.Sprintf("%v", counterMap[sql.Fingerprint]),
			LastReceiveText:      sql.RawText,
			LastReceiveTimestamp: sql.QueryAt.Format("2006-01-02 15:04:05"),
			QueryTimeAvg:         &queryTimeAvg,
			QueryTimeMax:         &queryTimeMax,
			DBUser:               sql.DBUser,
			Schema:               sql.Schema,
		})
	}

	err := c.UploadReq(scanner.PartialUpload, apName, reqBody)
	return err
}
