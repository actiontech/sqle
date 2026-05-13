package hive

import (
	"context"
	databaseDriver "database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	sqleDriver "github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/beltran/gohive"
	"github.com/sirupsen/logrus"
)

func init() {
	sqleDriver.BuiltInPluginProcessors[driverV2.DriverTypeHive] = &PluginProcessor{}
}

// PluginProcessor implements driver.PluginProcessor for Hive.
type PluginProcessor struct{}

// HiveDriverImpl implements driver.Plugin for Hive.
type HiveDriverImpl struct {
	log  *logrus.Entry
	dsn  *driverV2.DSN
	conn *gohive.Connection
}

func (p *PluginProcessor) GetDriverMetas() (*driverV2.DriverMetas, error) {
	return &driverV2.DriverMetas{
		PluginName:               driverV2.DriverTypeHive,
		DatabaseDefaultPort:      10000,
		Logo:                     logo,
		DatabaseAdditionalParams: additionalParams(),
		Rules:                    []*driverV2.Rule{},
		EnabledOptionalModule:    []driverV2.OptionalModule{},
	}, nil
}

func (p *PluginProcessor) Open(l *logrus.Entry, cfg *driverV2.Config) (sqleDriver.Plugin, error) {
	impl := &HiveDriverImpl{
		log: l,
	}
	if cfg.DSN != nil {
		impl.dsn = cfg.DSN
		conn, err := newHiveConnection(cfg.DSN)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to Hive: %v", err)
		}
		impl.conn = conn
	}
	return impl, nil
}

func (p *PluginProcessor) Stop() error {
	return nil
}

func additionalParams() params.Params {
	return params.Params{
		{
			Key:   "auth",
			Value: "NOSASL",
			Desc:  "authentication mode",
			Type:  params.ParamTypeString,
			Enums: []params.EnumsValue{
				{Value: "NOSASL", Desc: "No authentication"},
				{Value: "NONE", Desc: "No authentication (SASL)"},
				{Value: "LDAP", Desc: "LDAP authentication"},
				{Value: "KERBEROS", Desc: "Kerberos authentication"},
			},
		},
		{
			Key:   "transport_mode",
			Value: "binary",
			Desc:  "transport mode (binary or http)",
			Type:  params.ParamTypeString,
			Enums: []params.EnumsValue{
				{Value: "binary", Desc: "Binary transport (default)"},
				{Value: "http", Desc: "HTTP transport"},
			},
		},
	}
}

// newHiveConnection creates a gohive connection from DSN parameters.
// It reads host, port, user, password, database from DSN and auth/transport_mode
// from AdditionalParams. This follows the same approach as DMS-EE's NewHiveConn.
func newHiveConnection(dsn *driverV2.DSN) (*gohive.Connection, error) {
	port, err := strconv.Atoi(dsn.Port)
	if err != nil {
		return nil, fmt.Errorf("invalid port %q: %v", dsn.Port, err)
	}

	conf := gohive.NewConnectConfiguration()
	conf.Username = dsn.User
	conf.Password = dsn.Password
	if dsn.DatabaseName != "" {
		conf.Database = dsn.DatabaseName
	}

	auth := "NOSASL"
	if dsn.AdditionalParams != nil {
		if authParam := dsn.AdditionalParams.GetParam("auth"); authParam != nil {
			if v := authParam.String(); v != "" {
				auth = v
			}
		}
		if transportParam := dsn.AdditionalParams.GetParam("transport_mode"); transportParam != nil {
			if v := transportParam.String(); v != "" {
				conf.TransportMode = v
			}
		}
		if serviceParam := dsn.AdditionalParams.GetParam("service"); serviceParam != nil {
			if v := serviceParam.String(); v != "" {
				conf.Service = v
			}
		}
	}

	conn, err := gohive.Connect(dsn.Host, port, auth, conf)
	if err != nil {
		return nil, fmt.Errorf("gohive connect failed: %v", err)
	}
	return conn, nil
}

// Ping tests the connectivity to the Hive server by executing SELECT 1.
func (h *HiveDriverImpl) Ping(ctx context.Context) error {
	if h.conn == nil {
		return fmt.Errorf("hive connection is not initialized")
	}
	cursor := h.conn.Cursor()
	cursor.Exec(ctx, "SELECT 1")
	defer cursor.Close()
	if cursor.Err != nil {
		return fmt.Errorf("hive ping failed: %v", cursor.Err)
	}
	return nil
}

// Parse parses sqlText into Node array. It uses keyword prefix matching
// to classify SQL statements as DQL/DML/DDL.
func (h *HiveDriverImpl) Parse(ctx context.Context, sqlText string) ([]driverV2.Node, error) {
	sqls := splitSQL(sqlText)
	nodes := make([]driverV2.Node, 0, len(sqls))
	for _, sql := range sqls {
		sqlType := classifySQL(sql)
		nodes = append(nodes, driverV2.Node{
			Text:        sql,
			Type:        sqlType,
			Fingerprint: sql,
		})
	}
	return nodes, nil
}

