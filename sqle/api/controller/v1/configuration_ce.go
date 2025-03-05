//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	goGit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
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
	branches, err := getBranches(references)
	return c.JSON(http.StatusOK, &TestGitConnectionResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: TestGitConnectionResDataV1{
			IsConnectedSuccess: true,
			Branches:           branches,
		},
	})
}

func getBranches(references storer.ReferenceIter) ([]string, error) {
	branches := make([]string, 0)
	err := references.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			branches = append(branches, ref.Name().Short())
		}
		return nil
	})
	if err != nil {
		return branches, err
	}
	if len(branches) < 2 {
		return branches, nil
	}
	// 第一个元素确认了默认分支名，需要把可以checkout的默认分支提到第一个元素
	defaultBranch := "origin/" + branches[0]
	defaultBranchIndex := -1
	for i, branch := range branches {
		if branch == defaultBranch {
			defaultBranchIndex = i
			break
		}
	}
	resultBranches := make([]string, len(branches)-1)
	// 将默认分支提到第一个元素
	resultBranches = append(branches[:0], append([]string{branches[defaultBranchIndex]}, branches[1:defaultBranchIndex]...)...)
	if defaultBranchIndex+1 < len(branches)-1 {
		resultBranches = append(resultBranches, branches[defaultBranchIndex+1:]...)
	}
	return resultBranches, nil
}
