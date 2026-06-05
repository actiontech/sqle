package server

import (
	"context"
	_driver "database/sql/driver"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/agiledragon/gomonkey"
	"github.com/stretchr/testify/assert"
)

type degradeAuditPlugin struct {
	parseFn func(context.Context, string) ([]driverV2.Node, error)
	auditFn func(context.Context, []string) ([]*driverV2.AuditResults, error)
}

func (p *degradeAuditPlugin) Close(ctx context.Context)             {}
func (p *degradeAuditPlugin) Ping(ctx context.Context) error        { return nil }
func (p *degradeAuditPlugin) KillProcess(ctx context.Context) error { return nil }
func (p *degradeAuditPlugin) Exec(ctx context.Context, query string) (_driver.Result, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) ExecBatch(ctx context.Context, queries ...string) ([]_driver.Result, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) Tx(ctx context.Context, queries ...string) (*driverV2.TxResponse, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) Schemas(ctx context.Context) ([]string, error) { return nil, nil }
func (p *degradeAuditPlugin) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	if p.parseFn != nil {
		return p.parseFn(ctx, sqlText)
	}
	return []driverV2.Node{{Text: sqlText, Fingerprint: sqlText}}, nil
}
func (p *degradeAuditPlugin) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	if p.auditFn != nil {
		return p.auditFn(ctx, sqls)
	}
	results := make([]*driverV2.AuditResults, 0, len(sqls))
	for range sqls {
		result := driverV2.NewAuditResults()
		result.Add(driverV2.RuleLevelNormal, "normal_rule", i18nPkg.I18nStr{i18nPkg.DefaultLang: "normal"})
		results = append(results, result)
	}
	return results, nil
}
func (p *degradeAuditPlugin) GenRollbackSQL(ctx context.Context, sql string) (string, i18nPkg.I18nStr, error) {
	return "", nil, nil
}
func (p *degradeAuditPlugin) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) ExplainJSONFormat(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainJSONResult, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) GetTableMetaBySQL(ctx context.Context, conf *driver.GetTableMetaBySQLConf) (*driver.GetTableMetaBySQLResult, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabaseSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) Backup(ctx context.Context, backupStrategy string, sql string, backupMaxRows uint64) ([]string, string, error) {
	return nil, "", nil
}
func (p *degradeAuditPlugin) RecommendBackupStrategy(ctx context.Context, sql string) (*driver.RecommendBackupStrategyRes, error) {
	return nil, nil
}
func (p *degradeAuditPlugin) GetSelectivityOfSQLColumns(ctx context.Context, sql string) (map[string]map[string]float32, error) {
	return nil, nil
}

func TestBuildExecuteSQLsFromSQLToleratesEmptyAndParseFailure(t *testing.T) {
	plugin := &degradeAuditPlugin{parseFn: func(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
		return nil, errors.New("parse failed")
	}}

	empty, err := BuildExecuteSQLsFromSQL(context.Background(), plugin, " \n\t ", BuildExecuteSQLsOptions{})
	assert.NoError(t, err)
	assert.Empty(t, empty)

	executeSQLs, err := BuildExecuteSQLsFromSQL(context.Background(), plugin, " bad tdsql syntax ", BuildExecuteSQLsOptions{StartNumber: 7, SourceFile: "a.sql", StartLine: 3})
	assert.NoError(t, err)
	assert.Len(t, executeSQLs, 1)
	assert.Equal(t, uint(7), executeSQLs[0].Number)
	assert.Equal(t, "bad tdsql syntax", executeSQLs[0].Content)
	assert.Equal(t, "a.sql", executeSQLs[0].SourceFile)
	assert.Equal(t, uint64(3), executeSQLs[0].StartLine)
}

func TestHookAuditDegradesParseAndBatchAuditFailures(t *testing.T) {
	patches := gomonkey.ApplyMethod(reflect.TypeOf(&model.Storage{}), "GetSqlWhitelistByProjectId", func(_ *model.Storage, _ string) ([]model.SqlWhitelist, error) {
		return nil, nil
	})
	defer patches.Reset()

	auditCalls := make([][]string, 0)
	plugin := &degradeAuditPlugin{
		parseFn: func(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
			if strings.Contains(sqlText, "bad_parse") {
				return nil, errors.New("parse failed")
			}
			return []driverV2.Node{{Text: sqlText, Fingerprint: "fp:" + sqlText}}, nil
		},
		auditFn: func(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
			auditCalls = append(auditCalls, append([]string{}, sqls...))
			if len(sqls) > 1 {
				return nil, errors.New("batch audit failed")
			}
			if strings.Contains(sqls[0], "bad_audit") {
				return nil, errors.New("single audit failed")
			}
			result := driverV2.NewAuditResults()
			result.Add(driverV2.RuleLevelNormal, "normal_rule", i18nPkg.I18nStr{i18nPkg.DefaultLang: "normal"})
			return []*driverV2.AuditResults{result}, nil
		},
	}
	task := &model.Task{ExecuteSQLs: []*model.ExecuteSQL{
		{BaseSQL: model.BaseSQL{Content: "select 1"}},
		{BaseSQL: model.BaseSQL{Content: "bad_parse"}},
		{BaseSQL: model.BaseSQL{Content: "bad_audit"}},
	}}

	err := hookAudit(log.NewEntry(), task, plugin, &EmptyAuditHook{}, "project1", nil)
	assert.NoError(t, err)
	assert.Equal(t, model.TaskStatusAudited, task.Status)
	assert.Equal(t, string(driverV2.RuleLevelNormal), task.ExecuteSQLs[0].AuditLevel)
	assert.Equal(t, string(driverV2.RuleLevelWarn), task.ExecuteSQLs[1].AuditLevel)
	assert.Equal(t, string(driverV2.RuleLevelWarn), task.ExecuteSQLs[2].AuditLevel)
	assert.Len(t, auditCalls, 3)
	assert.Equal(t, []string{"select 1", "bad_audit"}, auditCalls[0])
	assert.Equal(t, []string{"select 1"}, auditCalls[1])
	assert.Equal(t, []string{"bad_audit"}, auditCalls[2])
}

func TestReplenishTaskStatisticsWithEmptyTask(t *testing.T) {
	task := &model.Task{}
	ReplenishTaskStatistics(task)
	assert.Equal(t, model.TaskStatusAudited, task.Status)
	assert.Equal(t, float64(1), task.PassRate)
	assert.Equal(t, string(driverV2.RuleLevelNull), task.AuditLevel)
	assert.Equal(t, int32(0), task.Score)
}
