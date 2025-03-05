package dms

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	dmsRegister "github.com/actiontech/dms/pkg/dms-common/register"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetAllUsers(ctx context.Context, dmsAddr string) ([]*model.User, error) {
	ret := make([]*model.User, 0)
	for pageIndex, pageSize := 1, 10; ; pageIndex++ {
		users, _, err := dmsobject.ListUsers(ctx, controller.GetDMSServerAddress(), dmsV1.ListUserReq{PageSize: uint32(pageSize), PageIndex: uint32(pageIndex), FilterDeletedUser: true})
		if err != nil {
			return nil, err
		}
		for _, user := range users {
			ret = append(ret, convertListUserToModel(user))
		}
		if len(users) < pageSize {
			break
		}
	}
	return ret, nil
}

func GetMapUsers(ctx context.Context, userUid []string, dmsAddr string) (map[string] /*user_id*/ *model.User, error) {
	users, err := GetUsers(ctx, userUid, dmsAddr)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]*model.User)
	for _, user := range users {
		ret[user.GetIDStr()] = user
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("cant't find any users")
	}
	return ret, nil
}

func GetUsers(ctx context.Context, userUid []string, dmsAddr string) ([]*model.User, error) {
	users, _, err := dmsobject.ListUsers(ctx, controller.GetDMSServerAddress(), dmsV1.ListUserReq{PageSize: 999, PageIndex: 1, FilterDeletedUser: true, FilterByUids: strings.Join(userUid, ",")})
	if err != nil {
		return nil, err
	}
	ret := make([]*model.User, 0)
	for _, user := range users {
		ret = append(ret, convertListUserToModel(user))
	}
	return ret, nil
}

func convertListUserToModel(user *dmsV1.ListUser) *model.User {
	id, _ := strconv.Atoi(user.UserUid)
	model_ := model.Model{ID: uint(id)}
	if user.IsDeleted {
		// 仅记录为已删除
		model_.DeletedAt = gorm.DeletedAt{Valid: true}
	}
	ret := &model.User{
		Model:    model_,
		Name:     user.Name,
		Email:    user.Email,
		Phone:    user.Phone,
		WeChatID: user.WxID,
	}
	if user.Stat != dmsV1.StatOK && user.Stat != dmsV1.StatOKEn {
		ret.Stat = 1
	}
	return ret
}

func GetUser(ctx context.Context, userUid string, dmsAddr string) (*model.User, error) {
	dmsUser, err := dmsobject.GetUser(ctx, userUid, dmsAddr)
	if err != nil {
		return nil, err
	}
	return convertUserToModel(dmsUser), nil

}

func GetCurrentUserLanguage(c echo.Context) string {
	user, err := controller.GetCurrentUser(c, GetUser)
	if err != nil {
		return ""
	}
	if user.Name == model.DefaultSysUser {
		// 系统用户直接通过请求头AcceptLanguage确定语言
		return i18nPkg.GetLangByAcceptLanguage(c)
	}
	return user.Language
}

