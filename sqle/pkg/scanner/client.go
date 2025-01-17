package scanner

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	v2 "github.com/actiontech/sqle/sqle/api/controller/v2"
)

// %s = audit plan name
const (
	// Post
	LoginUri = "/sqle/v1/login"
	// Post
	TriggerAudit = "/sqle/v1/projects/%v/audit_plans/%s/trigger"
	// Post
	FullUpload = "/sqle/v2/projects/%v/audit_plans/%s/sqls/full"
	// Post
	PartialUpload = "/sqle/v2/projects/%v/audit_plans/%s/sqls/partial"
	// Get										                  %v=report_id
	GetAuditReport = "/sqle/v1/projects/%v/audit_plans/%s/reports/%v/sqls?page_index=%d&page_size=%d"
	// 创建sql审核记录
	CreateSqlAudit = "/sqle/v1/projects/%v/sql_audit_records"
	// 获取sql审核记录
	GetSqlAudit = "/sqle/v1/projects/%v/sql_audit_records/%v/"
	// 获取task sqls
	GetTaskSQLs = "/sqle/v2/tasks/audits/%v/sqls?page_index=%d&page_size=%d"
	// 获取所有项目
	GetAllProjects = "/v1/dms/projects?page_index=%d&page_size=%d"
)

// %s = project name
// %s = instance audit plan id
// %s = instance audit plan type
const (
	// Post
	UploadSQL = "/sqle/v2/projects/%v/audit_plans/%s/sqls/upload"
)

type (
	BaseRes                     = controller.BaseRes
	GetAuditPlanReportSQLsRes   = v1.GetAuditPlanReportSQLsResV1
	AuditPlanSQLReq             = v2.AuditPlanSQLReqV2
	FullSyncAuditPlanSQLsReq    = v2.FullSyncAuditPlanSQLsReqV2
	PartialSyncAuditPlanSQLsReq = v1.PartialSyncAuditPlanSQLsReqV1
	TriggerAuditPlanRes         = v1.TriggerAuditPlanResV1
	CreateSqlAuditReq           = v1.CreateSQLAuditRecordReqV1
	CreateSqlAuditResp          = v1.CreateSQLAuditRecordResV1
	GetSqlAuditResp             = v1.GetSQLAuditRecordResV1
	GetAuditTaskSqls            = v2.GetAuditTaskSQLsResV2
)

type Client struct {
	baseURL    string
	httpClient *client
	token      string
	project    string
}

func NewSQLEClient(timeout time.Duration, host, port string) *Client {
	baseURL := fmt.Sprintf("%s:%v", host, port)
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "http://" + baseURL
	}

	client := &Client{
		baseURL:    baseURL,
		httpClient: newClient(timeout, nil),
	}

	return client
}

func (sc *Client) WithToken(token string) *Client {
	sc.token = token
	sc2 := *sc
	return &sc2
}

func (sc *Client) WithProject(project string) *Client {
	sc.project = project
	sc2 := *sc
	return &sc2
}

