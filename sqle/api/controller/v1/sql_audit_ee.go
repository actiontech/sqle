//go:build enterprise
// +build enterprise

package v1

import (
	"context"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/common"
	dms "github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/knowledge_base"
	"github.com/labstack/echo/v4"
)

func directGetSQLAnalysis(c echo.Context) error {
	req := new(GetSQLAnalysisReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	projectUid, err := dms.GetPorjectUIDByName(context.TODO(), req.ProjectName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	inst, exist, err := dms.GetInstanceInProjectByName(context.Background(), projectUid, req.InstanceName)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, ErrInstanceNoAccess)
	}

	can, err := CheckCurrentUserCanOpInstances(c.Request().Context(), projectUid, controller.GetUserID(c), []*model.Instance{inst})
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

		metaData := convertExplainAndMetaDataToRes(c.Request().Context(), explainResult, explainMessage, metaDataResult, sql)
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

// TODO 放到knowledge中
func getRuleKnowledge(c echo.Context) error {
	ruleName := c.Param("rule_name")
	dbType := c.Param("db_type")

	knowledge, err := knowledge_base.GetRuleWithKnowledge(c.Request().Context(), ruleName, dbType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}

	return c.JSON(http.StatusOK, GetRuleKnowledgeResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: RuleKnowledgeResV1{
			Rule: RuleInfo{
				Desc:       knowledge.Description,
				Annotation: knowledge.Title,
			},
			KnowledgeContent: knowledge.Content,
		},
	})
}

// TODO 放到knowledge中
func updateRuleKnowledge(c echo.Context) error {
	req := new(UpdateRuleKnowledgeReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ruleName := c.Param("rule_name")
	dbType := c.Param("db_type")
	if req.KnowledgeContent == nil {
		return c.JSON(http.StatusOK, controller.JSONBaseErrorReq(c, nil))
	}

	ctx := c.Request().Context()
	if err := knowledge_base.UpdateRuleKnowledgeContent(ctx, ruleName, dbType, *req.KnowledgeContent); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}

func getCustomRuleKnowledge(c echo.Context) error {
	ruleName := c.Param("rule_name")
	dbType := c.Param("db_type")
	knowledge, err := knowledge_base.GetCustomRuleWithKnowledge(c.Request().Context(), ruleName, dbType)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return c.JSON(http.StatusOK, GetRuleKnowledgeResV1{
		BaseRes: controller.NewBaseReq(nil),
		Data: RuleKnowledgeResV1{
			Rule: RuleInfo{
				Desc:       knowledge.Description,
				Annotation: knowledge.Title,
			},
			KnowledgeContent: knowledge.Content,
		},
	})
}

// TODO 放到knowledge中
func updateCustomRuleKnowledge(c echo.Context) error {
	req := new(UpdateRuleKnowledgeReq)
	if err := controller.BindAndValidateReq(c, req); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	ruleName := c.Param("rule_name")
	dbType := c.Param("db_type")
	if req.KnowledgeContent == nil {
		return c.JSON(http.StatusOK, controller.JSONBaseErrorReq(c, nil))
	}

	ctx := c.Request().Context()
	if err := knowledge_base.UpdateCustomRuleKnowledgeContent(ctx, ruleName, dbType, *req.KnowledgeContent); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	return controller.JSONBaseErrorReq(c, nil)
}
