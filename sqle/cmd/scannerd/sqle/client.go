package sqle

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"actiontech.cloud/sqle/sqle/sqle/api/controller"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/config"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/utils"
	"actiontech.cloud/sqle/sqle/sqle/cmd/scannerd/utils/httpc"

	v1 "actiontech.cloud/sqle/sqle/sqle/api/controller/v1"
)

// %s = audit plan name
const (
	// Post
	LoginUri = "/v1/login"
	// Post
	TriggerAudit = "/v1/audit_plans/%s/trigger"
	// Post
	FullUpload = "/v1/audit_plans/%s/sqls/full"
	// Post
	PartialUpload = "/v1/audit_plans/%s/sqls/partial"
	// Get										%v=report_id
	GetAuditReport = "/v1/audit_plans/%s/report/%v/?page_index=%d&page_size=%d"
)

type (
	BaseRes                     = controller.BaseRes
	GetAuditPlanReportSQLsRes   = v1.GetAuditPlanReportSQLsResV1
	AuditPlanSQLReq             = v1.AuditPlanSQLReqV1
	FullSyncAuditPlanSQLsReq    = v1.FullSyncAuditPlanSQLsReqV1
	PartialSyncAuditPlanSQLsReq = v1.PartialSyncAuditPlanSQLsReqV1
	TriggerAuditPlanRes         = v1.TriggerAuditPlanResV1
)

type Client struct {
	baseURL    string
	httpClient *httpc.Client
	token      string
}

func NewSQLEClient(timeout time.Duration, cfg *config.Config) *Client {
	baseURL := fmt.Sprintf("%s:%v", cfg.Host, cfg.Port)
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	client := &Client{
		baseURL:    baseURL,
		httpClient: httpc.NewHTTPClient(timeout, nil),
	}

	return client
}
func (sc *Client) WithToken(token string) *Client {
	sc.token = token
	sc2 := *sc
	return &sc2
}

func (sc *Client) UploadReq(uri string, auditPlanName string, sqlList []AuditPlanSQLReq) error {
	url := sc.baseURL + fmt.Sprintf(uri, auditPlanName)

	reqBody := &FullSyncAuditPlanSQLsReq{
		SQLs: sqlList,
	}
	body, err := utils.JSONMarshal(reqBody)
	if err != nil {
		return err
	}

	resBody, err := sc.httpClient.SendRequest(context.TODO(), url, http.MethodPost, sc.token, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	baseRes := new(BaseRes)
	err = json.Unmarshal(resBody, baseRes)
	if err != nil {
		return err
	}
	if baseRes.Code != 0 {
		return fmt.Errorf("failed to request %s", url)
	}
	return nil
}

func (sc *Client) TriggerAuditReq(auditPlanName string) (string, error) {
	url := sc.baseURL + fmt.Sprintf(TriggerAudit, auditPlanName)

	resBody, err := sc.httpClient.SendRequest(context.TODO(), url, http.MethodPost, sc.token, nil)
	if err != nil {
		return "", err
	}

	triggerRes := new(TriggerAuditPlanRes)
	err = json.Unmarshal(resBody, triggerRes)
	if err != nil {
		return "", err
	}
	if triggerRes.Code != 0 {
		return "", fmt.Errorf("failed to request %s", url)
	}
	return triggerRes.Data.Id, nil
}

func (sc *Client) GetAuditReportReq(auditPlanName string, reportID string) error {
	var pageIndex, pageSize, cursor uint64
	pageIndex, pageSize = 1, 10
	cursor = pageIndex * pageSize

	for {
		url := sc.baseURL + fmt.Sprintf(GetAuditReport, auditPlanName, reportID, pageIndex, pageSize)
		resBody, err := sc.httpClient.SendRequest(context.TODO(), url, http.MethodGet, sc.token, nil)
		if err != nil {
			return err
		}

		auditRes := new(GetAuditPlanReportSQLsRes)
		err = json.Unmarshal(resBody, auditRes)
		if err != nil {
			return err
		}
		if auditRes.Code != 0 {
			return fmt.Errorf("failed to request %s", url)
		}
		for _, res := range auditRes.Data {
			fmt.Println(res.LastReceiveText)
			fmt.Println(res.AuditResult)
			if strings.Contains(res.AuditResult, "[error]") {
				return fmt.Errorf("audit result error, stopped")
			}
		}

		if cursor < auditRes.TotalNums {
			pageIndex++
			cursor = pageIndex * pageSize
		} else {
			break
		}
	}

	return nil
}
