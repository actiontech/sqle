package dms

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	base "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
	"github.com/actiontech/sqle/sqle/errors"
)

// OpsType 与 AC-B06 / dms-common OpsType 契约对齐（uid + name）。
type OpsType struct {
	UID  string `json:"uid,omitempty"`
	Name string `json:"name"`
}

type listOpsTypesReply struct {
	Data  []*OpsType `json:"data"`
	Total int64      `json:"total_nums"`
	base.GenericResp
}

// ListOpsTypes 调用 DMS 项目运维类型字典 list（与 dmsobject.ListOpsTypes / AC-B06 同一 HTTP 契约）。
// 本仓 vendor 尚未同步 dms-common 新符号时，在此复用既有 pkgHttp + DefaultDMSToken 封装。
func ListOpsTypes(ctx context.Context, projectUID string, pageIndex, pageSize uint32) ([]*OpsType, int64, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/dms/projects/%s/ops_types", GetDMSServerAddress(), projectUID))
	if err != nil {
		return nil, 0, fmt.Errorf("failed to parse ops types URL: %v", err)
	}
	query := url.Values{}
	query.Set("page_size", fmt.Sprintf("%d", pageSize))
	query.Set("page_index", fmt.Sprintf("%d", pageIndex))
	baseURL.RawQuery = query.Encode()

	reply := &listOpsTypesReply{}
	if err := pkgHttp.Get(ctx, baseURL.String(), header, nil, reply); err != nil {
		return nil, 0, fmt.Errorf("failed to list ops types from %v: %v", baseURL.String(), err)
	}
	if reply.Code != 0 {
		return nil, 0, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}
	return reply.Data, reply.Total, nil
}

// ErrOpsTypeNotBelongToProject 所选运维类型不属于本项目字典（用户可读）。
var ErrOpsTypeNotBelongToProject = errors.New(errors.DataInvalid, fmt.Errorf("所选运维类型不属于本项目字典"))

// ValidateOpsTypeBelongToProject 空串视为未设置（通过）；非空须属于本项目字典。
func ValidateOpsTypeBelongToProject(ctx context.Context, projectUID, opsTypeUID string) error {
	opsTypeUID = strings.TrimSpace(opsTypeUID)
	if opsTypeUID == "" {
		return nil
	}
	items, _, err := ListOpsTypes(ctx, projectUID, 1, 1000)
	if err != nil {
		return err
	}
	for _, item := range items {
		if item != nil && item.UID == opsTypeUID {
			return nil
		}
	}
	return ErrOpsTypeNotBelongToProject
}

// BuildOpsTypeNameMap 按项目一次 ListOpsTypes，构建 uid→name 内存 map（AC-B17 / D8）。
// 拉取失败或空项目 → 空 map（展示侧 omitempty，不向上抛错）。
func BuildOpsTypeNameMap(ctx context.Context, projectUID string) map[string]string {
	nameByUID := map[string]string{}
	if strings.TrimSpace(projectUID) == "" {
		return nameByUID
	}
	items, _, err := ListOpsTypes(ctx, projectUID, 1, 1000)
	if err != nil {
		return nameByUID
	}
	for _, item := range items {
		if item != nil && item.UID != "" {
			nameByUID[item.UID] = item.Name
		}
	}
	return nameByUID
}

// ResolveOpsTypeFromMap 从已构建的 uid→name map 回填展示对象。
// 未设置或字典已删/不命中 → nil（omitempty）。
func ResolveOpsTypeFromMap(opsTypeUID string, nameByUID map[string]string) *OpsType {
	opsTypeUID = strings.TrimSpace(opsTypeUID)
	if opsTypeUID == "" || nameByUID == nil {
		return nil
	}
	name, ok := nameByUID[opsTypeUID]
	if !ok {
		return nil
	}
	return &OpsType{UID: opsTypeUID, Name: name}
}

// ResolveOpsTypeDisplay 详情单条解析：内部走 BuildOpsTypeNameMap + ResolveOpsTypeFromMap（与列表共用）。
// 未设置、字典拉取失败、或字典项已删/不命中 → nil（前端「-」/省略）。
func ResolveOpsTypeDisplay(ctx context.Context, projectUID, opsTypeUID string) *OpsType {
	return ResolveOpsTypeFromMap(opsTypeUID, BuildOpsTypeNameMap(ctx, projectUID))
}
