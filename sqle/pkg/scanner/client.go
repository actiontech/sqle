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
	LoginUri = "/v1/login"
	// Post
	TriggerAudit = "/v1/projects/%v/audit_plans/%s/trigger"
	// Post
	FullUpload = "/v2/projects/%v/audit_plans/%s/sqls/full"
	// Post
	PartialUpload = "/v2/projects/%v/audit_plans/%s/sqls/partial"
	// Get										%v=report_id
	GetAuditReport = "/v1/projects/%v/audit_plans/%s/reports/%v/sqls?page_index=%d&page_size=%d"
)

type (
	BaseRes                     = controller.BaseRes
	GetAuditPlanReportSQLsRes   = v1.GetAuditPlanReportSQLsResV1
	AuditPlanSQLReq             = v2.AuditPlanSQLReqV2
	FullSyncAuditPlanSQLsReq    = v2.FullSyncAuditPlanSQLsReqV2
	PartialSyncAuditPlanSQLsReq = v1.PartialSyncAuditPlanSQLsReqV1
	TriggerAuditPlanRes         = v1.TriggerAuditPlanResV1
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

func (sc *Client) UploadReq(uri string, auditPlanName string, sqlList []*AuditPlanSQLReq) error {
	bodyBuf := &bytes.Buffer{}
	encoder := json.NewEncoder(bodyBuf)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(&FullSyncAuditPlanSQLsReq{
		SQLs: sqlList,
	})
	if err != nil {
		return err
	}

	url := sc.baseURL + fmt.Sprintf(uri, sc.project, auditPlanName)
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