func convertUserToModel(user *dmsV1.GetUser) *model.User {
	id, _ := strconv.Atoi(user.UserUid)
	model_ := model.Model{ID: uint(id)}
	ret := &model.User{
		Model:              model_,
		Name:               user.Name,
		Email:              user.Email,
		Phone:              user.Phone,
		WeChatID:           user.WxID,
		Language:           user.Language,
		ThirdPartyUserInfo: user.ThirdPartyUserInfo,
	}
	if user.Stat != dmsV1.StatOK && user.Stat != dmsV1.StatOKEn {
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

// !dms-todo: 注意脚本服务调用该接口，接口修改会导致脚本服务调用原接口失败，需要通知相关开发修改接口调用
// dms-todo: 临时方案
func GetPorjectUIDByName(ctx context.Context, projectName string, needActive ...bool) (projectUID string, err error) {
	project, err := GetPorjectByName(ctx, projectName)
	if err != nil {
		return "", err
	}

	if len(needActive) == 1 && needActive[0] && project.Archived {
		return "", fmt.Errorf("project is archived")
	}
	return project.ProjectUid, nil
}
func GetProjectByID(ProjectId string) (project dmsV1.ListProject, err error) {
	ret, _, err := dmsobject.ListProjects(context.TODO(), GetDMSServerAddress(), dmsV1.ListProjectReq{
		PageSize:    1,
		PageIndex:   1,
		FilterByUID: ProjectId,
	})
	if err != nil {
		return project, err
	}
	if len(ret) > 0 && ret[0] != nil {
		project = *ret[0]
	}
	return project, nil
}

func GetPorjectByName(ctx context.Context, projectName string) (project *dmsV1.ListProject, err error) {
	ret, total, err := dmsobject.ListProjects(ctx, controller.GetDMSServerAddress(), dmsV1.ListProjectReq{
		PageSize:     1,
		PageIndex:    1,
		FilterByName: projectName,
	})
	if err != nil {
		return nil, err
	}
	if total == 0 || len(ret) == 0 {
		return nil, fmt.Errorf("project %s not found", projectName)
	}

	return ret[0], nil
}

func GetProjects() ([]string, error) {
	projectIds := make([]string, 0)
	projects, _, err := dmsobject.ListProjects(context.Background(), controller.GetDMSServerAddress(), dmsV1.ListProjectReq{
		PageSize:  9999,
		PageIndex: 1,
	})
	if err != nil {
		return nil, err
	}
	for _, namespce := range projects {
		projectIds = append(projectIds, namespce.ProjectUid)
	}
	return projectIds, nil
}

// TODO 这里没有考虑到sqled开启https的情况
func RegisterAsDMSTarget(sqleConfig *config.SqleOptions) error {
	controller.InitDMSServerAddress(sqleConfig.DMSServerAddress)
	InitDMSServerAddress(sqleConfig.DMSServerAddress)
	ctx := context.Background()

	// 向DMS注册反向代理
	if err := dmsRegister.RegisterDMSProxyTarget(ctx, controller.GetDMSServerAddress(), "sqle", fmt.Sprintf("http://%v:%v", sqleConfig.APIServiceOpts.Addr, sqleConfig.APIServiceOpts.Port) /* TODO https的处理*/, config.Version, []string{"/sqle/v"}, dmsV1.ProxyScenarioInternalService); nil != err {
		return fmt.Errorf("failed to register dms proxy target: %v", err)
	}
	// 注册校验接口
	if err := dmsRegister.RegisterDMSPlugin(ctx, controller.GetDMSServerAddress(), &dmsV1.Plugin{
		Name:                         "sqle",
		OperateDataResourceHandleUrl: fmt.Sprintf("http://%s:%d/%s/%s", sqleConfig.APIServiceOpts.Addr, sqleConfig.APIServiceOpts.Port, "v1", "data_resource/handle"),
		GetDatabaseDriverOptionsUrl:  fmt.Sprintf("http://%s:%d/%s/%s", sqleConfig.APIServiceOpts.Addr, sqleConfig.APIServiceOpts.Port, "v1", "database_driver_options"),
		GetDatabaseDriverLogosUrl:    fmt.Sprintf("http://%s:%d/%s/%s", sqleConfig.APIServiceOpts.Addr, sqleConfig.APIServiceOpts.Port, "v1", "database_driver_logos"),
	}); err != nil {
		return fmt.Errorf("failed to register dms plugin for operation data source handle")
	}

	return nil
}

func ListProjectUserTips(ctx context.Context, projectUid string) (users []*model.User, err error) {
	dmsUsers, _, err := dmsobject.ListMembersInProject(ctx, controller.GetDMSServerAddress(), dmsV1.ListMembersForInternalReq{
		PageSize:   999,
		PageIndex:  1,
		ProjectUid: projectUid,
	})
	if err != nil {
		return nil, fmt.Errorf("get user from dms error: %v", err)
	}

	for _, dmsUser := range dmsUsers {
		id, err := strconv.Atoi(dmsUser.User.Uid)
		if err != nil {
			return nil, err
		}
		model_ := model.Model{ID: uint(id)}
		users = append(users, &model.User{
			Model: model_,
			Name:  dmsUser.User.Name,
		})
	}
	return users, nil
}
