//go:build enterprise
// +build enterprise

package differ

import (
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/utils"
)

// DiffType enumerates possible ways that two objects differ
type DiffType int

// Constants representing the types of diff operations.
const (
	DiffTypeNone DiffType = iota
	DiffTypeCreate
	DiffTypeDrop
	DiffTypeAlter
	DiffTypeRename
)

// ObjectDiff is an interface allowing generic handling of differences between
// two objects.
type ObjectDiff interface {
	ObjectKeyer
	DiffType() DiffType
	Statement(StatementModifiers) (string, error)
}

// NextAutoIncMode enumerates various ways of handling AUTO_INCREMENT
// discrepancies between two tables.
type NextAutoIncMode uint8

// Constants for how to handle next-auto-inc values in table diffs. Usually
// these are ignored in diffs entirely, but in some cases they are included.
const (
	NextAutoIncIgnore      NextAutoIncMode = iota // omit auto-inc value changes in diff
	NextAutoIncIfIncreased                        // only include auto-inc value if the "from" side is less than the "to" side
	NextAutoIncIfAlready                          // only include auto-inc value if the "from" side is already greater than 1
	NextAutoIncAlways                             // always include auto-inc value in diff
)

// PartitioningMode enumerates ways of handling partitioning status -- that is,
// presence or lack of a PARTITION BY clause.
type PartitioningMode uint8

// Constants for how to handle partitioning status differences.
const (
	PartitioningPermissive PartitioningMode = iota // don't negate any partitioning-related clauses
	PartitioningRemove                             // negate PARTITION BY clauses from DDL
	PartitioningKeep                               // negate REMOVE PARTITIONING clauses from ALTERs
)

// StatementModifiers are options that may be applied to adjust the DDL emitted
// for a particular table, and/or generate errors if certain clauses are
// present.
type StatementModifiers struct {
	NextAutoInc            NextAutoIncMode  // How to handle differences in next-auto-inc values
	Partitioning           PartitioningMode // How to handle differences in partitioning status
	AllowUnsafe            bool             // Whether to allow potentially-destructive DDL (drop table, drop column, modify col type, etc)
	LockClause             string           // Include a LOCK=[value] clause in generated ALTER TABLE
	AlgorithmClause        string           // Include an ALGORITHM=[value] clause in generated ALTER TABLE
	StrictIndexOrder       bool             // If true, maintain index order even in cases where there is no functional difference
	StrictCheckConstraints bool             // If true, maintain check constraint definition even if differences are cosmetic (name change; relative order of check definitions in MariaDB)
	StrictForeignKeyNaming bool             // If true, maintain foreign key definition even if differences are cosmetic (name change, RESTRICT vs NO ACTION, etc)
	StrictColumnDefinition bool             // If true, maintain column properties that are purely cosmetic (only affects MySQL 8)
	LaxColumnOrder         bool             // If true, don't modify columns if they only differ by position
	LaxComments            bool             // If true, don't modify tables/columns/indexes/routines if they only differ by comment clauses
	CompareMetadata        bool             // If true, compare creation-time sql_mode and db collation for stored programs
	VirtualColValidation   bool             // If true, add WITH VALIDATION clause for ALTER TABLE affecting virtual columns
	SkipPreDropAlters      bool             // If true, skip ALTERs that were only generated to make DROP TABLE faster
	Flavor                 Flavor           // Adjust generated DDL to match vendor/version. Zero value is FlavorUnknown which makes no adjustments.
}

///// SchemaDiff ///////////////////////////////////////////////////////////////

// SchemaDiff represents a set of differences between two database schemas,
// encapsulating diffs of various different object types.
type SchemaDiff struct {
	FromSchema   *Schema
	ToSchema     *Schema
	TableDiffs   []*TableDiff   // a set of statements that, if run, would turn tables in FromSchema into ToSchema
	RoutineDiffs []*RoutineDiff // " but for funcs and procs
}

// NewSchemaDiff computes the set of differences between two database schemas.
func NewSchemaDiff(from, to *Schema) *SchemaDiff {
	result := &SchemaDiff{
		FromSchema: from,
		ToSchema:   to,
	}

	if from == nil && to == nil {
		return result
	}

	result.TableDiffs = compareTables(from, to)
	result.RoutineDiffs = compareRoutines(from, to)
	return result
}

func compareTables(from, to *Schema) []*TableDiff {
	var tableDiffs, addFKAlters []*TableDiff
	fromByName := from.TablesByName()
	toByName := to.TablesByName()

	for name, fromTable := range fromByName {
		toTable, stillExists := toByName[name]
		if !stillExists {
			tableDiffs = append(tableDiffs, PreDropAlters(fromTable)...)
			tableDiffs = append(tableDiffs, NewDropTable(fromTable))
			continue
		}
		td := NewAlterTable(fromTable, toTable)
		if td != nil {
			otherAlter, addFKAlter := td.SplitAddForeignKeys()
			alters := otherAlter.SplitConflicts()
			tableDiffs = append(tableDiffs, alters...)
			if addFKAlter != nil {
				addFKAlters = append(addFKAlters, addFKAlter)
			}
		}
	}
	for name, toTable := range toByName {
		if _, alreadyExists := fromByName[name]; !alreadyExists {
			tableDiffs = append(tableDiffs, NewCreateTable(toTable))
		}
	}

	// We put ALTER TABLEs containing ADD FOREIGN KEY last, since the FKs may rely
	// on tables, columns, or indexes that are being newly created earlier in the
	// diff. (This is not a comprehensive solution yet though, since FKs can refer
	// to other schemas, and NewSchemaDiff only operates within one schema.)
	tableDiffs = append(tableDiffs, addFKAlters...)
	return tableDiffs
}

