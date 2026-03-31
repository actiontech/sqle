package v1

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/actiontech/sqle/sqle/utils"
)

// BuildAuditReportData 从 Task 和数据库查询构建报告数据。
// 该函数放在 controller 层而非 utils 层，因为 utils 被 model 引用，
// 若 utils 反向引用 model 会产生循环依赖。
func BuildAuditReportData(task *model.Task, s *model.Storage, noDuplicate bool, ctx context.Context) (*utils.AuditReportData, error) {
	// 1. 获取 SQL 列表
	data := map[string]interface{}{
		"task_id":      fmt.Sprintf("%d", task.ID),
		"no_duplicate": noDuplicate,
	}

	taskSQLsDetail, _, err := s.GetTaskSQLsByReq(data)
	if err != nil {
		return nil, fmt.Errorf("get task SQLs failed: %w", err)
	}

	// 2. 获取回滚 SQL 映射
	rollbackSqlMap, err := server.BackupService{}.GetRollbackSqlsMap(task.ID)
	if err != nil {
		return nil, fmt.Errorf("get rollback SQLs failed: %w", err)
	}

	// 3. 构建 SQL 列表和统计数据
	levelDist := make(map[string]int)
	ruleHits := make(map[string]int)
	var sqlList []utils.AuditSQLItem
	var problemSQLs []utils.AuditSQLItem

	for _, td := range taskSQLsDetail {
		// 构造临时 ExecuteSQL 对象以复用状态描述方法
		tempSQL := &model.ExecuteSQL{
			AuditResults: td.AuditResults,
			AuditStatus:  td.AuditStatus,
		}
		tempSQL.ExecStatus = td.ExecStatus

		// 提取规则名称和审核建议
		ruleName, suggestion := extractRuleInfo(td.AuditResults, ctx)

		item := utils.AuditSQLItem{
			Number:      td.Number,
			SQL:         td.ExecSQL,
			AuditLevel:  td.AuditLevel,
			AuditStatus: tempSQL.GetAuditStatusDesc(ctx),
			AuditResult: tempSQL.GetAuditResultDesc(ctx),
			ExecStatus:  tempSQL.GetExecStatusDesc(ctx),
			ExecResult:  td.ExecResult,
			RollbackSQL: strings.Join(rollbackSqlMap[td.Id], "\n"),
			Description: td.Description,
			RuleName:    ruleName,
			Suggestion:  suggestion,
		}
		sqlList = append(sqlList, item)

		// 统计等级分布
		level := td.AuditLevel
		if level == "" {
			level = "normal"
		}
		levelDist[level]++

		// 区分问题 SQL（AuditLevel 不是 normal 且不为空）
		if level != "normal" {
			problemSQLs = append(problemSQLs, item)
		}

		// 统计规则命中
		for _, ar := range td.AuditResults {
			if ar.RuleName != "" {
				ruleHits[ar.RuleName]++
			}
		}
	}

	// 4. 构建国际化标签（当前使用 locale 包提供的 i18n 标签）
	labels := buildReportLabels(ctx)

	// 5. 构建审核时间
	auditTime := time.Now().Format("2006-01-02 15:04:05")
	if task.CreatedAt.Year() > 1 {
		auditTime = task.CreatedAt.Format("2006-01-02 15:04:05")
	}

	return &utils.AuditReportData{
		TaskID:       uint64(task.ID),
		Title:        locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelTitle),
		InstanceName: task.InstanceName(),
		Schema:       task.Schema,
		GeneratedAt:  time.Now(),
		Lang:         locale.Bundle.GetLangTagFromCtx(ctx).String(),
		LogoBase64:   "",
		Summary: utils.AuditSummary{
			AuditTime:    auditTime,
			InstanceName: task.InstanceName(),
			Schema:       task.Schema,
			TotalSQL:     len(sqlList),
			PassRate:     task.PassRate,
			Score:        task.Score,
			AuditLevel:   task.AuditLevel,
		},
		Statistics: utils.AuditStatistics{
			LevelDistribution: toLevelCounts(levelDist),
			RuleHits:          toRuleHits(ruleHits),
		},
		SQLList:     sqlList,
		ProblemSQLs: problemSQLs,
		Labels:      labels,
	}, nil
}

