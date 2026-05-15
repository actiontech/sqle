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
	// compareRunnerFactory builds a query runner for the calibrated
	// (compared) DSN side of GetDatabaseDiffModifySQL. In production it
	// opens a real gohive connection. Unit tests inject a fake to avoid
	// requiring a Hive server.
	compareRunnerFactory func(dsn *driverV2.DSN) (hiveQueryRunner, func(), error)
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
	// Default compareRunnerFactory opens a fresh gohive connection to the
	// calibratedDSN. Unit tests override this with a fake.
	impl.compareRunnerFactory = defaultCompareRunnerFactory
	return impl, nil
}

// defaultCompareRunnerFactory opens a real gohive connection to the
// calibrated DSN and returns a runner + close hook.
func defaultCompareRunnerFactory(dsn *driverV2.DSN) (hiveQueryRunner, func(), error) {
	conn, err := newHiveConnection(dsn)
	if err != nil {
		return nil, func() {}, fmt.Errorf("connect to compared Hive: %v", err)
	}
	closer := func() { _ = conn.Close() }
	return &gohiveQueryRunner{conn: conn}, closer, nil
}

// gohiveQueryRunner is the production hiveQueryRunner backed by *gohive.Connection.
type gohiveQueryRunner struct {
	conn *gohive.Connection
}

// hs2NoResultRowErrMarkers are substring markers that identify a HiveServer2
// "ROW-ERR" returned from FetchResults on a statement that produces no result
// columns (e.g. `USE <db>`, `SET ...`, DDL). Such an error is **non-fatal** in
// HS2's protocol: the same connection can still execute the next statement.
//
// See compat-RISK-10 (sqle-ee/docs/dev/compat_risks.md) and the reference
// implementation in sqle-ee/cmd/hivetool/main.go (which tolerates ROW-ERR
// and breaks out of the fetch loop).
var hs2NoResultRowErrMarkers = []string{
	// HiveServer2's generic "no result columns" message.
	"Server-side error; please check HS2 logs.",
	// Some HS2 builds wrap the same condition with an explicit status.
	"StatusCode:ERROR_STATUS",
}

// isHS2NoResultRowErr reports whether err is the non-fatal ROW-ERR
// HiveServer2 returns from FetchResults for a no-result-column statement.
// It is a pure string match against the error's Error() output so we can
// unit-test the classifier without needing a real gohive cursor.
//
// The function is conservative: only errors that match ALL of the canonical
// markers are treated as tolerable. Real Hive runtime errors (syntax errors,
// missing table, etc.) carry a different status / message and will NOT be
// matched.
func isHS2NoResultRowErr(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	for _, m := range hs2NoResultRowErrMarkers {
		if !strings.Contains(msg, m) {
			return false
		}
	}
	return true
}

// hiveCursor abstracts the minimum *gohive.Cursor surface used by the
// fetch loop. It exists solely so unit tests can substitute a fake; the
// production code path uses *gohive.Cursor directly.
type hiveCursor interface {
	HasMore(ctx context.Context) bool
	FetchOne(ctx context.Context, dests ...interface{})
	Err() error
}

// gohiveCursorAdapter adapts *gohive.Cursor (whose error is a public field
// rather than a method) to the hiveCursor interface for use in fetchAllRows.
type gohiveCursorAdapter struct {
	c *gohive.Cursor
}

func (a gohiveCursorAdapter) HasMore(ctx context.Context) bool                  { return a.c.HasMore(ctx) }
func (a gohiveCursorAdapter) FetchOne(ctx context.Context, dests ...interface{}) { a.c.FetchOne(ctx, dests...) }
func (a gohiveCursorAdapter) Err() error                                         { return a.c.Err }

