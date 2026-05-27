// Tests for the Hive driver's compare-capability surface introduced by
// Issue #2872. The 16 `Test_*` functions defined here cover every row of
// the design.md §3.4 variant SQL matrix and §5.2.4 unit-test table, plus
// the metas declaration required by compat-RISK-1.
//
// Design choices:
//   - Each test uses the map-case style required by dev.md (key = scenario,
//     value = inputs + expectations); table iteration is deterministic via
//     sorted keys when output order matters.
//   - No network or database is touched: the Hive cursor is replaced with
//     fakeQueryRunner that satisfies the hiveQueryRunner interface.
//   - WARNING text assertions use strings.Contains rather than full-string
//     equality so cosmetic tweaks (whitespace, trailing newline) do not
//     break tests.
package hive

import (
	"context"
	"sort"
	"strings"
	"testing"

	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/sirupsen/logrus"
)

// fakeQueryRunner implements hiveQueryRunner using a preloaded map of
// query string -> rows. The match is performed by exact equality first,
// then by prefix (so callers can register a "SHOW CREATE TABLE tbl_order"
// response without having to anticipate USE-statement prefixes).
//
// When a query is registered with errOnQuery (non-empty), runSingleStringQuery
// returns that as an error to simulate Hive failures.
type fakeQueryRunner struct {
	rows       map[string][]string
	errOnQuery map[string]string
	// log captures every query the code under test issued, in order.
	log []string
}

func newFakeRunner() *fakeQueryRunner {
	return &fakeQueryRunner{
		rows:       map[string][]string{},
		errOnQuery: map[string]string{},
	}
}

func (f *fakeQueryRunner) on(query string, rows ...string) *fakeQueryRunner {
	f.rows[query] = rows
	return f
}

func (f *fakeQueryRunner) fail(query, errMsg string) *fakeQueryRunner {
	f.errOnQuery[query] = errMsg
	return f
}

func (f *fakeQueryRunner) runSingleStringQuery(_ context.Context, q string) ([]string, error) {
	f.log = append(f.log, q)
	if msg, ok := f.errOnQuery[q]; ok {
		return nil, fakeError(msg)
	}
	if rows, ok := f.rows[q]; ok {
		return rows, nil
	}
	// Default: pretend the table/view does not exist by returning an error
	// (matches how SHOW CREATE TABLE behaves for missing objects).
	return nil, fakeError("missing object for query: " + q)
}

type fakeError string

func (e fakeError) Error() string { return string(e) }

// newDriverWithRunners builds a HiveDriverImpl wired up with provided
// fake runners for base and compared sides. The compareRunnerFactory
// returns the supplied compareRunner regardless of DSN.
func newDriverWithRunners(base, compared hiveQueryRunner) *HiveDriverImpl {
	return &HiveDriverImpl{
		log:    logrus.NewEntry(logrus.New()),
		runner: base,
		compareRunnerFactory: func(_ *driverV2.DSN) (hiveQueryRunner, func(), error) {
			return compared, func() {}, nil
		},
	}
}

// --------------------------------------------------------------------- //
// 1. Test_Metas_EnabledOptional
// --------------------------------------------------------------------- //

