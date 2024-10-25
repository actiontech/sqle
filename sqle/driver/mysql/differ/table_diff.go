package differ

import (
	"errors"
	"fmt"
	"strings"
)

// TableDiff represents a difference between two tables.
type TableDiff struct {
	Type         DiffType
	From         *Table
	To           *Table
	alterClauses []TableAlterClause
	supported    bool
}

// ObjectKey returns a value representing the type and name of the table being
// diff'ed. The name will be the From side table, unless the diffType is
// DiffTypeCreate, in which case the To side table name is used.
func (td *TableDiff) ObjectKey() ObjectKey {
	if td == nil {
		return ObjectKey{}
	}
	if td.Type == DiffTypeCreate {
		return td.To.ObjectKey()
	}
	return td.From.ObjectKey()
}

// DiffType returns the type of diff operation.
func (td *TableDiff) DiffType() DiffType {
	if td == nil {
		return DiffTypeNone
	}
	return td.Type
}

// NewCreateTable returns a *TableDiff representing a CREATE TABLE statement,
// i.e. a table that only exists in the "to" side schema in a diff.
func NewCreateTable(table *Table) *TableDiff {
	return &TableDiff{
		Type:      DiffTypeCreate,
		To:        table,
		supported: true,
	}
}

// NewAlterTable returns a *TableDiff representing an ALTER TABLE statement,
// i.e. a table that exists in the "from" and "to" side schemas but with one
// or more differences. If the supplied tables are identical, nil will be
// returned instead of a TableDiff.
func NewAlterTable(from, to *Table) *TableDiff {
	clauses, supported := from.Diff(to)
	if supported && len(clauses) == 0 {
		return nil
	}
	return &TableDiff{
		Type:         DiffTypeAlter,
		From:         from,
		To:           to,
		alterClauses: clauses,
		supported:    supported,
	}
}

// NewDropTable returns a *TableDiff representing a DROP TABLE statement,
// i.e. a table that only exists in the "from" side schema in a diff.
func NewDropTable(table *Table) *TableDiff {
	return &TableDiff{
		Type:      DiffTypeDrop,
		From:      table,
		supported: true,
	}
}

// PreDropAlters returns a slice of *TableDiff to run prior to dropping a
// table. For tables partitioned with RANGE or LIST partitioning, this returns
// ALTERs to drop all partitions but one. In all other cases, this returns nil.
func PreDropAlters(table *Table) []*TableDiff {
	if table.Partitioning == nil || table.Partitioning.SubMethod != "" {
		return nil
	}
	// Only RANGE, RANGE COLUMNS, LIST, LIST COLUMNS support ALTER TABLE...DROP
	// PARTITION clause
	if !strings.HasPrefix(table.Partitioning.Method, "RANGE") && !strings.HasPrefix(table.Partitioning.Method, "LIST") {
		return nil
	}

	fakeTo := &Table{}
	*fakeTo = *table
	fakeTo.Partitioning = nil
	var result []*TableDiff
	for _, p := range table.Partitioning.Partitions[0 : len(table.Partitioning.Partitions)-1] {
		clause := ModifyPartitions{
			Drop:         []*Partition{p},
			ForDropTable: true,
		}
		result = append(result, &TableDiff{
			Type:         DiffTypeAlter,
			From:         table,
			To:           fakeTo,
			alterClauses: []TableAlterClause{clause},
			supported:    true,
		})
	}
	return result
}

