package dms

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	dmsRegister "github.com/actiontech/dms/pkg/dms-common/register"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
)

func GetMapUsers(ctx context.Context, userUid []string, dmsAddr string) (map[string] /*user_id*/ *model.User, error) {
	users, _, err := dmsobject.ListUsers(ctx, controller.GetDMSServerAddress(), dmsV1.ListUserReq{PageSize: 999, PageIndex: 1, FilterDeletedUser: true, FilterByUids: strings.Join(userUid, ",")})
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*model.User)
	for _, user := range users {
		userModel := convertListUserToModel(user)
		ret[userModel.GetIDStr()] = userModel
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("cant't find any users")
	}
	return ret, nil
}

func convertListUserToModel(user *dmsV1.ListUser) *model.User {
	id, _ := strconv.Atoi(user.UserUid)
	model_ := model.Model{ID: uint(id)}
	if user.IsDeleted {
		// 仅记录为已删除
		model_.DeletedAt = &time.Time{}
	}
	ret := &model.User{
		Model:    model_,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		WeChatID: user.WxID,
	}
	if user.Stat != dmsV1.StatOK {
		ret.Stat = 1
	}
	return ret
}

func GetUser(ctx context.Context, userUid string, dmsAddr string) (*model.User, error) {
	dmsUser, err := dmsobject.GetUser(ctx, userUid, controller.GetDMSServerAddress())
	if err != nil {
		return nil, err
	}
	return convertUserToModel(dmsUser), nil

}

func convertUserToModel(user *dmsV1.GetUser) *model.User {
	id, _ := strconv.Atoi(user.UserUid)
	model_ := model.Model{ID: uint(id)}
	ret := &model.User{
		Model:    model_,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		WeChatID: user.WxID,
	}
	if user.Stat != dmsV1.StatOK {
		ret.Stat = 1
	}
	return ret
}

// dms-todo: 1. 缓存 user 信息；2. 后续考虑所有需要name的接口返回 user id + name 组合的形式
func GetUserNameWithDelTag(userId string) string {
	if userId == "" {
		return ""
	}
	users, _, err := dmsobject.ListUsers(context.TODO(), controller.GetDMSServerAddress(), dmsV1.ListUserReq{PageSize: 1, PageIndex: 1, FilterDeletedUser: true, FilterByUids: userId})
	if err != nil {
		log.NewEntry().WithField("user_id", userId).Errorln("fail to get user from dms")
		return ""
	}
	if len(users) == 0 {
		return ""
	}
	user := users[0]

	if user.IsDeleted {
		return fmt.Sprintf("%s[x]", user.Name)
	}
	return user.Name
}

// dms-todo: 临时方案
func GetPorjectUIDByName(ctx context.Context, projectName string) (projectUID string, err error) {
	ret, total, err := dmsobject.ListNamespaces(ctx, controller.GetDMSServerAddress(), dmsV1.ListNamespaceReq{
		PageSize:     1,
		PageIndex:    1,
		FilterByName: projectName,
	})
	if err != nil {
		return "", err
	}
	if total == 0 || len(ret) == 0 {
		return "", fmt.Errorf("namespace %s not found", projectName)
	}
	if ret[0].Archived {
		return "", fmt.Errorf("project is archived")
	}
	return ret[0].NamespaceUid, nil
}

func GetPorjectByName(ctx context.Context, projectName string) (project *dmsV1.ListNamespace, err error) {
	ret, total, err := dmsobject.ListNamespaces(ctx, controller.GetDMSServerAddress(), dmsV1.ListNamespaceReq{
		PageSize:     1,
		PageIndex:    1,
		FilterByName: projectName,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 || len(ret) == 0 {
		return nil, fmt.Errorf("namespace %s not found", projectName)
	}
	if ret[0].Archived {
		return nil, fmt.Errorf("project is archived")
	}
	return ret[0], nil
}

func GetProjects() ([]string, error) {
	projectIds := make([]string, 0)
	namespaces, _, err := dmsobject.ListNamespaces(context.Background(), controller.GetDMSServerAddress(), dmsV1.ListNamespaceReq{
		PageSize:  9999,
		PageIndex: 1,
	})
	if err != nil {
		return nil, err
	}
	for _, namespce := range namespaces {
		projectIds = append(projectIds, namespce.NamespaceUid)
	}
	return projectIds, nil
}

func RegisterAsDMSTarget(sqleConfig config.SqleConfig) error {
	controller.InitDMSServerAddress(sqleConfig.DMSServerAddress)
	ctx := context.Background()

	// 向DMS注册反向代理
	if err := dmsRegister.RegisterDMSProxyTarget(ctx, controller.GetDMSServerAddress(), "sqle", fmt.Sprintf("http://%v:%v", sqleConfig.SqleServerHost, sqleConfig.SqleServerPort) /* TODO https的处理*/, config.Version, []string{"/sqle/v"}); nil != err {
		return fmt.Errorf("failed to register dms proxy target: %v", err)
	}
	// 注册校验接口
	if err := dmsRegister.RegisterDMSPlugin(ctx, controller.GetDMSServerAddress(), &dmsV1.Plugin{
		Name:                         "sqle",
		OperateDataResourceHandleUrl: fmt.Sprintf("http://%s:%d/%s/%s", sqleConfig.SqleServerHost, sqleConfig.SqleServerPort, "v1", "data_resource/handle"),
	}); err != nil {
		return fmt.Errorf("failed to register dms plugin for operation data source handle")
	}

	return nil
}
