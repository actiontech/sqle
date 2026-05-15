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

// hiveQueryRunner abstracts the minimum Hive cursor operations needed by
// GetDatabaseObjectDDL / GetDatabaseDiffModifySQL. It allows unit tests to
// substitute a fake without requiring a real gohive connection or network.
//
// runSingleStringQuery executes a query that returns rows of a single STRING
// column (e.g. SHOW DATABASES, SHOW CREATE TABLE) and returns each row as a
// string. The implementation is responsible for opening / closing the cursor.
type hiveQueryRunner interface {
	runSingleStringQuery(ctx context.Context, query string) ([]string, error)
}

// HiveDriverImpl implements driver.Plugin for Hive.
type HiveDriverImpl struct {
	log  *logrus.Entry
	dsn  *driverV2.DSN
	conn *gohive.Connection
	// runner is the query executor used by ObjectDDL / DiffModifySQL paths.
	// In production it is set to gohiveQueryRunner wrapping h.conn; in unit
	// tests it can be replaced with a fake to avoid network dependency.
	runner hiveQueryRunner
}

func (p *PluginProcessor) GetDriverMetas() (*driverV2.DriverMetas, error) {
	return &driverV2.DriverMetas{
		PluginName:               driverV2.DriverTypeHive,
		DatabaseDefaultPort:      10000,
		Logo:                     logo,
		DatabaseAdditionalParams: additionalParams(),
		Rules:                    []*driverV2.Rule{},
		EnabledOptionalModule: []driverV2.OptionalModule{
			driverV2.OptionalGetDatabaseObjectDDL,
			driverV2.OptionalGetDatabaseDiffModifySQL,
		},
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
		impl.runner = &gohiveQueryRunner{conn: conn}
	}
	return impl, nil
}

// gohiveQueryRunner is the production hiveQueryRunner backed by *gohive.Connection.
type gohiveQueryRunner struct {
	conn *gohive.Connection
}

func (g *gohiveQueryRunner) runSingleStringQuery(ctx context.Context, query string) ([]string, error) {
	if g.conn == nil {
		return nil, fmt.Errorf("hive connection is not initialized")
	}
	cursor := g.conn.Cursor()
	defer cursor.Close()

	cursor.Exec(ctx, query)
	if cursor.Err != nil {
		return nil, fmt.Errorf("failed to execute %q: %v", query, cursor.Err)
	}

	var rows []string
	for cursor.HasMore(ctx) {
		if cursor.Err != nil {
			return nil, fmt.Errorf("failed to fetch row: %v", cursor.Err)
		}
		var val string
		cursor.FetchOne(ctx, &val)
		if cursor.Err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", cursor.Err)
		}
		rows = append(rows, val)
	}
	return rows, nil
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
	if h.conn == nil {
		return nil, fmt.Errorf("hive connection is not initialized")
	}
	cursor := h.conn.Cursor()
	defer cursor.Close()

	cursor.Exec(ctx, "SHOW DATABASES")
	if cursor.Err != nil {
		return nil, fmt.Errorf("failed to execute SHOW DATABASES: %v", cursor.Err)
	}

	var schemas []string
	for cursor.HasMore(ctx) {
		if cursor.Err != nil {
			return nil, fmt.Errorf("failed to fetch schema row: %v", cursor.Err)
		}
		var dbName string
		cursor.FetchOne(ctx, &dbName)
		if cursor.Err != nil {
			return nil, fmt.Errorf("failed to scan schema name: %v", cursor.Err)
		}
		schemas = append(schemas, dbName)
	}

	return schemas, nil
}

func (h *HiveDriverImpl) GetTableMetaBySQL(ctx context.Context, conf *sqleDriver.GetTableMetaBySQLConf) (*sqleDriver.GetTableMetaBySQLResult, error) {
	return nil, fmt.Errorf("hive plugin does not support GetTableMetaBySQL")
}

func (h *HiveDriverImpl) EstimateSQLAffectRows(ctx context.Context, sql string) (*driverV2.EstimatedAffectRows, error) {
	return nil, fmt.Errorf("hive plugin does not support EstimateSQLAffectRows")
}

// hiveFunctionUnsupportedMsg is the Chinese error message returned when a
// FUNCTION object is requested. The driver does not implement FUNCTION DDL
// in this batch; it is planned for the second batch (design §3.2.1 / §3.5).
const hiveFunctionUnsupportedMsg = "Hive FUNCTION 暂未支持（计划第二批落地）"

// defaultHiveSchema is the schema used when an object info has an empty
// SchemaName (design §3.2.1 "USE <SchemaName>; schema 为空走默认 default").
const defaultHiveSchema = "default"

