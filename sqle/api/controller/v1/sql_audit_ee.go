//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
)

func directGetSQLAnalysis(c echo.Context) error {
	req := new(GetSQLAnalysisReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	s := model.GetStorage()
	inst, exist, err := s.GetInstanceByNameAndProjectName(req.InstanceName, req.ProjectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	can, err := checkCurrentUserCanAccessInstance(c, inst)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !can {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	plugin, err := common.NewDriverManagerWithoutAudit(log.NewEntry(), inst, req.SchemaName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	defer plugin.Close(context.TODO())

	nodes, err := plugin.Parse(context.TODO(), req.Sql)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	var resp []*SqlAnalysisResDataV1
	for _, node := range nodes {
		sql := node.Text
		explainResult, explainMessage, metaDataResult, err := getSQLAnalysisResultFromDriver(log.NewEntry(), req.SchemaName, sql, inst)
		if err != nil {
			return controller.JSONBaseErrorReq(c, err)
		}

		metaData := convertExplainAndMetaDataToRes(explainResult, explainMessage, metaDataResult, sql)
		resp = append(resp, &SqlAnalysisResDataV1{
			SQLExplain: metaData.SQLExplain,
			TableMetas: metaData.TableMetas,
		})
	}

	return c.JSON(http.StatusOK, DirectGetSQLAnalysisResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data:    resp,
	})
}
