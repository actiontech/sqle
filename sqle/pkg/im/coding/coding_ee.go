//go:build enterprise
// +build enterprise

package coding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/actiontech/sqle/sqle/log"
	"io"
	"net/http"
	"time"
)

type CodingClient struct {
	url   string
	token string
}

func NewCodingClient(url string, token string) *CodingClient {
	return &CodingClient{url: url, token: token}
}

type CreateIssueRequestBody struct {
	Name        string `json:"Name"`
	Priority    string `json:"Priority"`
	ProjectName string `json:"ProjectName"`
	Type        string `json:"Type"`
	Description string `json:"Description"`
	// coding必填项，默认8小时
	WorkingHours float32 `json:"WorkingHours"`
}

// CreateIssue https://coding.net/help/openapi#/operations/CreateIssue
func (codingClient *CodingClient) CreateIssue(createIssueRequestBody CreateIssueRequestBody) (*CreateIssueResponseBody, error) {
	jsonData, err := json.Marshal(createIssueRequestBody)
	if err != nil {
		return nil, err
	}
	createIssueUrl := fmt.Sprintf("%s/open-api/CreateIssue?Action=CreateIssue", codingClient.url)
	request, err := http.NewRequest("POST", createIssueUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %v", codingClient.token))
	client := &http.Client{
		Timeout: 60 * time.Second,
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Logger().Errorf("failed to invoke coding create issue api error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Logger().Errorf("failed to read coding create issue api response: %v", err)
		return nil, err
	}
	var createIssueResponse CreateIssueResponseBody
	err = json.Unmarshal(body, &createIssueResponse)
	if err != nil || createIssueResponse.Response == nil {
		return nil, fmt.Errorf("failed to create coding issue, maybe blocked by coding")
	}
	if createIssueResponse.Response.Error != nil {
		return nil, errors.New(createIssueResponse.Response.Error.Message)
	}
	return &createIssueResponse, nil
}

type CreateIssueResponseBody struct {
	Response *CreateIssueResp `json:"Response"`
}

type CreateIssueResp struct {
	Error     *CreateIssueError `json:"Error"`
	Issue     *interface{}      `json:"Issue"`
	RequestId string            `json:"RequestId"`
}

type CreateIssueError struct {
	Message string `json:"Message"`
	Code    string `json:"Code"`
}
