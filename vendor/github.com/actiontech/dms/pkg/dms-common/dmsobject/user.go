package dmsobject

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

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
		return nil, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, nil
}

func GetMemberGroup(ctx context.Context, memberGroupUid, projectUid string, dmsAddr string) (*dmsV1.GetMemberGroup, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.GetMemberGroupReply{}

	url := fmt.Sprintf("%v%v", dmsAddr, dmsV1.GetMemberGroupRouter(memberGroupUid, projectUid))

	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, fmt.Errorf("failed to get member group from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, nil
}

func GetUserOpPermission(ctx context.Context, projectUid, userUid, dmsAddr string) (ret []dmsV1.OpPermissionItem, isAdmin bool, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reqBody := struct {
		UserOpPermission *dmsV1.UserOpPermission `json:"user_op_permission"`
	}{
		UserOpPermission: &dmsV1.UserOpPermission{ProjectUid: projectUid},
	}

	reply := &dmsV1.GetUserOpPermissionReply{}

	url := fmt.Sprintf("%v%v", dmsAddr, dmsV1.GetUserOpPermissionRouter(userUid))

	if err := pkgHttp.Get(ctx, url, header, reqBody, reply); err != nil {
		return nil, false, fmt.Errorf("failed to get user op permission from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, false, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data.OpPermissionList, reply.Data.IsAdmin, nil

}

// GetUserOpPermissionWithBWP is like GetUserOpPermission but also returns the BusinessWritePermission field.
func GetUserOpPermissionWithBWP(ctx context.Context, projectUid, userUid, dmsAddr string) (ret []dmsV1.OpPermissionItem, isAdmin bool, businessWritePermission bool, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reqBody := struct {
		UserOpPermission *dmsV1.UserOpPermission `json:"user_op_permission"`
	}{
		UserOpPermission: &dmsV1.UserOpPermission{ProjectUid: projectUid},
	}

	reply := &dmsV1.GetUserOpPermissionReply{}

	url := fmt.Sprintf("%v%v", dmsAddr, dmsV1.GetUserOpPermissionRouter(userUid))

	if err := pkgHttp.Get(ctx, url, header, reqBody, reply); err != nil {
		return nil, false, true, fmt.Errorf("failed to get user op permission from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, false, true, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data.OpPermissionList, reply.Data.IsAdmin, reply.Data.BusinessWritePermission, nil
}

func ListMembersInProject(ctx context.Context, dmsAddr string, req dmsV1.ListMembersForInternalReq) ([]*dmsV1.ListMembersForInternalItem, int64, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.ListMembersForInternalReply{}

	url := fmt.Sprintf("%v%v?page_size=%v&page_index=%v", dmsAddr, dmsV1.GetListMembersForInternalRouter(req.ProjectUid), req.PageSize, req.PageIndex)

	if err := pkgHttp.Get(ctx, url, header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to get member from %v: %v", url, err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil
}

func ListUsers(ctx context.Context, dmsAddr string, req dmsV1.ListUserReq) (ret []*dmsV1.ListUser, total int64, err error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	reply := &dmsV1.ListUserReply{}

	// 构建查询参数
	params := url.Values{}
	params.Set("page_size", strconv.FormatUint(uint64(req.PageSize), 10))
	if req.PageIndex > 0 {
		params.Set("page_index", strconv.FormatUint(uint64(req.PageIndex), 10))
	}
	if req.OrderBy != "" {
		params.Set("order_by", string(req.OrderBy))
	}
	if req.FilterByName != "" {
		params.Set("filter_by_name", req.FilterByName)
	}
	if req.FilterByUids != "" {
		params.Set("filter_by_uids", req.FilterByUids)
	}
	if req.FilterDeletedUser {
		params.Set("filter_del_user", "true")
	}
	if req.FuzzyKeyword != "" {
		params.Set("fuzzy_keyword", req.FuzzyKeyword)
	}
	if req.FilterByEmail != "" {
		params.Set("filter_by_email", req.FilterByEmail)
	}
	if req.FilterByPhone != "" {
		params.Set("filter_by_phone", req.FilterByPhone)
	}
	if req.FilterByStat != "" {
		params.Set("filter_by_stat", string(req.FilterByStat))
	}
	if req.FilterByAuthenticationType != "" {
		params.Set("filter_by_authentication_type", string(req.FilterByAuthenticationType))
	}
	if req.FilterBySystem != "" {
		params.Set("filter_by_system", string(req.FilterBySystem))
	}

	requestURL := fmt.Sprintf("%v%v?%v", dmsAddr, dmsV1.GetUsersRouter(), params.Encode())

	if err := pkgHttp.Get(ctx, requestURL, header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to list users from %v: %v", requestURL, err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply.Data, reply.Total, nil

}
