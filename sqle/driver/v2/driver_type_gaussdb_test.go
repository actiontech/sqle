package driverV2

import (
	"testing"
)

// TestDriverTypeGaussDB_literal hard-asserts the cross-repo literal contract:
//   - sqle-ee   DriverTypeGaussDB    = "GaussDB"
//   - dms-ee    DBTypeGaussDB        = "GaussDB"   (Task-D04 whitelist)
//   - odc       DialectType.GAUSSDB.name() == "GAUSSDB"
//     (dms-ee convertDBType normalizes "GaussDB" -> "GAUSSDB" at the boundary)
//   - odc-client ConnectType.GAUSSDB = 'GAUSSDB'   (Task-D03)
//
// Breaking the literal here breaks every cross-repo routing path.
func TestDriverTypeGaussDB_literal(t *testing.T) {
	cases := map[string]struct {
		got      string
		expected string
	}{
		"DriverTypeGaussDB equals \"GaussDB\"": {
			got:      DriverTypeGaussDB,
			expected: "GaussDB",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.got != tc.expected {
				t.Errorf("got %q, want %q (cross-repo literal contract)", tc.got, tc.expected)
			}
		})
	}
}

// TestDriverTypeGaussDB_distinct_from_PostgreSQL covers CR-7 / decision-3:
// GaussDB and PostgreSQL must remain independent constants so the SQLE plugin
// router never folds GaussDB into the PostgreSQL branch.
func TestDriverTypeGaussDB_distinct_from_PostgreSQL(t *testing.T) {
	if DriverTypeGaussDB == DriverTypePostgreSQL {
		t.Errorf("DriverTypeGaussDB (%q) must be distinct from DriverTypePostgreSQL (%q)",
			DriverTypeGaussDB, DriverTypePostgreSQL)
	}
	if DriverTypePostgreSQL != "PostgreSQL" {
		t.Errorf("DriverTypePostgreSQL = %q, want %q (PG regression guard)",
			DriverTypePostgreSQL, "PostgreSQL")
	}
}

// TestDriverTypeGaussDB_distinct_from_other_drivers is a map-case
// anti-collision check: GaussDB must not accidentally equal any other
// existing DriverType constant.
func TestDriverTypeGaussDB_distinct_from_other_drivers(t *testing.T) {
	others := map[string]string{
		"MySQL":          DriverTypeMySQL,
		"TiDB":           DriverTypeTiDB,
		"SQL Server":     DriverTypeSQLServer,
		"Oracle":         DriverTypeOracle,
		"DB2":            DriverTypeDB2,
		"OceanBase":      DriverTypeOceanBase,
		"TDSQLForInnoDB": DriverTypeTDSQLForInnoDB,
		"TBase":          DriverTypeTBase,
		"HANA":           DriverTypeHANA,
	}
	for name, val := range others {
		t.Run(name, func(t *testing.T) {
			if val == DriverTypeGaussDB {
				t.Errorf("DriverType%s (%q) collides with DriverTypeGaussDB (%q)",
					name, val, DriverTypeGaussDB)
			}
		})
	}
}

// TestDriverTypeGaussDB_constant_table_consistency enumerates the entire
// DriverType* constant table (util.go L27-39) and asserts the literal
// "GaussDB" appears exactly once. Adding another constant that shadows
// "GaussDB" -- or removing DriverTypeGaussDB -- breaks this test.
func TestDriverTypeGaussDB_constant_table_consistency(t *testing.T) {
	all := []string{
		DriverTypeMySQL,
		DriverTypePostgreSQL,
		DriverTypeTiDB,
		DriverTypeSQLServer,
		DriverTypeOracle,
		DriverTypeDB2,
		DriverTypeOceanBase,
		DriverTypeTDSQLForInnoDB,
		DriverTypeTBase,
		DriverTypeHANA,
		DriverTypeGaussDB,
	}
	count := 0
	for _, d := range all {
		if d == "GaussDB" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("DriverType=GaussDB occurrence: got %d, want 1", count)
	}
}
