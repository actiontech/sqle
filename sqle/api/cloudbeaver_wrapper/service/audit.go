package service

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver"
	"github.com/actiontech/sqle/sqle/log"
	sqleModel "github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/server"
)

func AuditSQL(sql string, connectionID string) (auditSuccess bool, result string, err error) {
	if sql == "" || connectionID == "" {
		return true, "", nil
	}

	// 获取SQLE实例
	s := sqleModel.GetStorage()
	cache, err := s.GetCloudBeaverInstanceCacheByCBInstIDs([]string{connectionID})
	if err != nil {
		return false, "", err
	}

	// 找不到这个实例的缓存表示这个实例不在SQLE管理范围内
	if len(cache) == 0 {
		return true, "", nil
	}

	// 找不到sqle实例可能是实例被删除后没更新的脏数据
	inst, exist, err := s.GetInstanceById(fmt.Sprintf("%v", cache[0].SQLEInstanceID))
	if err != nil || !exist {
		return false, "", err
	}

	// 不用审核的实例跳过审核
	if !inst.SqlQueryConfig.AuditEnabled {
		return true, "", nil
	}

	task, err := server.AuditSQLByDBType(log.NewEntry(), sql, inst.DbType, "")
	if err != nil {
		return false, "", err
	}

	if driver.RuleLevel(task.AuditLevel).MoreOrEqual(driver.RuleLevel(inst.SqlQueryConfig.AllowQueryWhenLessThanAuditLevel)) {
		return false, generateAuditResult(task), nil
	}

	return true, "", nil
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
