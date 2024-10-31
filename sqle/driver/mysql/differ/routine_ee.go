//go:build enterprise
// +build enterprise

package differ

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/actiontech/sqle/sqle/driver/mysql/executor"
	"github.com/jmoiron/sqlx"
)

// Routine represents a stored procedure or function.
type Routine struct {
	Name              string     `json:"name"`
	Type              ObjectType `json:"type"`                     // Will be ObjectTypeProcedure or ObjectTypeFunction
	Body              string     `json:"body"`                     // Has correct escaping despite I_S mutilating it
	ParamString       string     `json:"paramString"`              // Formatted as per original CREATE
	ReturnDataType    string     `json:"returnDataType,omitempty"` // Includes charset/collation when relevant
	Definer           string     `json:"definer"`
	DatabaseCollation string     `json:"dbCollation"` // from creation time
	Comment           string     `json:"comment,omitempty"`
	Deterministic     bool       `json:"deterministic,omitempty"`
	SQLDataAccess     string     `json:"sqlDataAccess,omitempty"`
	SecurityType      string     `json:"securityType"`
	SQLMode           string     `json:"sqlMode"`    // sql_mode in effect at creation time
	CreateStatement   string     `json:"showCreate"` // complete SHOW CREATE obtained from an instance
}

// ObjectKey returns a value useful for uniquely refering to a Routine within a
// single Schema, for example as a map key.
func (r *Routine) ObjectKey() ObjectKey {
	if r == nil {
		return ObjectKey{}
	}
	return ObjectKey{
		Type: r.Type,
		Name: r.Name,
	}
}

// Def returns the routine's CREATE statement as a string.
func (r *Routine) Def() string {
	return r.CreateStatement
}

// Definition generates and returns a canonical CREATE PROCEDURE or CREATE
// FUNCTION statement based on the Routine's Go field values.
func (r *Routine) Definition(flavor Flavor) string {
	return fmt.Sprintf("%s%s", r.head(flavor), r.Body)
}

// DefinerClause returns the routine's DEFINER, quoted/escaped in a way
// consistent with SHOW CREATE.
func (r *Routine) DefinerClause() string {
	if atPos := strings.LastIndex(r.Definer, "@"); atPos >= 0 {
		return fmt.Sprintf("DEFINER=%s@%s", EscapeIdentifier(r.Definer[0:atPos]), EscapeIdentifier(r.Definer[atPos+1:]))
	}
	return fmt.Sprintf("DEFINER=%s", r.Definer)
}

// head returns the portion of a CREATE statement prior to the body.
func (r *Routine) head(_ Flavor) string {
	var definer, returnClause, characteristics string

	if r.Definer != "" {
		definer = r.DefinerClause() + " "
	}
	if r.Type == ObjectTypeFunc {
		returnClause = fmt.Sprintf(" RETURNS %s", r.ReturnDataType)
	}

	clauses := make([]string, 0)
	if r.SQLDataAccess != "CONTAINS SQL" {
		clauses = append(clauses, fmt.Sprintf("    %s\n", r.SQLDataAccess))
	}
	if r.Deterministic {
		clauses = append(clauses, "    DETERMINISTIC\n")
	}
	if r.SecurityType != "DEFINER" {
		clauses = append(clauses, fmt.Sprintf("    SQL SECURITY %s\n", r.SecurityType))
	}
	if r.Comment != "" {
		clauses = append(clauses, fmt.Sprintf("    COMMENT '%s'\n", EscapeValueForCreateTable(r.Comment)))
	}
	characteristics = strings.Join(clauses, "")

	return fmt.Sprintf("CREATE %s%s %s(%s)%s\n%s",
		definer,
		r.Type.Caps(),
		EscapeIdentifier(r.Name),
		r.ParamString,
		returnClause,
		characteristics)
}

// Equals returns true if two routines are identical, false otherwise.
func (r *Routine) Equals(other *Routine) bool {
	// shortcut if both nil pointers, or both pointing to same underlying struct
	if r == other {
		return true
	}
	// if one is nil, but the two pointers aren't equal, then one is non-nil
	if r == nil || other == nil {
		return false
	}

	// All fields are simple scalars, so we can just use equality check once we
	// know neither is nil
	return *r == *other
}

