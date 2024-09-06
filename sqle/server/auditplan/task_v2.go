package auditplan

import (
	"context"

	"github.com/actiontech/sqle/sqle/dms"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/actiontech/sqle/sqle/server"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sirupsen/logrus"
)

type AuditResultResp struct {
	AuditPlanID uint64
	Task        *model.Task
	AuditedSqls []*model.SQLManageRecord
}

type Head struct {
	Name     string
	Desc     *i18n.Message
	Type     string
	Sortable bool
}

// checkAndGetOrderByName: 对传入的order by 进行校验, 如果非定义的类型则返回空
func checkAndGetOrderByName(head []Head, orderByName string) string {
	for _, v := range head {
		if v.Sortable && v.Name == orderByName {
			return orderByName
		}
	}
	return ""
}

type FilterInputType string

const (
	FilterInputTypeInt      FilterInputType = "int"
	FilterInputTypeString   FilterInputType = "string"
	FilterInputTypeDateTime FilterInputType = "date_time"
)

type FilterOpType string

const (
	FilterOpTypeEqual   = "equal"
	FilterOpTypeBetween = "between"
)

type FilterMeta struct {
	Name            string
	Desc            *i18n.Message
	FilterInputType FilterInputType
	FilterOpType    FilterOpType
	FilterTips      []FilterTip
}

type FilterTip struct {
	Value string `json:"value"`
	Desc  string `json:"desc"`
	Group string `json:"group"`
}

type Filter struct {
	Name                  string             `json:"filter_name"`
	FilterComparisonValue string             `json:"filter_compare_value"`
	FilterBetweenValue    FilterBetweenValue `json:"filter_between_value"`
}

type FilterBetweenValue struct {
	From string
	To   string
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
	HighPriorityParams() params.ParamsWithOperator
	Metrics() []string
}

type AuditPlanCollector interface {
	ExtractSQL(logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) ([]*SQLV2, error) // 处理 server 采集的数据
}

type AuditPlanHandler interface {
	AggregateSQL(cache SQLV2Cacher, sql *SQLV2) error // 数据聚合
	Audit([]*model.SQLManageRecord) (*AuditResultResp, error)
	// GetSQLs(ap *AuditPlan, persist *model.Storage, args map[string]interface{}) ([]Head, []map[string] /* head name */ string, uint64, error) // todo: 弃用

	// todo: 放到 AuditPlanMeta 里, 原因是meta里未保存 AuditPlanMeta, 其他开发者也在改造等合并后在处理。
	Head(ap *AuditPlan) []Head

	// todo: 放到 AuditPlanMeta 里, 原因是meta里未保存 AuditPlanMeta, 其他开发者也在改造等合并后在处理。
	Filters(ctx context.Context, logger *logrus.Entry, ap *AuditPlan, persist *model.Storage) []FilterMeta
	GetSQLData(ap *AuditPlan, persist *model.Storage, filters []Filter, orderBy string, isAsc bool, limit, offset int) ([]map[string] /* head name */ string, uint64, error)
}

func auditSQLs(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	logger := log.NewEntry()
	persist := model.GetStorage()
	if len(sqls) == 0 {
		return nil, ErrNoSQLNeedToBeAudited
	}
	// 同一批sql都属于同一个任务
	instAuditPlanID := sqls[0].SourceId
	auditPlanType := sqls[0].Source
	auditPlan, err := dms.GetAuditPlansWithInstanceV2(instAuditPlanID, auditPlanType, persist.GetAuditPlanDetailByInstAuditPlanIdAndType)
	if err != nil {
		return nil, err
	}
	schema := sqls[0].SchemaName

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
		AuditPlanID: uint64(auditPlan.ID),
		Task:        task,
		AuditedSqls: sqls,
	}, nil
}
