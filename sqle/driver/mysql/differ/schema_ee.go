//go:build enterprise
// +build enterprise

package differ

// Schema represents a database schema.
type Schema struct {
	Name      string     `json:"databaseName"`
	CharSet   string     `json:"defaultCharSet"`
	Collation string     `json:"defaultCollation"`
	Tables    []*Table   `json:"tables,omitempty"`
	Routines  []*Routine `json:"routines,omitempty"`
}

// ObjectKey returns a value useful for uniquely refering to a Schema, for
// example as a map key.
func (s *Schema) ObjectKey() ObjectKey {
	if s == nil {
		return ObjectKey{}
	}
	return ObjectKey{
		Type: ObjectTypeDatabase,
		Name: s.Name,
	}
}

// Def returns the schema's CREATE statement as a string.
func (s *Schema) Def() string {
	return s.CreateStatement()
}

// TablesByName returns a mapping of table names to Table struct pointers, for
// all tables in the schema.
func (s *Schema) TablesByName() map[string]*Table {
	if s == nil {
		return map[string]*Table{}
	}
	result := make(map[string]*Table, len(s.Tables))
	for _, t := range s.Tables {
		result[t.Name] = t
	}
	return result
}

// HasTable returns true if a table with the given name exists in the schema.
// Callers should be careful to supply a name that takes into account the
// server's lower_case_table_names setting.
func (s *Schema) HasTable(name string) bool {
	return s != nil && s.Table(name) != nil
}

// Table returns a table by name.
// Callers should be careful to supply a name that takes into account the
// server's lower_case_table_names setting.
func (s *Schema) Table(name string) *Table {
	if s != nil {
		for _, t := range s.Tables {
			if t.Name == name {
				return t
			}
		}
	}
	return nil
}

// ProceduresByName returns a mapping of stored procedure names to Routine
// struct pointers, for all stored procedures in the schema.
func (s *Schema) ProceduresByName() map[string]*Routine {
	return s.routinesByNameAndType(ObjectTypeProc)
}

// FunctionsByName returns a mapping of function names to Routine struct
// pointers, for all functions in the schema.
func (s *Schema) FunctionsByName() map[string]*Routine {
	return s.routinesByNameAndType(ObjectTypeFunc)
}

func (s *Schema) routinesByNameAndType(ot ObjectType) map[string]*Routine {
	if s == nil {
		return map[string]*Routine{}
	}
	result := make(map[string]*Routine, len(s.Routines))
	for _, r := range s.Routines {
		if r.Type == ot {
			result[r.Name] = r
		}
	}
	return result
}

// Objects returns DefKeyers for all objects in the schema, excluding the schema
// itself. The result is a map, keyed by ObjectKey (type+name).
func (s *Schema) Objects() map[ObjectKey]DefKeyer {
	if s == nil {
		return nil
	}
	dict := make(map[ObjectKey]DefKeyer, len(s.Tables)+len(s.Routines))
	for _, table := range s.Tables {
		dict[table.ObjectKey()] = table
	}
	for _, routine := range s.Routines {
		dict[routine.ObjectKey()] = routine
	}
	return dict
}

// StripMatches removes objects from s if they match any supplied pattern. The
// in-memory representation of the schema is modified in-place. This does not
// affect any actual database instances.
func (s *Schema) StripMatches(removePatterns []ObjectPattern) {
	if s == nil {
		return
	}
	for _, pattern := range removePatterns {
		switch pattern.Type {
		case ObjectTypeTable:
			s.Tables = stripMatchingObjects(s.Tables, pattern)
		case ObjectTypeProc, ObjectTypeFunc:
			s.Routines = stripMatchingObjects(s.Routines, pattern)
		}
	}
}

func stripMatchingObjects[T ObjectKeyer](s []T, pattern ObjectPattern) (result []T) {
	for _, obj := range s {
		if !pattern.Match(obj) {
			result = append(result, obj)
		}
	}
	return
}

// Diff returns the set of differences between this schema and another schema.
func (s *Schema) Diff(other *Schema) *SchemaDiff {
	return NewSchemaDiff(s, other)
}

// DropStatement returns a SQL statement that, if run, would drop this schema.
func (s *Schema) DropStatement() string {
	return "DROP DATABASE " + EscapeIdentifier(s.Name)
}

// CreateStatement returns a SQL statement that, if run, would create this
// schema.
func (s *Schema) CreateStatement() string {
	var charSet, collate string
	if s.CharSet != "" {
		charSet = " CHARACTER SET " + s.CharSet
	}
	if s.Collation != "" {
		collate = " COLLATE " + s.Collation
	}
	return "CREATE DATABASE " + EscapeIdentifier(s.Name) + charSet + collate
}

// AlterStatement returns a SQL statement that, if run, would alter this
// schema's default charset and/or collation to the supplied values.
// If charSet is "" and collation isn't, only the collation will be changed.
// If collation is "" and charSet isn't, the default collation for charSet is
// used automatically.
// If both params are "", or if values equal to the schema's current charSet
// and collation are supplied, an empty string is returned.
func (s *Schema) AlterStatement(charSet, collation string) string {
	var charSetClause, collateClause string
	if s.CharSet != charSet && charSet != "" {
		charSetClause = " CHARACTER SET " + charSet
	}
	if s.Collation != collation && collation != "" {
		collateClause = " COLLATE " + collation
	}
	if charSetClause == "" && collateClause == "" {
		return ""
	}
	return "ALTER DATABASE " + EscapeIdentifier(s.Name) + charSetClause + collateClause
}