// SplitAddForeignKeys looks through a TableDiff's alterClauses and pulls out
// any AddForeignKey clauses into a separate TableDiff. The first returned
// TableDiff is guaranteed to contain no AddForeignKey clauses, and the second
// returned value is guaranteed to only consist of AddForeignKey clauses. If
// the receiver contained no AddForeignKey clauses, the first return value will
// be the receiver, and the second will be nil. If the receiver contained only
// AddForeignKey clauses, the first return value will be nil, and the second
// will be the receiver.
// This method is useful for several reasons: it is desirable to only add FKs
// after other alters have been made (since FKs rely on indexes on both sides);
// it is illegal to drop and re-add an FK with the same name in the same ALTER;
// some versions of MySQL recommend against dropping and adding FKs in the same
// ALTER even if they have different names.
func (td *TableDiff) SplitAddForeignKeys() (*TableDiff, *TableDiff) {
	if td.Type != DiffTypeAlter || !td.supported || len(td.alterClauses) == 0 {
		return td, nil
	}

	addFKClauses := make([]TableAlterClause, 0)
	otherClauses := make([]TableAlterClause, 0, len(td.alterClauses))
	for _, clause := range td.alterClauses {
		if _, ok := clause.(AddForeignKey); ok {
			addFKClauses = append(addFKClauses, clause)
		} else {
			otherClauses = append(otherClauses, clause)
		}
	}
	if len(addFKClauses) == 0 {
		return td, nil
	} else if len(otherClauses) == 0 {
		return nil, td
	}
	result1 := &TableDiff{
		Type:         DiffTypeAlter,
		From:         td.From,
		To:           td.To,
		alterClauses: otherClauses,
		supported:    true,
	}
	result2 := &TableDiff{
		Type:         DiffTypeAlter,
		From:         td.From,
		To:           td.To,
		alterClauses: addFKClauses,
		supported:    true,
	}
	return result1, result2
}

// SplitConflicts looks through a TableDiff's alterClauses and pulls out any
// clauses that need to be placed into a separate TableDiff in order to yield
// legal or error-free DDL, due to DDL edge-cases. This includes attempts to add
// multiple FULLTEXT indexes in a single ALTER, and attempts to rename an index
// while also changing its visibility/ignored status.
// This method returns a slice of TableDiffs. The first element will be
// equivalent to the receiver (td) with any conflicting clauses removed;
// subsequent slice elements, if any, will be separate TableDiffs each
// consisting of individual conflicting clauses.
// This method does not interact with AddForeignKey clauses; see dedicated
// method SplitAddForeignKeys for that logic.
func (td *TableDiff) SplitConflicts() (result []*TableDiff) {
	if td == nil {
		return nil
	} else if td.Type != DiffTypeAlter || !td.supported || len(td.alterClauses) == 0 {
		return []*TableDiff{td}
	}

	var seenAddFulltext bool
	keepClauses := make([]TableAlterClause, 0, len(td.alterClauses))
	separateClauses := make([]TableAlterClause, 0)
	for _, clause := range td.alterClauses {
		if addIndex, ok := clause.(AddIndex); ok && addIndex.Index.Type == "FULLTEXT" {
			if seenAddFulltext {
				separateClauses = append(separateClauses, clause)
				continue
			}
			seenAddFulltext = true
		} else if mi, ok := clause.(ModifyIndex); ok && mi.FromIndex.Equivalent(mi.ToIndex) && mi.FromIndex.Name != mi.ToIndex.Name && mi.FromIndex.Invisible != mi.ToIndex.Invisible {
			// Put an AlterIndex into separateClauses so that we run that clause in its
			// own separate ALTER TABLE, or skipped if StatementModifiers cause the
			// original ModifyIndex to be handled as a DROP/re-ADD.
			separateClauses = append(separateClauses, AlterIndex{
				Name:         mi.ToIndex.Name, // use the post-rename name, since the rename happens first!
				Invisible:    mi.ToIndex.Invisible,
				linkedRename: &mi,
			})
		}
		keepClauses = append(keepClauses, clause)
	}

	result = append(result, &TableDiff{
		Type:         DiffTypeAlter,
		From:         td.From,
		To:           td.To,
		alterClauses: keepClauses,
		supported:    true,
	})
	for n := range separateClauses {
		result = append(result, &TableDiff{
			Type:         DiffTypeAlter,
			From:         td.From,
			To:           td.To,
			alterClauses: []TableAlterClause{separateClauses[n]},
			supported:    true,
		})
	}
	return result
}

