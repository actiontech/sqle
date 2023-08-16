package dmsobject

import (
	"context"
	"fmt"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

func GetUser(ctx context.Context, userUid string, dmsAddr string) (*dmsV1.GetUser, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.GetUserReply{}

	url := fmt.Sprintf("%v%v", dmsAddr, dmsV1.GetUserRouter(userUid))

	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, fmt.Errorf("failed to get user from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Msg)
	}

	return reply.Payload.User, nil
}

func GetUserOpPermission(ctx context.Context, namespaceUid, userUid, dmsAddr string) (ret []dmsV1.OpPermissionItem, isAdmin bool, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reqBody := struct {
		UserOpPermission *dmsV1.UserOpPermission `json:"user_op_permission"`
	}{
		UserOpPermission: &dmsV1.UserOpPermission{NamespaceUid: namespaceUid},
	}

	reply := &dmsV1.GetUserOpPermissionReply{}

	url := fmt.Sprintf("%v%v", dmsAddr, dmsV1.GetUserOpPermissionRouter(userUid))

	if err := pkgHttp.Get(ctx, url, header, reqBody, reply); err != nil {
		return nil, false, fmt.Errorf("failed to get user op permission from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, false, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Msg)
	}

	return reply.Payload.OpPermissionList, reply.Payload.IsAdmin, nil

}

func ListMembersInNamespace(ctx context.Context, dmsAddr string, req dmsV1.ListMembersForInternalReq) ([]*dmsV1.ListMembersForInternalItem, int64, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.ListMembersForInternalReply{}

	url := fmt.Sprintf("%v%v?page_size=%v&page_index=%v&namespace_uid=%v",
		dmsAddr, dmsV1.GetListMembersForInternalRouter(), req.PageSize, req.PageIndex, req.NamespaceUid)

	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to get member from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Msg)
	}

	return reply.Payload.Members, reply.Payload.Total, nil
}

func ListUsers(ctx context.Context, dmsAddr string, req dmsV1.ListUserReq) (ret []*dmsV1.ListUser, total int64, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.ListUserReply{}

	url := fmt.Sprintf("%v%v?page_size=%v&page_index=%v&filter_del_user=%v&filter_by_uids=%v", dmsAddr, dmsV1.GetUsersRouter(), req.PageSize, req.PageIndex, req.FilterDeletedUser, req.FilterByUids)

	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to list users from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Msg)
	}

	return reply.Payload.Users, reply.Payload.Total, nil

}
