package differ

import (
	"fmt"
)

// Check represents a single check constraint in a table.
type Check struct {
	Name     string `json:"name"`
	Clause   string `json:"clause"`
	Enforced bool   `json:"enforced"` // Always true in MariaDB
}

// Definition returns this Check's definition clause, for use as part of a DDL
// statement.
func (cc *Check) Definition(flavor Flavor) string {
	var notEnforced string
	if !cc.Enforced {
		notEnforced = " /*!80016 NOT ENFORCED */"
	}
	return fmt.Sprintf("CONSTRAINT %s CHECK (%s)%s", EscapeIdentifier(cc.Name), cc.Clause, notEnforced)
}
