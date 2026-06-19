package server

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
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

func TestSkipSQLRuleExceptionResultsOnlySkipsMatchedRule(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	model.InitMockStorage(mockDB)

	results := driverV2.NewAuditResults()
	results.Add(driverV2.RuleLevelError, "rule_error_excepted", i18nPkg.ConvertStr2I18nAsDefaultLang("excepted rule"))
	results.Add(driverV2.RuleLevelWarn, "rule_warn_kept", i18nPkg.ConvertStr2I18nAsDefaultLang("kept rule"))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_rule_exception` WHERE (project_id = ? AND sql_fingerprint IN (?,?) AND rule_name IN (?,?)) AND `sql_rule_exception`.`deleted_at` IS NULL")).
		WithArgs(model.ProjectUID("700300"), "select * from ?", "select * from t1", "rule_error_excepted", "rule_warn_kept").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "instance_id", "sql_fingerprint", "rule_name", "reason"}).
			AddRow(1, "700300", 1, "select * from ?", "rule_error_excepted", "business exception"))

	filteredResults, skippedResults, err := skipSQLRuleExceptionResults(log.NewEntry(), model.GetStorage(), "700300", 1, "select * from ?", "select * from t1", results)
	assert.NoError(t, err)
	if assert.Len(t, filteredResults.Results, 1) {
		assert.Equal(t, "rule_warn_kept", filteredResults.Results[0].RuleName)
		assert.Equal(t, driverV2.RuleLevelWarn, filteredResults.Level())
	}
	if assert.Len(t, skippedResults, 1) {
		assert.Equal(t, "rule_error_excepted", skippedResults[0].RuleName)
		assert.Equal(t, string(driverV2.RuleLevelError), skippedResults[0].Level)
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSkipSQLRuleExceptionResultsFallsBackToSQLTextFingerprint(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	model.InitMockStorage(mockDB)

	results := driverV2.NewAuditResults()
	results.Add(driverV2.RuleLevelError, "SQLE00008", i18nPkg.ConvertStr2I18nAsDefaultLang("table must have primary key"))
	results.Add(driverV2.RuleLevelWarn, "SQLE00033", i18nPkg.ConvertStr2I18nAsDefaultLang("kept rule"))
	results.Add(driverV2.RuleLevelWarn, "SQLE00061", i18nPkg.ConvertStr2I18nAsDefaultLang("kept rule"))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_rule_exception` WHERE (project_id = ? AND sql_fingerprint IN (?) AND rule_name IN (?,?,?)) AND `sql_rule_exception`.`deleted_at` IS NULL")).
		WithArgs(model.ProjectUID("700300"), "create table rule_exc_e2e_std_1781803945 (id int, name varchar(20))", "SQLE00008", "SQLE00033", "SQLE00061").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "instance_id", "sql_fingerprint", "rule_name", "reason"}).
			AddRow(24, "700300", 1, "create table rule_exc_e2e_std_1781803945 (id int, name varchar(20))", "SQLE00008", "business exception"))

	filteredResults, skippedResults, err := skipSQLRuleExceptionResults(log.NewEntry(), model.GetStorage(), "700300", 1, "", " create table rule_exc_e2e_std_1781803945 (id int, name varchar(20)) ", results)
	assert.NoError(t, err)
	if assert.Len(t, filteredResults.Results, 2) {
		assert.Equal(t, "SQLE00033", filteredResults.Results[0].RuleName)
		assert.Equal(t, "SQLE00061", filteredResults.Results[1].RuleName)
	}
	if assert.Len(t, skippedResults, 1) {
		assert.Equal(t, "SQLE00008", skippedResults[0].RuleName)
		assert.Equal(t, string(driverV2.RuleLevelError), skippedResults[0].Level)
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHookAuditWholeSQLWhitelistSkipsAllRulesBeforeRuleException(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	model.InitMockStorage(mockDB)

	plugin := &auditFallbackPlugin{nodes: []driverV2.Node{{Text: "select * from t1", Fingerprint: "select * from ?"}}}
	task := &model.Task{
		InstanceId:  1,
		Instance:    &model.Instance{ProjectId: "700300"},
		ExecuteSQLs: []*model.ExecuteSQL{{BaseSQL: model.BaseSQL{Content: "select * from t1"}}},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_whitelist` WHERE sql_whitelist.project_id = ? AND `sql_whitelist`.`deleted_at` IS NULL")).
		WithArgs("700300").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "value", "match_type"}).
			AddRow(9, "700300", "select * from t1", model.SQLWhitelistExactMatch))
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("UPDATE `sql_whitelist` SET `last_matched_time`=?,`matched_count`=matched_count + ? WHERE sql_whitelist.id = ? AND `sql_whitelist`.`deleted_at` IS NULL")).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = hookAudit(log.NewEntry(), task, plugin, &EmptyAuditHook{}, "700300", nil)
	assert.NoError(t, err)
	assert.Empty(t, plugin.auditCalls)
	if assert.Len(t, task.ExecuteSQLs, 1) {
		assert.Equal(t, string(driverV2.RuleLevelNormal), task.ExecuteSQLs[0].AuditLevel)
		assert.Len(t, task.ExecuteSQLs[0].AuditResults, 1)
		assert.Empty(t, task.ExecuteSQLs[0].SkippedAuditResults)
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestHookAuditWholeSQLWhitelistMissThenRuleExceptionSkipsMatchedRuleOnly(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err)
	mock.ExpectQuery("SELECT VERSION\\(\\)").WillReturnRows(sqlmock.NewRows([]string{"VERSION()"}).AddRow("5.7"))
	model.InitMockStorage(mockDB)

	results := driverV2.NewAuditResults()
	results.Add(driverV2.RuleLevelError, "rule_error_excepted", i18nPkg.ConvertStr2I18nAsDefaultLang("excepted rule"))
	results.Add(driverV2.RuleLevelWarn, "rule_warn_kept", i18nPkg.ConvertStr2I18nAsDefaultLang("kept rule"))
	plugin := &auditFallbackPlugin{
		nodes:        []driverV2.Node{{Text: "select * from t1", Fingerprint: "select * from ?"}},
		auditResults: []*driverV2.AuditResults{results},
	}
	task := &model.Task{
		InstanceId:  1,
		Instance:    &model.Instance{ProjectId: "700300"},
		ExecuteSQLs: []*model.ExecuteSQL{{BaseSQL: model.BaseSQL{Content: "select * from t1"}}},
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_whitelist` WHERE sql_whitelist.project_id = ? AND `sql_whitelist`.`deleted_at` IS NULL")).
		WithArgs("700300").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "value", "match_type"}).
			AddRow(9, "700300", "select * from t2", model.SQLWhitelistExactMatch))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sql_rule_exception` WHERE (project_id = ? AND sql_fingerprint IN (?,?) AND rule_name IN (?,?)) AND `sql_rule_exception`.`deleted_at` IS NULL")).
		WithArgs(model.ProjectUID("700300"), "select * from ?", "select * from t1", "rule_error_excepted", "rule_warn_kept").
		WillReturnRows(sqlmock.NewRows([]string{"id", "project_id", "instance_id", "sql_fingerprint", "rule_name", "reason"}).
			AddRow(1, "700300", 1, "select * from ?", "rule_error_excepted", "business exception"))

	err = hookAudit(log.NewEntry(), task, plugin, &EmptyAuditHook{}, "700300", nil)
	assert.NoError(t, err)
	assert.Equal(t, [][]string{{"select * from t1"}}, plugin.auditCalls)
	if assert.Len(t, task.ExecuteSQLs, 1) {
		assert.Equal(t, string(driverV2.RuleLevelWarn), task.ExecuteSQLs[0].AuditLevel)
		if assert.Len(t, task.ExecuteSQLs[0].AuditResults, 1) {
			assert.Equal(t, "rule_warn_kept", task.ExecuteSQLs[0].AuditResults[0].RuleName)
		}
		if assert.Len(t, task.ExecuteSQLs[0].SkippedAuditResults, 1) {
			assert.Equal(t, "rule_error_excepted", task.ExecuteSQLs[0].SkippedAuditResults[0].RuleName)
		}
	}

	mock.ExpectClose()
	assert.NoError(t, mockDB.Close())
	assert.NoError(t, mock.ExpectationsWereMet())
}