func (sc *Client) UploadReq(uri, auditPlanID, errorMessage string, sqlList []*AuditPlanSQLReq) error {
	bodyBuf := &bytes.Buffer{}
	encoder := json.NewEncoder(bodyBuf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(&FullSyncAuditPlanSQLsReq{
		SQLs:         sqlList,
		ErrorMessage: errorMessage,
	})
	if err != nil {
		return err
	}

	url := sc.baseURL + fmt.Sprintf(uri, sc.project, auditPlanID)
	resBody, err := sc.httpClient.sendRequest(context.TODO(), url, http.MethodPost, sc.token, bytes.NewBuffer(bodyBuf.Bytes()))
	if err != nil {
		return err
	}

	baseRes := new(BaseRes)
	err = json.Unmarshal(resBody, baseRes)
	if err != nil {
		return err
	}
	if baseRes.Code != 0 {
		return fmt.Errorf("failed to request %s, error:%s", url, baseRes.Message)
	}
	return nil
}

func (sc *Client) TriggerAuditReq(auditPlanName string) (string, error) {
	url := sc.baseURL + fmt.Sprintf(TriggerAudit, sc.project, auditPlanName)

	resBody, err := sc.httpClient.sendRequest(context.TODO(), url, http.MethodPost, sc.token, nil)
	if err != nil {
		return "", err
	}

	triggerRes := new(TriggerAuditPlanRes)
	err = json.Unmarshal(resBody, triggerRes)
	if err != nil {
		return "", err
	}
	if triggerRes.Code != 0 {
		return "", fmt.Errorf("failed to request %s, error:%s", url, triggerRes.Message)
	}
	return triggerRes.Data.Id, nil
}

func (sc *Client) DirectAudit(ctx context.Context, sqlAuditReq *CreateSqlAuditReq) error {
	createSqlAuditResp, err := sc.CreateSqlAudit(ctx, sqlAuditReq)
	if err != nil {
		return fmt.Errorf("failed to create sql audit record, error: %s", err)
	}

	sqlAudit, err := sc.GetSqlAudit(ctx, createSqlAuditResp.Data.Id)
	if err != nil {
		return fmt.Errorf("failed to get sql audit record, error: %s", err)
	}

	projectID, err := sc.GetProjectUidByName(ctx)
	if err != nil {
		return fmt.Errorf("failed to get project uid by name, error: %s", err)
	}

	err = sc.GetTaskSQLs(ctx, sqlAudit.Data.Task.Id, sqlAudit.Data.SQLAuditRecordId, projectID)
	if err != nil {
		return err
	}

	return nil
}

func (sc *Client) GetProjectUidByName(ctx context.Context) (string, error) {
	projectList, err := sc.GetProjectList(ctx)
	if err != nil {
		return "", err
	}

	var projectID string
	for _, project := range projectList.Data {
		if project.Name == sc.project {
			projectID = project.ProjectUid
			break
		}
	}

	return projectID, nil
}

func (sc *Client) CreateSqlAudit(ctx context.Context, sqlAuditReq *CreateSqlAuditReq) (*CreateSqlAuditResp, error) {
	url := sc.baseURL + fmt.Sprintf(CreateSqlAudit, sc.project)
	bodyBuf := &bytes.Buffer{}
	encoder := json.NewEncoder(bodyBuf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(&sqlAuditReq)
	if err != nil {
		return nil, err
	}

	resBody, err := sc.httpClient.sendRequest(ctx, url, http.MethodPost, sc.token, bytes.NewBuffer(bodyBuf.Bytes()))
	if err != nil {
		return nil, err
	}

	resp := new(CreateSqlAuditResp)
	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("failed to request %s, error:%s", url, resp.Message)
	}

	return resp, nil
}

func (sc *Client) GetSqlAudit(ctx context.Context, sqlAuditRecordID string) (*GetSqlAuditResp, error) {
	url := sc.baseURL + fmt.Sprintf(GetSqlAudit, sc.project, sqlAuditRecordID)
	resBody, err := sc.httpClient.sendRequest(ctx, url, http.MethodGet, sc.token, nil)
	if err != nil {
		return nil, err
	}

	resp := new(GetSqlAuditResp)
	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("failed to request %s message: %s", url, resp.Message)
	}

	return resp, nil
}

func (sc *Client) GetTaskSQLs(ctx context.Context, taskID uint, sqlAuditRecordID, projectID string) error {
	var pageIndex, pageSize, cursor uint64
	pageIndex, pageSize = 1, 10
	cursor = pageIndex * pageSize
	var finalErr error
	auditError := errors.New("audit result error")
	var totalCount, errorCount, warningCount int

	for {
		url := sc.baseURL + fmt.Sprintf(GetTaskSQLs, taskID, pageIndex, pageSize)
		resp, err := sc.httpClient.sendRequest(ctx, url, http.MethodGet, sc.token, nil)
		if err != nil {
			return err
		}

		taskSqlList := new(GetAuditTaskSqls)
		err = json.Unmarshal(resp, taskSqlList)
		if err != nil {
			return err
		}
		if taskSqlList.Code != 0 {
			return fmt.Errorf("failed to request %s,message: %s", url, taskSqlList.Message)
		}
		fmt.Println("---------------------------------------------------------")
		for _, sql := range taskSqlList.Data {
			totalCount++
			fmt.Println(sql.ExecSQL)
			for _, result := range sql.AuditResult {
				fmt.Printf("[%s]%s\n", result.Level, result.Message)
			}
			fmt.Println("---------------------------------------------------------")

			if sql.AuditLevel == "error" {
				errorCount++
				finalErr = auditError
			}

			if sql.AuditLevel == "warn" {
				warningCount++
			}
		}

		if cursor < taskSqlList.TotalNums {
			pageIndex++
			cursor = pageIndex * pageSize
		} else {
			break
		}
	}

	fmt.Printf("total sqls: %d, error sqls: %d, warning sqls: %d, if you want to view the details, visit the link %s\n",
		totalCount, errorCount, warningCount, sc.baseURL+fmt.Sprintf("/sqle/project/%s/sql-audit/detail/%s", projectID, sqlAuditRecordID))

	return finalErr
}

