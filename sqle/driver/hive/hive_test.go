package hive

import (
	"context"
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

	// Verify EnabledOptionalModule is empty
	if len(metas.EnabledOptionalModule) != 0 {
		t.Errorf("expected empty EnabledOptionalModule, got %d modules", len(metas.EnabledOptionalModule))
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
