//go:build enterprise
// +build enterprise

package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/log"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/xen0n/go-workwx"
)

const (
	getTokenUrl       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?"
	createTemplateUrl = "https://qyapi.weixin.qq.com/cgi-bin/oa/approval/create_template?"
	httpTimeOut       = 60

	templateName      = "SQLE审核"
	language          = "zh_CN"
	textControl       = "Text"
	tableControl      = "Table"
	projectNameComp   = "项目名称"
	workflowNameComp  = "工单名称"
	workflowLinkComp  = "工单链接"
	oaTypeComp        = "审核操作"
	dataSourceComp    = "数据源"
	auditScoreComp    = "审核得分"
	auditPassRateComp = "审核通过率"
	sqlTextComp       = "SQL详情"
	sqlDetailComp     = "审核详情"

	tableCompId         = "Table-1712652040429"
	projectNameCompId   = "Text-1712467954299"
	workflowNameCompId  = "Text-1712467972770"
	workflowLinkCompId  = "Text-1712653256570"
	oaTypeCompId        = "Text-1712653346567"
	dataSourceCompId    = "Text-1712652061059"
	auditScoreCompId    = "Text-1712652144998"
	auditPassRateCompId = "Text-1712652891077"
	sqlTextCompId       = "Text-1712653324978"
)

type wechatClient struct {
	client *workwx.WorkwxApp
}

func NewWechatClient(corpId, secret string) *wechatClient {
	client := workwx.New(corpId)
	// WithApp方法需要传应用密钥以及agentID，但是这里不需要使用agentID，所以使用`1`代替
	wxClient := client.WithApp(secret, 1)
	return &wechatClient{client: wxClient}
}

