package hive

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sqleDriver "github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/sirupsen/logrus"
)

func TestHivePluginRegistered(t *testing.T) {
	processor, ok := sqleDriver.BuiltInPluginProcessors[driverV2.DriverTypeHive]
	if !ok {
		t.Fatalf("expected BuiltInPluginProcessors to contain key %q", driverV2.DriverTypeHive)
	}
	if processor == nil {
		t.Fatal("expected BuiltInPluginProcessors[Hive] to be non-nil")
	}
}

func TestGetDriverMetas(t *testing.T) {
	p := &PluginProcessor{}
	metas, err := p.GetDriverMetas()
	if err != nil {
		t.Fatalf("GetDriverMetas() returned error: %v", err)
	}

	// Verify PluginName
	if metas.PluginName != "Hive" {
		t.Errorf("expected PluginName=%q, got %q", "Hive", metas.PluginName)
	}

	// Verify DefaultPort
	if metas.DatabaseDefaultPort != 10000 {
		t.Errorf("expected DatabaseDefaultPort=10000, got %d", metas.DatabaseDefaultPort)
	}

	// Verify Rules is empty
	if len(metas.Rules) != 0 {
		t.Errorf("expected empty Rules, got %d rules", len(metas.Rules))
	}

	// Verify EnabledOptionalModule declares structure-compare capabilities (compat-RISK-1).
	// The set must contain OptionalGetDatabaseObjectDDL and OptionalGetDatabaseDiffModifySQL
	// for the controller/server capability check whitelist to accept Hive.
	expectedModules := map[driverV2.OptionalModule]bool{
		driverV2.OptionalGetDatabaseObjectDDL:     false,
		driverV2.OptionalGetDatabaseDiffModifySQL: false,
	}
	for _, m := range metas.EnabledOptionalModule {
		if _, ok := expectedModules[m]; ok {
			expectedModules[m] = true
		}
	}
	for m, seen := range expectedModules {
		if !seen {
			t.Errorf("expected EnabledOptionalModule to contain %v", m)
		}
	}

	// Verify additionalParams: auth
	authParam := metas.DatabaseAdditionalParams.GetParam("auth")
	if authParam == nil {
		t.Fatal("expected additionalParams to contain 'auth' param")
	}
	if authParam.Value != "NONE" {
		t.Errorf("expected auth default value=%q, got %q", "NONE", authParam.Value)
	}
	expectedAuthEnums := []string{"NONE", "NOSASL", "LDAP", "KERBEROS"}
	if len(authParam.Enums) != len(expectedAuthEnums) {
		t.Fatalf("expected %d auth enums, got %d", len(expectedAuthEnums), len(authParam.Enums))
	}
	for i, expected := range expectedAuthEnums {
		if authParam.Enums[i].Value != expected {
			t.Errorf("auth enum[%d]: expected %q, got %q", i, expected, authParam.Enums[i].Value)
		}
	}

	// Verify additionalParams: transport_mode
	transportParam := metas.DatabaseAdditionalParams.GetParam("transport_mode")
	if transportParam == nil {
		t.Fatal("expected additionalParams to contain 'transport_mode' param")
	}
	if transportParam.Value != "binary" {
		t.Errorf("expected transport_mode default value=%q, got %q", "binary", transportParam.Value)
	}
	expectedTransportEnums := []string{"binary", "http"}
	if len(transportParam.Enums) != len(expectedTransportEnums) {
		t.Fatalf("expected %d transport_mode enums, got %d", len(expectedTransportEnums), len(transportParam.Enums))
	}
	for i, expected := range expectedTransportEnums {
		if transportParam.Enums[i].Value != expected {
			t.Errorf("transport_mode enum[%d]: expected %q, got %q", i, expected, transportParam.Enums[i].Value)
		}
	}

	// Verify additionalParams: service
	serviceParam := metas.DatabaseAdditionalParams.GetParam("service")
	if serviceParam == nil {
		t.Fatal("expected additionalParams to contain 'service' param")
	}
	if serviceParam.Value != "" {
		t.Errorf("expected service default value=%q, got %q", "", serviceParam.Value)
	}
}

