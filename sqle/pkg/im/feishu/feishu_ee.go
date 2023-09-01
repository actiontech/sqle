package feishu

import (
	"context"
	"fmt"

	"github.com/larksuite/oapi-sdk-go/v3/service/approval/v4"
)

const (
	FormContent = `
[
  {
    "id": "1",
    "type": "input",
    "name": "@i18n@text1"
  },
  {
    "id": "2",
    "type": "input",
    "name": "@i18n@text2"
  },
  {
    "id": "3",
    "type": "input",
    "name": "@i18n@text3"
  },
  {
    "id": "4",
    "type": "input",
    "name": "@i18n@text4"
  },
  {
    "id": "5",
    "type": "fieldList",
    "name": "@i18n@text5",
    "value": [
      {
        "id": "6",
        "name": "@i18n@text6",
        "type": "input",
        "required": true
      },
      {
        "id": "7",
        "name": "@i18n@text7",
        "type": "input",
        "required": true
      },
      {
        "id": "8",
        "name": "@i18n@text8",
        "type": "input",
        "required": true
      }
    ],
    "option": {
      "inputType": "FORM",
      "printType": "FORM"
    }
  }
]
`

	CreateInstanceForm = `
[
  {
    "id": "1",
    "type": "input",
    "value": "%s"
  },
  {
    "id": "2",
    "type": "input",
    "value": "%s"
  },
  {
    "id": "3",
    "type": "input",
    "value": "%s"
  },
  {
    "id": "4",
    "type": "input",
    "value": "%s"
  },
  {
    "id": "5",
    "type": "fieldList",
    "value": [%s]
  }
]
`
)

// CreateApprovalTemplate 创建审批定义
// https://open.feishu.cn/document/server-docs/approval-v4/approval/create
func (f *FeishuClient) CreateApprovalTemplate(ctx context.Context) (approvalCode *string, err error) {
	req := larkapproval.NewCreateApprovalReqBuilder().
		ApprovalCreate(larkapproval.NewApprovalCreateBuilder().
			ApprovalName(`@i18n@approval_name`).
			ApprovalCode(``).
			Viewers([]*larkapproval.ApprovalCreateViewers{
				larkapproval.NewApprovalCreateViewersBuilder().
					ViewerType(`NONE`).
					Build(),
			}).
			Form(larkapproval.NewApprovalFormBuilder().
				FormContent(FormContent).
				Build()).
			NodeList([]*larkapproval.ApprovalNode{
				larkapproval.NewApprovalNodeBuilder().
					Id(`START`).Build(),
				larkapproval.NewApprovalNodeBuilder().
					Id(`approve`).
					Name(`@i18n@node_name`).
					NodeType(`OR`).
					Approver([]*larkapproval.ApprovalApproverCcer{
						larkapproval.NewApprovalApproverCcerBuilder().
							Type(`Free`).Build()}).Build(),
				larkapproval.NewApprovalNodeBuilder().
					Id(`END`).
					Build(),
			}).
			I18nResources([]*larkapproval.I18nResource{
				larkapproval.NewI18nResourceBuilder().
					IsDefault(true).
					Locale(`zh-CN`).
					Texts([]*larkapproval.I18nResourceText{
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@approval_name`).
							Value(`sqle审批`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@node_name`).
							Value(`Approval`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text1`).
							Value(`项目名称`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text2`).
							Value(`工单名称`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text3`).
							Value(`工单描述`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text4`).
							Value(`工单链接`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text5`).
							Value(`审核结果`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text6`).
							Value(`数据源`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text7`).
							Value(`审核得分`).
							Build(),
						larkapproval.NewI18nResourceTextBuilder().
							Key(`@i18n@text8`).
							Value(`审核通过率`).
							Build(),
					}).Build(),
			}).Build()).
		Build()

	resp, err := f.client.Approval.Approval.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("create approval instance failed: respCode=%v, respMsg=%v, respRequestId=%v", resp.Code, resp.Msg, resp.RequestId())
	}

	return resp.Data.ApprovalCode, nil
}

// CreateApprovalInstance 创建审批实例
// https://open.feishu.cn/document/server-docs/approval-v4/instance/create?appId=cli_a4668286c92ed013
func (f *FeishuClient) CreateApprovalInstance(ctx context.Context, approvalCode, workflowName string, originUserId string, approveUserIds []string, auditResult, projectName, desc, workflowUrl string) (*string, error) {
	form := fmt.Sprintf(CreateInstanceForm, projectName, workflowName, desc, workflowUrl, auditResult)

	req := larkapproval.NewCreateInstanceReqBuilder().
		InstanceCreate(larkapproval.NewInstanceCreateBuilder().
			ApprovalCode(approvalCode).
			OpenId(originUserId).
			Form(form).
			NodeApproverOpenIdList([]*larkapproval.NodeApprover{
				larkapproval.NewNodeApproverBuilder().
					Key(`approve`).
					Value(approveUserIds).
					Build(),
			}).Build()).
		Build()

	resp, err := f.client.Approval.Instance.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf("create approval instance failed: respCode=%v, respMsg=%v, respRequestId=%v", resp.Code, resp.Msg, resp.RequestId())
	}

	return resp.Data.InstanceCode, nil
}
