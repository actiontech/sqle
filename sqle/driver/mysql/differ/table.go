package differ

import (
	"fmt"
	"regexp"
	"strings"
)

// Table represents a single database table.
type Table struct {
	Name              string             `json:"name"`
	Engine            string             `json:"storageEngine"`
	CharSet           string             `json:"defaultCharSet"`
	Collation         string             `json:"defaultCollation"`
	ShowCollation     bool               `json:"showCollation,omitempty"` // Include default COLLATE in SHOW CREATE TABLE: logic differs by flavor
	CreateOptions     string             `json:"createOptions,omitempty"` // row_format, stats_persistent, stats_auto_recalc, etc
	Columns           []*Column          `json:"columns"`
	PrimaryKey        *Index             `json:"primaryKey,omitempty"`
	SecondaryIndexes  []*Index           `json:"secondaryIndexes,omitempty"`
	ForeignKeys       []*ForeignKey      `json:"foreignKeys,omitempty"`
	Checks            []*Check           `json:"checks,omitempty"`
	Comment           string             `json:"comment,omitempty"`
	Tablespace        string             `json:"tablespace,omitempty"`
	NextAutoIncrement uint64             `json:"nextAutoIncrement,omitempty"`
	Partitioning      *TablePartitioning `json:"partitioning,omitempty"`       // nil if table isn't partitioned
	UnsupportedDDL    bool               `json:"unsupportedForDiff,omitempty"` // If true, cannot diff this table or auto-generate its CREATE TABLE
	CreateStatement   string             `json:"showCreateTable"`              // complete SHOW CREATE TABLE obtained from an instance
}

// ObjectKey returns a value useful for uniquely refering to a Table within a
// single Schema, for example as a map key.
func (t *Table) ObjectKey() ObjectKey {
	if t == nil {
		return ObjectKey{}
	}
	return ObjectKey{
		Type: ObjectTypeTable,
		Name: t.Name,
	}
}

// Def returns the table's CREATE statement as a string.
func (t *Table) Def() string {
	return t.CreateStatement
}

// AlterStatement returns the prefix to a SQL "ALTER TABLE" statement.
func (t *Table) AlterStatement() string {
	return fmt.Sprintf("ALTER TABLE %s", EscapeIdentifier(t.Name))
}

// DropStatement returns a SQL statement that, if run, would drop this table.
func (t *Table) DropStatement() string {
	return fmt.Sprintf("DROP TABLE %s", EscapeIdentifier(t.Name))
}

// GeneratedCreateStatement generates a CREATE TABLE statement based on the
// Table's Go field values. If t.UnsupportedDDL is false, this will match
// the output of MySQL's SHOW CREATE TABLE statement. But if t.UnsupportedDDL
// is true, this means the table uses MySQL features that does not yet
// support, and so the output of this method will differ from MySQL.
func (t *Table) GeneratedCreateStatement(flavor Flavor) string {
	defs := make([]string, len(t.Columns), len(t.Columns)+len(t.SecondaryIndexes)+len(t.ForeignKeys)+len(t.Checks)+1)
	for n, c := range t.Columns {
		defs[n] = c.Definition(flavor)
	}
	if t.PrimaryKey != nil {
		defs = append(defs, t.PrimaryKey.Definition(flavor))
	}
	for _, idx := range t.SecondaryIndexes {
		defs = append(defs, idx.Definition(flavor))
	}
	for _, fk := range t.ForeignKeys {
		defs = append(defs, fk.Definition(flavor))
	}
	for _, cc := range t.Checks {
		defs = append(defs, cc.Definition(flavor))
	}
	var tablespaceClause string
	if t.Tablespace != "" {
		tablespaceClause = fmt.Sprintf(" /*!50100 TABLESPACE %s */", EscapeIdentifier(t.Tablespace))
	}
	var autoIncClause string
	if t.NextAutoIncrement > 1 {
		autoIncClause = fmt.Sprintf(" AUTO_INCREMENT=%d", t.NextAutoIncrement)
	}
	charSet := t.CharSet
	// MySQL 8.0.24+ uses "utf8mb3" for table default charset in SHOW CREATE TABLE,
	// but still "utf8" for cols there, and "utf8" everywhere in I_S
	if charSet == "utf8" && flavor.MinMySQL(8, 0, 24) {
		charSet = "utf8mb3"
	}
	var collate string
	if t.ShowCollation {
		collate = " COLLATE=" + t.Collation
	}
	var createOptions string
	if t.CreateOptions != "" {
		createOptions = " " + t.CreateOptions
	}
	var comment string
	if t.Comment != "" {
		comment = fmt.Sprintf(" COMMENT='%s'", EscapeValueForCreateTable(t.Comment))
	}
	result := fmt.Sprintf("CREATE TABLE %s (\n  %s\n)%s ENGINE=%s%s DEFAULT CHARSET=%s%s%s%s%s",
		EscapeIdentifier(t.Name),
		strings.Join(defs, ",\n  "),
		tablespaceClause,
		t.Engine,
		autoIncClause,
		charSet,
		collate,
		createOptions,
		comment,
		t.Partitioning.Definition(flavor),
	)
	return result
}