func TestClassifySQL(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected string
	}{
		// DQL cases
		"SELECT uppercase":           {input: "SELECT * FROM t", expected: driverV2.SQLTypeDQL},
		"select lowercase":           {input: "select id from t", expected: driverV2.SQLTypeDQL},
		"WITH CTE":                   {input: "WITH cte AS (SELECT 1) SELECT * FROM cte", expected: driverV2.SQLTypeDQL},
		"SHOW TABLES":                {input: "SHOW TABLES", expected: driverV2.SQLTypeDQL},
		"DESCRIBE table":             {input: "DESCRIBE my_table", expected: driverV2.SQLTypeDQL},
		"DESC table":                 {input: "DESC my_table", expected: driverV2.SQLTypeDQL},
		"EXPLAIN query":              {input: "EXPLAIN SELECT 1", expected: driverV2.SQLTypeDQL},
		"leading whitespace SELECT":  {input: "  SELECT 1", expected: driverV2.SQLTypeDQL},

		// DML cases
		"INSERT":                     {input: "INSERT INTO t VALUES (1)", expected: driverV2.SQLTypeDML},
		"UPDATE":                     {input: "UPDATE t SET a=1", expected: driverV2.SQLTypeDML},
		"DELETE":                     {input: "DELETE FROM t WHERE id=1", expected: driverV2.SQLTypeDML},
		"MERGE":                      {input: "MERGE INTO t USING s ON t.id=s.id", expected: driverV2.SQLTypeDML},
		"LOAD":                       {input: "LOAD DATA INPATH '/path' INTO TABLE t", expected: driverV2.SQLTypeDML},
		"EXPORT":                     {input: "EXPORT TABLE t TO '/path'", expected: driverV2.SQLTypeDML},

		// DDL cases (default)
		"CREATE TABLE":               {input: "CREATE TABLE t (id INT)", expected: driverV2.SQLTypeDDL},
		"ALTER TABLE":                {input: "ALTER TABLE t ADD COLUMNS (col STRING)", expected: driverV2.SQLTypeDDL},
		"DROP TABLE":                 {input: "DROP TABLE t", expected: driverV2.SQLTypeDDL},
		"GRANT":                      {input: "GRANT SELECT ON t TO user", expected: driverV2.SQLTypeDDL},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := classifySQL(tc.input)
			if got != tc.expected {
				t.Errorf("classifySQL(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestSplitSQL(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected []string
	}{
		"single SQL": {
			input:    "SELECT 1",
			expected: []string{"SELECT 1"},
		},
		"multiple SQLs": {
			input:    "SELECT 1; SELECT 2; SELECT 3",
			expected: []string{"SELECT 1", "SELECT 2", "SELECT 3"},
		},
		"trailing semicolon": {
			input:    "SELECT 1;",
			expected: []string{"SELECT 1"},
		},
		"empty input": {
			input:    "",
			expected: []string{},
		},
		"whitespace only": {
			input:    "   ;  ;  ",
			expected: []string{},
		},
		"mixed with whitespace": {
			input:    "  SELECT 1 ;  INSERT INTO t VALUES(1)  ;  ",
			expected: []string{"SELECT 1", "INSERT INTO t VALUES(1)"},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := splitSQL(tc.input)
			if len(got) != len(tc.expected) {
				t.Fatalf("splitSQL(%q): got %d results, want %d", tc.input, len(got), len(tc.expected))
			}
			for i, s := range got {
				if s != tc.expected[i] {
					t.Errorf("splitSQL(%q)[%d] = %q, want %q", tc.input, i, s, tc.expected[i])
				}
			}
		})
	}
}

