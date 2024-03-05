package v2

import (
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	v1 "github.com/actiontech/sqle/sqle/api/controller/v1"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/labstack/echo/v4"
)

type ExtractSQLOperationsReqV2 struct {
	InstanceType string `json:"instance_type" form:"instance_type" example:"MySQL" valid:"required"`
	SQLContent   string `json:"sql_content" form:"sql_content" example:"select * from t1; select * from t2;" valid:"required"`
}

type ExtractSQLOperationsResV2 struct {
	controller.BaseRes
	Data []ExtractSQLOperationsDataV2 `json:"data"`
}

type ExtractSQLOperationsDataV2 struct {
	DbName     string `json:"db_name"`
	SchemaName string `json:"schema_name"`
	TableName  string `json:"table_name"`
	Operation  string `json:"operation"`
}

// ExtractSQLOperations
// @Summary 获取sql操作对象
// @Description extract sql operations
// @Id GetSQLOperations
// @Tags extract_sql_operations
// @Security ApiKeyAuth
// @Param req body v2.ExtractSQLOperationsReqV2 true "extract sql operations"
// @Success 200 {object} v2.ExtractSQLOperationsResV2
// @router /v2/sql_operations [post]
func ExtractSQLOperations(c echo.Context) error {
	req := new(ExtractSQLOperationsReqV2)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}

	l := log.NewEntry().WithField(c.Path(), "extract sql operations")

	ops, err := server.GetSQLOperations(l, req.SQLContent, req.InstanceType)
	if err != nil {
		l.Errorf("extract sql operations failed: %v", err)
		return controller.JSONBaseErrorReq(c, v1.ErrDirectAudit)
	}

	data := make([]ExtractSQLOperationsDataV2, 0, len(ops))
	for _, op := range ops {
		for _, item := range op.ObjectOps {
			data = append(data, ExtractSQLOperationsDataV2{
				DbName:     item.Object.DatabaseName,
				SchemaName: item.Object.SchemaName,
				TableName:  item.Object.TableName,
				Operation:  string(item.Op),
			})
		}
	}

	return c.JSON(http.StatusOK, ExtractSQLOperationsResV2{
		BaseRes: controller.BaseRes{},
		Data:    data,
	})
}