// Statement returns the full DDL statement corresponding to the TableDiff. A
// blank string may be returned if the mods indicate the statement should be
// skipped. If the mods indicate the statement should be disallowed, it will
// still be returned as-is, but the error will be non-nil. Be sure not to
// ignore the error value of this method.
func (td *TableDiff) Statement(mods StatementModifiers) (string, error) {
	if td == nil {
		return "", nil
	}

	var err error
	switch td.Type {
	case DiffTypeCreate:
		stmt := td.To.CreateStatement
		if td.To.Partitioning != nil && mods.Partitioning == PartitioningRemove {
			stmt = td.To.UnpartitionedCreateStatement(mods.Flavor)
		}
		if td.To.HasAutoIncrement() && (mods.NextAutoInc == NextAutoIncIgnore || mods.NextAutoInc == NextAutoIncIfAlready) {
			stmt, _ = ParseCreateAutoInc(stmt)
		}
		return stmt, nil
	case DiffTypeAlter:
		return td.alterStatement(mods)
	case DiffTypeDrop:
		stmt := td.From.DropStatement()
		if !mods.AllowUnsafe {
			err = &UnsafeDiffError{
				Reason: "Desired drop of table " + EscapeIdentifier(td.From.Name) + " would cause all of its data to be lost.",
			}
		}
		return stmt, err
	default: // DiffTypeRename not supported yet
		return "", fmt.Errorf("unsupported diff type %d", td.Type)
	}
}

// Clauses returns the body of the statement represented by the table diff.
// For DROP statements, this will be an empty string. For CREATE statements,
// it will be everything after "CREATE TABLE [name] ". For ALTER statements,
// it will be everything after "ALTER TABLE [name] ".
func (td *TableDiff) Clauses(mods StatementModifiers) (string, error) {
	stmt, err := td.Statement(mods)
	if stmt == "" {
		return stmt, err
	}
	switch td.Type {
	case DiffTypeCreate:
		prefix := fmt.Sprintf("CREATE TABLE %s ", EscapeIdentifier(td.To.Name))
		return strings.Replace(stmt, prefix, "", 1), err
	case DiffTypeAlter:
		prefix := fmt.Sprintf("%s ", td.From.AlterStatement())
		return strings.Replace(stmt, prefix, "", 1), err
	case DiffTypeDrop:
		return "", err
	default: // DiffTypeRename not supported yet
		return "", fmt.Errorf("unsupported diff type %d", td.Type)
	}
}

func (td *TableDiff) alterStatement(mods StatementModifiers) (string, error) {
	// Force StrictIndexOrder to be enabled for InnoDB tables that have no primary
	// key and at least one unique index with non-nullable columns
	if !mods.StrictIndexOrder && td.To.Engine == "InnoDB" && td.To.ClusteredIndexKey() != td.To.PrimaryKey {
		mods.StrictIndexOrder = true
	}

	clauseStrings := make([]string, 0, len(td.alterClauses))
	var unsafeReasons []string
	var partitionClauseString string
	var changingComment bool
	for _, clause := range td.alterClauses {
		if !mods.AllowUnsafe {
			if clause, ok := clause.(Unsafer); ok {
				if unsafe, reason := clause.Unsafe(mods); unsafe {
					unsafeReasons = append(unsafeReasons, reason)
				}
			}
		}
		if clauseString := clause.Clause(mods); clauseString != "" {
			switch clause.(type) {
			case PartitionBy, RemovePartitioning:
				// Adding or removing partitioning must occur at the end of the ALTER
				// TABLE, and oddly *without* a preceeding comma
				partitionClauseString = clauseString
				continue // do NOT append to clauseStrings
			case ModifyPartitions:
				// Other partitioning-related clauses cannot appear alongside any other
				// clauses, including ALGORITHM or LOCK clauses
				mods.LockClause = ""
				mods.AlgorithmClause = ""
			case ChangeComment:
				// Track this for LaxComments modifier
				changingComment = true
			}
			clauseStrings = append(clauseStrings, clauseString)
		}
	}

	// Determine any errors: unsafe, unsupported, or both.
	// The "both" situation happens when the table uses unsupported features but
	// we're still able to generate at least a partial diff, and that partial diff
	// is unsafe. In that case, the UnsupportedDiffError wraps the UnsafeDiffError
	// (instead of vice versa) for purposes of using the unsupported error message
	// as the primary error message.
	var err error
	if len(unsafeReasons) > 0 {
		err = &UnsafeDiffError{
			Reason: "Desired alteration for " + td.ObjectKey().String() + " is not safe: " + strings.Join(unsafeReasons, "; ") + ".",
		}
	}
	if !td.supported {
		if td.To.UnsupportedDDL {
			subjectAndVerb := `The desired state ("to" side of diff) contains `
			if td.From.UnsupportedDDL {
				subjectAndVerb = "Both sides of the diff contain "
			}
			err = fmt.Errorf("unsupported %s", subjectAndVerb)
		}
	}
	if len(clauseStrings) == 0 && partitionClauseString == "" {
		return "", err
	}

	// LaxComments means "only change the comment if some other non-comment thing
	// is also being changed"
	if mods.LaxComments && len(clauseStrings) == 1 && partitionClauseString == "" && changingComment {
		return "", err
	}

	if mods.LockClause != "" {
		lockClause := fmt.Sprintf("LOCK=%s", strings.ToUpper(mods.LockClause))
		clauseStrings = append([]string{lockClause}, clauseStrings...)
	}
	if mods.AlgorithmClause != "" {
		algorithmClause := fmt.Sprintf("ALGORITHM=%s", strings.ToUpper(mods.AlgorithmClause))
		clauseStrings = append([]string{algorithmClause}, clauseStrings...)
	}
	if mods.VirtualColValidation {
		var canValidate bool
		for _, clause := range td.alterClauses {
			switch clause := clause.(type) {
			case AddColumn:
				canValidate = canValidate || clause.Column.Virtual
			case ModifyColumn:
				canValidate = canValidate || clause.NewColumn.Virtual
			}
		}
		if canValidate {
			clauseStrings = append(clauseStrings, "WITH VALIDATION")
		}
	}

	var spacer string
	if len(clauseStrings) > 0 && partitionClauseString != "" {
		spacer = " "
	}
	return td.From.AlterStatement() + " " + strings.Join(clauseStrings, ", ") + spacer + partitionClauseString, err
}