func TestAuditReturnsEmptyResults(t *testing.T) {
	p := &PluginProcessor{}
	plugin, err := p.Open(logrus.NewEntry(logrus.New()), &driverV2.Config{})
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	impl := plugin.(*HiveDriverImpl)

	sqls := []string{"SELECT 1", "INSERT INTO t VALUES(1)", "CREATE TABLE t(id INT)"}
	results, err := impl.Audit(context.Background(), sqls)
	if err != nil {
		t.Fatalf("Audit() returned error: %v", err)
	}
	if len(results) != len(sqls) {
		t.Fatalf("Audit() returned %d results, want %d", len(results), len(sqls))
	}
	for i, r := range results {
		if r == nil {
			t.Errorf("Audit() result[%d] is nil", i)
		}
	}
}

func TestParse(t *testing.T) {
	p := &PluginProcessor{}
	plugin, err := p.Open(logrus.NewEntry(logrus.New()), &driverV2.Config{})
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	impl := plugin.(*HiveDriverImpl)

	cases := map[string]struct {
		input         string
		expectedCount int
		expectedTypes []string
	}{
		"single DQL": {
			input:         "SELECT 1",
			expectedCount: 1,
			expectedTypes: []string{driverV2.SQLTypeDQL},
		},
		"multiple mixed": {
			input:         "SELECT 1; INSERT INTO t VALUES(1); CREATE TABLE t(id INT)",
			expectedCount: 3,
			expectedTypes: []string{driverV2.SQLTypeDQL, driverV2.SQLTypeDML, driverV2.SQLTypeDDL},
		},
		"empty input": {
			input:         "",
			expectedCount: 0,
			expectedTypes: []string{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			nodes, err := impl.Parse(context.Background(), tc.input)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tc.input, err)
			}
			if len(nodes) != tc.expectedCount {
				t.Fatalf("Parse(%q) returned %d nodes, want %d", tc.input, len(nodes), tc.expectedCount)
			}
			for i, node := range nodes {
				if node.Type != tc.expectedTypes[i] {
					t.Errorf("Parse(%q) node[%d].Type = %q, want %q", tc.input, i, node.Type, tc.expectedTypes[i])
				}
			}
		})
	}
}

func TestPingWithNilDSN(t *testing.T) {
	p := &PluginProcessor{}
	plugin, err := p.Open(logrus.NewEntry(logrus.New()), &driverV2.Config{})
	if err != nil {
		t.Fatalf("Open() returned error: %v", err)
	}
	impl := plugin.(*HiveDriverImpl)

	err = impl.Ping(context.Background())
	if err == nil {
		t.Error("expected Ping() to return error when dsn is nil")
	}
}

func TestOpenWithNilDSN(t *testing.T) {
	p := &PluginProcessor{}
	plugin, err := p.Open(logrus.NewEntry(logrus.New()), &driverV2.Config{DSN: nil})
	if err != nil {
		t.Fatalf("Open() with nil DSN should succeed in offline mode, got error: %v", err)
	}
	if plugin == nil {
		t.Fatal("Open() returned nil plugin")
	}
}

func TestPingWithNilConn(t *testing.T) {
	// When Open is called with nil DSN (offline audit mode), conn is nil.
	// Ping should return an error indicating uninitialized connection.
	impl := &HiveDriverImpl{
		log: logrus.NewEntry(logrus.New()),
	}
	err := impl.Ping(context.Background())
	if err == nil {
		t.Error("expected Ping() to return error when conn is nil")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected error to mention 'not initialized', got: %v", err)
	}
}

func TestCloseWithNilConn(t *testing.T) {
	// Close should not panic when conn is nil (offline audit mode).
	impl := &HiveDriverImpl{
		log: logrus.NewEntry(logrus.New()),
	}
	// Should not panic
	impl.Close(context.Background())
}

// fakeHiveCursor lets tests drive fetchAllRows through fully-scripted
// HasMore / FetchOne / Err sequences. It mimics how HiveServer2 surfaces a
// non-fatal ROW-ERR after a no-result-column statement (USE, SET, DDL).
//
// Each step in `steps` represents one HasMore tick. Setting HasMore=false
// signals end of stream; setting Err to a non-nil value triggers the
// tolerance branch in fetchAllRows. If FetchValue is non-empty, FetchOne
// will assign it to the destination string before the next HasMore tick.
type fakeHiveCursor struct {
	steps    []fakeCursorStep
	idx      int
	err      error
	fetchVal string
}