// equalsIgnoringCharacteristics returns true if two routines are identical,
// or only differ by characteristics which can be adjusted in-place using ALTER:
// SQLDataAccess, SecurityType, or Comment.
func (r *Routine) equalsIgnoringCharacteristics(other *Routine) bool {
	// shortcut if both nil pointers, or both pointing to same underlying struct
	if r == other {
		return true
	}
	// if one is nil, but the two pointers aren't equal, then one is non-nil
	if r == nil || other == nil {
		return false
	}

	if r.Name != other.Name || r.Type != other.Type || r.Body != other.Body || r.Definer != other.Definer {
		return false
	}
	if r.Deterministic != other.Deterministic {
		return false // arguably a characteristic, but nonetheless not supported for ALTER...
	}
	if r.ParamString != other.ParamString || r.ReturnDataType != other.ReturnDataType {
		return false
	}
	if r.DatabaseCollation != other.DatabaseCollation || r.SQLMode != other.SQLMode {
		return false
	}
	return true
}

// DropStatement returns a SQL statement that, if run, would drop this routine.
func (r *Routine) DropStatement() string {
	return fmt.Sprintf("DROP %s %s", r.Type.Caps(), EscapeIdentifier(r.Name))
}

// parseCreateStatement populates Body, ParamString, and ReturnDataType by
// parsing CreateStatement. It is used during introspection of routines in
// situations where the mysql.proc table is unavailable or does not exist.
func (r *Routine) parseCreateStatement(flavor Flavor, schema string) error {
	// Find matching parens around arg list
	argStart := strings.IndexRune(r.CreateStatement, '(')
	var argEnd int
	nestCount := 1
	for pos, r := range r.CreateStatement {
		if nestCount == 0 {
			argEnd = pos
			break
		} else if pos <= argStart {
			continue
		} else if r == '(' {
			nestCount++
		} else if r == ')' {
			nestCount--
		}
	}
	if argStart <= 0 || argEnd <= 0 {
		return fmt.Errorf("failed to parse show create %s %s.%s: %s", r.Type.Caps(), EscapeIdentifier(schema), EscapeIdentifier(r.Name), r.CreateStatement)
	}
	r.ParamString = r.CreateStatement[argStart+1 : argEnd-1]

	if r.Type == ObjectTypeFunc {
		retStart := argEnd + len(" RETURNS ")
		retEnd := retStart + strings.IndexRune(r.CreateStatement[retStart:], '\n')
		if retEnd <= 0 {
			return fmt.Errorf("failed to parse show create %s %s.%s: %s", r.Type.Caps(), EscapeIdentifier(schema), EscapeIdentifier(r.Name), r.CreateStatement)
		}
		r.ReturnDataType = r.CreateStatement[retStart:retEnd]
	}

	// Attempt to replace r.Body with one that doesn't have character conversion problems
	if header := r.head(flavor); strings.HasPrefix(r.CreateStatement, header) {
		r.Body = r.CreateStatement[len(header):]
	}
	return nil
}

///// Diff logic ///////////////////////////////////////////////////////////////

// RoutineDiff represents a difference between two routines. For diffs modifying
// an existing routine, if it is a characteristic-only change, this will be
// represented as a single RoutineDiff with DiffTypeAlter. Otherwise a
// modification including non-characteristic changes will be represented as
// two separate RoutineDiffs: one DiffTypeDrop and one DiffTypeCreate. This is
// needed to handle flavors which don't support CREATE OR REPLACE syntax.
// Flavors that *do* support CREATE OR REPLACE will simply blank-out the DROP
// portion of the pair.
type RoutineDiff struct {
	Type DiffType
	From *Routine
	To   *Routine
}

// ObjectKey returns a value representing the type and name of the routine being
// diff'ed. The type will be either ObjectTypeFunc or ObjectTypeProc. The name
// will be the From side routine, unless this is a Create, in which case the To
// side routine name is used.
func (rd *RoutineDiff) ObjectKey() ObjectKey {
	if rd != nil && rd.From != nil {
		return rd.From.ObjectKey()
	} else if rd != nil && rd.To != nil {
		return rd.To.ObjectKey()
	}
	return ObjectKey{}
}

