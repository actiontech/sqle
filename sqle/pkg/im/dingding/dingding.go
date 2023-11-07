package dingding

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dingTalkOauth "github.com/alibabacloud-go/dingtalk/oauth2_1_0"
	dingTalkWorkflow "github.com/alibabacloud-go/dingtalk/workflow_1_0"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

const (
	dingTalkOpenApi = "https://oapi.dingtalk.com/topapi"
	timeout         = 30 * time.Second

	projectNameComp           = "项目名称"
	workflowLinkComp          = "工单链接"
	workflowNameComp          = "工单名称"
	workDescComp              = "工单描述"
	auditResultComp           = "审核结果"
	approvalTemplateModelName = "sqle审批"
)

type DingTalk struct {
	Id          uint
	AppKey      string
	AppSecret   string
	ProcessCode string
}

func GetToken(key, secret string) (string, error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	client, err := dingTalkOauth.NewClient(config)
	if err != nil {
		return "", fmt.Errorf("get dingtalk client error: %v", err)
	}

	getAccessTokenRequest := &dingTalkOauth.GetAccessTokenRequest{
		AppKey:    tea.String(key),
		AppSecret: tea.String(secret),
	}

	result, err := client.GetAccessToken(getAccessTokenRequest)
	if err != nil {
		return "", fmt.Errorf("get dingtalk token error: %v", err)
	}

	return *result.Body.AccessToken, nil
}

// CreateApprovalTemplate
// https://open.dingtalk.com/document/orgapp-server/create-an-approval-form-template
func (d *DingTalk) CreateApprovalTemplate() error {
	var processCode string

	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return fmt.Errorf("get token error: %v", err)
	}

	client, err := newWorkflowClient()
	if err != nil {
		return fmt.Errorf("get workflow client error: %v", err)
	}

	formCreateHeaders := &dingTalkWorkflow.FormCreateHeaders{}
	formCreateHeaders.XAcsDingtalkAccessToken = tea.String(token)

	projectNameComponent := &dingTalkWorkflow.FormComponent{
		ComponentType: tea.String("TextField"),
		Props: &dingTalkWorkflow.FormComponentProps{
			ComponentId: tea.String("TextField_17EZKEGEIBOC0"),
			Label:       tea.String(projectNameComp),
		},
	}

	workflowLinkComponent := &dingTalkWorkflow.FormComponent{
		ComponentType: tea.String("TextField"),
		Props: &dingTalkWorkflow.FormComponentProps{
			ComponentId: tea.String("TextField_17EZKEGE908C0"),
			Label:       tea.String(workflowLinkComp),
		},
	}

	workflowNameComponent := &dingTalkWorkflow.FormComponent{
		ComponentType: tea.String("TextField"),
		Props: &dingTalkWorkflow.FormComponentProps{
			ComponentId: tea.String("TextField_17EZKEGSOCTC0"),
			Label:       tea.String(workflowNameComp),
		},
	}

	workflowDescComponent := &dingTalkWorkflow.FormComponent{
		ComponentType: tea.String("TextField"),
		Props: &dingTalkWorkflow.FormComponentProps{
			ComponentId: tea.String("TextField_17wdcEGSOCTC0"),
			Label:       tea.String(workDescComp),
		},
	}

	tableComponent := &dingTalkWorkflow.FormComponent{
		ComponentType: tea.String("TableField"),
		Children: []*dingTalkWorkflow.FormComponent{
			{
				ComponentType: tea.String("TextField"),
				Props: &dingTalkWorkflow.FormComponentProps{
					Label:       tea.String("数据源"),
					ComponentId: tea.String("TextField_1712EGE908C0"),
				},
			}, {
				ComponentType: tea.String("TextField"),
				Props: &dingTalkWorkflow.FormComponentProps{
					Label:       tea.String("审核得分"),
					ComponentId: tea.String("TextField_17EZ3GSOCTC0"),
				},
			}, {
				ComponentType: tea.String("TextField"),
				Props: &dingTalkWorkflow.FormComponentProps{
					Label:       tea.String("审核通过率"),
					ComponentId: tea.String("TextField_17EZK43OCTC0"),
				},
			},
		},
		Props: &dingTalkWorkflow.FormComponentProps{
			TableViewMode: tea.String("table"),
			VerticalPrint: tea.Bool(true),
			Label:         tea.String(auditResultComp),
			ComponentId:   tea.String("TableField_17EZKEG9OCTC0"),
		},
	}

	formCreateRequest := &dingTalkWorkflow.FormCreateRequest{
		Name:           tea.String(approvalTemplateModelName),
		FormComponents: []*dingTalkWorkflow.FormComponent{projectNameComponent, workflowNameComponent, workflowDescComponent, workflowLinkComponent, tableComponent},
	}

	// 存在 processCode 则更新模版
	if d.ProcessCode != "" {
		formCreateRequest.ProcessCode = tea.String(d.ProcessCode)
	}

	resp, err := client.FormCreateWithOptions(formCreateRequest, formCreateHeaders, &util.RuntimeOptions{})
	if err != nil {
		if sdkErr, ok := err.(*tea.SDKError); ok {
			// 如果用户更换新的审批应用，这时候用数据库中的 ProcessCode 去更新审批模版，会报 processcode.error（processCode对应的审批流程不存在）错误
			// 去掉 ProcessCode ，不去更新，而是在新的审批应用上创建一个新的审批模版
			// 错误码： https://open.dingtalk.com/document/orgapp-server/create-an-approval-form-template
			if !tea.BoolValue(util.Empty(sdkErr.Code)) && sdkErr.Code == tea.String("processcode.error") {
				formCreateRequest.ProcessCode = tea.String("")
				//nolint:staticcheck
				resp, err = client.FormCreateWithOptions(formCreateRequest, formCreateHeaders, &util.RuntimeOptions{})
				if err != nil {
					return fmt.Errorf("second attempt create approval template error: %v", err)
				}
				if resp.Body.Result.ProcessCode != nil {
					processCode = *resp.Body.Result.ProcessCode
				}
				goto End
				// https://github.com/actiontech/sqle/issues/1487
			} else if !tea.BoolValue(util.Empty(sdkErr.Code)) && *sdkErr.Code == "formName.error" {
				getProcessCodeByNameHeaders := &dingTalkWorkflow.GetProcessCodeByNameHeaders{}
				getProcessCodeByNameHeaders.XAcsDingtalkAccessToken = tea.String(token)

				getProcessCodeByNameRequest := &dingTalkWorkflow.GetProcessCodeByNameRequest{
					Name: tea.String(approvalTemplateModelName),
				}

				getProcessCodeResp, err := client.GetProcessCodeByNameWithOptions(getProcessCodeByNameRequest, getProcessCodeByNameHeaders, &util.RuntimeOptions{})
				if err != nil {
					return fmt.Errorf("get Process Code error: %v", err)
				}
				if getProcessCodeResp.Body.Result.ProcessCode != nil {
					processCode = *getProcessCodeResp.Body.Result.ProcessCode
				}

				goto End
			}
		}

		return fmt.Errorf("create approval template error: %v", err)
	} else {
		if resp.Body.Result.ProcessCode != nil {
			processCode = *resp.Body.Result.ProcessCode
		}
	}