// DatabaseDiff returns an object representing database-level DDL (CREATE
// DATABASE, ALTER DATABASE, DROP DATABASE), or nil if no database-level DDL
// is necessary.
func (sd *SchemaDiff) DatabaseDiff() *DatabaseDiff {
	dd := &DatabaseDiff{From: sd.FromSchema, To: sd.ToSchema}
	if dd.DiffType() == DiffTypeNone {
		return nil
	}
	return dd
}

// ObjectDiffs returns a slice of all ObjectDiffs in the SchemaDiff. The results
// are returned in a sorted order, such that the diffs' Statements are legal.
// For example, if a CREATE DATABASE is present, it will occur in the slice
// prior to any table-level DDL in that schema.
func (sd *SchemaDiff) ObjectDiffs() []ObjectDiff {
	result := make([]ObjectDiff, 0)
	dd := sd.DatabaseDiff()
	if dd != nil {
		result = append(result, dd)
	}
	for _, td := range sd.TableDiffs {
		result = append(result, td)
	}
	for _, rd := range sd.RoutineDiffs {
		result = append(result, rd)
	}
	return result
}

// String returns the set of differences between two schemas as a single string.
// In building this string representation, note that no statement modifiers are
// applied, and any errors from Statement() are ignored. This means the returned
// string may contain destructive statements, and should only be used for
// display purposes, not for DDL execution.
func (sd *SchemaDiff) String() string {
	allDiffs := sd.ObjectDiffs()
	diffStatements := make([]string, len(allDiffs))
	for n, diff := range allDiffs {
		stmt, _ := diff.Statement(StatementModifiers{})
		diffStatements[n] = fmt.Sprintf("%s;\n", stmt)
	}
	return strings.Join(diffStatements, "")
}

///// DatabaseDiff /////////////////////////////////////////////////////////////

// DatabaseDiff represents differences of schema characteristics (default
// character set or default collation), or a difference in the existence of the
// the schema.
type DatabaseDiff struct {
	From *Schema
	To   *Schema
}

// ObjectKey returns a value representing the type and name of the schema being
// diff'ed. The name will be the From side schema, unless it is nil (CREATE
// DATABASE), in which case the To side schema name is returned.
func (dd *DatabaseDiff) ObjectKey() ObjectKey {
	if dd == nil || (dd.From == nil && dd.To == nil) {
		return ObjectKey{}
	}
	if dd.From == nil {
		return dd.To.ObjectKey()
	}
	return dd.From.ObjectKey()
}

// DiffType returns the type of diff operation.
func (dd *DatabaseDiff) DiffType() DiffType {
	if dd == nil || (dd.From == nil && dd.To == nil) {
		return DiffTypeNone
	} else if dd.From == nil && dd.To != nil {
		return DiffTypeCreate
	} else if dd.From != nil && dd.To == nil {
		// 忽略drop database类型差异：https://github.com/actiontech/sqle-ee/issues/1949#issuecomment-2453922579
		// return DiffTypeDrop
		return DiffTypeNone
	}

	if dd.From.CharSet != dd.To.CharSet || dd.From.Collation != dd.To.Collation {
		return DiffTypeAlter
	}
	return DiffTypeNone
}

// Statement returns a DDL statement corresponding to the DatabaseDiff. A blank
// string may be returned if there is no statement to execute.
func (dd *DatabaseDiff) Statement(mods StatementModifiers) (string, error) {
	if dd == nil {
		return "", nil
	}
	switch dd.DiffType() {
	case DiffTypeCreate:
		// fmt.Sprintf("%s;\n use %s\n", dd.To.Name, dd.To.CreateStatement())
		return fmt.Sprintf("%s;\nUSE %s", dd.To.CreateStatement(), utils.SupplementalQuotationMarks(dd.To.Name)), nil
	case DiffTypeDrop:
		stmt := dd.From.DropStatement()
		var err error
		if !mods.AllowUnsafe {
			if len(dd.From.Objects()) > 0 {
				err = &UnsafeDiffError{
					Reason: "Desired drop of " + dd.ObjectKey().String() + " would cause data loss.",
				}
			}
		}
		return stmt, err
	case DiffTypeAlter:
		return dd.From.AlterStatement(dd.To.CharSet, dd.To.Collation), nil
	}
	return "", nil
}

///// Errors ///////////////////////////////////////////////////////////////////

// UnsafeDiffError can be returned by ObjectDiff.Statement when the supplied
// statement modifiers do not permit the generated ObjectDiff to be used in this
// situation.
type UnsafeDiffError struct {
	Reason string
}

// Error satisfies the builtin error interface.
func (e *UnsafeDiffError) Error() string {
	return e.Reason
}