// MarkSupported provides a mechanism for callers to vouch for the correctness
// of a TableDiff that was automatically marked as unsupported. This should only
// be used in cases where a table with UnsupportedDDL is being altered in a way
// which either doesn't interact with the unsupported features, or easily
// removes those features. It is the caller's responsibility to first verify
// that the TableDiff's Statement() returns accurate, non-empty SQL.
func (td *TableDiff) MarkSupported() error {
	if td == nil || len(td.alterClauses) == 0 {
		return errors.New("cannot mark TableDiff as supported: no alter clauses were generated")
	} else if td.supported {
		return errors.New("cannot mark TableDiff as supported: supported is already true")
	}
	td.supported = true
	return nil
}

func diffTables(from, to *Table) (clauses []TableAlterClause, supported bool) {
	if from.Name != to.Name {
		// Table renaming not yet supported
		return []TableAlterClause{}, false
	}

	// If both tables have same output for SHOW CREATE TABLE, we know they're the same.
	// We do this check prior to the UnsupportedDDL check so that we only emit the
	// warning if the tables actually changed.
	if from.CreateStatement != "" && from.CreateStatement == to.CreateStatement {
		return []TableAlterClause{}, true
	}

	// If we're attempting to alter a supported table into an unsupported table,
	// don't even bother attempting to generate clauses; we know with 100%
	// certainty that the emitted DDL will be incomplete or incorrect. (In other
	// cases, we still attempt to generate DDL, since the alter MAY just consist
	// of fully-supported alterations to otherwise-unsupported tables. For example:
	// a table is unsupported due to having a spatial index, but the alter is just
	// adding some unrelated column.)
	supported = !from.UnsupportedDDL && !to.UnsupportedDDL
	if !from.UnsupportedDDL && to.UnsupportedDDL {
		return nil, false
	}

	clauses = make([]TableAlterClause, 0)

	// Check for default charset or collation changes first, prior to looking at
	// column adds, to ensure the default change affects any new columns that don't
	// explicitly override the table default
	if from.CharSet != to.CharSet || from.Collation != to.Collation {
		clauses = append(clauses, ChangeCharSet{
			FromCharSet:   from.CharSet,
			FromCollation: from.Collation,
			ToCharSet:     to.CharSet,
			ToCollation:   to.Collation,
		})
	}

	// Process column drops, modifications, adds. Must be done in this specific order
	// so that column reordering works properly.
	cc := compareColumnExistence(from, to)
	clauses = append(clauses, cc.columnDrops()...)
	clauses = append(clauses, cc.columnModifications()...)
	clauses = append(clauses, cc.columnAdds()...)

	// Compare PK
	if !from.PrimaryKey.Equals(to.PrimaryKey) {
		if from.PrimaryKey != nil {
			clauses = append(clauses, DropIndex{Index: from.PrimaryKey})
		}
		if to.PrimaryKey != nil {
			clauses = append(clauses, AddIndex{Index: to.PrimaryKey})
		}
	}

	// Compare secondary indexes
	clauses = append(clauses, compareSecondaryIndexes(from, to)...)

	// Compare foreign keys. If only the name of an FK changes, we consider this
	// difference to be cosmetic, and suppress it at clause generation time unless
	// requested. (This is important for pt-osc support, since it renames FKs due
	// to their namespace being schema-wide.)
	fromForeignKeys := from.foreignKeysByName()
	toForeignKeys := to.foreignKeysByName()
	fkChangeCosmeticOnly := func(fk *ForeignKey, others []*ForeignKey) bool {
		for _, other := range others {
			if fk.Equivalent(other) {
				return true
			}
		}
		return false
	}
	for _, toFk := range toForeignKeys {
		if _, existedBefore := fromForeignKeys[toFk.Name]; !existedBefore {
			clauses = append(clauses, AddForeignKey{
				ForeignKey:   toFk,
				cosmeticOnly: fkChangeCosmeticOnly(toFk, from.ForeignKeys),
			})
		}
	}
	for _, fromFk := range fromForeignKeys {
		toFk, stillExists := toForeignKeys[fromFk.Name]
		if !stillExists {
			clauses = append(clauses, DropForeignKey{
				ForeignKey:   fromFk,
				cosmeticOnly: fkChangeCosmeticOnly(fromFk, to.ForeignKeys),
			})
		} else if !fromFk.Equals(toFk) {
			cosmeticOnly := fromFk.Equivalent(toFk) // e.g. just changes between RESTRICT and NO ACTION
			drop := DropForeignKey{
				ForeignKey:   fromFk,
				cosmeticOnly: cosmeticOnly,
			}
			add := AddForeignKey{
				ForeignKey:   toFk,
				cosmeticOnly: cosmeticOnly,
			}
			clauses = append(clauses, drop, add)
		}
	}

	// Compare check constraints. Although the order of check constraints has no
	// functional impact, ordering changes must nonetheless must be detected, as
	// MariaDB lists checks in creation order for I_S and SHOW CREATE. And similar
	// to FKs, we must detect naming-only changes for OSC tool compatibility.
	fromChecks := from.checksByName()
	toChecks := to.checksByName()
	checkChangeNameOnly := func(cc *Check, others []*Check) bool {
		for _, other := range others {
			if cc.Clause == other.Clause && cc.Enforced == other.Enforced {
				return true
			}
		}
		return false
	}
	var fromCheckStillExist []*Check // ordered list of checks from "from" that still exist in "to"
	for _, fromCheck := range from.Checks {
		if _, stillExists := toChecks[fromCheck.Name]; stillExists {
			fromCheckStillExist = append(fromCheckStillExist, fromCheck)
		} else {
			clauses = append(clauses, DropCheck{
				Check:      fromCheck,
				renameOnly: checkChangeNameOnly(fromCheck, to.Checks),
			})
		}
	}
	var reorderChecks bool
	for n, toCheck := range to.Checks {
		if fromCheck, existedBefore := fromChecks[toCheck.Name]; !existedBefore {
			clauses = append(clauses, AddCheck{
				Check:      toCheck,
				renameOnly: checkChangeNameOnly(toCheck, from.Checks),
			})
			reorderChecks = true
		} else if fromCheck.Clause != toCheck.Clause {
			clauses = append(clauses, DropCheck{Check: fromCheck}, AddCheck{Check: toCheck})
			reorderChecks = true
		} else if fromCheck.Enforced != toCheck.Enforced {
			// Note: if MariaDB ever supports NOT ENFORCED, this will need extra logic
			// similar to how AlterIndex.alsoReordering works!
			clauses = append(clauses, AlterCheck{Check: fromCheck, NewEnforcement: toCheck.Enforced})
		} else if reorderChecks {
			clauses = append(clauses,
				DropCheck{Check: fromCheck, reorderOnly: true},
				AddCheck{Check: toCheck, reorderOnly: true})
		} else if fromCheckStillExist[n].Name != toCheck.Name {
			// If we get here, reorderChecks was previously false, meaning anything
			// *before* this position was identical on both sides. We can therefore leave
			// *this* check alone and just reorder anything that now comes *after* it.
			reorderChecks = true
		}
	}

	// Compare storage engine
	if from.Engine != to.Engine {
		clauses = append(clauses, ChangeStorageEngine{NewStorageEngine: to.Engine})
	}

	// Compare next auto-inc value
	if from.NextAutoIncrement != to.NextAutoIncrement && to.HasAutoIncrement() {
		cai := ChangeAutoIncrement{
			NewNextAutoIncrement: to.NextAutoIncrement,
			OldNextAutoIncrement: from.NextAutoIncrement,
		}
		clauses = append(clauses, cai)
	}

	// Compare create options
	if from.CreateOptions != to.CreateOptions {
		cco := ChangeCreateOptions{
			OldCreateOptions: from.CreateOptions,
			NewCreateOptions: to.CreateOptions,
		}
		clauses = append(clauses, cco)
	}

	// Compare comment
	if from.Comment != to.Comment {
		clauses = append(clauses, ChangeComment{NewComment: to.Comment})
	}

	// Compare tablespace
	if from.Tablespace != to.Tablespace {
		clauses = append(clauses, ChangeTablespace{NewTablespace: to.Tablespace})
	}

	// Compare partitioning. This must be performed last due to a MySQL requirement
	// of PARTITION BY / REMOVE PARTITIONING occurring last in a multi-clause ALTER
	// TABLE.
	// Note that some partitioning differences aren't supported yet, and others are
	// intentionally ignored.
	partClauses, partSupported := from.Partitioning.Diff(to.Partitioning)
	clauses = append(clauses, partClauses...)
	if !partSupported {
		supported = false
	}

	// If the SHOW CREATE TABLE output differed between the two tables, but we
	// did not generate any clauses, this indicates some aspect of the change is
	// unsupported (even though the two tables are individually supported). This
	// normally shouldn't happen, but could be possible given differences between
	// MySQL versions, vendors, storage engines, etc.
	if len(clauses) == 0 && from.CreateStatement != "" && to.CreateStatement != "" {
		supported = false
	}

	return
}

