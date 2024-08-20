package auditplan

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/model"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
)

type TiDBAuditLogTaskV2 struct {
	*DefaultTaskV2
}

func NewTiDBAuditLogTaskV2Fn() func() interface{} {
	return func() interface{} {
		return &TiDBAuditLogTaskV2{}
	}
}

/*
	{
		Type:         TypeTiDBAuditLog,
		Desc:         "TiDB审计日志",
		InstanceType: InstanceTypeTiDB,
		CreateTask:   NewTaskWrap(&TiDBAuditLogTaskV2{}),
		Handler:      &TiDBAuditLogTaskV2{},
		Params: []*params.Param{
			{
				Key:   paramKeyAuditSQLsScrappedInLastPeriodMinute,
				Desc:  "审核过去时间段内抓取的SQL（分钟）",
				Value: "0",
				Type:  params.ParamTypeInt,
			},
		},
	},
*/

func (at *TiDBAuditLogTaskV2) InstanceType() string {
	return InstanceTypeTiDB
}

func (at *TiDBAuditLogTaskV2) Params(instanceId ...string) params.Params {
	return []*params.Param{}
}

func (at *TiDBAuditLogTaskV2) HighPriorityParams() params.ParamsWithOperator {
	return []*params.ParamWithOperator{
		defaultAuditLevelOperateParams,
	}
}

// todo: tidb 审核部分与其他类型的不太一样
func (at *TiDBAuditLogTaskV2) Audit(sqls []*model.SQLManageRecord) (*AuditResultResp, error) {
	return nil, nil
	// var task *model.Task
	// if at.ap.InstanceName == "" {
	// 	task = &model.Task{
	// 		DBType: at.ap.DBType,
	// 	}
	// } else {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	// 	defer cancel()

	// 	instance, _, err := dms.GetInstanceInProjectByName(ctx, string(at.ap.ProjectId), at.ap.InstanceName)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	task = &model.Task{
	// 		Instance: instance,
	// 		Schema:   at.ap.InstanceDatabase,
	// 		DBType:   at.ap.DBType,
	// 	}
	// }

	// auditPlanSQLs, err := at.persist.GetAuditPlanSQLsV2Unaudit(at.ap.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// if len(auditPlanSQLs) == 0 {
	// 	return nil, errNoSQLInAuditPlan
	// }

	// filteredSqls, err := filterSQLsByPeriodV2(at.ap.Params, auditPlanSQLs)
	// if err != nil {
	// 	return nil, err
	// }

	// if len(filteredSqls) == 0 {
	// 	return nil, ErrNoSQLNeedToBeAudited
	// }

	// for i, sql := range filteredSqls {
	// 	schema := ""
	// 	info, _ := sql.Info.OriginValue()
	// 	if schemaStr, ok := info[server.AuditSchema].(string); ok {
	// 		schema = schemaStr
	// 	}

	// 	executeSQL := &model.ExecuteSQL{
	// 		BaseSQL: model.BaseSQL{
	// 			Number:  uint(i),
	// 			Content: sql.SqlText,
	// 			Schema:  schema,
	// 		},
	// 	}

	// 	task.ExecuteSQLs = append(task.ExecuteSQLs, executeSQL)
	// }

	// hook := &TiDBAuditHook{
	// 	originalSQLs: map[*model.ExecuteSQL]string{},
	// }
	// projectId := model.ProjectUID(at.ap.ProjectId)
	// err = server.HookAudit(at.logger, task, hook, &projectId, at.ap.RuleTemplateName)
	// if err != nil {
	// 	return nil, err
	// }

	// return &AuditResultResp{
	// 	AuditPlanID:  uint64(at.ap.ID),
	// 	Task:         task,
	// 	FilteredSqls: filteredSqls,
	// }, nil
}

// 审核前填充上缺失的schema, 审核后还原被审核SQL, 并添加注释说明sql在哪个库执行的
type TiDBAuditHook struct {
	originalSQLs map[*model.ExecuteSQL]string
}

func (t *TiDBAuditHook) BeforeAudit(sql *model.ExecuteSQL) {
	if sql.Schema == "" {
		return
	}
	t.originalSQLs[sql] = sql.Content
	newSQL, err := tidbCompletionSchema(sql.Content, sql.Schema)
	if err != nil {
		return
	}
	sql.Content = newSQL
}

func (t *TiDBAuditHook) AfterAudit(sql *model.ExecuteSQL) {
	if sql.Schema == "" {
		return
	}
	if o, ok := t.originalSQLs[sql]; ok {
		sql.Content = fmt.Sprintf("%v -- current schema: %v", o, sql.Schema)
	}
}

// 填充sql缺失的schema
func tidbCompletionSchema(sql, schema string) (string, error) {
	stmts, _, err := parser.New().PerfectParse(sql, "", "")
	if err != nil {
		return "", err
	}
	if len(stmts) != 1 {
		return "", parser.ErrSyntax
	}

	stmts[0].Accept(&completionSchemaVisitor{schema: schema})
	buf := new(bytes.Buffer)
	restoreCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, buf)
	err = stmts[0].Restore(restoreCtx)
	return buf.String(), err
}

// completionSchemaVisitor implements ast.Visitor interface.
type completionSchemaVisitor struct {
	schema string
}

func (g *completionSchemaVisitor) Enter(n ast.Node) (node ast.Node, skipChildren bool) {
	if stmt, ok := n.(*ast.TableName); ok {
		if stmt.Schema.L == "" {
			stmt.Schema.L = strings.ToLower(g.schema)
			stmt.Schema.O = g.schema
		}
	}
	return n, false
}

func (g *completionSchemaVisitor) Leave(n ast.Node) (node ast.Node, ok bool) {
	return n, true
}
