package state

// Known PostgresVersion values - use these for checks in version-dependent code
const (
	PostgresVersion10 = 100000
	PostgresVersion11 = 110000
	PostgresVersion12 = 120000
	PostgresVersion13 = 130000
	PostgresVersion14 = 140000
	PostgresVersion15 = 150000

	// MinRequiredPostgresVersion - We require PostgreSQL 10 or newer
	MinRequiredPostgresVersion = PostgresVersion10
)

// PostgresVersion - Identifying information about the PostgreSQL server version and build details
type PostgresVersion struct {
	Full    string `json:"full"`    // e.g. "PostgreSQL 9.5.1 on x86_64-pc-linux-gnu, compiled by gcc (Debian 4.9.2-10) 4.9.2, 64-bit"
	Short   string `json:"short"`   // e.g. "9.5.1"
	Numeric int    `json:"numeric"` // e.g. 90501

	// For collector use only, to avoid calling functions that don't work
	IsAwsAurora bool
	IsCitus     bool
}
