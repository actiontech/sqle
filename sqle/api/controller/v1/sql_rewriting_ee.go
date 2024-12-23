//go:build enterprise
// +build enterprise

package v1

import (
	"fmt"
	"net/http"

	"github.com/actiontech/sqle/sqle/api/controller"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/errors"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server/fillsql"
	"github.com/actiontech/sqle/sqle/server/sqlrewriting"
	"github.com/labstack/echo/v4"
)

type RewriteType string

const (
	TypeStatement RewriteType = "statement" // 语句级重写
	TypeStructure RewriteType = "structure" // 结构级重写
	TypeOther     RewriteType = "other"     // 其他
)

func getRewriteSQLData(c echo.Context) error {
	req := new(RewriteSQLReq)
	err := controller.BindAndValidateReq(c, req)
	if err != nil {
		return err
	}
	enableStructure := req.EnableStructureType

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
		if res.TableMetaResultErr == driverV2.ErrSQLIsNotSupported {
			res.TableMetaResult = &driver.GetTableMetaBySQLResult{}
		}
		return controller.JSONBaseErrorReq(c, fmt.Errorf("get table meta failed: %v", res.TableMetaResultErr))
	}
	// TODO: 需要Explain和PerformanceStatistics
	taskDbType, err := s.GetTaskDbTypeByID(taskID)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	params := &sqlrewriting.SQLRewritingParams{
		DBType:              taskDbType,
		SQL:                 taskSql,
		TableStructures:     res.TableMetaResult.TableMetas,
		Explain:             nil,
		EnableStructureType: enableStructure,
	}
	// 进行重写
	rewrittenRes, err := sqlrewriting.SQLRewriting(c.Request().Context(), params)
	if err != nil {
		return controller.JSONBaseErrorReq(c, err)
	}
	// 以下开始构造返回结果
	ret := &RewriteSQLData{
		BusinessDesc:              rewrittenRes.BusinessDesc,
		LogicDesc:                 rewrittenRes.LogicDesc,
		BusinessNonEquivalentDesc: rewrittenRes.BusinessNonEquivalent,
		RewrittenSQLBusinessDesc:  rewrittenRes.BusinessDescAfterOptimize,
		RewrittenSQLLogicDesc:     rewrittenRes.LogicDescAfterOptimize,
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
		ret.Suggestions = append(ret.Suggestions, &RewriteSuggestion{
			RuleName:     r.I18nRuleInfo.GetRuleInfoByLangTag(lang).Desc,
			AuditLevel:   r.Level,
			Type:         string(sqlRewriteSuggestionTypeConvert(suggestion.Type)),
			Desc:         suggestion.Description,
			RewrittenSQL: suggestion.RewrittenSql,
			DDL_DCL_desc: suggestion.DDLDCLDesc,
			DDL_DCL:      suggestion.DDLDCL,
		})
	}
	ret.RewrittenSQL = lastRewrittenSQL
	return c.JSON(http.StatusOK, &RewriteSQLRes{
		BaseRes: controller.NewBaseReq(nil),
		Data:    ret,
	})
}

func sqlRewriteSuggestionTypeConvert(typ string) RewriteType {
	switch typ {
	case "语句级优化":
		return TypeStatement
	case "结构级优化":
		return TypeStructure
	default:
		return TypeOther
	}
}