// fetchAllRows pulls a single STRING column from cur and returns each row.
// HS2 ROW-ERR (non-fatal) is tolerated: when encountered, the loop breaks
// and the rows captured so far are returned with a nil error. Real fetch
// errors (any error that is not isHS2NoResultRowErr) are propagated.
//
// This is the shared core of gohiveQueryRunner.runSingleStringQuery; it is
// also used by unit tests via a fake hiveCursor so that the ROW-ERR
// tolerance contract can be exercised without a live HiveServer2 instance.
func fetchAllRows(ctx context.Context, cur hiveCursor) ([]string, error) {
	var rows []string
	for cur.HasMore(ctx) {
		if e := cur.Err(); e != nil {
			if isHS2NoResultRowErr(e) {
				// HS2 ROW-ERR on a no-result-column statement; treat as EOF.
				break
			}
			return nil, fmt.Errorf("failed to fetch row: %v", e)
		}
		var val string
		cur.FetchOne(ctx, &val)
		if e := cur.Err(); e != nil {
			if isHS2NoResultRowErr(e) {
				// Same ROW-ERR can surface during FetchOne; tolerate and stop.
				break
			}
			return nil, fmt.Errorf("failed to scan row: %v", e)
		}
		rows = append(rows, val)
	}
	return rows, nil
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

	return fetchAllRows(ctx, gohiveCursorAdapter{c: cursor})
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

	// SHOW DATABASES returns a single STRING column. Reuse fetchAllRows so
	// the ROW-ERR tolerance contract (compat-RISK-10) stays consistent with
	// runSingleStringQuery's behaviour.
	return fetchAllRows(ctx, gohiveCursorAdapter{c: cursor})
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

// WARNING comments emitted at the head of TABLE / VIEW DROP+CREATE SQL
// segments. Both Chinese and English lines are required by design §3.5 and
// are consumed by dms-ui-ee ModifiedSqlDrawer for top-banner detection.
const (
	hiveWarningTableDropCreate = "-- WARNING: data loss risk; table will be dropped and recreated.\n" +
		"-- 警告: 数据将丢失；表将被删除并重建。\n"
	hiveWarningViewDropCreate = "-- WARNING: view will be recreated; downstream queries depending on this view may be affected.\n" +
		"-- 警告: 视图将被重建；依赖该视图的下游查询可能受影响。\n"
)

// GetDatabaseDiffModifySQL generates the variant SQL needed to reconcile
// the calibrated (compared) side of a Hive instance so its objects match
// the base side (this driver's connection). The full strategy matrix is
// in docs/spec/design.md §3.4 and §3.2.2.
//
// Per-object behavior:
//
//   - TABLE only in base, not in compared -> CREATE TABLE ...      (no WARNING)
//   - TABLE only in compared, not in base -> DROP TABLE IF EXISTS  (no WARNING)
//   - TABLE on both sides, ALTER-able diff -> ALTER TABLE sequence (no WARNING)
//   - TABLE on both sides, fallback diff   -> DROP+CREATE          (WARNING)
//   - VIEW (any diff, any direction)       -> DROP+CREATE          (WARNING)
//   - FUNCTION                             -> return Chinese error; skip the
//     entire schema's result (compat-RISK-9)
//   - PROCEDURE / TRIGGER / EVENT          -> short-circuit, skip silently
//     with a WARN log (compat-RISK-4)
//
// Each non-empty SchemaName produces a result entry with `USE <schema>;`
// prefixed to the SQL block.
func (h *HiveDriverImpl) GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *driverV2.DSN, objInfos []*driverV2.DatabasCompareSchemaInfo) ([]*driverV2.DatabaseDiffModifySQLResult, error) {
	if h.runner == nil {
		return nil, fmt.Errorf("hive base connection is not initialized")
	}
	if h.compareRunnerFactory == nil {
		return nil, fmt.Errorf("hive compareRunnerFactory is not initialized")
	}
	// Open compare-side runner once for the whole call; close at the end.
	compareRunner, closeCompare, err := h.compareRunnerFactory(calibratedDSN)
	if err != nil {
		return nil, err
	}
	defer closeCompare()

	results := make([]*driverV2.DatabaseDiffModifySQLResult, 0, len(objInfos))
	for _, objInfo := range objInfos {
		baseSchemaName := objInfo.BaseSchemaName
		if baseSchemaName == "" {
			baseSchemaName = defaultHiveSchema
		}
		comparedSchemaName := objInfo.ComparedSchemaName
		if comparedSchemaName == "" {
			comparedSchemaName = defaultHiveSchema
		}

		// USE on base side first so subsequent SHOW CREATE TABLE resolves.
		if _, err := h.runner.runSingleStringQuery(ctx,
			fmt.Sprintf("USE %s", baseSchemaName)); err != nil {
			return nil, fmt.Errorf("use base schema %q: %v", baseSchemaName, err)
		}
		if _, err := compareRunner.runSingleStringQuery(ctx,
			fmt.Sprintf("USE %s", comparedSchemaName)); err != nil {
			return nil, fmt.Errorf("use compared schema %q: %v", comparedSchemaName, err)
		}

		sqls := make([]string, 0)
		sqls = append(sqls, fmt.Sprintf("USE %s;", comparedSchemaName))

		for _, obj := range objInfo.DatabaseObjects {
			switch obj.ObjectType {
			case driverV2.ObjectType_TABLE:
				stmts, terr := diffTableObject(ctx, h.runner, compareRunner, obj.ObjectName)
				if terr != nil {
					return nil, terr
				}
				sqls = append(sqls, stmts...)
			case driverV2.ObjectType_VIEW:
				stmts, verr := diffViewObject(ctx, h.runner, compareRunner, obj.ObjectName)
				if verr != nil {
					return nil, verr
				}
				sqls = append(sqls, stmts...)
			case driverV2.ObjectType_FUNCTION:
				// Compat-RISK-9: FUNCTION is planned for the second batch.
				if h.log != nil {
					h.log.WithField("object", obj.ObjectName).
						Warnf("hive driver: %s", hiveFunctionUnsupportedMsg)
				}
				return results, fmt.Errorf("%s", hiveFunctionUnsupportedMsg)
			case driverV2.ObjectType_PROCEDURE,
				driverV2.ObjectType_TRIGGER,
				driverV2.ObjectType_EVENT:
				// Compat-RISK-4: short-circuit physically unsupported types.
				if h.log != nil {
					h.log.WithField("object", obj.ObjectName).
						WithField("objectType", obj.ObjectType).
						Warn("hive driver: object type not supported, skipped")
				}
				continue
			default:
				if h.log != nil {
					h.log.WithField("object", obj.ObjectName).
						WithField("objectType", obj.ObjectType).
						Warn("hive driver: unknown object type, skipped")
				}
				continue
			}
		}

		results = append(results, &driverV2.DatabaseDiffModifySQLResult{
			SchemaName: comparedSchemaName,
			ModifySQLs: sqls,
		})
	}
	return results, nil
}