type fakeCursorStep struct {
	HasMore       bool
	FetchValue    string
	ErrBeforeNext error // err returned by Err() *after* this step's HasMore
}

func (f *fakeHiveCursor) HasMore(ctx context.Context) bool {
	if f.idx >= len(f.steps) {
		return false
	}
	step := f.steps[f.idx]
	f.fetchVal = step.FetchValue
	// Error surfaces immediately on entering the next iteration so the
	// fetchAllRows branch can pick it up before FetchOne.
	f.err = step.ErrBeforeNext
	return step.HasMore
}

func (f *fakeHiveCursor) FetchOne(ctx context.Context, dests ...interface{}) {
	if f.idx < len(f.steps) && len(dests) > 0 {
		if p, ok := dests[0].(*string); ok {
			*p = f.fetchVal
		}
	}
	f.idx++
}

func (f *fakeHiveCursor) Err() error { return f.err }

func Test_IsHS2NoResultRowErr(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil",
			err:  nil,
			want: false,
		},
		{
			name: "tolerable_HS2_row_err",
			// Real-world payload from HiveServer2 FetchResults on a USE statement.
			err: fmt.Errorf("TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.] SqlState:<nil> ErrorCode:<nil> ErrorMessage:<nil>})"),
			want: true,
		},
		{
			name: "syntax_error_not_tolerable",
			err:  fmt.Errorf("FAILED: ParseException line 1:0 cannot recognize input near 'SELEKT'"),
			want: false,
		},
		{
			name: "missing_marker_status_only",
			err:  fmt.Errorf("StatusCode:ERROR_STATUS but message differs"),
			want: false,
		},
		{
			name: "missing_status_marker",
			err:  fmt.Errorf("Server-side error; please check HS2 logs."),
			want: false,
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := isHS2NoResultRowErr(tc.err)
			if got != tc.want {
				t.Errorf("isHS2NoResultRowErr(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// Test_FetchAllRows_RowErrTolerant exercises the HS2 ROW-ERR tolerance
// contract for compat-RISK-10. The three scenarios verify:
//
//  1. USE-like statement returns ROW-ERR with zero rows -> fetchAllRows
//     yields (nil, nil) (no hard error; caller can keep executing).
//  2. SHOW TABLES yields two rows then a trailing ROW-ERR -> rows are
//     returned and the ROW-ERR is treated as EOF (compat-RISK-10 must
//     not drop already-fetched rows).
//  3. A real Hive runtime error (syntax error) is propagated -> caller
//     receives the wrapped failure as before the fix.
func Test_FetchAllRows_RowErrTolerant(t *testing.T) {
	hs2RowErr := fmt.Errorf("TStatus({StatusCode:ERROR_STATUS InfoMessages:[Server-side error; please check HS2 logs.]})")
	syntaxErr := fmt.Errorf("FAILED: ParseException syntax error")

	cases := map[string]struct {
		steps    []fakeCursorStep
		wantRows []string
		wantErr  bool
		errMatch string
	}{
		"USE_statement_row_err_tolerated": {
			// First HasMore tick surfaces ROW-ERR with no row -> loop breaks.
			steps: []fakeCursorStep{
				{HasMore: true, ErrBeforeNext: hs2RowErr},
			},
			wantRows: nil,
			wantErr:  false,
		},
		"SHOW_TABLES_rows_then_trailing_row_err": {
			// Two real rows arrive; the third HasMore tick is the HS2
			// terminator with ROW-ERR. Tolerated -> existing rows preserved.
			steps: []fakeCursorStep{
				{HasMore: true, FetchValue: "t_base_only"},
				{HasMore: true, FetchValue: "t_diff_only"},
				{HasMore: true, ErrBeforeNext: hs2RowErr},
			},
			wantRows: []string{"t_base_only", "t_diff_only"},
			wantErr:  false,
		},
		"real_syntax_error_propagates": {
			// A non-ROW-ERR is a genuine failure and must surface.
			steps: []fakeCursorStep{
				{HasMore: true, ErrBeforeNext: syntaxErr},
			},
			wantRows: nil,
			wantErr:  true,
			errMatch: "ParseException",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			cur := &fakeHiveCursor{steps: tc.steps}
			rows, err := fetchAllRows(context.Background(), cur)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil (rows=%v)", rows)
				}
				if tc.errMatch != "" && !strings.Contains(err.Error(), tc.errMatch) {
					t.Errorf("err = %v, want substring %q", err, tc.errMatch)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(rows) != len(tc.wantRows) {
				t.Fatalf("rows len = %d, want %d (rows=%v)", len(rows), len(tc.wantRows), rows)
			}
			for i, r := range rows {
				if r != tc.wantRows[i] {
					t.Errorf("rows[%d] = %q, want %q", i, r, tc.wantRows[i])
				}
			}
		})
	}
}

// Test_RunSingleStringQuery_NilConn verifies the early-return guard
// continues to work after the refactor.
func Test_RunSingleStringQuery_NilConn(t *testing.T) {
	g := &gohiveQueryRunner{conn: nil}
	_, err := g.runSingleStringQuery(context.Background(), "USE default")
	if err == nil {
		t.Fatal("expected error when conn is nil")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("err = %v, want substring %q", err, "not initialized")
	}
}

// scriptedRunner is a hiveQueryRunner whose runSingleStringQuery returns
// pre-baked rows / errors keyed by exact query text. Used to verify the
// listAllSchemaObjects helper and the "default discovery" behaviour of
// GetDatabaseObjectDDL without needing a live HiveServer2.
type scriptedRunner struct {
	scripts map[string]scriptedReply
	calls   []string
}

type scriptedReply struct {
	rows []string
	err  error
}

func (s *scriptedRunner) runSingleStringQuery(ctx context.Context, query string) ([]string, error) {
	s.calls = append(s.calls, query)
	r, ok := s.scripts[query]
	if !ok {
		return nil, fmt.Errorf("scriptedRunner: no script for %q", query)
	}
	return r.rows, r.err
}

func Test_ListAllSchemaObjects(t *testing.T) {
	t.Run("tables_and_views_classified", func(t *testing.T) {
		runner := &scriptedRunner{scripts: map[string]scriptedReply{
			"SHOW TABLES": {rows: []string{"t_base_only", "v_user_summary", "t_alter_widen"}},
			"SHOW VIEWS":  {rows: []string{"v_user_summary"}},
		}}
		objs, err := listAllSchemaObjects(context.Background(), runner)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if len(objs) != 3 {
			t.Fatalf("len=%d want 3 (%v)", len(objs), objs)
		}
		got := map[string]string{}
		for _, o := range objs {
			got[o.ObjectName] = o.ObjectType
		}
		if got["v_user_summary"] != driverV2.ObjectType_VIEW {
			t.Errorf("v_user_summary type = %q want VIEW", got["v_user_summary"])
		}
		if got["t_base_only"] != driverV2.ObjectType_TABLE {
			t.Errorf("t_base_only type = %q want TABLE", got["t_base_only"])
		}
		if got["t_alter_widen"] != driverV2.ObjectType_TABLE {
			t.Errorf("t_alter_widen type = %q want TABLE", got["t_alter_widen"])
		}
	})
	t.Run("show_views_failure_degrades_to_all_table", func(t *testing.T) {
		// Older Hive (< 2.2) does not support SHOW VIEWS; we must not blow up.
		runner := &scriptedRunner{scripts: map[string]scriptedReply{
			"SHOW TABLES": {rows: []string{"t1", "v1"}},
			"SHOW VIEWS":  {err: fmt.Errorf("unsupported on old HS2")},
		}}
		objs, err := listAllSchemaObjects(context.Background(), runner)
		if err != nil {
			t.Fatalf("expected tolerant fallback, got err: %v", err)
		}
		if len(objs) != 2 {
			t.Fatalf("len=%d want 2", len(objs))
		}
		for _, o := range objs {
			if o.ObjectType != driverV2.ObjectType_TABLE {
				t.Errorf("%s degraded type = %q want TABLE", o.ObjectName, o.ObjectType)
			}
		}
	})
	t.Run("show_tables_error_propagates", func(t *testing.T) {
		runner := &scriptedRunner{scripts: map[string]scriptedReply{
			"SHOW TABLES": {err: fmt.Errorf("real fetch error")},
		}}
		_, err := listAllSchemaObjects(context.Background(), runner)
		if err == nil {
			t.Fatal("expected err to propagate")
		}
		if !strings.Contains(err.Error(), "show tables") {
			t.Errorf("err = %v want substring %q", err, "show tables")
		}
	})
	t.Run("empty_names_skipped", func(t *testing.T) {
		runner := &scriptedRunner{scripts: map[string]scriptedReply{
			"SHOW TABLES": {rows: []string{"", "t1", ""}},
			"SHOW VIEWS":  {rows: []string{}},
		}}
		objs, err := listAllSchemaObjects(context.Background(), runner)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if len(objs) != 1 || objs[0].ObjectName != "t1" {
			t.Errorf("expected single t1, got %v", objs)
		}
	})
}

func Test_GetDatabaseObjectDDL_DefaultDiscovery(t *testing.T) {
	// When caller passes objInfo.DatabaseObjects = nil, the driver must
	// auto-enumerate via listAllSchemaObjects before producing DDLs.
	runner := &scriptedRunner{scripts: map[string]scriptedReply{
		"USE default":                    {rows: nil},
		"SHOW TABLES":                    {rows: []string{"t_base_only", "v_user_summary"}},
		"SHOW VIEWS":                     {rows: []string{"v_user_summary"}},
		"SHOW CREATE TABLE t_base_only":  {rows: []string{"CREATE TABLE `t_base_only` (`id` int)"}},
		"SHOW CREATE TABLE v_user_summary": {rows: []string{"CREATE VIEW `v_user_summary` AS SELECT 1"}},
	}}
	impl := &HiveDriverImpl{runner: runner}

	res, err := impl.GetDatabaseObjectDDL(context.Background(), []*driverV2.DatabaseSchemaInfo{
		{SchemaName: "default"}, // DatabaseObjects intentionally nil
	})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(res) != 1 {
		t.Fatalf("len(res)=%d want 1", len(res))
	}
	if res[0].SchemaName != "default" {
		t.Errorf("schema = %q want default", res[0].SchemaName)
	}
	if len(res[0].DatabaseObjectDDLs) != 2 {
		t.Fatalf("ddls len = %d want 2 (auto-discovery): %v", len(res[0].DatabaseObjectDDLs), res[0].DatabaseObjectDDLs)
	}
	gotTypes := map[string]string{}
	for _, d := range res[0].DatabaseObjectDDLs {
		gotTypes[d.DatabaseObject.ObjectName] = d.DatabaseObject.ObjectType
	}
	if gotTypes["t_base_only"] != driverV2.ObjectType_TABLE {
		t.Errorf("t_base_only type = %q want TABLE", gotTypes["t_base_only"])
	}
	if gotTypes["v_user_summary"] != driverV2.ObjectType_VIEW {
		t.Errorf("v_user_summary type = %q want VIEW", gotTypes["v_user_summary"])
	}
}

// fakeExecRunner records every Exec invocation and lets tests script the
// per-query error response. It is the unit-test injection point for the
// Exec / ExecBatch contract on HiveDriverImpl.
type fakeExecRunner struct {
	calls   []string
	errs    map[string]error
	defaultErr error
}

func (f *fakeExecRunner) exec(ctx context.Context, query string) error {
	f.calls = append(f.calls, query)
	if e, ok := f.errs[query]; ok {
		return e
	}
	return f.defaultErr
}

func Test_StripSQLTerminator(t *testing.T) {
	cases := map[string]struct {
		in   string
		want string
	}{
		"plain":              {in: "SELECT 1", want: "SELECT 1"},
		"trailing semicolon": {in: "SELECT 1;", want: "SELECT 1"},
		"trailing whitespace_and_semicolons": {in: " SELECT 1 ;;\n ", want: "SELECT 1"},
		"only_semicolons":    {in: ";;;", want: ""},
		"empty":              {in: "", want: ""},
		"whitespace_only":    {in: "   \n  ", want: ""},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := stripSQLTerminator(tc.in)
			if got != tc.want {
				t.Errorf("stripSQLTerminator(%q) = %q want %q", tc.in, got, tc.want)
			}
		})
	}
}

