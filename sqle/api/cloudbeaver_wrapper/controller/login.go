package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/common"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
	"github.com/actiontech/sqle/sqle/log"
)

// ActiveUser is the resolver for the activeUser field.
func (r *QueryResolverImpl) ActiveUser(ctx context.Context) (*model.UserInfo, error) {
	data, err := r.Next(r.Ctx)
	if err != nil {
		return nil, err
	}

	resp := &struct {
		Data struct {
			User *model.UserInfo `json:"user"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}

	if resp.Data.User != nil && resp.Data.User.DisplayName != nil {
		*resp.Data.User.DisplayName = service.RestoreFromCloudBeaverUserName(*resp.Data.User.DisplayName)
	}

	return resp.Data.User, err
}

// AuthLogout is the resolver for the authLogout field.
func (r *QueryResolverImpl) AuthLogout(ctx context.Context, provider *string, configuration *string) (*bool, error) {
	bySqle := r.Ctx.Request().Header.Get(common.InvokedBySqleKey)
	// 不允许从cloudbeaver页面登出
	if bySqle != common.InvokedBySqleValue {
		errMsg := "please logout cloudbeaver from SQLE"
		log.NewEntry().Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	data, err := r.Next(r.Ctx)
	if err != nil {
		return nil, err
	}
	resp := &struct {
		Data struct {
			AuthLogout *bool `json:"authLogout"`
		} `json:"data"`
	}{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, err
	}
	return resp.Data.AuthLogout, nil
}
