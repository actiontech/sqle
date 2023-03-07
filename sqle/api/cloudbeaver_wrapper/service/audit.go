package service

import (
	"fmt"
	"strings"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	sqleModel "github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
)

type AuditResult struct {
	Result     string
	LimitLevel string
	AuditLevel string
}

func AuditSQL(sql string, connectionID string) (auditSuccess bool, result *AuditResult, err error) {
	if sql == "" || connectionID == "" {
		return true, nil, nil
	}

	// 获取SQLE实例
	s := sqleModel.GetStorage()
	cache, err := s.GetCloudBeaverInstanceCacheByCBInstIDs([]string{connectionID})
	if err != nil {
		return false, nil, err
	}

	// 找不到这个实例的缓存表示这个实例不在SQLE管理范围内
	if len(cache) == 0 {
		return true, nil, nil
	}

	// 找不到sqle实例可能是实例被删除后没更新的脏数据
	inst, exist, err := s.GetInstanceById(fmt.Sprintf("%v", cache[0].SQLEInstanceID))
	if err != nil || !exist {
		return false, nil, err
	}

	// 不用审核的实例跳过审核
	if !inst.SqlQueryConfig.AuditEnabled {
		return true, nil, nil
	}

	ruleTemplate, exist, err := s.GetRuleTemplatesByInstanceNameAndProjectId(inst.Name, inst.ProjectId)
	if err != nil {
		return false, nil, err
	}
	ruleTemplateName := ""
	var projectId uint
	if exist {
		ruleTemplateName = ruleTemplate.Name
		projectId = ruleTemplate.ProjectId
	}

	task, err := server.AuditSQLByDBType(log.NewEntry(), sql, inst.DbType, &projectId, ruleTemplateName)
	if err != nil {
		return false, nil, err
	}

	if driverV2.RuleLevel(task.AuditLevel).More(driverV2.RuleLevel(inst.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel)) {
		return false, &AuditResult{
			Result:     generateAuditResult(task),
			LimitLevel: inst.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel,
			AuditLevel: task.AuditLevel,
		}, nil
	}

	return true, nil, nil
}

func generateAuditResult(task *sqleModel.Task) string {
	builder := strings.Builder{}
	for _, executeSQL := range task.ExecuteSQLs {
		builder.WriteString(strings.TrimSpace(executeSQL.Content))
		builder.WriteString("\n------------------------------------------------------------------------------------------------\n")
		builder.WriteString(executeSQL.AuditResult)
		builder.WriteString("\n------------------------------------------------------------------------------------------------\n\n")
	}
	return builder.String()
}