End:
	if processCode == "" {
		return fmt.Errorf("create approval template error: %v", resp.Body.Result)
	}

	s := model.GetStorage()
	if err := s.UpdateImConfigById(d.Id, map[string]interface{}{"process_code": processCode}); err != nil {
		return fmt.Errorf("update process code error: %v", err)
	}

	return nil
}

// CreateApprovalInstance
// https://open.dingtalk.com/document/orgapp-server/create-an-approval-instance
func (d *DingTalk) CreateApprovalInstance(workflowName string, workflowId string, originUserId *string, userIds []*string, auditResult, projectName, desc, workflowUrl string) error {
	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return fmt.Errorf("get token error: %v", err)
	}

	client, err := newWorkflowClient()
	if err != nil {
		return fmt.Errorf("get workflow client error: %v", err)
	}

	startProcessInstanceHeaders := &dingTalkWorkflow.StartProcessInstanceHeaders{}
	startProcessInstanceHeaders.XAcsDingtalkAccessToken = tea.String(token)

	actionType := "NONE"
	if len(userIds) > 1 {
		actionType = "OR"
	}

	var startProcessInstanceRequestApprovers []*dingTalkWorkflow.StartProcessInstanceRequestApprovers
	startProcessInstanceRequestApprovers = append(startProcessInstanceRequestApprovers, &dingTalkWorkflow.StartProcessInstanceRequestApprovers{
		ActionType: tea.String(actionType),
		UserIds:    userIds,
	})

	var startProcessInstanceRequestFormComponentValues []*dingTalkWorkflow.StartProcessInstanceRequestFormComponentValues
	startProcessInstanceRequestFormComponentValues = append(startProcessInstanceRequestFormComponentValues,
		&dingTalkWorkflow.StartProcessInstanceRequestFormComponentValues{
			Name:  tea.String(projectNameComp),
			Value: tea.String(projectName)},
		&dingTalkWorkflow.StartProcessInstanceRequestFormComponentValues{
			Name:  tea.String(workflowNameComp),
			Value: tea.String(workflowName)},
		&dingTalkWorkflow.StartProcessInstanceRequestFormComponentValues{
			Name:  tea.String(workDescComp),
			Value: tea.String(desc)},
		&dingTalkWorkflow.StartProcessInstanceRequestFormComponentValues{
			Name:  tea.String(workflowLinkComp),
			Value: tea.String(workflowUrl)},
		&dingTalkWorkflow.StartProcessInstanceRequestFormComponentValues{
			Name:  tea.String(auditResultComp),
			Value: tea.String(auditResult),
		})

	startProcessInstanceRequest := &dingTalkWorkflow.StartProcessInstanceRequest{
		OriginatorUserId:    originUserId,
		ProcessCode:         tea.String(d.ProcessCode),
		Approvers:           startProcessInstanceRequestApprovers,
		FormComponentValues: startProcessInstanceRequestFormComponentValues,
	}

	resp, err := client.StartProcessInstanceWithOptions(startProcessInstanceRequest, startProcessInstanceHeaders, &util.RuntimeOptions{})
	if err != nil {
		return fmt.Errorf("create approval instance error: %v", err)
	}

	approvalDetail, err := d.GetApprovalDetail(*resp.Body.InstanceId)
	if err != nil {
		return fmt.Errorf("get approval detail error: %v", err)
	}

	taskID := *approvalDetail.Tasks[0].TaskId
	dingTalkInstance := model.DingTalkInstance{ApproveInstanceCode: *resp.Body.InstanceId, WorkflowId: workflowId, TaskID: int64(uint(taskID))}
	s := model.GetStorage()
	if err := s.Save(&dingTalkInstance); err != nil {
		return fmt.Errorf("save dingtalk instance error: %v", err)
	}

	return nil
}