// UnpartitionedCreateStatement returns the table's CREATE statement without
// its PARTITION BY clause. Supplying an accurate flavor improves performance,
// but is not required; FlavorUnknown still works correctly.
func (t *Table) UnpartitionedCreateStatement(flavor Flavor) string {
	if t.Partitioning == nil {
		return t.CreateStatement
	}
	if partClause := t.Partitioning.Definition(flavor); strings.HasSuffix(t.CreateStatement, partClause) {
		return t.CreateStatement[0 : len(t.CreateStatement)-len(partClause)]
	}
	base, _ := ParseCreatePartitioning(t.CreateStatement)
	return base
}

// ColumnsByName returns a mapping of column names to Column value pointers,
// for all columns in the table.
func (t *Table) ColumnsByName() map[string]*Column {
	result := make(map[string]*Column, len(t.Columns))
	for _, c := range t.Columns {
		result[c.Name] = c
	}
	return result
}

// SecondaryIndexesByName returns a mapping of index names to Index value
// pointers, for all secondary indexes in the table.
func (t *Table) SecondaryIndexesByName() map[string]*Index {
	result := make(map[string]*Index, len(t.SecondaryIndexes))
	for _, idx := range t.SecondaryIndexes {
		result[idx.Name] = idx
	}
	return result
}

// foreignKeysByName returns a mapping of foreign key names to ForeignKey value
// pointers, for all foreign keys in the table.
func (t *Table) foreignKeysByName() map[string]*ForeignKey {
	result := make(map[string]*ForeignKey, len(t.ForeignKeys))
	for _, fk := range t.ForeignKeys {
		result[fk.Name] = fk
	}
	return result
}

// checksByName returns a mapping of check constraint names to Check value
// pointers, for all check constraints in the table.
func (t *Table) checksByName() map[string]*Check {
	result := make(map[string]*Check, len(t.Checks))
	for _, cc := range t.Checks {
		result[cc.Name] = cc
	}
	return result
}

// HasAutoIncrement returns true if the table contains an auto-increment column,
// or false otherwise.
func (t *Table) HasAutoIncrement() bool {
	for _, c := range t.Columns {
		if c.AutoIncrement {
			return true
		}
	}
	return false
}

// ClusteredIndexKey returns which index is used for an InnoDB table's clustered
// index. This will be the primary key if one exists; otherwise, it will be the
// first unique key made of only non-nullable, non-expression columns. If there
// is no such key, or if the table's engine isn't InnoDB, this method returns
// nil.
func (t *Table) ClusteredIndexKey() *Index {
	if t.Engine != "InnoDB" {
		return nil
	}
	if t.PrimaryKey != nil {
		return t.PrimaryKey
	}
	cols := t.ColumnsByName()
	nullable := func(index *Index) bool {
		for _, part := range index.Parts {
			if col := cols[part.ColumnName]; col == nil || col.Nullable {
				return true
			}
		}
		return false
	}
	for _, index := range t.SecondaryIndexes {
		if index.Unique && !index.Functional() && !nullable(index) {
			return index
		}
	}
	return nil
}

var reTableRowFormatClause = regexp.MustCompile(`ROW_FORMAT=(\w+)`)

// RowFormatClause returns the table's ROW_FORMAT clause, if one was explicitly
// specified in the table's creation options. If no ROW_FORMAT clause was
// specified, but a KEY_BLOCK_SIZE is, "COMPRESSED" will be returned since MySQL
// applies this automatically. If no ROW_FORMAT or KEY_BLOCK_SIZE was specified,
// a blank string is returned.
// This method does not query an instance to determine if the table's actual
// ROW_FORMAT differs from what was requested in creation options; nor does it
// query the default row format if none was specified.
func (t *Table) RowFormatClause() string {
	matches := reTableRowFormatClause.FindStringSubmatch(t.CreateOptions)
	if matches != nil {
		return matches[1]
	}
	if strings.Contains(t.CreateOptions, "KEY_BLOCK_SIZE") {
		return "COMPRESSED"
	}
	return ""
}

// UniqueConstraintsWithColumn returns a slice of Indexes which have uniqueness
// constraints (primary key or unique secondary index) and include col as one
// of the index parts. If col is not part of any uniqueness constraints, a nil
// slice is returned.
func (t *Table) UniqueConstraintsWithColumn(col *Column) []*Index {
	var result []*Index
	if t.PrimaryKey != nil && indexHasColumn(t.PrimaryKey, col) {
		result = append(result, t.PrimaryKey)
	}
	for _, idx := range t.SecondaryIndexes {
		if idx.Unique && indexHasColumn(idx, col) {
			result = append(result, idx)
		}
	}
	return result
}

func indexHasColumn(idx *Index, col *Column) bool {
	for _, part := range idx.Parts {
		if part.ColumnName == col.Name {
			return true
		}
	}
	return false
}

// Diff returns a set of differences between this table and another table. Some
// edge cases are not supported, such as sub-partitioning, spatial indexes,
// MariaDB application time periods, or various non-InnoDB table features; in
// this case, supported will be false and clauses MAY OR MAY NOT be empty. Any
// returned clauses in that case must be carefully verified for correctness.
func (t *Table) Diff(to *Table) (clauses []TableAlterClause, supported bool) {
	return diffTables(t, to)
}