// Audit performs SQL audit. Currently returns empty results (no audit rules)
// as per design requirement TC-02.
func (h *HiveDriverImpl) Audit(ctx context.Context, sqls []string) ([]*driverV2.AuditResults, error) {
	results := make([]*driverV2.AuditResults, len(sqls))
	for i := range sqls {
		results[i] = &driverV2.AuditResults{}
	}
	return results, nil
}

func (h *HiveDriverImpl) Close(ctx context.Context) {
	if h.conn != nil {
		h.conn.Close()
	}
}

func (h *HiveDriverImpl) Exec(ctx context.Context, query string) (databaseDriver.Result, error) {
	return nil, fmt.Errorf("hive plugin does not support Exec")
}

func (h *HiveDriverImpl) ExecBatch(ctx context.Context, sqls ...string) ([]databaseDriver.Result, error) {
	return nil, fmt.Errorf("hive plugin does not support ExecBatch")
}

func (h *HiveDriverImpl) Tx(ctx context.Context, queries ...string) (*driverV2.TxResponse, error) {
	return nil, fmt.Errorf("hive plugin does not support Tx")
}

func (h *HiveDriverImpl) Query(ctx context.Context, sql string, conf *driverV2.QueryConf) (*driverV2.QueryResult, error) {
	return nil, fmt.Errorf("hive plugin does not support Query")
}

func (h *HiveDriverImpl) Explain(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainResult, error) {
	return nil, fmt.Errorf("hive plugin does not support Explain")
}

func (h *HiveDriverImpl) ExplainJSONFormat(ctx context.Context, conf *driverV2.ExplainConf) (*driverV2.ExplainJSONResult, error) {
	return nil, fmt.Errorf("hive plugin does not support ExplainJSONFormat")
}

func (h *HiveDriverImpl) GenRollbackSQL(ctx context.Context, sql string) (string, i18nPkg.I18nStr, error) {
	return "", nil, nil
}

func (h *HiveDriverImpl) KillProcess(ctx context.Context) error {
	return fmt.Errorf("hive plugin does not support KillProcess")
}

func (h *HiveDriverImpl) Schemas(ctx context.Context) ([]string, error) {
	return nil, fmt.Errorf("hive plugin does not support Schemas")
}

func (h *HiveDriverImpl) GetTableMetaBySQL(ctx context.Context, conf *sqleDriver.GetTableMetaBySQLConf) (*sqleDriver.GetTableMetaBySQLResult, error) {
	return nil, fmt.Errorf("hive plugin does not support GetTableMetaBySQL")
}

func (h *HiveDriverImpl) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, fmt.Errorf("hive plugin does not support EstimateSQLAffectRows")
}

func (h *HiveDriverImpl) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabaseSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	return nil, fmt.Errorf("hive plugin does not support GetDatabaseObjectDDL")
}

func (h *HiveDriverImpl) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	return nil, fmt.Errorf("hive plugin does not support GetDatabaseDiffModifySQL")
}

func (h *HiveDriverImpl) Backup(ctx context.Context, backupStrategy string, sql string, backupMaxRows uint64) ([]string, string, error) {
	return nil, "", fmt.Errorf("hive plugin does not support Backup")
}

func (h *HiveDriverImpl) RecommendBackupStrategy(ctx context.Context, sql string) (*sqleDriver.RecommendBackupStrategyRes, error) {
	return nil, fmt.Errorf("hive plugin does not support RecommendBackupStrategy")
}

func (h *HiveDriverImpl) GetSelectivityOfSQLColumns(ctx context.Context, sql string) (map[string]map[string]float32, error) {
	return nil, fmt.Errorf("hive plugin does not support GetSelectivityOfSQLColumns")
}

// splitSQL splits SQL text by semicolons and filters out empty statements.
func splitSQL(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// classifySQL classifies a SQL statement by its keyword prefix.
// Returns driverV2.SQLTypeDQL, driverV2.SQLTypeDML, or driverV2.SQLTypeDDL.
func classifySQL(sql string) string {
	upper := strings.ToUpper(strings.TrimSpace(sql))

	switch {
	case strings.HasPrefix(upper, "SELECT"),
		strings.HasPrefix(upper, "WITH"),
		strings.HasPrefix(upper, "SHOW"),
		strings.HasPrefix(upper, "DESCRIBE"),
		strings.HasPrefix(upper, "DESC"),
		strings.HasPrefix(upper, "EXPLAIN"):
		return driverV2.SQLTypeDQL
	case strings.HasPrefix(upper, "INSERT"),
		strings.HasPrefix(upper, "UPDATE"),
		strings.HasPrefix(upper, "DELETE"),
		strings.HasPrefix(upper, "MERGE"),
		strings.HasPrefix(upper, "LOAD"),
		strings.HasPrefix(upper, "EXPORT"):
		return driverV2.SQLTypeDML
	default:
		return driverV2.SQLTypeDDL
	}
}
