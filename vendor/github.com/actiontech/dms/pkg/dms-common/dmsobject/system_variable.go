package dmsobject

import (
	"context"
	"fmt"
	"net/url"

	baseV1 "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

// GetSystemVariables 获取系统变量配置
func GetSystemVariables(ctx context.Context, dmsAddr string) (*dmsV1.GetSystemVariablesReply, error) {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	// 构建基础 URL
	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/dms/configurations/system_variables", dmsAddr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %v", err)
	}

	// 调用 HTTP GET 请求
	reply := &dmsV1.GetSystemVariablesReply{}
	if err := pkgHttp.Get(ctx, baseURL.String(), header, nil, reply); err != nil {
		return nil, err
	}
	if reply.Code != 0 {
		return nil, fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return reply, nil
}

// UpdateSystemVariables 更新系统变量配置
func UpdateSystemVariables(ctx context.Context, dmsAddr string, req *dmsV1.UpdateSystemVariablesReqV1) error {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}

	// 构建基础 URL
	baseURL, err := url.Parse(fmt.Sprintf("%s/v1/dms/configurations/system_variables", dmsAddr))
	if err != nil {
		return fmt.Errorf("failed to parse base URL: %v", err)
	}

	// 调用 HTTP PATCH 请求
	reply := &baseV1.GenericResp{}
	if err := pkgHttp.Call(ctx, "PATCH", baseURL.String(), header, req, reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Message)
	}

	return nil
}
