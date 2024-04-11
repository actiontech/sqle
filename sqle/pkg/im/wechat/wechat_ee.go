//go:build enterprise
// +build enterprise

package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/xen0n/go-workwx"
)

const (
	getTokenUrl       = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?"
	createTemplateUrl = "https://qyapi.weixin.qq.com/cgi-bin/oa/approval/create_template?"

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

	resp, err := http.Get(apiURL)
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
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
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
