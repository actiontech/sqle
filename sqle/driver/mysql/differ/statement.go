package differ

import (
	"fmt"
	"strings"
)

// StatementType indicates the type of a SQL statement found in a SQLFile.
// Parsing of types is very rudimentary, which can be advantageous for linting
// purposes. Otherwise, SQL errors or typos would prevent type detection.
type StatementType int

// Constants enumerating different types of statements
const (
	StatementTypeUnknown StatementType = iota
	StatementTypeNoop                  // entirely whitespace and/or comments
	StatementTypeCommand               // currently just USE or DELIMITER
	StatementTypeCreate
	StatementTypeCreateUnsupported // edge cases like CREATE...SELECT
	StatementTypeAlter             // not actually ever parsed yet
	// Other types will be added once they are supported by the package
)

// Statement represents a logical instruction in a file, consisting of either
// an SQL statement, a command (e.g. "USE some_database"), or whitespace and/or
// comments between two separate statements or commands.
type Statement struct {
	File            string
	LineNo          int
	CharNo          int
	Text            string // includes trailing Delimiter and newline
	DefaultDatabase string // only populated if an explicit USE command was encountered
	Type            StatementType
	ObjectType      ObjectType
	ObjectName      string
	ObjectQualifier string
	Delimiter       string // delimiter in use at the time of statement; not necessarily present in Text though
	Compound        bool   // if true, this is a compound statement (stored program with a BEGIN block, requiring alternative delimiter)
	nameClause      string // raw version, potentially with schema name qualifier and/or surrounding backticks
}

// Location returns the file, line number, and character number where the
// statement was obtained from
func (stmt *Statement) Location() string {
	if stmt.File == "" && stmt.LineNo == 0 && stmt.CharNo == 0 {
		return ""
	}
	if stmt.File == "" {
		return fmt.Sprintf("unknown:%d:%d", stmt.LineNo, stmt.CharNo)
	}
	return fmt.Sprintf("%s:%d:%d", stmt.File, stmt.LineNo, stmt.CharNo)
}

// ObjectKey returns an ObjectKey for the object affected by this
// statement.
func (stmt *Statement) ObjectKey() ObjectKey {
	return ObjectKey{
		Type: stmt.ObjectType,
		Name: stmt.ObjectName,
	}
}

// Schema returns the schema name that this statement impacts.
func (stmt *Statement) Schema() string {
	if stmt.ObjectQualifier != "" {
		return stmt.ObjectQualifier
	}
	return stmt.DefaultDatabase
}

// Body returns the Statement's Text, without any trailing delimiter,
// whitespace, or qualified schema name.
func (stmt *Statement) Body() string {
	body, _ := stmt.SplitTextBody()
	if stmt.ObjectQualifier == "" || stmt.nameClause == "" {
		return body
	}
	return strings.Replace(body, stmt.nameClause, EscapeIdentifier(stmt.ObjectName), 1)
}

// SplitTextBody returns Text with its trailing delimiter and whitespace (if
// any) separated out into a separate string.
func (stmt *Statement) SplitTextBody() (body string, suffix string) {
	if stmt == nil {
		return "", ""
	}
	body = strings.TrimRight(stmt.Text, "\n\r\t ")
	if stmt.Delimiter != "" && stmt.Delimiter != "\000" {
		body = strings.TrimSuffix(body, stmt.Delimiter)
		body = strings.TrimRight(body, "\n\r\t ")
	}
	return body, stmt.Text[len(body):]
}

// NormalizeTrailer ensures the statement text ends in a delimiter (if required
// based on the statement type) and newline. This method modifies stmt in-place.
func (stmt *Statement) NormalizeTrailer() {
	body, trailer := stmt.SplitTextBody()

	// If delimiter isn't known, or this is a DELIMITER command line, or a noop/
	// comment line, or unknown statement: just ensure there's a trailing newline
	if stmt.Delimiter == "" || stmt.Delimiter == "\000" || stmt.Type == StatementTypeUnknown || stmt.Type == StatementTypeNoop {
		if !strings.Contains(trailer, "\n") {
			stmt.Text += "\n"
		}
		return
	}

	if !strings.Contains(trailer, stmt.Delimiter) {
		stmt.Text = body + stmt.Delimiter + trailer
	}
	if !strings.Contains(trailer, "\n") {
		stmt.Text += "\n"
	}
}

// Compounder is implemented by types that have the ability to represent
// compound statements, requiring special delimiter handling.
type Compounder interface {
	IsCompoundStatement() bool
}

// IsCompoundStatement returns true if stmt is a compound statement.
func (stmt *Statement) IsCompoundStatement() bool {
	return stmt != nil && stmt.Compound
}