func Test_IsAllCommentLines(t *testing.T) {
	cases := map[string]struct {
		in   string
		want bool
	}{
		"empty":                {in: "", want: true},
		"single_comment":       {in: "-- hello", want: true},
		"multiline_comments":   {in: "-- WARNING: data loss risk\n-- second comment", want: true},
		"blank_lines_and_comment": {in: "\n  -- only comment\n  \n", want: true},
		"mixed":                {in: "-- WARNING\nDROP TABLE t", want: false},
		"sql_only":             {in: "DROP TABLE t", want: false},
	}
	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			got := isAllCommentLines(tc.in)
			if got != tc.want {
				t.Errorf("isAllCommentLines(%q) = %v want %v", tc.in, got, tc.want)
			}
		})
	}
}

// Test_Exec_SingleStatement covers the happy path: a DDL is forwarded to
// the runner without alteration; the trailing semicolon (if present) is
// stripped before submission; the returned hiveExecResult exposes the
// "not supported" contract for LastInsertId / RowsAffected.
func Test_Exec_SingleStatement(t *testing.T) {
	runner := &fakeExecRunner{}
	impl := &HiveDriverImpl{
		execRunnerFactory: func(_ *HiveDriverImpl) hiveExecRunner { return runner },
	}
	res, err := impl.Exec(context.Background(), "DROP TABLE IF EXISTS sqle_compare_test.t_diff_only;")
	if err != nil {
		t.Fatalf("Exec returned err: %v", err)
	}
	if res == nil {
		t.Fatal("Exec returned nil result")
	}
	if _, err := res.LastInsertId(); err == nil {
		t.Error("expected LastInsertId to return an error (hive not supported)")
	}
	if _, err := res.RowsAffected(); err == nil {
		t.Error("expected RowsAffected to return an error (hive not supported)")
	}
	if len(runner.calls) != 1 {
		t.Fatalf("expected 1 runner.calls, got %d (%v)", len(runner.calls), runner.calls)
	}
	if runner.calls[0] != "DROP TABLE IF EXISTS sqle_compare_test.t_diff_only" {
		t.Errorf("expected trailing-; stripped, got %q", runner.calls[0])
	}
}

