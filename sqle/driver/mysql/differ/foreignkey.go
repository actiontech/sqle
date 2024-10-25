package differ

import (
	"fmt"
	"strings"
)

// ForeignKey represents a single foreign key constraint in a table. Note that
// the "referenced" side of the FK is tracked as strings, rather than *Schema,
// *Table, *[]Column to avoid potentially having to introspect multiple schemas
// in a particular order. Also, the referenced side is not gauranteed to exist,
// especially if foreign_key_checks=0 has been used at any point in the past.
type ForeignKey struct {
	Name                  string   `json:"name"`
	ColumnNames           []string `json:"columnNames"`
	ReferencedSchemaName  string   `json:"referencedSchemaName,omitempty"` // will be empty string if same schema
	ReferencedTableName   string   `json:"referencedTableName"`
	ReferencedColumnNames []string `json:"referencedColumnNames"` // slice length always identical to len(ColumnNames)
	UpdateRule            string   `json:"updateRule"`
	DeleteRule            string   `json:"deleteRule"`
}

// Definition returns this ForeignKey's definition clause, for use as part of a DDL
// statement.
func (fk *ForeignKey) Definition(flavor Flavor) string {
	colParts := make([]string, len(fk.ColumnNames))
	for n, colName := range fk.ColumnNames {
		colParts[n] = EscapeIdentifier(colName)
	}
	childCols := strings.Join(colParts, ", ")

	referencedTable := EscapeIdentifier(fk.ReferencedTableName)
	if fk.ReferencedSchemaName != "" {
		referencedTable = fmt.Sprintf("%s.%s", EscapeIdentifier(fk.ReferencedSchemaName), referencedTable)
	}

	for n, col := range fk.ReferencedColumnNames {
		colParts[n] = EscapeIdentifier(col)
	}
	parentCols := strings.Join(colParts, ", ")

	// MySQL 8 omits NO ACTION clauses, but includes RESTRICT clauses. In all other
	// flavors the opposite is true. (Even though NO ACTION and RESTRICT are
	// completely equivalent...)
	var hiddenRule, deleteRule, updateRule string
	if flavor.MinMySQL(8) {
		hiddenRule = "NO ACTION"
	} else {
		hiddenRule = "RESTRICT"
	}
	if fk.DeleteRule != hiddenRule {
		deleteRule = fmt.Sprintf(" ON DELETE %s", fk.DeleteRule)
	}
	if fk.UpdateRule != hiddenRule {
		updateRule = fmt.Sprintf(" ON UPDATE %s", fk.UpdateRule)
	}

	return fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)%s%s", EscapeIdentifier(fk.Name), childCols, referencedTable, parentCols, deleteRule, updateRule)
}

// Equals returns true if two ForeignKeys are completely identical (even in
// terms of cosmetic differences), false otherwise.
func (fk *ForeignKey) Equals(other *ForeignKey) bool {
	if fk == nil || other == nil {
		return fk == other // only equal if BOTH are nil
	}
	return fk.Name == other.Name && fk.UpdateRule == other.UpdateRule && fk.DeleteRule == other.DeleteRule && fk.Equivalent(other)
}

// Equivalent returns true if two ForeignKeys are functionally equivalent,
// regardless of whether or not they have the same names.
func (fk *ForeignKey) Equivalent(other *ForeignKey) bool {
	if fk == nil || other == nil {
		return fk == other // only equivalent if BOTH are nil
	}

	if fk.ReferencedSchemaName != other.ReferencedSchemaName || fk.ReferencedTableName != other.ReferencedTableName {
		return false
	}
	if fk.normalizedUpdateRule() != other.normalizedUpdateRule() || fk.normalizedDeleteRule() != other.normalizedDeleteRule() {
		return false
	}
	if len(fk.ColumnNames) != len(other.ColumnNames) {
		return false
	}
	for n := range fk.ColumnNames {
		if fk.ColumnNames[n] != other.ColumnNames[n] || fk.ReferencedColumnNames[n] != other.ReferencedColumnNames[n] {
			return false
		}
	}
	return true
}

func (fk *ForeignKey) normalizedUpdateRule() string {
	// MySQL and MariaDB both treat RESTRICT, NO ACTION, and lack of a rule
	// equivalently in terms of functionality.
	if fk.UpdateRule == "RESTRICT" || fk.UpdateRule == "NO ACTION" {
		return ""
	}
	return fk.UpdateRule
}

func (fk *ForeignKey) normalizedDeleteRule() string {
	// MySQL and MariaDB both treat RESTRICT, NO ACTION, and lack of a rule
	// equivalently in terms of functionality.
	if fk.DeleteRule == "RESTRICT" || fk.DeleteRule == "NO ACTION" {
		return ""
	}
	return fk.DeleteRule
}
