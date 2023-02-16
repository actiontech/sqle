package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	gqlClient "github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/client"
	"github.com/actiontech/sqle/sqle/api/cloudbeaver_wrapper/graph/model"
	sqleModel "github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/utils"
)

func SyncCurrentUser(cloudBeaverUser string) error {
	ctx := context.Background()

	// 获取SQLE缓存中的CloudBeaver用户信息和SQLE用户信息
	s := sqleModel.GetStorage()
	cache, cbExist, err := s.GetCloudBeaverUserCacheByCBUserID(cloudBeaverUser)
	if err != nil {
		return err
	}

	sqleUserName := RestoreFromCloudBeaverUserName(cloudBeaverUser)
	sqleUser, exist, err := s.GetUserByName(sqleUserName)
	if err != nil {
		return err
	}
	if !exist { // SQLE用户不存在有可能是用户使用自行添加的用户导致的, 此用户因为与SQLE无关, 所以忽略
		return nil
	}

	// 如果缓存存在且指纹校验通过, 则认为用户同步过且当前缓存为最新缓存
	// 反之则需要触发同步并更新缓存
	if cbExist && sqleUser.FingerPrint() == cache.SQLEFingerprint {
		return nil
	}

	if IsReserved(cloudBeaverUser) {
		return fmt.Errorf("this username cannot be used")
	}

	// 使用管理员身份登录
	client, err := GetGQLClientWithRootUser()
	if err != nil {
		return err
	}

	checkExistReq := gqlClient.NewRequest(QueryGQL.IsUserExistQuery(), map[string]interface{}{
		"userId": cloudBeaverUser,
	})

	type UserList struct {
		ListUsers []struct {
			UserID string `json:"userID"`
		} `json:"listUsers"`
	}
	users := UserList{}

	err = client.Run(ctx, checkExistReq, &users)
	if err != nil {
		return fmt.Errorf("check cloudbeaver user exist failed: %v", err)
	}

	// 用户不存在则创建CloudBeaver用户
	if len(users.ListUsers) == 0 {
		// 创建用户
		createUserReq := gqlClient.NewRequest(QueryGQL.CreateUserQuery(), map[string]interface{}{
			"userId": cloudBeaverUser,
		})
		err = client.Run(ctx, createUserReq, &UserList{})
		if err != nil {
			return fmt.Errorf("create cloudbeaver user failed: %v", err)
		}

		// 授予角色(不授予角色的用户无法登录)
		grantUserRoleReq := gqlClient.NewRequest(QueryGQL.GrantUserRoleQuery(), map[string]interface{}{
			"userId": cloudBeaverUser,
			"roleId": CBUserRole,
			"teamId": CBUserRole,
		})
		err = client.Run(ctx, grantUserRoleReq, nil)
		if err != nil {
			return fmt.Errorf("create cloudbeaver user failed: %v", err)
		}
	}

	// 更新CloudBeaver用户密码
	updatePasswordReq := gqlClient.NewRequest(QueryGQL.UpdatePasswordQuery(), map[string]interface{}{
		"userId": cloudBeaverUser,
		"credentials": model.JSON{
			"password": strings.ToUpper(utils.Md5(sqleUser.Password)),
		},
	})
	err = client.Run(ctx, updatePasswordReq, nil)
	if err != nil {
		return fmt.Errorf("update cloudbeaver user failed: %v", err)
	}

	// 更新SQLE缓存
	return s.UpdateCloudBeaverUserCache(sqleUser.ID, cloudBeaverUser)
}

// IsReserved 会检查用户名是否为无法使用的名称, 比如admin和user是角色名, 按照CloudBeaver的要求角色名无法作为用户名
func IsReserved(name string) bool {
	_, ok := map[string]struct{}{
		"admin": {},
		"user":  {},
	}[name]
	return ok
}

const CBNamePrefix = "sqle-"

func GenerateCloudBeaverUserName(name string) string {
	return CBNamePrefix + name
}

func RestoreFromCloudBeaverUserName(name string) string {
	return strings.TrimPrefix(name, CBNamePrefix)
}

// LoginToCBServer 的登录请求会直接被转发, 不会被中间件拦截处理
func LoginToCBServer(user, pwd string) (cookie []*http.Cookie, err error) {
	client := gqlClient.NewClient(GetGqlServerURI(), gqlClient.WithHttpResHandler(
		func(response *http.Response) {
			if response != nil {
				cookie = response.Cookies()
			}
		}))
	req := gqlClient.NewRequest(QueryGQL.LoginQuery(), map[string]interface{}{
		"credentials": model.JSON{
			"user":     user,
			"password": strings.ToUpper(utils.Md5(pwd)), // the password is an all-caps md5-32 string
		},
	})

	res := struct {
		AuthInfo struct {
			AuthId interface{} `json:"authId"`
		} `json:"authInfo"`
	}{}
	if err := client.Run(context.TODO(), req, &res); err != nil {
		return cookie, fmt.Errorf("cloudbeaver login failed: %v", err)
	}

	return cookie, nil
}