// UpdateApprovalStatus
// https://open.dingtalk.com/document/orgapp-server/approve-or-reject-the-approval-task
func (d *DingTalk) UpdateApprovalStatus(workflowId string, status, userId, reason string) error {
	s := model.GetStorage()
	dingTalkInstance, exist, err := s.GetDingTalkInstanceByWorkflowID(workflowId)
	if err != nil {
		return fmt.Errorf("get dingtalk instance error: %v", err)
	}
	if !exist {
		return errors.New("ding talk instance not exist")
	}

	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return fmt.Errorf("get token error: %v", err)
	}

	client, err := newWorkflowClient()
	if err != nil {
		return fmt.Errorf("get workflow client error: %v", err)
	}

	executeProcessInstanceHeaders := &dingTalkWorkflow.ExecuteProcessInstanceHeaders{}
	executeProcessInstanceHeaders.XAcsDingtalkAccessToken = tea.String(token)

	executeProcessInstanceRequest := &dingTalkWorkflow.ExecuteProcessInstanceRequest{
		ProcessInstanceId: tea.String(dingTalkInstance.ApproveInstanceCode),
		Result:            tea.String(status),
		ActionerUserId:    tea.String(userId),
		TaskId:            tea.Int64(dingTalkInstance.TaskID),
	}

	if reason != "" {
		executeProcessInstanceRequest.Remark = tea.String(reason)
	}

	_, err = client.ExecuteProcessInstanceWithOptions(executeProcessInstanceRequest, executeProcessInstanceHeaders, &util.RuntimeOptions{})
	if err != nil {
		return fmt.Errorf("update approval status error: %v", err)
	}

	dingTalkInstance.Status = status
	if err := s.Save(&dingTalkInstance); err != nil {
		return fmt.Errorf("save dingtalk instance error: %v", err)
	}

	return nil
}

// GetApprovalDetail
// https://open.dingtalk.com/document/orgapp-server/obtains-the-details-of-a-single-approval-instance-pop
func (d *DingTalk) GetApprovalDetail(approveInstanceCode string) (*dingTalkWorkflow.GetProcessInstanceResponseBodyResult, error) {
	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return nil, fmt.Errorf("get token error: %v", err)
	}

	client, err := newWorkflowClient()
	if err != nil {
		return nil, fmt.Errorf("get workflow client error: %v", err)
	}

	getProcessInstanceHeaders := &dingTalkWorkflow.GetProcessInstanceHeaders{}
	getProcessInstanceHeaders.XAcsDingtalkAccessToken = tea.String(token)
	getProcessInstanceRequest := &dingTalkWorkflow.GetProcessInstanceRequest{
		ProcessInstanceId: tea.String(approveInstanceCode),
	}

	resp, err := client.GetProcessInstanceWithOptions(getProcessInstanceRequest, getProcessInstanceHeaders, &util.RuntimeOptions{})
	if err != nil {
		return nil, fmt.Errorf("get approval status error: %v", err)
	}

	return resp.Body.Result, nil
}