func Test_Exec_EmptyAndCommentStatementsAreNoOp(t *testing.T) {
	runner := &fakeExecRunner{defaultErr: fmt.Errorf("runner must not be called")}
	impl := &HiveDriverImpl{
		execRunnerFactory: func(_ *HiveDriverImpl) hiveExecRunner { return runner },
	}
	cases := []string{
		"",
		"   ",
		";",
		";;\n;",
		"-- WARNING: data loss risk",
		"-- 警告: 数据将丢失\n-- second comment",
	}
	for _, q := range cases {
		q := q
		t.Run(fmt.Sprintf("noop_%q", q), func(t *testing.T) {
			res, err := impl.Exec(context.Background(), q)
			if err != nil {
				t.Fatalf("Exec(%q) returned err: %v", q, err)
			}
			if res == nil {
				t.Fatalf("Exec(%q) returned nil result", q)
			}
		})
	}
	if len(runner.calls) != 0 {
		t.Errorf("expected runner to receive zero calls for no-op queries, got %v", runner.calls)
	}
}

func Test_Exec_PropagatesRunnerError(t *testing.T) {
	runner := &fakeExecRunner{
		errs: map[string]error{
			"CREATE TABLE x(a INT)": fmt.Errorf("FAILED: SemanticException [Error 10001]: Table x already exists"),
		},
	}
	impl := &HiveDriverImpl{
		execRunnerFactory: func(_ *HiveDriverImpl) hiveExecRunner { return runner },
	}
	_, err := impl.Exec(context.Background(), "CREATE TABLE x(a INT)")
	if err == nil {
		t.Fatal("expected Exec to return runner err")
	}
	if !strings.Contains(err.Error(), "Table x already exists") {
		t.Errorf("expected runner err to be wrapped, got: %v", err)
	}
}