// Secondary indexes may be added, dropped, renamed, or have visibility changes.
// Although relative order of indexes is usually irrelevant, we still support
// dropping/re-adding indexes to result in a desired ordering, requiring extra
// bookkeeping.
// This code is relatively complex because some old flavors don't support
// renaming, in which case we must add/drop... which then has further
// implications if strict relative ordering is also requested.
func compareSecondaryIndexes(from, to *Table) (clauses []TableAlterClause) {
	fromIndexes := from.SecondaryIndexesByName()               // indexes in "from", keyed by name but later adjusted to use new name in case of rename
	toIndexes := to.SecondaryIndexesByName()                   // indexes in "to", keyed by name
	fromIndexStillExist := make([]*Index, 0, len(fromIndexes)) // ordered list of indexes from "from" that still exist in "to"

	// Determine which indexes have been fully dropped
	for _, fromIndex := range from.SecondaryIndexes {
		stillExists := (toIndexes[fromIndex.Name] != nil)

		// Determine if the index is "missing" on To side due to a rename: compare
		// all seemingly-new indexes to this one
		if !stillExists {
			for toName, toIndex := range toIndexes {
				// This comparison intentionally doesn't examine the Comment or Invisible
				// fields; logic elsewhere handles that appropriately
				if fromIndexes[toName] == nil && toIndex.Equivalent(fromIndex) {
					// Adjust the mapping of indexes in "from" to now be keyed by its new name.
					// This makes later lookup/comparison easier, and also "claims" the index
					// so that it won't erroneously be used as a candidate in multiple renames!
					delete(fromIndexes, fromIndex.Name)
					fromIndexes[toName] = fromIndex
					stillExists = true
					break
				}
			}
		}

		// Either some corresponding index still exists on "to" side, or the index has
		// been dropped entirely
		if stillExists {
			fromIndexStillExist = append(fromIndexStillExist, fromIndex)
		} else {
			clauses = append(clauses, DropIndex{Index: fromIndex})
		}
	}

	var reorderDueToClause []TableAlterClause // if non-empty and using StrictIndexOrder, *may* need to re-order, depending what SQL these clauses emit
	var reorderDueToMove bool                 // if true and using StrictIndexOrder, definitely must re-order
	for n, toIndex := range to.SecondaryIndexes {
		fromIndex, existedBefore := fromIndexes[toIndex.Name]

		// Entirely new index, not a modification to existing index. This also
		// means any pre-existing "To" side indexes after this must be re-ordered
		// if caller uses StatementModifiers.StrictIndexOrder.
		if !existedBefore {
			clause := AddIndex{Index: toIndex}
			clauses = append(clauses, clause)
			reorderDueToMove = true
			continue
		}

		// This index has changed, and/or a previous index change potentially
		// requires dropping/re-adding subsequent indexes to maintain the requested
		// order of secondary index definitions in the CREATE TABLE.
		if !fromIndex.Equals(toIndex) || reorderDueToMove || len(reorderDueToClause) > 0 {
			clause := ModifyIndex{
				FromIndex:          fromIndex,
				ToIndex:            toIndex,
				reorderDueToClause: reorderDueToClause,
				reorderDueToMove:   reorderDueToMove,
			}
			clauses = append(clauses, clause)
			reorderDueToClause = append(reorderDueToClause, clause)
		}

		// The relative order of pre-existing indexes has changed. With strict
		// ordering, all indexes *after* this one must move. This one can stay in
		// place; the other ones' moves will result in the desired order.
		if !reorderDueToMove && fromIndexStillExist[n].Name != fromIndex.Name {
			reorderDueToMove = true
		}
	}

	return clauses
}

