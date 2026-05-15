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
	if authParam.Value != "NOSASL" {
		t.Errorf("expected auth default value=%q, got %q", "NOSASL", authParam.Value)
	}
	expectedAuthEnums := []string{"NOSASL", "NONE", "LDAP", "KERBEROS"}
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
