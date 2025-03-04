//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	goGit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	goGitTransport "github.com/go-git/go-git/v5/plumbing/transport/http"
	"net/http"
	"os"

	"github.com/actiontech/sqle/sqle/errors"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

var (
	errCommunityEditionNotSupportFeishuAudit         = errors.New(errors.EnterpriseEditionFeatures, e.New("feishu audit is enterprise version feature"))
	errCommunityEditionNotSupportDingDingAudit       = errors.New(errors.EnterpriseEditionFeatures, e.New("dingding audit is enterprise version feature"))
	errCommunityEditionNotSupportWechatAudit         = errors.New(errors.EnterpriseEditionFeatures, e.New("wechat audit is enterprise version feature"))
	errCommunityEditionDoesNotSupportScheduledNotify = errors.New(errors.EnterpriseEditionFeatures, e.New("community edition does not support workflow scheduled notify"))
)

func updateFeishuAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportFeishuAudit)
}

func getFeishuAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportFeishuAudit)
}

func testFeishuAuditConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportFeishuAudit)
}

func getDingTalkConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportDingDingAudit)
}

func updateDingTalkConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportDingDingAudit)
}

func testDingTalkConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportDingDingAudit)
}

func getWechatAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func updateWechatAuditConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func testWechatAuditConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func getCodingConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func updateCodingConfigurationV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func testCodingAuditConfigV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionNotSupportWechatAudit)
}

func getScheduledTaskDefaultOptionV1(c echo.Context) error {
	return controller.JSONBaseErrorReq(c, errCommunityEditionDoesNotSupportScheduledNotify)
}

func testGitConnectionV1(c echo.Context) error {
	request := new(TestGitConnectionReqV1)
	if err := controller.BindAndValidateReq(c, request); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	directory, err := os.MkdirTemp("./", "git-repo-")
	defer os.RemoveAll(directory)
	cloneOpts := &goGit.CloneOptions{
		URL: request.GitHttpUrl,
	}
	// public repository do not require an user name and password
	userName := c.FormValue(GitUserName)
	password := c.FormValue(GitPassword)
	if userName != "" {
		cloneOpts.Auth = &goGitTransport.BasicAuth{
			Username: userName,
			Password: password,
		}
	}
	repository, err := goGit.PlainCloneContext(c.Request().Context(), directory, false, cloneOpts)
	if err != nil {
		return c.JSON(http.StatusOK, &TestGitConnectionResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestGitConnectionResDataV1{
				IsConnectedSuccess: false,
				ErrorMessage:       err.Error(),
			},
		})
	}
	references, err := repository.References()
	if err != nil {
		return c.JSON(http.StatusOK, &TestGitConnectionResV1{
			BaseRes: controller.NewBaseReq(nil),
			Data: TestGitConnectionResDataV1{
				IsConnectedSuccess: false,
				ErrorMessage:       err.Error(),
			},
		})
	}
	branches := make([]string, 0)
	err = references.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			branches = append(branches, ref.Name().Short())
		}
		return nil
	})

	return c.JSON(http.StatusOK, &TestGitConnectionResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestGitConnectionResDataV1{
			IsConnectedSuccess: true,
			Branches:           branches,
		},
	})
}