func compareColumnExistence(self, other *Table) columnsComparison {
	cc := columnsComparison{
		fromTable:           self,
		toTable:             other,
		fromColumnsByName:   self.ColumnsByName(),
		fromStillPresent:    make([]bool, len(self.Columns)),
		toAlreadyExisted:    make([]bool, len(other.Columns)),
		fromOrderCommonCols: make([]*Column, 0, len(self.Columns)),
		toOrderCommonCols:   make([]*Column, 0, len(other.Columns)),
	}
	toColumnsByName := other.ColumnsByName()
	for n, col := range self.Columns {
		if _, existsInOther := toColumnsByName[col.Name]; existsInOther {
			cc.fromStillPresent[n] = true
			cc.fromOrderCommonCols = append(cc.fromOrderCommonCols, col)
		}
	}
	for n, col := range other.Columns {
		if _, existsInSelf := cc.fromColumnsByName[col.Name]; existsInSelf {
			cc.toAlreadyExisted[n] = true
			cc.toOrderCommonCols = append(cc.toOrderCommonCols, col)
			if !cc.commonColumnsMoved && col.Name != cc.fromOrderCommonCols[len(cc.toOrderCommonCols)-1].Name {
				cc.commonColumnsMoved = true
			}
		}
	}
	return cc
}