// DiffType returns the type of diff operation.
func (rd *RoutineDiff) DiffType() DiffType {
	if rd == nil {
		return DiffTypeNone
	}
	return rd.Type
}

// Statement returns the full DDL statement corresponding to the RoutineDiff. A
// blank string may be returned if the mods indicate the statement should be
// skipped. If the mods indicate the statement should be disallowed, it will
// still be returned as-is, but the error will be non-nil. Be sure not to
// ignore the error value of this method.
func (rd *RoutineDiff) Statement(mods StatementModifiers) (stmt string, err error) {
	if rd == nil {
		return "", nil
	}

	// MySQL and MariaDB both support ALTER only for a limited set of changes.
	// Handle this first since it's the simplest case.
	if rd.Type == DiffTypeAlter {
		return rd.alterStatement(mods)
	}

	// It's not an ALTER, so it's either a DROP or CREATE. This may be a related
	// pair if it represents a non-characteristic modification to an existing
	// routine. Detect some special-case types of replacements.
	var metadataOnlyReplace, mariaReplace, clearCommentReplace bool
	if rd.From != nil && rd.To != nil { // related pair for a replacement
		if rd.From.CreateStatement == rd.To.CreateStatement {
			// If we're replacing a routine only because its creation-time sql_mode or
			// db collation has changed, only proceed if mods indicate we should. (This
			// type of replacement is effectively opt-in because it is counter-intuitive
			// and obscure.)
			if !mods.CompareMetadata {
				return "", nil
			}
			metadataOnlyReplace = true
		} else if rd.From.Comment != rd.To.Comment && rd.From.equalsIgnoringCharacteristics(rd.To) {
			// Setting a comment to a blank string requires a DROP/CREATE pair in MySQL
			// 8.0+ due to a server bug, so compareRoutines() always emits a DROP/CREATE
			// pair since the flavor is not known at that time. For non-MySQL8+ flavors,
			// we then convert this pair back into a single ALTER.
			clearCommentReplace = true

			// However, if *only* the comment has changed, suppress the diff entirely
			// if mods indicate not to generate comment-only changes
			if mods.LaxComments && rd.From.SQLDataAccess == rd.To.SQLDataAccess && rd.From.SecurityType == rd.To.SecurityType {
				return "", nil
			}
		}

		// MariaDB can use CREATE OR REPLACE to modify routines in a single statement
		mariaReplace = mods.Flavor.IsMariaDB()
	}

	if rd.Type == DiffTypeDrop {
		// Omit the DROP part of the pair entirely in cases where we're doing an atomic replacement or alter
		if mariaReplace || (clearCommentReplace && !mods.Flavor.MinMySQL(8)) {
			return "", nil
		}
		stmt = rd.From.DropStatement()
		if metadataOnlyReplace {
			stmt = "# Dropping and re-creating " + rd.ObjectKey().String() + " to update metadata\n" + stmt
		}
		if !mods.AllowUnsafe {
			if rd.To == nil { // pure DROP, always unsafe
				err = &UnsafeDiffError{
					Reason: "Desired drop of " + rd.ObjectKey().String() + " is risky, since you must first ensure that it is not used in any application queries, or referenced by other routines.",
				}
			} else { // DROP just ahead of re-CREATE to replace routine in MySQL
				err = &UnsafeDiffError{
					Reason: "Desired modification to " + rd.ObjectKey().String() + " requires dropping and re-creating it, and application queries may fail if they attempt to call the routine during the brief moment after the DROP but before the re-CREATE.",
				}
			}
		}
		return stmt, err

	} else if rd.Type == DiffTypeCreate {
		if clearCommentReplace && !mods.Flavor.MinMySQL(8) {
			return rd.alterStatement(mods)
		}
		stmt = rd.To.CreateStatement
		if mariaReplace {
			stmt = strings.Replace(stmt, "CREATE ", "CREATE OR REPLACE ", 1)
			if metadataOnlyReplace {
				stmt = "# Replacing " + rd.ObjectKey().String() + " to update metadata\n" + stmt
			}
		}

		// If modifying a routine to adjust the params or return, mark the CREATE as
		// unsafe, even in MariaDB. In MySQL, this intentionally overwrites the
		// general-purpose Reason set above.
		if rd.From != nil && !mods.AllowUnsafe && (rd.From.ParamString != rd.To.ParamString || rd.From.ReturnDataType != rd.To.ReturnDataType) {
			err = &UnsafeDiffError{
				Reason: "Desired modification to " + rd.ObjectKey().String() + " affects its parameters or return type, which may break call-sites in application queries, or in other routines. There is no way to simultaneously deploy application and routine changes in an atomic fashion.",
			}
		}
		return stmt, err
	}

	// DiffTypeRename not used, no equivalent syntax
	return "", fmt.Errorf("unsupported diff type %d", rd.DiffType())
}