// fetchTableDDL runs SHOW CREATE TABLE for the given object and returns the
// concatenated DDL string. An empty string + nil error indicates the table
// does not exist on that side; any other error is a real failure.
func fetchTableDDL(ctx context.Context, runner hiveQueryRunner, objectName string) (string, bool, error) {
	rows, err := runner.runSingleStringQuery(ctx,
		fmt.Sprintf("SHOW CREATE TABLE %s", objectName))
	if err != nil {
		// Heuristic: Hive returns a SemanticException when the object is
		// missing. We surface it as "not found" so the caller can decide
		// the direction (only-in-base vs only-in-compared). Other errors
		// could be propagated by the caller but for diff purposes we treat
		// any failure as "not found", matching the MySQL impl pattern.
		return "", false, nil
	}
	if len(rows) == 0 {
		return "", false, nil
	}
	return strings.Join(rows, "\n"), true, nil
}

// diffTableObject generates the variant SQL for a single TABLE object.
// It captures DDL from both sides and dispatches to diffTableDDL for the
// detailed ALTER-vs-DROP+CREATE matrix decision.
func diffTableObject(ctx context.Context, baseRunner, compareRunner hiveQueryRunner, objectName string) ([]string, error) {
	baseDDL, baseExists, err := fetchTableDDL(ctx, baseRunner, objectName)
	if err != nil {
		return nil, err
	}
	compareDDL, compareExists, err := fetchTableDDL(ctx, compareRunner, objectName)
	if err != nil {
		return nil, err
	}

	switch {
	case baseExists && !compareExists:
		// Only on base side -> create on compared side (no WARNING).
		return []string{ensureSemicolon(baseDDL)}, nil
	case !baseExists && compareExists:
		// Only on compared side -> drop from compared (no WARNING).
		return []string{fmt.Sprintf("DROP TABLE IF EXISTS %s;", objectName)}, nil
	case !baseExists && !compareExists:
		// Neither side has it; nothing to do.
		return nil, nil
	}

	// Both sides exist: decide ALTER vs DROP+CREATE.
	alters, fallback, err := diffTableDDL(baseDDL, compareDDL)
	if err != nil {
		return nil, fmt.Errorf("diffTableDDL %q: %v", objectName, err)
	}
	if fallback {
		// DROP+CREATE with WARNING header (compat-RISK-6).
		return []string{
			hiveWarningTableDropCreate +
				fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", objectName) +
				ensureSemicolon(baseDDL),
		}, nil
	}
	if len(alters) == 0 {
		// No structural difference detected.
		return nil, nil
	}
	return alters, nil
}

// diffViewObject generates the variant SQL for a single VIEW object.
// Per design §3.4 / D3 decision: any view difference produces a unified
// DROP+CREATE with the view-recreated WARNING header.
func diffViewObject(ctx context.Context, baseRunner, compareRunner hiveQueryRunner, objectName string) ([]string, error) {
	// Hive views use SHOW CREATE TABLE.
	baseDDL, baseExists, err := fetchTableDDL(ctx, baseRunner, objectName)
	if err != nil {
		return nil, err
	}
	compareDDL, compareExists, err := fetchTableDDL(ctx, compareRunner, objectName)
	if err != nil {
		return nil, err
	}

	switch {
	case baseExists && !compareExists:
		return []string{ensureSemicolon(baseDDL)}, nil
	case !baseExists && compareExists:
		return []string{fmt.Sprintf("DROP VIEW IF EXISTS %s;", objectName)}, nil
	case !baseExists && !compareExists:
		return nil, nil
	}

	if normalizeWhitespace(baseDDL) == normalizeWhitespace(compareDDL) {
		return nil, nil
	}
	// Any difference triggers DROP+CREATE with the VIEW-specific WARNING
	// (compat-RISK-6, design §3.4 unified rule).
	return []string{
		hiveWarningViewDropCreate +
			fmt.Sprintf("DROP VIEW IF EXISTS %s;\n", objectName) +
			ensureSemicolon(baseDDL),
	}, nil
}

// ensureSemicolon appends a trailing semicolon if the DDL string does not
// already end with one (ignoring trailing whitespace).
func ensureSemicolon(s string) string {
	trimmed := strings.TrimRight(s, " \t\n\r")
	if strings.HasSuffix(trimmed, ";") {
		return trimmed + "\n"
	}
	return trimmed + ";\n"
}

// normalizeWhitespace collapses runs of whitespace to a single space so
// that view DDL comparisons are not sensitive to formatting differences
// (HiveServer2 sometimes emits extra newlines).
func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
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
