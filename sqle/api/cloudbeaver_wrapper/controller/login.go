package controller

// import (
// 	"context"
// 	"encoding/json"

// 	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
// 	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/service"
// )

// // ActiveUser is the resolver for the activeUser field.
// func (r *QueryResolverImpl) ActiveUser(ctx context.Context) (*model.UserInfo, error) {
// 	data, err := r.Next(r.Ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	resp := &struct {
// 		Data struct {
// 			User *model.UserInfo `json:"user"`
// 		} `json:"data"`
// 	}{}

// 	err = json.Unmarshal(data, resp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.Data.User != nil && resp.Data.User.DisplayName != nil {
// 		*resp.Data.User.DisplayName = service.RestoreFromCloudBeaverUserName(*resp.Data.User.DisplayName)
// 	}

// 	return resp.Data.User, err
// }
