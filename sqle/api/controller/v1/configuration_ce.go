//go:build !enterprise
// +build !enterprise

package v1

import (
	e "errors"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/utils"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/labstack/echo/v4"
	"net/http"
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
	repository, _, cleanup, err := utils.CloneGitRepository(c.Request().Context(), request.GitHttpUrl, request.GitUserName, request.GitUserPassword)
	defer cleanup()
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
	if defaultBranchIndex == -1 {
		return branches, nil
	}
	//1. 根据第一个元素，找到其余元素中的默认分支
	//2. 如果找到，将找到的默认分支名移到第一个元素，并且删除原来的第一个元素。
	branches[0], branches[defaultBranchIndex] = branches[defaultBranchIndex], branches[0]
	branches = append(branches[:defaultBranchIndex], branches[defaultBranchIndex+1:]...)
	return branches, nil
}
