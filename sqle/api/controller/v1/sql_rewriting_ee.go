//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/fillsql"
	"github.com/actiontech/sqle/sqle/server/sqlrewriting"
	"github.com/labstack/echo/v4"
)

func getTaskRewrittenSQLData(c echo.Context) error {
	taskID := c.Param("task_id")
	sqlNumber := c.Param("number")
	s := model.GetStorage()
	task, err := GetTaskById(c.Request().Context(), taskID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if err := CheckCurrentUserCanViewTask(c, task); err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	taskSql, exist, err := s.GetTaskSQLByNumber(taskID, sqlNumber)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if !exist {
		return controller.JSONBaseErrorReq(c, errors.NewDataNotExistErr("sql number not found"))
	}
	sqlContent, err := fillsql.FillingSQLWithParamMarker(taskSql.Content, task)
	if err != nil {
		log.NewEntry().Errorf("fill param marker sql failed: %v", err)
		sqlContent = taskSql.Content
	}
	res, err := GetSQLAnalysisResult(log.NewEntry(), task.Instance, task.Schema, sqlContent)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	if res.TableMetaResultErr != nil {
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get table meta failed: %v", res.TableMetaResultErr))
	}
	// TODO: 需要Explain和PerformanceStatistics
	taskDbType, err := s.GetTaskDbTypeByID(taskID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	params := &sqlrewriting.SQLRewritingParams{
		DBType:          taskDbType,
		SQL:             taskSql,
		TableStructures: res.TableMetaResult.TableMetas,
		Explain:         nil,
	}
	// 进行重写
	rewrittenRes, err := sqlrewriting.SQLRewriting(c.Request().Context(), params)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 以下开始构造返回结果
	// 考虑将来可能也会需要ExplainResult和PerformanceStatistics，这里直接复用convertSQLAnalysisResultToRes
	analysisResult := convertSQLAnalysisResultToRes(c.Request().Context(), res, taskSql.Content)
	ret := &TaskRewrittenSQLData{
		TableMetas:                analysisResult.TableMetas,
		BusinessNonEquivalentDesc: rewrittenRes.BusinessNonEquivalent,
	}
	var lastRewrittenSQL string
	lang := locale.Bundle.GetLangTagFromCtx(c.Request().Context())
	for _, suggestion := range rewrittenRes.Suggestions {
		ruleName := sqlrewriting.ConvertRuleIDToRuleName(suggestion.RuleID)
		r, exist, err := s.GetRule(ruleName, taskDbType)
		if err != nil {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("get rule failed: %v", err))
		}
		if !exist {
			return controller.JSONBaseErrorReq(c, fmt.Errorf("rule not found: %s", ruleName))
		}
		if suggestion.RewrittenSql != "" {
			lastRewrittenSQL = suggestion.RewrittenSql
		}
		ret.Suggestions = append(ret.Suggestions, &SQLRewrittenSuggestion{
			RuleName:     r.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc,
			AuditLevel:   r.Level,
			Type:         suggestion.Type,
			Desc:         suggestion.Description,
			RewrittenSQL: suggestion.RewrittenSql,
			DDL_DCL_desc: suggestion.DDLDCLDesc,
			DDL_DCL:      suggestion.DDLDCL,
		})
	}
	ret.RewrittenSQL = lastRewrittenSQL
	return c.JSON(http.StatusOK, &GetTaskRewrittenSQLRes{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ret,
	})
}
