package state

import "github.com/guregu/null"

// PostgresRole - A role in the PostgreSQL system - note that this includes users (Login=true)
type PostgresRole struct {
	Oid                Oid       // ID of role
	Name               string    // Role name
	Inherit            bool      // Role automatically inherits privileges of roles it is a member of
	Login              bool      // Role can log in. That is, this role can be given as the initial session authorization identifier
	CreateDb           bool      // Role can create databases
	CreateRole         bool      // Role can create more roles
	SuperUser          bool      // Role has superuser privileges
	Replication        bool      // Role can initiate streaming replication and put the system in and out of backup mode.
	BypassRLS          bool      // Role bypasses every row level security policy, see https://www.postgresql.org/docs/9.5/static/ddl-rowsecurity.html
	ConnectionLimit    int32     // For roles that can log in, this sets maximum number of concurrent connections this role can make. -1 means no limit.
	PasswordValidUntil null.Time // Password expiry time (only used for password authentication); null if no expiration
	Config             []string  // Role-specific defaults for run-time configuration variables
	MemberOf           []Oid     // List of roles that this role is a member of (i.e. whose permissions it inherits)
}