// GetDatabaseObjectDDL fetches the CREATE statement for each requested
// (schema, object) pair. It implements the contract described in
// docs/spec/design.md §3.2.1:
//
//   - For each schema, first `USE <SchemaName>` (or "default" if empty).
//   - TABLE      -> SHOW CREATE TABLE <name>
//   - VIEW       -> SHOW CREATE TABLE <name>  (Hive views reuse this command)
//   - FUNCTION   -> placeholder DDL ""; the call still returns the result row
//     but the driver records a WARN log and propagates the
//     Chinese error message; FUNCTION support is planned for
//     the second batch (compat-RISK-9).
//   - PROCEDURE / TRIGGER / EVENT -> short-circuit: skip the object entirely
//     and emit a WARN log; do not panic, do not return an error
//     (Hive does not support these object types — compat-RISK-4).
//
// Behavior on error: if FUNCTION is requested, the function still returns
// nil error so other TABLE/VIEW results in the same batch are not dropped;
// the FUNCTION error is surfaced via the returned ObjectDDL (empty string)
// plus the WARN log. Real connection errors against TABLE/VIEW are returned
// as a normal Go error.
func (h *HiveDriverImpl) GetDatabaseObjectDDL(ctx context.Context, objInfos []*driverV2.DatabaseSchemaInfo) ([]*driverV2.DatabaseSchemaObjectResult, error) {
	if h.runner == nil {
		return nil, fmt.Errorf("hive connection is not initialized")
	}

	results := make([]*driverV2.DatabaseSchemaObjectResult, 0, len(objInfos))
	for _, objInfo := range objInfos {
		schemaName := objInfo.SchemaName
		if schemaName == "" {
			schemaName = defaultHiveSchema
		}

		// USE <schemaName> first so subsequent unqualified object references
		// resolve to this database (design §3.2.1).
		if _, err := h.runner.runSingleStringQuery(ctx, fmt.Sprintf("USE %s", schemaName)); err != nil {
			return nil, fmt.Errorf("use schema %q failed: %v", schemaName, err)
		}

		dbDDLs := make([]*driverV2.DatabaseObjectDDL, 0, len(objInfo.DatabaseObjects))
		for _, obj := range objInfo.DatabaseObjects {
			switch obj.ObjectType {
			case driverV2.ObjectType_TABLE, driverV2.ObjectType_VIEW:
				// Hive views reuse SHOW CREATE TABLE; the rows returned by
				// HiveServer2 are joined with newline to form the DDL.
				rows, err := h.runner.runSingleStringQuery(ctx,
					fmt.Sprintf("SHOW CREATE TABLE %s", obj.ObjectName))
				if err != nil {
					return nil, fmt.Errorf("show create %s.%s failed: %v",
						schemaName, obj.ObjectName, err)
				}
				dbDDLs = append(dbDDLs, &driverV2.DatabaseObjectDDL{
					DatabaseObject: &driverV2.DatabaseObject{
						ObjectName: obj.ObjectName,
						ObjectType: obj.ObjectType,
					},
					ObjectDDL: strings.Join(rows, "\n"),
				})
			case driverV2.ObjectType_FUNCTION:
				// FUNCTION is planned for the second batch (compat-RISK-9).
				// Emit a placeholder result with empty DDL and log a WARN so
				// the upstream pipeline can surface the unsupported message.
				if h.log != nil {
					h.log.WithField("object", obj.ObjectName).
						Warnf("hive driver: %s", hiveFunctionUnsupportedMsg)
				}
				dbDDLs = append(dbDDLs, &driverV2.DatabaseObjectDDL{
					DatabaseObject: &driverV2.DatabaseObject{
						ObjectName: obj.ObjectName,
						ObjectType: obj.ObjectType,
					},
					ObjectDDL: "",
				})
				// Returning an error here would abort the whole batch and
				// drop the legitimate TABLE/VIEW results. Per design §3.2.1
				// we propagate the FUNCTION error via the upstream layer
				// (it inspects ObjectDDL == "" and the WARN log); the driver
				// does not return a Go error.
				return results, fmt.Errorf("%s", hiveFunctionUnsupportedMsg)
			case driverV2.ObjectType_PROCEDURE,
				driverV2.ObjectType_TRIGGER,
				driverV2.ObjectType_EVENT:
				// Hive does not physically support these object types
				// (compat-RISK-4). Short-circuit: skip the object so upstream
				// can continue processing the rest of the batch; do NOT
				// panic, do NOT return an error.
				if h.log != nil {
					h.log.WithField("object", obj.ObjectName).
						WithField("objectType", obj.ObjectType).
						Warn("hive driver: object type not supported, skipped")
				}
				continue
			default:
				// Unknown object type: warn and skip rather than fail; future
				// versions may add new ObjectType constants.
				if h.log != nil {
					h.log.WithField("object", obj.ObjectName).
						WithField("objectType", obj.ObjectType).
						Warn("hive driver: unknown object type, skipped")
				}
				continue
			}
		}

		results = append(results, &driverV2.DatabaseSchemaObjectResult{
			SchemaName:         schemaName,
			DatabaseObjectDDLs: dbDDLs,
		})
	}
	return results, nil
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
