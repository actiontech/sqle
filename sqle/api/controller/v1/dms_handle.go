package v1

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	baseV1 "github.com/actiontech/dms/pkg/dms-common/api/base/v1"
	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/iancoleman/strcase"
	"github.com/labstack/echo/v4"
)

var operateHandlers map[string]OperationHanler = make(map[string]OperationHanler)

// NOTE:
// The implemented structure must be named[CamelCase] by the combination of DataResourceType, OperationType, and OperationTimingType
type OperationHanler interface {
	Hanle(ctx context.Context, currentUserId string, dataResourceId string) error
}

func init() {
	for _, v := range []OperationHanler{AfterDeleteNamespace{}, BeforeDeleteNamespace{}, BeforeArvhiveNamespace{}, AfterCreateNamespace{}} {
		structName := strcase.ToSnake(reflect.TypeOf(v).Name())
		operateHandlers[structName] = v
	}
}

// 内部接口
func OperateDataResourceHandle(c echo.Context) error {
	req := new(dmsV1.OperateDataResourceHandleReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	h, ok := operateHandlers[fmt.Sprintf("%s_%s_%s", req.OperationTiming, req.OperationType, req.DataResourceType)]
	if ok {
		err := h.Hanle(c.Request().Context(), "", req.DataResourceUid)
		if err != nil {
			return c.JSON(http.StatusOK, dmsV1.OperateDataResourceHandleReply{GenericResp: baseV1.GenericResp{Code: http.StatusBadRequest, Msg: err.Error()}})
		}
	}
	return c.JSON(http.StatusOK, dmsV1.OperateDataResourceHandleReply{GenericResp: baseV1.GenericResp{Msg: "OK"}})
}

type AfterDeleteNamespace struct {
}

type BeforeArvhiveNamespace struct {
}

type BeforeDeleteNamespace struct {
}
type AfterCreateNamespace struct {
}