type GetUserIDByPhoneRep struct {
	Mobile string `json:"mobile"`
}

type GetUserIDByPhoneResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		Userid string `json:"userid"`
	}
}

// GetUserIDByPhone
// https://open.dingtalk.com/document/orgapp-server/query-users-by-phone-number
func (d *DingTalk) GetUserIDByPhone(phone string) (*string, error) {
	// todo : token cache
	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return nil, fmt.Errorf("get token error: %v", err)
	}

	url := fmt.Sprintf("%s/v2/user/getbymobile?access_token=%s", dingTalkOpenApi, token)

	newEntry := log.NewEntry()

	getUserIDByPhoneRep := &GetUserIDByPhoneRep{
		Mobile: phone,
	}

	body, err := json.Marshal(getUserIDByPhoneRep)
	if err != nil {
		return nil, fmt.Errorf("marshal req error: %v", err)
	}

	resp, err := Requester(url, http.MethodPost, token, body)
	if err != nil {
		return nil, fmt.Errorf("get user id by mobile error: %v", err)
	}

	var user GetUserIDByPhoneResp
	if err := json.Unmarshal(resp, &user); err != nil {
		return nil, fmt.Errorf("unmarshal error: %v", err)
	}

	if user.ErrCode != 0 {
		newEntry.Errorf("get user id by mobile error,code: %v errMsg: %v", user.ErrCode, user.ErrMsg)
	}

	return &user.Result.Userid, nil
}

type GetUserByUserIdReq struct {
	UserId string `json:"userid"`
}

type GetUserByUserIdResp struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
	Result  struct {
		Mobile string `json:"mobile"`
	}
}

// GetMobileByUserID
// https://open.dingtalk.com/document/orgapp-server/query-user-details
func (d *DingTalk) GetMobileByUserID(userId string) (string, error) {
	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return "", fmt.Errorf("get token error: %v", err)
	}

	url := fmt.Sprintf("%s/v2/user/get?access_token=%s", dingTalkOpenApi, token)

	body, err := json.Marshal(&GetUserByUserIdReq{
		UserId: userId,
	})
	if err != nil {
		return "", fmt.Errorf("marshal req error: %v", err)
	}

	resp, err := Requester(url, http.MethodPost, token, body)
	if err != nil {
		return "", fmt.Errorf("get user by user id error: %v", err)
	}

	var user GetUserByUserIdResp
	if err := json.Unmarshal(resp, &user); err != nil {
		return "", fmt.Errorf("unmarshal error: %v", err)
	}

	if user.ErrCode != 0 {
		return "", fmt.Errorf("get user by user id error,code: %v errMsg: %v", user.ErrCode, user.ErrMsg)
	}

	return user.Result.Mobile, nil
}

// CancelApprovalInstance
// https://open.dingtalk.com/document/orgapp-server/revoke-an-approval-instance
func (d *DingTalk) CancelApprovalInstance(instanceCode string) error {
	token, err := GetToken(d.AppKey, d.AppSecret)
	if err != nil {
		return fmt.Errorf("get token error: %v", err)
	}

	client, err := newWorkflowClient()
	if err != nil {
		return fmt.Errorf("get dingtalk client error: %v", err)
	}

	terminateProcessInstanceHeaders := &dingTalkWorkflow.TerminateProcessInstanceHeaders{}
	terminateProcessInstanceHeaders.XAcsDingtalkAccessToken = tea.String(token)
	terminateProcessInstanceRequest := &dingTalkWorkflow.TerminateProcessInstanceRequest{
		ProcessInstanceId: tea.String(instanceCode),
		IsSystem:          tea.Bool(true),
		Remark:            tea.String("工单已关闭"),
	}

	_, err = client.TerminateProcessInstanceWithOptions(terminateProcessInstanceRequest, terminateProcessInstanceHeaders, &util.RuntimeOptions{})
	if err != nil {
		return fmt.Errorf("cancel approval instance error: %v, instanceCode: %v", err, instanceCode)
	}

	return nil
}

func Requester(url, method, token string, body []byte) ([]byte, error) {
	reader := bytes.NewReader(body)
	client := &http.Client{Timeout: timeout}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func newWorkflowClient() (*dingTalkWorkflow.Client, error) {
	config := &openapi.Config{}
	config.Protocol = tea.String("https")
	config.RegionId = tea.String("central")
	client, err := dingTalkWorkflow.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("get dingtalk client error: %v", err)
	}
	return client, nil
}
