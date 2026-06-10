package server

import (
	"context"
	"database/sql/driver"
	"errors"
	"testing"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	sqleDriver "github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/stretchr/testify/assert"
)

type auditFallbackPlugin struct {
	nodes        []driverV2.Node
	parseErr     error
	auditErr     error
	auditResults []*driverV2.AuditResults
	auditCalls   [][]string
}

func (m *auditFallbackPlugin) Close(ctx context.Context) {}
func (m *auditFallbackPlugin) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}
	return m.nodes, nil
}
func (m *auditFallbackPlugin) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	m.auditCalls = append(m.auditCalls, append([]string{}, sqls...))
	if m.auditErr != nil {
		return nil, m.auditErr
	}
	return m.auditResults, nil
}
func (m *auditFallbackPlugin) GenRollbackSQL(ctx context.Context, sql string) (string, i18nPkg.I18nStr, error) {
	return "", nil, nil
}
func (m *auditFallbackPlugin) Ping(ctx context.Context) error { return nil }
func (m *auditFallbackPlugin) Exec(ctx context.Context, query string) (driver.Result, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) ExecBatch(ctx context.Context, sqls ...string) ([]driver.Result, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) Tx(ctx context.Context, queries ...string) (*driverV2.TxResponse, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) ExplainJSONFormat(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainJSONResult, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) KillProcess(ctx context.Context) error         { return nil }
func (m *auditFallbackPlugin) Schemas(ctx context.Context) ([]string, error) { return nil, nil }
func (m *auditFallbackPlugin) GetTableMetaBySQL(ctx context.Context, conf *sqleDriver.GetTableMetaBySQLConf) (*sqleDriver.GetTableMetaBySQLResult, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabaseSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) Backup(ctx context.Context, backupStrategy string, sql string, backupMaxRows uint64) ([]string, string, error) {
	return nil, "", nil
}
func (m *auditFallbackPlugin) RecommendBackupStrategy(ctx context.Context, sql string) (*sqleDriver.RecommendBackupStrategyRes, error) {
	return nil, nil
}
func (m *auditFallbackPlugin) GetSelectivityOfSQLColumns(ctx context.Context, sql string) (map[string]map[string]float32, error) {
	return nil, nil
}

func TestConvertSQLsToTaskFallback(t *testing.T) {
	tests := map[string]struct {
		sql     string
		plugin  *auditFallbackPlugin
		wantSQL []string
	}{
		"blank returns empty task": {
			sql:    " \n\t ",
			plugin: &auditFallbackPlugin{parseErr: errors.New("must not matter")},
		},
		"parse error keeps raw sql": {
			sql:     "select from ;",
			plugin:  &auditFallbackPlugin{parseErr: errors.New("parse failed")},
			wantSQL: []string{"select from ;"},
		},
		"zero nodes keeps raw sql": {
			sql:     "select from ;",
			plugin:  &auditFallbackPlugin{},
			wantSQL: []string{"select from ;"},
		},
		"parse success uses nodes": {
			sql: "select 1; select 2;",
			plugin: &auditFallbackPlugin{nodes: []driverV2.Node{
				{Text: "select 1"},
				{Text: "select 2"},
			}},
			wantSQL: []string{"select 1", "select 2"},
		},
		"partial parse keeps unsupported fragments": {
			sql: "select 1; select from ; select 2;",
			plugin: &auditFallbackPlugin{nodes: []driverV2.Node{
				{Text: "select 1"},
				{Text: "select 2"},
			}},
			wantSQL: []string{"select 1", "select from", "select 2"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			task, err := convertSQLsToTask(tt.sql, tt.plugin)
			assert.NoError(t, err)
			assert.Len(t, task.ExecuteSQLs, len(tt.wantSQL))
			for i, wantSQL := range tt.wantSQL {
				assert.Equal(t, uint(i+1), task.ExecuteSQLs[i].Number)
				assert.Equal(t, wantSQL, task.ExecuteSQLs[i].Content)
			}
		})
	}
}

func TestAuditSQLsOneByOneFallback(t *testing.T) {
	plugin := &auditFallbackPlugin{
		auditErr: errors.New("batch failed"),
	}
	results := auditSQLsOneByOne(log.NewEntry(), plugin, []string{"select from ;", "select 1"})
	assert.Len(t, results, 2)
	assert.Equal(t, driverV2.RuleLevelWarn, results[0].Level())
	assert.Contains(t, results[0].Message(), unsupportedSQLWarnMessage.GetStrInLang(i18nPkg.DefaultLang))
	assert.Equal(t, driverV2.RuleLevelWarn, results[1].Level())
	assert.Equal(t, [][]string{{"select from ;"}, {"select 1"}}, plugin.auditCalls)
}

func TestReplenishTaskStatisticsEmptyTask(t *testing.T) {
	task := &model.Task{}
	ReplenishTaskStatistics(task)
	assert.Equal(t, model.TaskStatusAudited, task.Status)
	assert.Zero(t, task.Score)
	assert.Zero(t, task.PassRate)
	assert.Equal(t, string(driverV2.RuleLevelNull), task.AuditLevel)
}
