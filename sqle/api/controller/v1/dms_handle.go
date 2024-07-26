package v1

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	baseV1 "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func init() {
	dmsobject.InitOperateHandlers([]dmsobject.OperationHandler{
		AfterDeleteProject{},
		BeforeDeleteProject{},
		BeforeArchiveProject{},
		AfterCreateProject{},
		BeforeDeleteDbService{},
		AfterUpdateDbService{},
	})
}

// 内部接口
func OperateDataResourceHandle(c echo.Context) error {
	req := new(dmsV1.OperateDataResourceHandleReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	h := dmsobject.GetOperateHandle(fmt.Sprintf("%s_%s_%s", req.OperationTiming, req.OperationType, req.DataResourceType))

	if err := h.Handle(c.Request().Context(), "", req.DataResourceUid); err != nil {
		return c.JSON(http.StatusOK, dmsV1.OperateDataResourceHandleReply{GenericResp: baseV1.GenericResp{Code: http.StatusBadRequest, Message: err.Error()}})
	}

	return c.JSON(http.StatusOK, dmsV1.OperateDataResourceHandleReply{GenericResp: baseV1.GenericResp{Message: "OK"}})
}

type AfterDeleteProject struct {
}

type BeforeArchiveProject struct {
}

type BeforeDeleteProject struct {
}

type AfterCreateProject struct {
}

type BeforeDeleteDbService struct {
}

type AfterUpdateDbService struct {
}

func (h BeforeDeleteDbService) Handle(ctx context.Context, currentUserId string, instanceIdStr string) error {
	instanceId, err := strconv.ParseInt(instanceIdStr, 10, 64)
	if err != nil {
		return err
	}

	return common.CheckDeleteInstance(instanceId)
}

func (h AfterUpdateDbService) Handle(ctx context.Context, currentUserId string, instanceIdStr string) error {
	instanceId, err := strconv.ParseInt(instanceIdStr, 10, 64)
	if err != nil {
		return err
	}
	instance, _, err := dms.GetInstancesById(ctx, instanceIdStr)
	if err != nil {
		return err
	}
	s := model.GetStorage()
	insAuditPlan, exist, err := s.GetInstanceAuditPlanByInstanceID(instanceId)
	if err != nil {
		return err
	}
	if !exist {
		return nil
	}
	if needUpdate := h.isNeedUpdateService(instance, insAuditPlan); needUpdate {
		insAuditPlan.Business = instance.Business
		err := s.Save(insAuditPlan)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h AfterUpdateDbService) isNeedUpdateService(instance *model.Instance, insAuditPlan *model.InstanceAuditPlan) bool {
	needUpdate := false
	if instance.Business != insAuditPlan.Business {
		return true
	}
	return needUpdate
}