type columnsComparison struct {
	fromTable           *Table
	fromColumnsByName   map[string]*Column
	fromStillPresent    []bool
	fromOrderCommonCols []*Column
	toTable             *Table
	toAlreadyExisted    []bool
	toOrderCommonCols   []*Column
	commonColumnsMoved  bool
}

func (cc *columnsComparison) columnDrops() []TableAlterClause {
	clauses := make([]TableAlterClause, 0)

	// Loop through cols in "from" table, and process column drops
	for fromPos, stillPresent := range cc.fromStillPresent {
		if !stillPresent {
			clauses = append(clauses, DropColumn{
				Column: cc.fromTable.Columns[fromPos],
			})
		}
	}
	return clauses
}

func (cc *columnsComparison) columnAdds() []TableAlterClause {
	clauses := make([]TableAlterClause, 0)

	// Loop through cols in "to" table, and process column adds
	for toPos, alreadyExisted := range cc.toAlreadyExisted {
		if alreadyExisted {
			continue
		}
		add := AddColumn{
			Table:  cc.toTable,
			Column: cc.toTable.Columns[toPos],
		}

		// Determine if the new col was positioned in a specific place.
		// i.e. are there any pre-existing cols that come after it?
		var existingColsAfter bool
		for _, afterAlreadyExisted := range cc.toAlreadyExisted[toPos+1:] {
			if afterAlreadyExisted {
				existingColsAfter = true
				break
			}
		}
		if existingColsAfter {
			if toPos == 0 {
				add.PositionFirst = true
			} else {
				add.PositionAfter = cc.toTable.Columns[toPos-1]
			}
		}
		clauses = append(clauses, add)
	}
	return clauses
}

