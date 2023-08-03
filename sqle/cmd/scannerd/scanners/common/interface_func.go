package common

import (
	"context"
	"time"
	"fmt"

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