func (sc *Client) GetAuditReportReq(auditPlanName string, reportID string) error {
	var pageIndex, pageSize, cursor uint64
	pageIndex, pageSize = 1, 10
	cursor = pageIndex * pageSize
	var finalErr error
	auditError := errors.New("audit result error")

	for {
		url := sc.baseURL + fmt.Sprintf(GetAuditReport, sc.project, auditPlanName, reportID, pageIndex, pageSize)
		resBody, err := sc.httpClient.sendRequest(context.TODO(), url, http.MethodGet, sc.token, nil)
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
			fmt.Println(res.SQL)
			fmt.Println(res.AuditResult)
			if strings.Contains(res.AuditResult, "[error]") {
				finalErr = auditError
			}
		}

		if cursor < auditRes.TotalNums {
			pageIndex++
			cursor = pageIndex * pageSize
		} else {
			break
		}
	}

	return finalErr
}

type GetProjectListResp struct {
	controller.BaseRes
	Data  []ListProject `json:"data"`
	Total int64         `json:"total_nums"`
}

// A dms Project
type ListProject struct {
	// Project uid
	ProjectUid string `json:"uid"`
	// Project name
	Name string `json:"name"`
}

func (sc *Client) GetProjectList(ctx context.Context) (*GetProjectListResp, error) {
	url := sc.baseURL + GetAllProjects

	resBody, err := sc.httpClient.sendRequest(ctx, fmt.Sprintf(url, 1, 999999), http.MethodGet, sc.token, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to request %s, error: %s", url, err)
	}

	resp := new(GetProjectListResp)
	err = json.Unmarshal(resBody, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body, error: %s", err)
	}

	if resp.Code != 0 {
		return nil, fmt.Errorf("failed to request %s message: %s", url, resp.Message)
	}

	return resp, nil
}

const (
	DefaultTimeoutNum = 10
	DefaultTimeout    = time.Second * time.Duration(DefaultTimeoutNum)
)

// client is a wrap of http.Client
type client struct {
	*http.Client
}

// newClient returns a new HTTP client with timeout and HTTPS support
func newClient(timeout time.Duration, tlsCfg *tls.Config) *client {
	if timeout < time.Second {
		timeout = DefaultTimeout
	}
	tp := &http.Transport{
		TLSClientConfig: tlsCfg,
		Dial:            (&net.Dialer{Timeout: 3 * time.Second}).Dial,
	}
	return &client{&http.Client{
		Timeout:   timeout,
		Transport: tp,
	}}
}

func (c *client) sendRequest(ctx context.Context, url, method, token string, body io.Reader) ([]byte, error) {
	defer c.CloseIdleConnections()
	switch method {
	case http.MethodGet:
		return c.get(ctx, url, token)
	case http.MethodPost:
		return c.post(ctx, url, token, body)
	default:
		return nil, fmt.Errorf("invalid request method")
	}
}

// get fetch a URL with GET method and returns the response
func (c *client) get(ctx context.Context, url, token string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", token)

	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return checkHTTPResponse(res)
}

// post send a POST request to the url and returns the response
func (c *client) post(ctx context.Context, url, token string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return checkHTTPResponse(res)
}

// checkHTTPResponse checks if an HTTP response is with normal status codes
func checkHTTPResponse(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return body, fmt.Errorf("error requesting %s, response: %s, code %d", res.Request.URL, string(body), res.StatusCode)
	}
	return body, nil
}