// Test_Metas_EnabledOptional verifies compat-RISK-1: the metas exposed by
// the Hive driver must declare BOTH OptionalGetDatabaseObjectDDL and
// OptionalGetDatabaseDiffModifySQL so the structure-compare capability
// check whitelist accepts Hive.
func Test_Metas_EnabledOptional(t *testing.T) {
	cases := map[string]struct {
		want driverV2.OptionalModule
	}{
		"OptionalGetDatabaseObjectDDL":     {want: driverV2.OptionalGetDatabaseObjectDDL},
		"OptionalGetDatabaseDiffModifySQL": {want: driverV2.OptionalGetDatabaseDiffModifySQL},
	}
	p := &PluginProcessor{}
	metas, err := p.GetDriverMetas()
	if err != nil {
		t.Fatalf("GetDriverMetas: %v", err)
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if !metas.IsOptionalModuleEnabled(tc.want) {
				t.Errorf("metas missing %v; got modules=%v", tc.want, metas.EnabledOptionalModule)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 2. Test_GetDatabaseObjectDDL_TableHappy
// --------------------------------------------------------------------- //

// Test_GetDatabaseObjectDDL_TableHappy verifies that TABLE objects flow
// through SHOW CREATE TABLE and the rows are concatenated as the DDL.
func Test_GetDatabaseObjectDDL_TableHappy(t *testing.T) {
	cases := map[string]struct {
		schema string
		object string
		rows   []string
		want   string
	}{
		"single row": {
			schema: "sqle_compare_test",
			object: "tbl_order",
			rows:   []string{"CREATE TABLE tbl_order (id BIGINT) STORED AS ORC;"},
			want:   "CREATE TABLE tbl_order (id BIGINT) STORED AS ORC;",
		},
		"multi row joined by newline": {
			schema: "default",
			object: "users",
			rows:   []string{"CREATE TABLE users(", "  id INT", ") STORED AS TEXTFILE;"},
			want:   "CREATE TABLE users(\n  id INT\n) STORED AS TEXTFILE;",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			runner := newFakeRunner().
				on("USE "+tc.schema, "OK").
				on("SHOW CREATE TABLE "+tc.object, tc.rows...)
			h := &HiveDriverImpl{
				log:    logrus.NewEntry(logrus.New()),
				runner: runner,
			}
			results, err := h.GetDatabaseObjectDDL(context.Background(),
				[]*driverV2.DatabaseSchemaInfo{{
					SchemaName: tc.schema,
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results) != 1 || len(results[0].DatabaseObjectDDLs) != 1 {
				t.Fatalf("unexpected results shape: %+v", results)
			}
			got := results[0].DatabaseObjectDDLs[0].ObjectDDL
			if got != tc.want {
				t.Errorf("DDL mismatch:\n got=%q\nwant=%q", got, tc.want)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 3. Test_GetDatabaseObjectDDL_ViewHappy
// --------------------------------------------------------------------- //

// Test_GetDatabaseObjectDDL_ViewHappy verifies that VIEW objects ALSO go
// through SHOW CREATE TABLE (Hive reuses the same command for views).
func Test_GetDatabaseObjectDDL_ViewHappy(t *testing.T) {
	cases := map[string]struct {
		object string
		rows   []string
	}{
		"basic view": {
			object: "v_order_active",
			rows: []string{
				"CREATE VIEW v_order_active AS SELECT id, name FROM tbl_order WHERE 1=1;",
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			runner := newFakeRunner().
				on("USE default", "OK").
				on("SHOW CREATE TABLE "+tc.object, tc.rows...)
			h := &HiveDriverImpl{log: logrus.NewEntry(logrus.New()), runner: runner}
			results, err := h.GetDatabaseObjectDDL(context.Background(),
				[]*driverV2.DatabaseSchemaInfo{{
					SchemaName: "default",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_VIEW},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results[0].DatabaseObjectDDLs) != 1 {
				t.Fatalf("expected 1 DDL, got %d", len(results[0].DatabaseObjectDDLs))
			}
			// The driver must dispatch the same SHOW CREATE TABLE query for
			// VIEW; verify against the call log.
			found := false
			for _, q := range runner.log {
				if q == "SHOW CREATE TABLE "+tc.object {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("expected SHOW CREATE TABLE call for VIEW; got log=%v", runner.log)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 4. Test_GetDatabaseObjectDDL_FunctionRejected
// --------------------------------------------------------------------- //

// Test_GetDatabaseObjectDDL_FunctionRejected verifies compat-RISK-9 after
// the FIX-002 driver alignment (TC-HIVE-015): FUNCTION objects are silently
// skipped (no Go error, no placeholder DDL entry), matching the behaviour
// of the PROCEDURE/TRIGGER/EVENT short-circuit. The Chinese error message
// is emitted via the WARN log only — never as a returned Go error.
func Test_GetDatabaseObjectDDL_FunctionRejected(t *testing.T) {
	cases := map[string]struct {
		object string
	}{
		"function planned for batch 2": {
			object: "my_udf",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			runner := newFakeRunner().on("USE default", "OK")
			h := &HiveDriverImpl{log: logrus.NewEntry(logrus.New()), runner: runner}
			results, err := h.GetDatabaseObjectDDL(context.Background(),
				[]*driverV2.DatabaseSchemaInfo{{
					SchemaName: "default",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_FUNCTION},
					},
				}})
			if err != nil {
				t.Fatalf("expected nil error for FUNCTION skip, got %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 schema result entry, got %d", len(results))
			}
			if got := len(results[0].DatabaseObjectDDLs); got != 0 {
				t.Errorf("expected 0 DDLs (FUNCTION skipped) but got %d entries: %#v",
					got, results[0].DatabaseObjectDDLs)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 5. Test_GetDatabaseObjectDDL_UnsupportedTypeShortCircuit
// --------------------------------------------------------------------- //

// Test_GetDatabaseObjectDDL_UnsupportedTypeShortCircuit verifies
// compat-RISK-4: PROCEDURE / TRIGGER / EVENT are silently skipped (no
// error, no panic). One case per object type.
func Test_GetDatabaseObjectDDL_UnsupportedTypeShortCircuit(t *testing.T) {
	cases := map[string]string{
		"PROCEDURE": driverV2.ObjectType_PROCEDURE,
		"TRIGGER":   driverV2.ObjectType_TRIGGER,
		"EVENT":     driverV2.ObjectType_EVENT,
	}
	for name, objType := range cases {
		t.Run(name, func(t *testing.T) {
			runner := newFakeRunner().on("USE default", "OK")
			h := &HiveDriverImpl{log: logrus.NewEntry(logrus.New()), runner: runner}
			results, err := h.GetDatabaseObjectDDL(context.Background(),
				[]*driverV2.DatabaseSchemaInfo{{
					SchemaName: "default",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "x_obj", ObjectType: objType},
					},
				}})
			if err != nil {
				t.Fatalf("expected nil error for %s short-circuit, got %v", objType, err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 result entry, got %d", len(results))
			}
			if got := len(results[0].DatabaseObjectDDLs); got != 0 {
				t.Errorf("expected 0 DDLs (skipped) for %s, got %d", objType, got)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 6. Test_GetDatabaseObjectDDL_EmptyObjectList
// --------------------------------------------------------------------- //

// Test_GetDatabaseObjectDDL_EmptyObjectList verifies the **new contract**
// for empty DatabaseObjects (compat-RISK-10 secondary fix, Task-TEST-FIX-001):
//
//   - When the caller passes an empty DatabaseObjects slice, the driver
//     now auto-enumerates every TABLE/VIEW in the schema via SHOW TABLES
//     / SHOW VIEWS (mirroring MySQL driver mysql_ee.go::GetDatabaseObjectDDL
//     line 380-389).
//   - This is required because server/compare/database_compare_ee.go
//     ExecDatabaseCompare never populates DatabaseObjects itself — it only
//     forwards SchemaName, expecting the driver to discover the rest.
//
// The previous contract ("empty list = zero DDLs returned") was incompatible
// with the controller's actual call shape and would always yield "same".
func Test_GetDatabaseObjectDDL_EmptyObjectList(t *testing.T) {
	cases := map[string]struct {
		schema       string
		showTables   []string
		showViews    []string
		showCreates  map[string]string // object name -> DDL row
		expectCount  int
	}{
		"empty objects in default schema": {
			schema:     "default",
			showTables: []string{"t1", "v1"},
			showViews:  []string{"v1"},
			showCreates: map[string]string{
				"t1": "CREATE TABLE `t1` (`id` int)",
				"v1": "CREATE VIEW `v1` AS SELECT 1",
			},
			expectCount: 2,
		},
		"empty objects in named schema": {
			schema:     "sqle_compare_test",
			showTables: []string{"t_diff_only"},
			showViews:  []string{},
			showCreates: map[string]string{
				"t_diff_only": "CREATE TABLE `t_diff_only` (`a` string)",
			},
			expectCount: 1,
		},
		"empty schema name falls back default": {
			schema:     "",
			showTables: []string{},
			showViews:  []string{},
			expectCount: 0,
		},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			effective := tc.schema
			if effective == "" {
				effective = "default"
			}
			runner := newFakeRunner().
				on("USE "+effective, "OK")
			// Auto-discovery scripts: SHOW TABLES, SHOW VIEWS, then SHOW
			// CREATE TABLE for every discovered object. fakeRunner.on uses
			// a variadic rows... param so a slice expands via "...".
			runner.on("SHOW TABLES", tc.showTables...)
			runner.on("SHOW VIEWS", tc.showViews...)
			for obj, ddl := range tc.showCreates {
				runner.on("SHOW CREATE TABLE "+obj, ddl)
			}
			h := &HiveDriverImpl{log: logrus.NewEntry(logrus.New()), runner: runner}
			results, err := h.GetDatabaseObjectDDL(context.Background(),
				[]*driverV2.DatabaseSchemaInfo{{
					SchemaName:      tc.schema,
					DatabaseObjects: nil,
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 schema result, got %d", len(results))
			}
			if got := len(results[0].DatabaseObjectDDLs); got != tc.expectCount {
				t.Errorf("expected %d DDLs from auto-discovery, got %d", tc.expectCount, got)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 7. Test_GetDatabaseDiffModifySQL_OnlyBaseHasTable
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_OnlyBaseHasTable: design §3.4 row
// "TABLE / 只在基准侧" -> CREATE TABLE on compared side; no WARNING.
func Test_GetDatabaseDiffModifySQL_OnlyBaseHasTable(t *testing.T) {
	cases := map[string]struct {
		object string
		ddl    string
	}{
		"new table only in base": {
			object: "tbl_order",
			ddl:    "CREATE TABLE tbl_order (id BIGINT, name STRING) STORED AS ORC",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE "+tc.object, tc.ddl)
			compared := newFakeRunner().
				on("USE compared_schema", "OK")
				// SHOW CREATE TABLE on compared returns default error → not found
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 result, got %d", len(results))
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			if !strings.Contains(full, "USE compared_schema;") {
				t.Errorf("expected USE compared_schema; header; got=%q", full)
			}
			if !strings.Contains(full, "CREATE TABLE "+tc.object) {
				t.Errorf("expected CREATE TABLE %s; got=%q", tc.object, full)
			}
			if strings.Contains(full, "WARNING") {
				t.Errorf("expected NO WARNING for only-in-base case; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 8. Test_GetDatabaseDiffModifySQL_OnlyTargetHasTable
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_OnlyTargetHasTable: design §3.4 row
// "TABLE / 只在比对侧" -> DROP TABLE IF EXISTS on compared side; no WARNING.
func Test_GetDatabaseDiffModifySQL_OnlyTargetHasTable(t *testing.T) {
	cases := map[string]struct {
		object string
		ddl    string
	}{
		"orphan table only in compared": {
			object: "tbl_stale",
			ddl:    "CREATE TABLE tbl_stale (id INT)",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK")
				// SHOW CREATE TABLE on base returns default "not found"
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE "+tc.object, tc.ddl)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			if !strings.Contains(full, "DROP TABLE IF EXISTS "+tc.object) {
				t.Errorf("expected DROP TABLE IF EXISTS %s; got=%q", tc.object, full)
			}
			if strings.Contains(full, "WARNING") {
				t.Errorf("expected NO WARNING for only-in-compared case; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 9. Test_GetDatabaseDiffModifySQL_TableAddColumn
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_TableAddColumn: design §3.4 row
// "TABLE / 加列 (ADD COLUMN)" -> ALTER ADD COLUMNS; no WARNING.
func Test_GetDatabaseDiffModifySQL_TableAddColumn(t *testing.T) {
	cases := map[string]struct {
		object  string
		baseDDL string
		compDDL string
		// substrings the produced SQL block must contain
		wantContains []string
		// substrings the produced SQL block must NOT contain
		wantNotContains []string
	}{
		"add a single column at end": {
			object:  "tbl_order",
			baseDDL: "CREATE TABLE tbl_order (id BIGINT, name STRING, age INT) STORED AS ORC",
			compDDL: "CREATE TABLE tbl_order (id BIGINT, name STRING) STORED AS ORC",
			wantContains: []string{
				"ALTER TABLE tbl_order ADD COLUMNS (age INT)",
			},
			wantNotContains: []string{"WARNING", "DROP TABLE"},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE "+tc.object, tc.baseDDL)
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE "+tc.object, tc.compDDL)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			for _, want := range tc.wantContains {
				if !strings.Contains(full, want) {
					t.Errorf("expected SQL to contain %q; got=%q", want, full)
				}
			}
			for _, notWant := range tc.wantNotContains {
				if strings.Contains(full, notWant) {
					t.Errorf("expected SQL NOT to contain %q; got=%q", notWant, full)
				}
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 10. Test_GetDatabaseDiffModifySQL_TableChangeColumnType_Compat
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_TableChangeColumnType_Compat: design §3.4
// row "TABLE / 改列类型（兼容 widen）" -> ALTER CHANGE COLUMN; no WARNING.
func Test_GetDatabaseDiffModifySQL_TableChangeColumnType_Compat(t *testing.T) {
	cases := map[string]struct {
		baseDDL string
		compDDL string
	}{
		"int to bigint widening": {
			baseDDL: "CREATE TABLE t (id BIGINT, age INT) STORED AS ORC",
			compDDL: "CREATE TABLE t (id INT, age INT) STORED AS ORC",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE t", tc.baseDDL)
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE t", tc.compDDL)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "t", ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			if !strings.Contains(full, "ALTER TABLE t CHANGE COLUMN id id BIGINT") {
				t.Errorf("expected CHANGE COLUMN id id BIGINT; got=%q", full)
			}
			if strings.Contains(full, "WARNING") {
				t.Errorf("compatible widening should not emit WARNING; got=%q", full)
			}
			if strings.Contains(full, "DROP TABLE") {
				t.Errorf("compatible widening must not fall back to DROP+CREATE; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 11. Test_GetDatabaseDiffModifySQL_TableChangeColumnType_Incompat
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_TableChangeColumnType_Incompat: design §3.4
// row "TABLE / 改列类型（不兼容）" -> DROP+CREATE; WARNING (data loss).
func Test_GetDatabaseDiffModifySQL_TableChangeColumnType_Incompat(t *testing.T) {
	cases := map[string]struct {
		baseDDL string
		compDDL string
	}{
		"bigint to string is incompatible": {
			baseDDL: "CREATE TABLE t (id STRING) STORED AS ORC",
			compDDL: "CREATE TABLE t (id BIGINT) STORED AS ORC",
		},
		"int to timestamp is incompatible": {
			baseDDL: "CREATE TABLE t (ts TIMESTAMP) STORED AS ORC",
			compDDL: "CREATE TABLE t (ts INT) STORED AS ORC",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE t", tc.baseDDL)
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE t", tc.compDDL)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "t", ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			if !strings.Contains(full, "DROP TABLE IF EXISTS t") {
				t.Errorf("expected DROP TABLE IF EXISTS t; got=%q", full)
			}
			if !strings.Contains(full, "CREATE TABLE t") {
				t.Errorf("expected CREATE TABLE t; got=%q", full)
			}
			if !strings.Contains(full, "-- WARNING:") {
				t.Errorf("expected WARNING marker for data loss case; got=%q", full)
			}
			if !strings.Contains(full, "数据将丢失") {
				t.Errorf("expected Chinese data loss warning; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 12. Test_GetDatabaseDiffModifySQL_TableChangeStoredAs
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_TableChangeStoredAs: design §3.4 row
// "TABLE / 改存储格式 (STORED AS)" -> DROP+CREATE; WARNING.
func Test_GetDatabaseDiffModifySQL_TableChangeStoredAs(t *testing.T) {
	cases := map[string]struct {
		baseDDL string
		compDDL string
	}{
		"ORC to PARQUET storage change": {
			baseDDL: "CREATE TABLE t (id BIGINT) STORED AS PARQUET",
			compDDL: "CREATE TABLE t (id BIGINT) STORED AS ORC",
		},
		"TEXTFILE to ORC storage change": {
			baseDDL: "CREATE TABLE t (id BIGINT) STORED AS ORC",
			compDDL: "CREATE TABLE t (id BIGINT) STORED AS TEXTFILE",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE t", tc.baseDDL)
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE t", tc.compDDL)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "t", ObjectType: driverV2.ObjectType_TABLE},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			if !strings.Contains(full, "DROP TABLE IF EXISTS t") {
				t.Errorf("expected DROP TABLE IF EXISTS t; got=%q", full)
			}
			if !strings.Contains(full, "-- WARNING:") {
				t.Errorf("expected WARNING marker; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 13. Test_GetDatabaseDiffModifySQL_ViewDiff_AlwaysDropCreate
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_ViewDiff_AlwaysDropCreate: design §3.4 D3
// row "VIEW / 任意差异" -> 统一 DROP+CREATE; WARNING (view recreated).
func Test_GetDatabaseDiffModifySQL_ViewDiff_AlwaysDropCreate(t *testing.T) {
	cases := map[string]struct {
		baseDDL string
		compDDL string
	}{
		"select clause differs": {
			baseDDL: "CREATE VIEW v AS SELECT id, name FROM t WHERE 1=1",
			compDDL: "CREATE VIEW v AS SELECT id FROM t",
		},
		"trivial whitespace-only diff should be ignored": {
			baseDDL: "CREATE VIEW v AS SELECT id FROM t",
			compDDL: "CREATE  VIEW v AS SELECT id   FROM t",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE v", tc.baseDDL)
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE v", tc.compDDL)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "v", ObjectType: driverV2.ObjectType_VIEW},
					},
				}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			// Whitespace-only diffs must NOT emit DROP+CREATE.
			if name == "trivial whitespace-only diff should be ignored" {
				if strings.Contains(full, "DROP VIEW") {
					t.Errorf("trivial whitespace diff must not emit DROP VIEW; got=%q", full)
				}
				return
			}
			if !strings.Contains(full, "DROP VIEW IF EXISTS v") {
				t.Errorf("expected DROP VIEW IF EXISTS v; got=%q", full)
			}
			if !strings.Contains(full, "-- WARNING:") {
				t.Errorf("expected view-recreated WARNING marker; got=%q", full)
			}
			if !strings.Contains(full, "视图将被重建") {
				t.Errorf("expected Chinese view-recreated message; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 14. Test_GetDatabaseDiffModifySQL_FunctionRejected
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_FunctionRejected: design §3.2.2 line 239
// row "FUNCTION / 第二批". After FIX-002 (TC-HIVE-015 / 016): the driver
// silently skips the FUNCTION object — no Go error, no FUNCTION entry in
// the SQL block. Only the leading `USE compared_schema;` header remains;
// upstream layers can detect "all objects unsupported" by the empty body
// (compat-RISK-9 verified; aligned with PROCEDURE/TRIGGER/EVENT
// short-circuit at line 824).
func Test_GetDatabaseDiffModifySQL_FunctionRejected(t *testing.T) {
	cases := map[string]struct {
		object string
	}{
		"function returns no error and skips silently": {
			object: "my_udf",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().on("USE base_schema", "OK")
			compared := newFakeRunner().on("USE compared_schema", "OK")
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: tc.object, ObjectType: driverV2.ObjectType_FUNCTION},
					},
				}})
			if err != nil {
				t.Fatalf("expected nil error for FUNCTION skip, got %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 schema result, got %d", len(results))
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			// Only the USE header should be present; no FUNCTION-related output.
			if !strings.Contains(full, "USE compared_schema;") {
				t.Errorf("expected USE compared_schema; header in block; got=%q", full)
			}
			if strings.Contains(full, tc.object) {
				t.Errorf("expected FUNCTION object %q to be skipped from results; got block=%q",
					tc.object, full)
			}
			if strings.Contains(full, "DROP") ||
				strings.Contains(full, "CREATE") ||
				strings.Contains(full, "ALTER") {
				t.Errorf("expected FUNCTION to produce no DDL; got block=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 14b. Test_GetDatabaseDiffModifySQL_MixedFunctionAndTable
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_MixedFunctionAndTable: TC-HIVE-016 (compat-
// RISK-9). Verifies that a mixed batch containing both a TABLE and a
// FUNCTION lets the TABLE main path produce its ALTER SQL while the
// FUNCTION is silently skipped. Before FIX-002 the driver returned a hard
// error at the FUNCTION branch, dropping the TABLE result entirely; this
// test pins the fixed behaviour.
func Test_GetDatabaseDiffModifySQL_MixedFunctionAndTable(t *testing.T) {
	cases := map[string]struct {
		baseDDL string
		compDDL string
	}{
		"int to bigint widening alongside FUNCTION": {
			baseDDL: "CREATE TABLE t_alter_widen (amt BIGINT) STORED AS ORC",
			compDDL: "CREATE TABLE t_alter_widen (amt INT) STORED AS ORC",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().
				on("USE base_schema", "OK").
				on("SHOW CREATE TABLE t_alter_widen", tc.baseDDL)
			compared := newFakeRunner().
				on("USE compared_schema", "OK").
				on("SHOW CREATE TABLE t_alter_widen", tc.compDDL)
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "t_alter_widen", ObjectType: driverV2.ObjectType_TABLE},
						{ObjectName: "fake_fn", ObjectType: driverV2.ObjectType_FUNCTION},
					},
				}})
			if err != nil {
				t.Fatalf("expected nil error for mixed batch, got %v", err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 schema result, got %d", len(results))
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			// TABLE main path: ALTER produced.
			if !strings.Contains(full, "ALTER TABLE t_alter_widen CHANGE COLUMN amt amt BIGINT") {
				t.Errorf("expected TABLE ALTER CHANGE COLUMN amt amt BIGINT; got=%q", full)
			}
			// FUNCTION must be skipped — no fake_fn anywhere in the block.
			if strings.Contains(full, "fake_fn") {
				t.Errorf("expected FUNCTION fake_fn to be skipped; got block=%q", full)
			}
			// No WARNING / DROP fallback for a compatible widen.
			if strings.Contains(full, "WARNING") {
				t.Errorf("compatible widening should not emit WARNING; got=%q", full)
			}
			if strings.Contains(full, "DROP TABLE") {
				t.Errorf("compatible widening must not fall back to DROP+CREATE; got=%q", full)
			}
			// Header preserved.
			if !strings.Contains(full, "USE compared_schema;") {
				t.Errorf("expected USE compared_schema; header in block; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 14c. Test_GetDatabaseObjectDDL_MixedFunctionAndTable
// --------------------------------------------------------------------- //

// Test_GetDatabaseObjectDDL_MixedFunctionAndTable: TC-HIVE-016 sister
// coverage for the SHOW CREATE TABLE path. A mixed batch (TABLE +
// FUNCTION) returns the TABLE DDL and silently skips the FUNCTION —
// the result entry contains exactly one DatabaseObjectDDL for the TABLE.
func Test_GetDatabaseObjectDDL_MixedFunctionAndTable(t *testing.T) {
	runner := newFakeRunner().
		on("USE default", "OK").
		on("SHOW CREATE TABLE tbl_order",
			"CREATE TABLE tbl_order (id BIGINT) STORED AS ORC;")
	h := &HiveDriverImpl{log: logrus.NewEntry(logrus.New()), runner: runner}
	results, err := h.GetDatabaseObjectDDL(context.Background(),
		[]*driverV2.DatabaseSchemaInfo{{
			SchemaName: "default",
			DatabaseObjects: []*driverV2.DatabaseObject{
				{ObjectName: "tbl_order", ObjectType: driverV2.ObjectType_TABLE},
				{ObjectName: "fake_fn", ObjectType: driverV2.ObjectType_FUNCTION},
			},
		}})
	if err != nil {
		t.Fatalf("expected nil error for mixed batch, got %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 schema result, got %d", len(results))
	}
	ddls := results[0].DatabaseObjectDDLs
	if len(ddls) != 1 {
		t.Fatalf("expected exactly 1 DDL (TABLE only), got %d: %#v", len(ddls), ddls)
	}
	if ddls[0].DatabaseObject == nil ||
		ddls[0].DatabaseObject.ObjectType != driverV2.ObjectType_TABLE ||
		ddls[0].DatabaseObject.ObjectName != "tbl_order" {
		t.Errorf("expected TABLE tbl_order entry; got %#v", ddls[0])
	}
	if !strings.Contains(ddls[0].ObjectDDL, "CREATE TABLE tbl_order") {
		t.Errorf("expected SHOW CREATE TABLE output for tbl_order; got=%q", ddls[0].ObjectDDL)
	}
}

// --------------------------------------------------------------------- //
// 15. Test_GetDatabaseDiffModifySQL_UnsupportedTypeFiltered
// --------------------------------------------------------------------- //

// Test_GetDatabaseDiffModifySQL_UnsupportedTypeFiltered: design §3.2.2
// row PROCEDURE/TRIGGER/EVENT -> driver skips them; SQL block only
// contains the USE prefix (compat-RISK-4).
func Test_GetDatabaseDiffModifySQL_UnsupportedTypeFiltered(t *testing.T) {
	cases := map[string]string{
		"PROCEDURE": driverV2.ObjectType_PROCEDURE,
		"TRIGGER":   driverV2.ObjectType_TRIGGER,
		"EVENT":     driverV2.ObjectType_EVENT,
	}
	for name, objType := range cases {
		t.Run(name, func(t *testing.T) {
			base := newFakeRunner().on("USE base_schema", "OK")
			compared := newFakeRunner().on("USE compared_schema", "OK")
			h := newDriverWithRunners(base, compared)
			results, err := h.GetDatabaseDiffModifySQL(context.Background(),
				&driverV2.DSN{},
				[]*driverV2.DatabasCompareSchemaInfo{{
					BaseSchemaName:     "base_schema",
					ComparedSchemaName: "compared_schema",
					DatabaseObjects: []*driverV2.DatabaseObject{
						{ObjectName: "x_obj", ObjectType: objType},
					},
				}})
			if err != nil {
				t.Fatalf("expected nil error for short-circuited type %s; got %v", objType, err)
			}
			if len(results) != 1 {
				t.Fatalf("expected 1 schema result, got %d", len(results))
			}
			full := strings.Join(results[0].ModifySQLs, "\n")
			// Only the USE header should be present; no DROP/CREATE/ALTER.
			if strings.Contains(full, "DROP") ||
				strings.Contains(full, "CREATE") ||
				strings.Contains(full, "ALTER") {
				t.Errorf("expected %s to be filtered; got block=%q", objType, full)
			}
			if !strings.Contains(full, "USE compared_schema;") {
				t.Errorf("expected USE compared_schema; header; got=%q", full)
			}
		})
	}
}

// --------------------------------------------------------------------- //
// 16. Test_DiffTableDDL_Matrix — design §3.4 row-by-row coverage
// --------------------------------------------------------------------- //

// Test_DiffTableDDL_Matrix walks every row of the design §3.4 ALTER vs
// DROP+CREATE decision matrix. Each case keys on a label that matches a
// matrix row; the body lists base/compared DDLs and the expected outcome
// (ALTER substring(s) or fallback=true with a WARNING expectation).
//
// Coverage by matrix row:
//   - 加列                           -> case "add column"
//   - 删列                           -> case "drop column"
//   - 改列名                         -> case "rename column"
//   - 改列类型 (兼容 widen)          -> case "widen int to bigint"
//   - 改列类型 (不兼容)              -> case "incompat type"
//   - 改列注释                       -> case "change column comment"
//   - 改表注释                       -> case "change table comment"
//   - 改 TBLPROPERTIES (业务字段)   -> case "alter tblproperties"
//   - 改分区键定义                   -> case "change partition key"
//   - 改存储格式 (STORED AS)         -> case "change stored as"
//   - 改 ROW FORMAT / SerDe          -> case "change row format"
//   - 改 EXTERNAL/MANAGED            -> case "toggle external"
//   - 改 LOCATION                    -> case "change location"
//   - 多类差异组合                   -> case "combined incompatible"
//   - 运行时 TBLPROPERTIES 不触发    -> case "runtime properties ignored"
func Test_DiffTableDDL_Matrix(t *testing.T) {
	cases := map[string]struct {
		base    string
		target  string
		wantAlt []string // substrings that must appear in alterStmts joined
		wantFB  bool    // expected fallbackDropCreate
	}{
		"add column": {
			base:    "CREATE TABLE t (id BIGINT, name STRING, age INT) STORED AS ORC",
			target:  "CREATE TABLE t (id BIGINT, name STRING) STORED AS ORC",
			wantAlt: []string{"ALTER TABLE t ADD COLUMNS (age INT)"},
		},
		"drop column": {
			base:   "CREATE TABLE t (id BIGINT) STORED AS ORC",
			target: "CREATE TABLE t (id BIGINT, name STRING) STORED AS ORC",
			wantFB: true,
		},
		"rename column": {
			base:    "CREATE TABLE t (id BIGINT, full_name STRING) STORED AS ORC",
			target:  "CREATE TABLE t (id BIGINT, name STRING) STORED AS ORC",
			wantFB:  true, // rename detected as "name removed + full_name added" → deletion path
		},
		"widen int to bigint": {
			base:    "CREATE TABLE t (id BIGINT) STORED AS ORC",
			target:  "CREATE TABLE t (id INT) STORED AS ORC",
			wantAlt: []string{"ALTER TABLE t CHANGE COLUMN id id BIGINT"},
		},
		"incompat type": {
			base:   "CREATE TABLE t (ts TIMESTAMP) STORED AS ORC",
			target: "CREATE TABLE t (ts INT) STORED AS ORC",
			wantFB: true,
		},
		"change column comment": {
			base:    "CREATE TABLE t (id BIGINT COMMENT 'new note') STORED AS ORC",
			target:  "CREATE TABLE t (id BIGINT COMMENT 'old note') STORED AS ORC",
			wantAlt: []string{"ALTER TABLE t CHANGE COLUMN id id BIGINT COMMENT 'new note'"},
		},
		"change table comment": {
			base:    "CREATE TABLE t (id BIGINT) COMMENT 'updated' STORED AS ORC",
			target:  "CREATE TABLE t (id BIGINT) COMMENT 'old' STORED AS ORC",
			wantAlt: []string{"SET TBLPROPERTIES ('comment'='updated')"},
		},
		"alter tblproperties": {
			base:    "CREATE TABLE t (id BIGINT) STORED AS ORC TBLPROPERTIES ('biz.owner'='alice')",
			target:  "CREATE TABLE t (id BIGINT) STORED AS ORC TBLPROPERTIES ('biz.owner'='bob')",
			wantAlt: []string{"SET TBLPROPERTIES ('biz.owner'='alice')"},
		},
		"change partition key": {
			base:   "CREATE TABLE t (id BIGINT) PARTITIONED BY (dt STRING) STORED AS ORC",
			target: "CREATE TABLE t (id BIGINT) PARTITIONED BY (dt STRING, region STRING) STORED AS ORC",
			wantFB: true,
		},
		"change stored as": {
			base:   "CREATE TABLE t (id BIGINT) STORED AS PARQUET",
			target: "CREATE TABLE t (id BIGINT) STORED AS ORC",
			wantFB: true,
		},
		"change row format": {
			base:   "CREATE TABLE t (id BIGINT) ROW FORMAT SERDE 'org.apache.hadoop.hive.ql.io.orc.OrcSerde' STORED AS ORC",
			target: "CREATE TABLE t (id BIGINT) ROW FORMAT DELIMITED FIELDS TERMINATED BY ',' STORED AS ORC",
			wantFB: true,
		},
		"toggle external": {
			base:    "CREATE EXTERNAL TABLE t (id BIGINT) STORED AS ORC",
			target:  "CREATE TABLE t (id BIGINT) STORED AS ORC",
			wantAlt: []string{"SET TBLPROPERTIES ('EXTERNAL'='TRUE')"},
		},
		"change location": {
			base:    "CREATE TABLE t (id BIGINT) STORED AS ORC LOCATION 'hdfs://new/path'",
			target:  "CREATE TABLE t (id BIGINT) STORED AS ORC LOCATION 'hdfs://old/path'",
			wantAlt: []string{"SET LOCATION 'hdfs://new/path'"},
		},
		"combined incompatible": {
			base:   "CREATE TABLE t (id BIGINT, age INT) PARTITIONED BY (dt STRING) STORED AS ORC",
			target: "CREATE TABLE t (id BIGINT) PARTITIONED BY (dt STRING) STORED AS PARQUET",
			wantFB: true, // STORED AS change alone is enough; combined with extra column makes it stronger
		},
		"runtime properties ignored": {
			base:   "CREATE TABLE t (id BIGINT) STORED AS ORC TBLPROPERTIES ('transient_lastDdlTime'='100')",
			target: "CREATE TABLE t (id BIGINT) STORED AS ORC TBLPROPERTIES ('transient_lastDdlTime'='200', 'numFiles'='5')",
			// No ALTER expected: runtime keys are filtered, so the schemas
			// compare equal.
		},
	}

	// Iterate in sorted-key order so failures are deterministic.
	keys := make([]string, 0, len(cases))
	for k := range cases {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, name := range keys {
		tc := cases[name]
		t.Run(name, func(t *testing.T) {
			alters, fallback, err := diffTableDDL(tc.base, tc.target)
			if err != nil {
				t.Fatalf("diffTableDDL error: %v", err)
			}
			if fallback != tc.wantFB {
				t.Errorf("fallback mismatch: got=%v want=%v\nbase=%q\ntarget=%q\nalters=%v",
					fallback, tc.wantFB, tc.base, tc.target, alters)
				return
			}
			joined := strings.Join(alters, "\n")
			for _, want := range tc.wantAlt {
				if !strings.Contains(joined, want) {
					t.Errorf("expected ALTER substring %q in:\n%s", want, joined)
				}
			}
			if tc.wantFB && len(alters) != 0 {
				t.Errorf("fallback case should have empty alterStmts, got %v", alters)
			}
		})
	}
}
