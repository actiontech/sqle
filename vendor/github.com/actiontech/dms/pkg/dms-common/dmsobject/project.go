package dmsobject

import (
	"context"
	"fmt"
	"net/url"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	dmsV2 "github.com/actiontech/dms/pkg/dms-common/api/dms/v2"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

func ListProjects(ctx context.Context, dmsAddr string, req dmsV1.ListProjectReq) (ret []*dmsV1.ListProject, total int64, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.ListProjectReply{}

	baseURL, err := url.Parse(fmt.Sprintf("%v%v", dmsAddr, dmsV2.GetProjectsRouter()))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse base URL: %v", err)
	}
	// 构建查询参数
	query := url.Values{}
	query.Set("page_size", fmt.Sprintf("%v", req.PageSize))
	query.Set("page_index", fmt.Sprintf("%v", req.PageIndex))

	if req.FilterByName != "" {
		query.Set("filter_by_name", req.FilterByName)
	}
	if req.FilterByUID != "" {
		query.Set("filter_by_uid", req.FilterByUID)
	}
	if req.FilterByProjectPriority != "" {
		query.Set("filter_by_project_priority", string(req.FilterByProjectPriority))
	}
	for _, projectUid := range req.FilterByProjectUids {
		query.Add("filter_by_project_uids", projectUid)
	}

	baseURL.RawQuery = query.Encode()

	if err := pkgHttp.Get(ctx, baseURL.String(), header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to list project from %v: %v", baseURL.String(), err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil
}