// extractRuleInfo 从审核结果中提取规则名称和审核建议。
// 如果有多条规则命中，使用逗号分隔拼接。
func extractRuleInfo(auditResults model.AuditResults, ctx context.Context) (ruleName string, suggestion string) {
	if len(auditResults) == 0 {
		return "", ""
	}

	lang := locale.Bundle.GetLangTagFromCtx(ctx)
	var ruleNames []string
	var suggestions []string

	for _, ar := range auditResults {
		if ar.RuleName != "" {
			ruleNames = append(ruleNames, ar.RuleName)
		}
		msg := ar.GetAuditMsgByLangTag(lang)
		if msg != "" {
			suggestions = append(suggestions, msg)
		}
	}

	return strings.Join(ruleNames, ", "), strings.Join(suggestions, "; ")
}

// toLevelCounts 将等级分布 map 转换为有序的 LevelCount 切片。
// 按 error > warn > notice > normal 顺序排列。
func toLevelCounts(dist map[string]int) []utils.LevelCount {
	if len(dist) == 0 {
		return []utils.LevelCount{}
	}

	levelOrder := map[string]int{
		"error":  0,
		"warn":   1,
		"notice": 2,
		"normal": 3,
	}

	result := make([]utils.LevelCount, 0, len(dist))
	for level, count := range dist {
		result = append(result, utils.LevelCount{
			Level: level,
			Count: count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		oi, ok := levelOrder[result[i].Level]
		if !ok {
			oi = 99
		}
		oj, ok := levelOrder[result[j].Level]
		if !ok {
			oj = 99
		}
		return oi < oj
	})

	return result
}

// toRuleHits 将规则命中 map 转换为按命中次数降序排列的 RuleHit 切片。
func toRuleHits(hits map[string]int) []utils.RuleHit {
	if len(hits) == 0 {
		return []utils.RuleHit{}
	}

	result := make([]utils.RuleHit, 0, len(hits))
	for name, count := range hits {
		result = append(result, utils.RuleHit{
			RuleName: name,
			HitCount: count,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].HitCount > result[j].HitCount
	})

	return result
}

// buildReportLabels 构建报告中使用的国际化标签。
// 当前版本使用 locale 包已有的国际化消息和硬编码中文标签，
// 后续阶段 8 将接入 go-i18n 框架实现完整国际化。
func buildReportLabels(ctx context.Context) utils.ReportLabels {
	return utils.ReportLabels{
		AuditSummary:      locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelAuditSummary),
		ResultStatistics:  locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelResultStatistics),
		ProblemSQLList:    locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelProblemSQLList),
		RuleHitStatistics: locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelRuleHitStatistics),
		AuditTime:         locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelAuditTime),
		DataSource:        locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelDataSource),
		Schema:            locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelSchema),
		TotalSQL:          locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelTotalSQL),
		PassRate:          locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelPassRate),
		Score:             locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelScore),
		AuditLevel:        locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelAuditLevel),
		Number:            locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportIndex),
		SQL:               locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportSQL),
		AuditStatus:       locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportAuditStatus),
		AuditResult:       locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportAuditResult),
		ExecStatus:        locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportExecStatus),
		ExecResult:        locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportExecResult),
		RollbackSQL:       locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportRollbackSQL),
		RuleName:          locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelRuleName),
		Description:       locale.Bundle.LocalizeMsgByCtx(ctx, locale.TaskSQLReportDescription),
		Suggestion:        locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelSuggestion),
		Count:             locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelCount),
		HitCount:          locale.Bundle.LocalizeMsgByCtx(ctx, locale.ReportLabelHitCount),
	}
}