func Test_Exec_NilConnAndNoFactoryFails(t *testing.T) {
	impl := &HiveDriverImpl{}
	_, err := impl.Exec(context.Background(), "DROP TABLE t")
	if err == nil {
		t.Fatal("expected Exec to fail when both conn and factory are nil")
	}
	if !strings.Contains(err.Error(), "not initialized") {
		t.Errorf("expected 'not initialized', got: %v", err)
	}
}

// Test_ExecBatch_AllSucceed verifies the batch contract: every statement is
// forwarded in order, results are returned 1:1, and the per-statement
// trailing-; is stripped consistently with Exec.
func Test_ExecBatch_AllSucceed(t *testing.T) {
	runner := &fakeExecRunner{}
	impl := &HiveDriverImpl{
		execRunnerFactory: func(_ *HiveDriverImpl) hiveExecRunner { return runner },
	}
	sqls := []string{
		"USE sqle_compare_test;",
		"ALTER TABLE t_alter_widen CHANGE COLUMN amt amt BIGINT;",
		"-- WARNING: data loss risk", // comment-only -> no-op
		"DROP TABLE IF EXISTS t_diff_only;",
	}
	results, err := impl.ExecBatch(context.Background(), sqls...)
	if err != nil {
		t.Fatalf("ExecBatch returned err: %v", err)
	}
	if len(results) != len(sqls) {
		t.Fatalf("expected %d results, got %d", len(sqls), len(results))
	}
	for i, r := range results {
		if r == nil {
			t.Errorf("results[%d] is nil", i)
		}
	}
	// Comment-only statement is filtered before reaching runner.
	wantCalls := []string{
		"USE sqle_compare_test",
		"ALTER TABLE t_alter_widen CHANGE COLUMN amt amt BIGINT",
		"DROP TABLE IF EXISTS t_diff_only",
	}
	if len(runner.calls) != len(wantCalls) {
		t.Fatalf("runner.calls = %v want %v", runner.calls, wantCalls)
	}
	for i, want := range wantCalls {
		if runner.calls[i] != want {
			t.Errorf("runner.calls[%d] = %q want %q", i, runner.calls[i], want)
		}
	}
}

