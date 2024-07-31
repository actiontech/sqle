package auditplan

import (
	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/sirupsen/logrus"
)

type AuditResultResp struct {
	AuditPlanID  uint64
	Task         *model.Task
	FilteredSqls []*model.OriginManageSQL
}

type Head struct {
	Name string
	Desc string
	Type string
}

// todo: 弃用
type SQL struct {
	SQLContent  string
	Fingerprint string
	Schema      string
	Info        map[string]interface{}
}

type AuditPlanMeta interface {
	InstanceType() string
	Params(instanceId ...string) params.Params
	Metrics() []string
}

type AuditPlanCollector interface {
	ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) // 处理 server 采集的数据
}

type AuditPlanHandler interface {
	AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error // 数据聚合
	Audit([]*model.OriginManageSQL) (*AuditResultResp, error)
	GetSQLs(ap *AuditPlan, persist *model.Storage, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error)
}

func auditSQLs(sqls []*model.OriginManageSQL) (*AuditResultResp, error) {
	logger := log.NewEntry()
	persist := model.GetStorage()
	if len(sqls) == 0 {
		return nil, errNoSQLNeedToBeAudited
	}
	// 同一批sql都属于同一个任务
	auditPlanID := sqls[0].SourceId

	auditPlan, err := dms.GetAuditPlansWithInstanceV2(auditPlanID, persist.GetAuditPlanDetailByID)
	if err != nil {
		return nil, err
	}
	schema := sqls[0].SchemaName
	if schema == "" {
		schema = auditPlan.SchemaName
	}

	task := &model.Task{Instance: auditPlan.Instance, Schema: schema, DBType: auditPlan.DBType}

	for i, sql := range sqls {
		task.ExecuteSQLs = append(task.ExecuteSQLs, &model.ExecuteSQL{
			BaseSQL: model.BaseSQL{
				Number:  uint(i),
				Content: sql.SqlText,
			},
		})
	}
	projectId := model.ProjectUID(auditPlan.ProjectId)
	err = server.Audit(logger, task, &projectId, auditPlan.RuleTemplateName)
	if err != nil {
		return nil, err
	}

	// update sql audit result
	for i, sql := range task.ExecuteSQLs {
		sqls[i].AuditResults = sql.AuditResults
		sqls[i].AuditLevel = sql.AuditLevel
	}

	return &AuditResultResp{
		AuditPlanID:  uint64(auditPlan.ID),
		Task:         task,
		FilteredSqls: sqls,
	}, nil
}