type accessTokenResponse struct {
	ErrCode     int64  `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
}

func getAccessToekn(path string, CorpID string, CorpSecret string) (string, error) {
	params := url.Values{}
	params.Add("corpid", CorpID)
	params.Add("corpsecret", CorpSecret)

	apiURL := path + params.Encode()

	httpClient := &http.Client{
		Timeout: httpTimeOut * time.Second,
	}

	resp, err := httpClient.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var accessTokenResponse accessTokenResponse
	err = json.Unmarshal(body, &accessTokenResponse)
	if err != nil {
		return "", err
	}
	if accessTokenResponse.ErrCode != 0 {
		return "", errors.New(accessTokenResponse.ErrMsg)
	}

	return accessTokenResponse.AccessToken, nil
}

type WxControl struct {
	Control string
	ID      string
	Title   string
}

var WxControls = []WxControl{
	{
		Control: textControl,
		ID:      projectNameCompId,
		Title:   projectNameComp,
	},
	{
		Control: textControl,
		ID:      workflowNameCompId,
		Title:   workflowNameComp,
	},
	{
		Control: textControl,
		ID:      workflowLinkCompId,
		Title:   workflowLinkComp,
	},
	{
		Control: textControl,
		ID:      oaTypeCompId,
		Title:   oaTypeComp,
	},
}

var TableWxControls = []WxControl{
	{
		Control: textControl,
		ID:      dataSourceCompId,
		Title:   dataSourceComp,
	},
	{
		Control: textControl,
		ID:      auditScoreCompId,
		Title:   auditScoreComp,
	},
	{
		Control: textControl,
		ID:      auditPassRateCompId,
		Title:   auditPassRateComp,
	},
	{
		Control: textControl,
		ID:      sqlTextCompId,
		Title:   sqlTextComp,
	},
}

type OATemplateDetail struct {
	TemplateName []workwx.OAText `json:"template_name"`
	// TemplateContent 模板控件信息
	TemplateContent workwx.OATemplateControls `json:"template_content"`
}

type createTemplateResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	TemplateId string `json:"template_id"`
}

// 市面上的wechat sdk都不支持创建审批模板接口，所以需要手动实现
// https://developer.work.weixin.qq.com/document/path/97437
func (c *wechatClient) CreateApprovalTemplate(ctx context.Context) (approvalCode *string, err error) {
	accessToken, err := getAccessToekn(getTokenUrl, c.client.CorpID, c.client.CorpSecret)
	if err != nil {
		return
	}

	templateDetail := OATemplateDetail{
		TemplateName: []workwx.OAText{
			{
				Text: templateName,
				Lang: language,
			},
		},
		TemplateContent: workwx.OATemplateControls{
			Controls: []workwx.OATemplateControl{},
		},
	}

	for _, control := range WxControls {
		templateDetail.TemplateContent.Controls = append(templateDetail.TemplateContent.Controls, workwx.OATemplateControl{
			Property: workwx.OATemplateControlProperty{
				Control: workwx.OAControl(control.Control),
				ID:      control.ID,
				Title: []workwx.OAText{
					{
						Text: control.Title,
						Lang: language,
					},
				},
			},
		})
	}

	tableControl := workwx.OATemplateControl{
		Property: workwx.OATemplateControlProperty{
			Control: tableControl,
			ID:      tableCompId,
			Title: []workwx.OAText{
				{
					Text: sqlDetailComp,
					Lang: language,
				},
			},
		},
		Config: workwx.OATemplateControlConfig{
			Table: workwx.OATemplateControlConfigTable{
				Children: []workwx.OATemplateControl{},
			},
		},
	}

	for _, control := range TableWxControls {
		tableControl.Config.Table.Children = append(tableControl.Config.Table.Children, workwx.OATemplateControl{
			Property: workwx.OATemplateControlProperty{
				Control: workwx.OAControl(control.Control),
				ID:      control.ID,
				Title: []workwx.OAText{
					{
						Text: control.Title,
						Lang: language,
					},
				},
			},
		})
	}

	templateDetail.TemplateContent.Controls = append(templateDetail.TemplateContent.Controls, tableControl)

	params := url.Values{}
	params.Add("access_token", accessToken)

	// 构造请求地址
	apiURL := createTemplateUrl + params.Encode()
	jsonData, err := json.Marshal(templateDetail)

	if err != nil {
		return nil, err
	}

	httpClient := &http.Client{
		Timeout: httpTimeOut * time.Second,
	}

	resp, err := httpClient.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var createTemplateResp createTemplateResp
	err = json.Unmarshal(body, &createTemplateResp)
	if err != nil {
		return nil, err
	}
	if createTemplateResp.ErrCode != 0 {
		return nil, errors.New(createTemplateResp.ErrMsg)
	}
	return &createTemplateResp.TemplateId, nil
}

// CreateApprovalInstance 创建审批实例
// https://developer.work.weixin.qq.com/document/path/91853
func (c *wechatClient) CreateApprovalInstance(ctx context.Context, approvalCode, workflowName string, originUserId string,
	approveUserIds []string, projectName, workflowUrl, oaType string, auditResults []*model.WorkflowInstanceRecord) (string, error) {

	oaApplyEvent := workwx.OAApplyEvent{
		CreatorUserID: originUserId,
		TemplateID:    approvalCode,
		Approver: []workwx.OAApprover{
			{
				// Attr 节点审批方式：1-或签；2-会签，仅在节点为多人审批时有效
				// sqle审核节点只要有一人同意就可进入下一个流程，所以选择或签
				Attr:   1,
				UserID: approveUserIds,
			},
		},
	}
	oaContents := workwx.OAContents{
		Contents: []workwx.OAContent{
			{
				Control: textControl,
				ID:      projectNameCompId,
				Title:   []workwx.OAText{{Text: projectNameComp, Lang: language}},
				// 该sdk创建json数据的时候使用的是结构体，导致不赋值的结构体都会初始化为零值
				// OAContentBankAccount结构体中AccountType字段不需要进行赋值，会被初始化为0，导致接口报错，必须赋值
				// 临时方案：将AccountType字段值初始化为1，不影响审批功能使用
				Value: workwx.OAContentValue{Text: projectName, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
			},
			{
				Control: textControl,
				ID:      workflowNameCompId,
				Title:   []workwx.OAText{{Text: workflowNameComp, Lang: language}},
				Value:   workwx.OAContentValue{Text: workflowName, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
			},
			{
				Control: textControl,
				ID:      workflowLinkCompId,
				Title:   []workwx.OAText{{Text: workflowLinkComp, Lang: language}},
				Value:   workwx.OAContentValue{Text: workflowUrl, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
			},
			{
				Control: textControl,
				ID:      oaTypeCompId,
				Title:   []workwx.OAText{{Text: oaTypeComp, Lang: language}},
				Value:   workwx.OAContentValue{Text: oaType, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
			},
		},
	}

	s := model.GetStorage()

	tableList := []workwx.OAContentTableList{}
	for _, result := range auditResults {
		task, exist, err := s.GetTaskDetailById(fmt.Sprint(result.TaskId))
		if !exist {
			log.NewEntry().Infof("task has not detail, task id:%v", result.TaskId)
			continue
		}
		if err != nil {
			log.NewEntry().Infof("get task detail failed, err:%v", err)
			continue
		}

		sqls := task.ExecuteSQLs
		if len(sqls) < 1 {
			continue
		}
		sqlContent := sqls[0].Content
		// 避免显示过长，只取前50个字符
		if len(sqlContent) > 50 {
			sqlContent = fmt.Sprintf("%s...", sqlContent[:50])
		}
		// 企微通知不支持换行符
		// https://developer.work.weixin.qq.com/document/path/91853#%E9%99%841-%E6%96%87%E6%9C%AC%E5%A4%9A%E8%A1%8C%E6%96%87%E6%9C%AC%E6%8E%A7%E4%BB%B6%EF%BC%88control%E5%8F%82%E6%95%B0%E4%B8%BAtext%E6%88%96textarea%EF%BC%89
		sqlContent = strings.ReplaceAll(sqlContent, "\n", "")
		sqlContent = strings.ReplaceAll(sqlContent, "\r", "")
		tableList = append(tableList, workwx.OAContentTableList{
			List: []workwx.OAContent{
				{
					Control: textControl,
					ID:      dataSourceCompId,
					Title:   []workwx.OAText{{Text: dataSourceComp, Lang: language}},
					Value:   workwx.OAContentValue{Text: task.Schema, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
				},
				{
					Control: textControl,
					ID:      auditScoreCompId,
					Title:   []workwx.OAText{{Text: auditScoreComp, Lang: language}},
					Value:   workwx.OAContentValue{Text: fmt.Sprintf("%d", task.Score), BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
				},
				{
					Control: textControl,
					ID:      auditPassRateCompId,
					Title:   []workwx.OAText{{Text: auditPassRateComp, Lang: language}},
					Value:   workwx.OAContentValue{Text: fmt.Sprintf("%.2f%%", task.PassRate*100), BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
				},
				{
					Control: textControl,
					ID:      sqlTextCompId,
					Title:   []workwx.OAText{{Text: sqlTextComp, Lang: language}},
					Value:   workwx.OAContentValue{Text: sqlContent, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
				},
			},
		})
	}

	// sql详情组件
	if len(tableList) > 0 {
		oaContents.Contents = append(oaContents.Contents, workwx.OAContent{
			Control: tableControl,
			ID:      tableCompId,
			Title:   []workwx.OAText{{Text: sqlDetailComp, Lang: language}},
			Value:   workwx.OAContentValue{Table: tableList, BankAccount: workwx.OAContentBankAccount{AccountType: 1}},
		})
	}

	oaApplyEvent.ApplyData = oaContents

	spNo, err := c.client.ApplyOAEvent(oaApplyEvent)
	return spNo, err
}

// GetApprovalInstDetail 获取审批实例详情
// https://developer.work.weixin.qq.com/document/path/91983
func (c *wechatClient) GetApprovalRecordDetail(ctx context.Context, spNo string) (*workwx.OAApprovalDetail, error) {
	detail, err := c.client.GetOAApprovalDetail(spNo)
	if err != nil {
		return nil, err
	}
	return detail, nil
}
