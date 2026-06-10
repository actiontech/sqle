package v1

import (
	"context"
	"database/sql/driver"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	sqleDriver "github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockParsePlugin struct {
	nodes    []driverV2.Node
	parseErr error
}

func (m *mockParsePlugin) Close(ctx context.Context) {}
func (m *mockParsePlugin) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	if m.parseErr != nil {
		return nil, m.parseErr
	}
	return m.nodes, nil
}
func (m *mockParsePlugin) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	return nil, nil
}
func (m *mockParsePlugin) GenRollbackSQL(ctx context.Context, sql string) (string, i18nPkg.I18nStr, error) {
	return "", nil, nil
}
func (m *mockParsePlugin) Ping(ctx context.Context) error { return nil }
func (m *mockParsePlugin) Exec(ctx context.Context, query string) (driver.Result, error) {
	return nil, nil
}
func (m *mockParsePlugin) ExecBatch(ctx context.Context, sqls ...string) ([]driver.Result, error) {
	return nil, nil
}
func (m *mockParsePlugin) Tx(ctx context.Context, queries ...string) (*driverV2.TxResponse, error) {
	return nil, nil
}
func (m *mockParsePlugin) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, nil
}
func (m *mockParsePlugin) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return nil, nil
}
func (m *mockParsePlugin) ExplainJSONFormat(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainJSONResult, error) {
	return nil, nil
}
func (m *mockParsePlugin) KillProcess(ctx context.Context) error         { return nil }
func (m *mockParsePlugin) Schemas(ctx context.Context) ([]string, error) { return nil, nil }
func (m *mockParsePlugin) GetTableMetaBySQL(ctx context.Context, conf *sqleDriver.GetTableMetaBySQLConf) (*sqleDriver.GetTableMetaBySQLResult, error) {
	return nil, nil
}
func (m *mockParsePlugin) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, nil
}
func (m *mockParsePlugin) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabaseSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	return nil, nil
}
func (m *mockParsePlugin) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	return nil, nil
}
func (m *mockParsePlugin) Backup(ctx context.Context, backupStrategy string, sql string, backupMaxRows uint64) ([]string, string, error) {
	return nil, "", nil
}
func (m *mockParsePlugin) RecommendBackupStrategy(ctx context.Context, sql string) (*sqleDriver.RecommendBackupStrategyRes, error) {
	return nil, nil
}
func (m *mockParsePlugin) GetSelectivityOfSQLColumns(ctx context.Context, sql string) (map[string]map[string]float32, error) {
	return nil, nil
}

func TestAddSQLsFromFileToTasksFallback(t *testing.T) {
	tests := map[string]struct {
		sqls       GetSQLFromFileResp
		plugin     *mockParsePlugin
		wantSQLs   []string
		wantStarts []uint64
	}{
		"blank form input is ignored": {
			sqls:   GetSQLFromFileResp{SQLsFromFormData: " \n\t "},
			plugin: &mockParsePlugin{parseErr: errors.New("must not matter")},
		},
		"parse error keeps raw sql": {
			sqls:       GetSQLFromFileResp{SQLsFromFormData: "create table t1 (id int) dbpartition by hash(id);"},
			plugin:     &mockParsePlugin{parseErr: errors.New("parse failed")},
			wantSQLs:   []string{"create table t1 (id int) dbpartition by hash(id);"},
			wantStarts: []uint64{0},
		},
		"zero nodes keeps raw sql with source": {
			sqls: GetSQLFromFileResp{SQLsFromXMLs: []SQLFromXML{{
				FilePath:  "mapper.xml",
				StartLine: 12,
				SQL:       "select from ;",
			}}},
			plugin:     &mockParsePlugin{},
			wantSQLs:   []string{"select from ;"},
			wantStarts: []uint64{12},
		},
		"parse success uses parsed nodes": {
			sqls: GetSQLFromFileResp{SQLsFromFormData: "select 1; select 2;"},
			plugin: &mockParsePlugin{nodes: []driverV2.Node{
				{Text: "select 1", StartLine: 1, Type: "DQL"},
				{Text: "select 2", StartLine: 2, Type: "DQL"},
			}},
			wantSQLs:   []string{"select 1", "select 2"},
			wantStarts: []uint64{1, 2},
		},
		"partial parse keeps unsupported fragments": {
			sqls: GetSQLFromFileResp{SQLsFromFormData: "select 1; select from ; select 2;"},
			plugin: &mockParsePlugin{nodes: []driverV2.Node{
				{Text: "select 1", StartLine: 1, Type: "DQL"},
				{Text: "select 2", StartLine: 1, Type: "DQL"},
			}},
			wantSQLs:   []string{"select 1", "select from", "select 2"},
			wantStarts: []uint64{1, 0, 1},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			task := &model.Task{}
			err := addSQLsFromFileToTasks(tt.sqls, task, tt.plugin)
			assert.NoError(t, err)
			assert.Len(t, task.ExecuteSQLs, len(tt.wantSQLs))
			for i, wantSQL := range tt.wantSQLs {
				assert.Equal(t, uint(i+1), task.ExecuteSQLs[i].Number)
				assert.Equal(t, wantSQL, task.ExecuteSQLs[i].Content)
				assert.Equal(t, tt.wantStarts[i], task.ExecuteSQLs[i].StartLine)
			}
		})
	}
}

func TestGetSqlsFromGitEmptyURL(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	c := e.NewContext(req, httptest.NewRecorder())

	sqlFiles, javaFiles, xmls, exist, err := getSqlsFromGit(c)
	assert.NoError(t, err)
	assert.False(t, exist)
	assert.Empty(t, sqlFiles)
	assert.Empty(t, javaFiles)
	assert.Empty(t, xmls)
}