func (rd *RoutineDiff) alterStatement(mods StatementModifiers) (stmt string, err error) {
	var clauses []string
	if rd.From.SQLDataAccess != rd.To.SQLDataAccess {
		clauses = append(clauses, rd.To.SQLDataAccess)
	}
	if rd.From.SecurityType != rd.To.SecurityType {
		clauses = append(clauses, "SQL SECURITY "+rd.To.SecurityType)
	}
	if rd.From.Comment != rd.To.Comment && (len(clauses) > 0 || !mods.LaxComments) {
		clauses = append(clauses, fmt.Sprintf("COMMENT '%s'", EscapeValueForCreateTable(rd.To.Comment)))
	}
	if len(clauses) > 0 {
		stmt = "ALTER " + rd.To.Type.Caps() + " " + EscapeIdentifier(rd.To.Name) + " " + strings.Join(clauses, " ")
	}
	return stmt, nil
}

// IsCompoundStatement returns true if the diff is a compound CREATE statement,
// requiring special delimiter handling.
func (rd *RoutineDiff) IsCompoundStatement() bool {
	return rd.Type == DiffTypeCreate && ParseStatementInString(rd.To.CreateStatement).Compound
}

func compareRoutines(from, to *Schema) []*RoutineDiff {
	routineDiffs := compareRoutinesByName(from.ProceduresByName(), to.ProceduresByName())
	routineDiffs = append(routineDiffs, compareRoutinesByName(from.FunctionsByName(), to.FunctionsByName())...)
	return routineDiffs
}

// compareRoutinesByName is a helper function for comparing maps of procs or
// funcs, keyed by name. Both maps should only contain the same type of routine.
// In other words, both fromByName and toByName should only contain procs, or
// both only contain funcs. No validation of this is performed here.
func compareRoutinesByName(fromByName map[string]*Routine, toByName map[string]*Routine) (routineDiffs []*RoutineDiff) {
	for name, from := range fromByName {
		to, stillExists := toByName[name]
		if !stillExists {
			routineDiffs = append(routineDiffs, &RoutineDiff{Type: DiffTypeDrop, From: from})
		} else if !from.Equals(to) {
			// Determine if the only difference is in characteristics which can be
			// adjusted in-place using an ALTER. One special-case is needed to work
			// around a MySQL 8.0+ bug, where ALTER cannot be used to remove a COMMENT
			// clause; this means we must *always* avoid ALTER in that situation because
			// the DB flavor is not known at this point in time
			if from.equalsIgnoringCharacteristics(to) && (from.Comment == "" || to.Comment != "") {
				routineDiffs = append(routineDiffs, &RoutineDiff{Type: DiffTypeAlter, From: from, To: to})
			} else {
				routineDiffs = append(routineDiffs,
					&RoutineDiff{Type: DiffTypeDrop, From: from, To: to},
					&RoutineDiff{Type: DiffTypeCreate, From: from, To: to},
				)
			}
		}
	}
	for name, to := range toByName {
		if _, alreadyExists := fromByName[name]; !alreadyExists {
			routineDiffs = append(routineDiffs, &RoutineDiff{Type: DiffTypeCreate, To: to})
		}
	}
	return
}