// Test_ExecBatch_StopsOnFirstError mirrors MySQL driver behaviour: on the
// first error the batch returns the partial result set plus a wrapped err,
// without executing any subsequent statement.
func Test_ExecBatch_StopsOnFirstError(t *testing.T) {
	runner := &fakeExecRunner{
		errs: map[string]error{
			"DROP TABLE IF EXISTS t_diff_only": fmt.Errorf("permission denied"),
		},
	}
	impl := &HiveDriverImpl{
		execRunnerFactory: func(_ *HiveDriverImpl) hiveExecRunner { return runner },
	}
	sqls := []string{
		"USE sqle_compare_test",
		"DROP TABLE IF EXISTS t_diff_only", // fails here
		"CREATE TABLE never_runs (id INT)",
	}
	results, err := impl.ExecBatch(context.Background(), sqls...)
	if err == nil {
		t.Fatal("expected ExecBatch to return err on first failure")
	}
	if !strings.Contains(err.Error(), "permission denied") {
		t.Errorf("expected err to wrap runner err, got: %v", err)
	}
	// Two results: one for USE (nil result is hiveExecResult{}), one for failed DROP (also returned).
	if len(results) != 2 {
		t.Fatalf("expected 2 partial results, got %d", len(results))
	}
	// The CREATE TABLE statement must NOT be invoked after the failure.
	if len(runner.calls) != 2 {
		t.Errorf("runner.calls = %v, expected 2 (USE + DROP), CREATE must be skipped", runner.calls)
	}
}
