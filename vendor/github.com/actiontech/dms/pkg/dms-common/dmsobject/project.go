package dmsobject

import (
	"context"
	"fmt"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

func ListProjects(ctx context.Context, dmsAddr string, req dmsV1.ListProjectReq) (ret []*dmsV1.ListProject, total int64, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.ListProjectReply{}

	url := fmt.Sprintf("%v%v?page_size=%v&page_index=%v&filter_by_name=%v&filter_by_uid=%v", dmsAddr, dmsV1.GetProjectsRouter(), req.PageSize, req.PageIndex, req.FilterByName, req.FilterByUID)

	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to list project from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil
}