///// Introspection logic //////////////////////////////////////////////////////

func querySchemaRoutines(conn *executor.Executor, schema string, routineNames []string, routineType string, flavor Flavor) ([]*Routine, error) {
	// Obtain the routines in the schema
	// We completely exclude routines that the user can call, but not examine --
	// e.g. user has EXECUTE priv but missing other vital privs. In this case
	// routine_definition will be NULL.
	type rawRoutine struct {
		Name              string         `json:"routine_name"`
		Type              string         `json:"routine_type"`
		Body              sql.NullString `json:"routine_definition"`
		IsDeterministic   string         `json:"is_deterministic"`
		SQLDataAccess     string         `json:"sql_data_access"`
		SecurityType      string         `json:"security_type"`
		SQLMode           string         `json:"sql_mode"`
		Comment           string         `json:"routine_comment"`
		Definer           string         `json:"definer"`
		DatabaseCollation string         `json:"database_collation"`
	}
	// Note on this query: MySQL 8.0 changes information_schema column names to
	// come back from queries in all caps, so we need to explicitly use AS clauses
	// in order to get them back as lowercase and have sqlx Select() work
	query := `
		SELECT SQL_BUFFER_RESULT
		       r.routine_name AS routine_name, UPPER(r.routine_type) AS routine_type,
		       r.routine_definition AS routine_definition,
		       UPPER(r.is_deterministic) AS is_deterministic,
		       UPPER(r.sql_data_access) AS sql_data_access,
		       UPPER(r.security_type) AS security_type,
		       r.sql_mode AS sql_mode, r.routine_comment AS routine_comment,
		       r.definer AS definer, r.database_collation AS database_collation
		FROM   information_schema.routines r
		WHERE  r.routine_schema = ? AND routine_type = ? AND routine_name in(?) AND routine_definition IS NOT NULL`
	query, args, queryErr := sqlx.In(query, schema, routineType, routineNames)
	if queryErr != nil {
		return nil, queryErr
	}
	results, queryErr := conn.Db.Query(query, args...)
	if queryErr != nil {
		return nil, fmt.Errorf("Error querying information_schema.routines for schema %s: %s", schema, queryErr)
	}
	if len(results) == 0 {
		return []*Routine{}, nil
	}
	retRoutines := make([]*rawRoutine, len(results))
	for i, record := range results {
		retRoutines[i] = &rawRoutine{
			Name:              record["routine_name"].String,
			Type:              record["routine_type"].String,
			Body:              record["routine_definition"],
			IsDeterministic:   record["is_deterministic"].String,
			SQLDataAccess:     record["sql_data_access"].String,
			SecurityType:      record["security_type"].String,
			SQLMode:           record["sql_mode"].String,
			Comment:           record["routine_comment"].String,
			Definer:           record["definer"].String,
			DatabaseCollation: record["database_collation"].String,
		}
	}
	routines := make([]*Routine, len(retRoutines))
	dict := make(map[ObjectKey]*Routine, len(retRoutines))
	for n, rawRoutine := range retRoutines {
		routines[n] = &Routine{
			Name:              rawRoutine.Name,
			Type:              ObjectType(strings.ToLower(rawRoutine.Type)),
			Body:              rawRoutine.Body.String, // This contains incorrect formatting conversions; overwritten later
			Definer:           rawRoutine.Definer,
			DatabaseCollation: rawRoutine.DatabaseCollation,
			Comment:           rawRoutine.Comment,
			Deterministic:     rawRoutine.IsDeterministic == "YES",
			SQLDataAccess:     rawRoutine.SQLDataAccess,
			SecurityType:      rawRoutine.SecurityType,
			SQLMode:           rawRoutine.SQLMode,
		}
		if routines[n].Type != ObjectTypeProc && routines[n].Type != ObjectTypeFunc {
			return nil, fmt.Errorf("unsupported routine type %s found in %s.%s", rawRoutine.Type, schema, rawRoutine.Name)
		}
		key := ObjectKey{Type: routines[n].Type, Name: routines[n].Name}
		dict[key] = routines[n]
	}

	// Obtain param string, return type string, and full create statement:
	// We can't rely only on information_schema, since it doesn't have the param
	// string formatted in the same way as the original CREATE, nor does
	// routines.body handle strings/charsets correctly for re-runnable SQL.
	// In flavors without the new data dictionary, we first try querying mysql.proc
	// to bulk-fetch sufficient info to rebuild the CREATE without needing to run
	// a SHOW CREATE per routine.
	// If mysql.proc doesn't exist or that query fails, we then run a SHOW CREATE
	// per routine, using multiple goroutines for performance reasons.
	var alreadyObtained int
	if !flavor.MinMySQL(8) {
		type rawRoutineMeta struct {
			Name      string `json:"name"`
			Type      string `json:"type"`
			Body      string `json:"body"`
			ParamList string `json:"param_list"`
			Returns   string `json:"returns"`
		}
		query := `
			SELECT name, type, body, param_list, returns
			FROM   mysql.proc
			WHERE  db = ? AND type = ? AND name in(?)`
		query, args, queryErr := sqlx.In(query, schema, routineType, routineNames)
		if queryErr != nil {
			return nil, queryErr
		}
		// Errors here are non-fatal. No need to even check; slice will be empty which is fine
		metaResults, queryErr := conn.Db.Query(query, args...)
		if queryErr != nil {
			return nil, fmt.Errorf("Error querying mysql.proc for schema %s: %s", schema, queryErr)
		}
		retRoutineMeta := make([]*rawRoutineMeta, len(metaResults))
		for i, record := range metaResults {
			retRoutineMeta[i] = &rawRoutineMeta{
				Name:      record["name"].String,
				Type:      record["type"].String,
				Body:      record["body"].String,
				ParamList: record["param_list"].String,
				Returns:   record["returns"].String,
			}
		}

		for _, meta := range retRoutineMeta {
			key := ObjectKey{Type: ObjectType(strings.ToLower(meta.Type)), Name: meta.Name}
			if routine, ok := dict[key]; ok {
				routine.ParamString = strings.Replace(meta.ParamList, "\r\n", "\n", -1)
				routine.ReturnDataType = meta.Returns
				routine.Body = strings.Replace(meta.Body, "\r\n", "\n", -1)
				routine.CreateStatement = routine.Definition(flavor)
				alreadyObtained++
			}
		}
	}

	var err error
	if alreadyObtained < len(routines) {
		// g, subCtx := errgroup.WithContext(ctx)
		for n := range routines {
			r := routines[n] // avoid issues with goroutines and loop iterator values
			if r.CreateStatement == "" {
				// g.Go(func() (err error) {
				r.CreateStatement, err = showCreateRoutine(conn, r.Name, r.Type)
				if err == nil {
					r.CreateStatement = strings.Replace(r.CreateStatement, "\r\n", "\n", -1)
					err = r.parseCreateStatement(flavor, schema)
				} else {
					err = fmt.Errorf("error executing show create %s for %s.%s: %s", r.Type.Caps(), EscapeIdentifier(schema), EscapeIdentifier(r.Name), err)
				}
				// return err
				// })
			}
		}
		// err = g.Wait()
	}
	return routines, err
}

func showCreateRoutine(conn *executor.Executor, routine string, ot ObjectType) (create string, err error) {
	query := fmt.Sprintf("SHOW CREATE %s %s", ot.Caps(), EscapeIdentifier(routine))
	if ot == ObjectTypeProc {
		result, err := conn.Db.Query(query)
		if err == nil && len(result) != 1 {
			return "", fmt.Errorf("sql: no rows in result set")
		} else if err == nil {
			create = result[0]["Create Procedure"].String
		}
	} else if ot == ObjectTypeFunc {
		result, err := conn.Db.Query(query)
		if err == nil && len(result) != 1 {
			return "", fmt.Errorf("sql: no rows in result set")
		} else if err == nil {
			create = result[0]["Create Function"].String
		}
	} else {
		return "", fmt.Errorf("object type %s is not a routine", ot)
	}
	return
}
