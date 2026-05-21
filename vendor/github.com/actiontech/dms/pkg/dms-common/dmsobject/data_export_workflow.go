package dmsobject

import (
	"context"
	"fmt"
	"net/url"

	dmsCommonV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

func GetGlobalDataExportWorkflowsList(ctx context.Context, dmsAddr string, req dmsCommonV1.FilterGlobalDataExportWorkflowReq) ([]*dmsCommonV1.ListDataExportWorkflow, int64, error) {

	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	// 构建基础 URL
	baseURL, err := url.Parse(fmt.Sprintf("%s%s", dmsAddr, dmsCommonV1.GetGlobalDataExportWorkflowsRouter()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// 构建查询参数
	query := url.Values{}
	query.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	query.Set("page_index", fmt.Sprintf("%d", req.PageIndex))

	if req.FilterByCreateUserUid != "" {
		query.Set("filter_by_create_user_uid", req.FilterByCreateUserUid)
	}

	if len(req.FilterStatusList) > 0 {
		for _, v := range req.FilterStatusList {
			if v != "" {
				query.Add("filter_status_list", string(v))
			}
		}
	}

	if len(req.FilterProjectUids) > 0 {
		for _, v := range req.FilterProjectUids {
			if v != "" {
				query.Add("filter_project_uids", v)
			}
		}
	}

	if req.FilterProjectUid != "" {
		query.Set("filter_project_uid", req.FilterProjectUid)
	}

	if req.FilterDBServiceUid != "" {
		query.Set("filter_db_service_uid", req.FilterDBServiceUid)
	}

	if req.FilterCurrentStepAssigneeUserId != "" {
		query.Set("filter_current_step_assignee_user_id", req.FilterCurrentStepAssigneeUserId)
	}

	if req.FuzzyKeyword != "" {
		query.Set("fuzzy_keyword", req.FuzzyKeyword)
	}

	if req.CheckUserCanAccess {
		query.Set("check_user_can_access", "true")
	}

	if req.CurrentUserID != "" {
		query.Set("current_user_id", req.CurrentUserID)
	}

	if len(req.ViewableDBServiceUids) > 0 {
		for _, v := range req.ViewableDBServiceUids {
			if v != "" {
				query.Add("viewable_db_service_uids", v)
			}
		}
	}

	if req.FilterCreateTimeFrom != "" {
		query.Set("filter_create_time_from", req.FilterCreateTimeFrom)
	}

	if req.FilterCreateTimeTo != "" {
		query.Set("filter_create_time_to", req.FilterCreateTimeTo)
	}

	if req.FilterUpdateTimeFrom != "" {
		query.Set("filter_update_time_from", req.FilterUpdateTimeFrom)
	}

	if req.FilterUpdateTimeTo != "" {
		query.Set("filter_update_time_to", req.FilterUpdateTimeTo)
	}

	// 将查询参数附加到 URL
	baseURL.RawQuery = query.Encode()

	// 调用 HTTP GET 请求
	reply := &dmsCommonV1.ListDataExportWorkflowsReply{}
	if err := pkgHttp.Get(ctx, baseURL.String(), header, nil, reply); err != nil {
		return nil, 0, err
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil
}
