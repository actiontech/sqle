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

	placeholder := "%s%s?page_size=%d&page_index=%d&order_by=%v&filter_by_business=%s&filter_by_host=%s&filter_by_uid=%s&filter_by_port=%s&filter_by_db_type=%s&filter_by_name=%s"
	requestUri := fmt.Sprintf(placeholder, dmsAddr, dmsV1.GetDBServiceRouter(req.ProjectUid), req.PageSize, req.PageIndex, req.OrderBy, url.QueryEscape(req.FilterByBusiness), req.FilterByHost, req.FilterByUID, req.FilterByPort, url.QueryEscape(req.FilterByDBType), url.QueryEscape(req.FilterByName))

	reply := &dmsV1.ListDBServiceReply{}

	if err := pkgHttp.Get(ctx, requestUri, header, nil, reply); err != nil {
		return nil, 0, err
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil
}
