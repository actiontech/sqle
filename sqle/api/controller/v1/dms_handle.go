package v1

import (
	"fmt"
	"net/http"

	baseV1 "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/dms/pkg/dms-common/dmsobject"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/labstack/echo/v4"
)

func init() {
	dmsobject.InitOperateHandlers([]dmsobject.OperationHanler{AfterDeleteProject{}, BeforeDeleteProject{}, BeforeArvhiveProject{}, AfterCreateProject{}})
}

// 内部接口
func OperateDataResourceHandle(c echo.Context) error {
	req := new(dmsV1.OperateDataResourceHandleReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	h, ok := dmsobject.OperateHandlers[fmt.Sprintf("%s_%s_%s", req.OperationTiming, req.OperationType, req.DataResourceType)]
	if ok {
		err := h.Hanle(c.Request().Context(), "", req.DataResourceUid)
		if err != nil {
			return c.JSON(http.StatusOK, dmsV1.OperateDataResourceHandleReply{GenericResp: baseV1.GenericResp{Code: http.StatusBadRequest, Message: err.Error()}})
		}
	}
	return c.JSON(http.StatusOK, dmsV1.OperateDataResourceHandleReply{GenericResp: baseV1.GenericResp{Message: "OK"}})
}

type AfterDeleteProject struct {
}

type BeforeArvhiveProject struct {
}

type BeforeDeleteProject struct {
}
type AfterCreateProject struct {
}