func (cc *columnsComparison) columnModifications() []TableAlterClause {
	clauses := make([]TableAlterClause, 0)
	commonCount := len(cc.fromOrderCommonCols)
	if commonCount == 0 {
		// no common cols = no possible MODIFY COLUMN clauses
		return clauses
	} else if !cc.commonColumnsMoved {
		// If all common cols are at same position, efficient comparison is simpler
		for toPos, toCol := range cc.toOrderCommonCols {
			if fromCol := cc.fromOrderCommonCols[toPos]; !fromCol.Equals(toCol) {
				clauses = append(clauses, ModifyColumn{
					Table:              cc.toTable,
					OldColumn:          fromCol,
					NewColumn:          toCol,
					InUniqueConstraint: cc.colInUniqueConstraint(fromCol, toCol),
				})
			}
		}
		return clauses
	}

	// If one or more common columns were re-positioned, identify the longest
	// increasing subsequence in the "from" side, to determine which columns can
	// stay put vs which ones need to be repositioned.
	toColPos := make(map[string]int, commonCount)
	for toPos, col := range cc.toOrderCommonCols {
		toColPos[col.Name] = toPos
	}
	fromIndexToPos := make([]int, commonCount)
	for fromPos, fromCol := range cc.fromOrderCommonCols {
		fromIndexToPos[fromPos] = toColPos[fromCol.Name]
	}
	stayPut := make([]bool, commonCount)
	for _, toPos := range longestIncreasingSubsequence(fromIndexToPos) {
		stayPut[toPos] = true
	}

	// For each common column (relative to the "to" order), emit a MODIFY COLUMN
	// clause if the col was reordered or modified.
	for toPos, toCol := range cc.toOrderCommonCols {
		fromCol := cc.fromColumnsByName[toCol.Name]
		if moved := !stayPut[toPos]; moved || !fromCol.Equals(toCol) {
			modify := ModifyColumn{
				Table:              cc.toTable,
				OldColumn:          fromCol,
				NewColumn:          toCol,
				PositionFirst:      moved && toPos == 0,
				InUniqueConstraint: cc.colInUniqueConstraint(fromCol, toCol),
			}
			if moved && toPos > 0 {
				modify.PositionAfter = cc.toOrderCommonCols[toPos-1]
			}
			clauses = append(clauses, modify)
		}
	}
	return clauses
}

// colInUniqueConstraint returns true if the old and new versions of the column
// are in at least one unique constraint that existed in both old and new
// versions of the table. This information is useful for determining if a
// collation change is unsafe due to affecting string equality for one or more
// indexes.
func (cc *columnsComparison) colInUniqueConstraint(fromCol, toCol *Column) bool {
	fromUniques := cc.fromTable.UniqueConstraintsWithColumn(fromCol)
	if len(fromUniques) == 0 {
		return false
	}
	toUniques := cc.toTable.UniqueConstraintsWithColumn(toCol)
	if len(toUniques) == 0 {
		return false
	}
	for _, fromIdx := range fromUniques {
		for _, toIdx := range toUniques {
			if toIdx.Name == fromIdx.Name {
				return true
			}
		}
	}
	return false
}
