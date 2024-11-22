package dmsobject

import (
	"context"
	"fmt"
	"net/url"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

func ListDbServices(ctx context.Context, dmsAddr string, req dmsV1.ListDBServiceReq) ([]*dmsV1.ListDBService, int64, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	// 构建基础 URL
	baseURL, err := url.Parse(fmt.Sprintf("%s%s", dmsAddr, dmsV1.GetDBServiceRouter(req.ProjectUid)))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// 构建查询参数
	query := url.Values{}
	query.Set("page_size", fmt.Sprintf("%d", req.PageSize))
	query.Set("page_index", fmt.Sprintf("%d", req.PageIndex))

	if req.OrderBy != "" {
		query.Set("order_by", fmt.Sprintf("%v", req.OrderBy))
	}
	if req.FilterByBusiness != "" {
		query.Set("filter_by_business", req.FilterByBusiness)
	}
	if req.FilterByHost != "" {
		query.Set("filter_by_host", req.FilterByHost)
	}
	if req.FilterByUID != "" {
		query.Set("filter_by_uid", req.FilterByUID)
	}
	if req.FilterByName != "" {
		query.Set("filter_by_name", req.FilterByName)
	}
	if req.FilterByPort != "" {
		query.Set("filter_by_port", req.FilterByPort)
	}
	if req.FilterByDBType != "" {
		query.Set("filter_by_db_type", req.FilterByDBType)
	}
	if req.FuzzyKeyword != "" {
		query.Set("fuzzy_keyword", req.FuzzyKeyword)
	}
	if req.IsEnableMasking != nil {
		query.Set("is_enable_masking", fmt.Sprintf("%t", *req.IsEnableMasking))
	}

	for _, id := range req.FilterByDBServiceIds {
		query.Add("filter_by_db_service_ids", id)
	}

	// 将查询参数附加到 URL
	baseURL.RawQuery = query.Encode()

	// 调用 HTTP GET 请求
	reply := &dmsV1.ListDBServiceReply{}
	if err := pkgHttp.Get(ctx, baseURL.String(), header, nil, reply); err != nil {
		return nil, 0, err
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil
}
