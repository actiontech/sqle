// Package hive provides the built-in Hive driver for sqle-ee. This file
// implements diffTableDDL, the helper that compares two Hive CREATE TABLE
// statements and decides whether the difference can be applied with ALTER
// statements or must fall back to DROP + CREATE.
//
// The implementation follows docs/spec/design.md §3.4 (variant SQL generation
// matrix, D3 decision) and the Hive ALTER capability matrix in
// docs/spec/exploration.md §4. It addresses compat-RISK-6: when a difference
// cannot be safely expressed with ALTER (partition key changes, storage
// format changes, SerDe changes, EXTERNAL flag changes, incompatible column
// type changes, column deletion), the caller must emit DROP+CREATE with a
// data-loss WARNING.
//
// Per project convention (dev.md "非必要少用正则表达式"), this parser is
// section-based and operates on trimmed/uppercased keywords rather than
// regular expressions. It is robust to whitespace and case differences.
package hive

import (
	"fmt"
	"sort"
	"strings"
)

// runtimeTBLPropertyKeys lists TBLPROPERTIES keys that Hive maintains
// internally and that change frequently as data is written. These must be
// filtered out before comparing TBLPROPERTIES sets so that the diff matrix
// does not produce false-positive DROP+CREATE results (compat-RISK-6).
//
// Source: design.md §3.4 + exploration.md §4.
var runtimeTBLPropertyKeys = map[string]struct{}{
	"transient_lastDdlTime":  {},
	"numFiles":               {},
	"numRows":                {},
	"rawDataSize":            {},
	"totalSize":              {},
	"COLUMN_STATS_ACCURATE":  {},
}

// hiveTableSchema captures the structural elements of a Hive table that
// matter for diff decisions. The struct deliberately ignores presentation
// concerns (quoting, whitespace, comment line position) so that two DDLs
// that are semantically identical compare equal.
type hiveTableSchema struct {
	// columns is the ordered list of column definitions. Order matters
	// because Hive ALTER TABLE ADD COLUMNS appends at the end.
	columns []hiveColumn
	// partitionedBy is the ordered list of partition columns. A non-empty
	// difference always forces DROP+CREATE because Hive cannot change
	// partition keys via ALTER (design §3.4).
	partitionedBy []hiveColumn
	// storedAs is the file format (TEXTFILE, ORC, PARQUET, AVRO, ...). An
	// empty string means the DDL did not declare STORED AS explicitly.
	storedAs string
	// rowFormat captures ROW FORMAT clauses including SERDE class and SERDE
	// properties. A difference forces DROP+CREATE.
	rowFormat string
	// tblProperties contains TBLPROPERTIES with runtime keys filtered out.
	tblProperties map[string]string
	// external is true when the CREATE TABLE statement is CREATE EXTERNAL
	// TABLE. Toggling EXTERNAL is allowed via ALTER per design §3.4 but
	// only via TBLPROPERTIES SET; the parser captures it here so the diff
	// can detect the toggle.
	external bool
	// location is the LOCATION clause value. A difference is treated as an
	// ALTER-able change (ALTER TABLE ... SET LOCATION).
	location string
	// tableName is the unqualified table name extracted from the DDL.
	// Used to compose ALTER statements.
	tableName string
	// comment captures the COMMENT 'xxx' table-level comment (separate from
	// TBLPROPERTIES). Hive supports changing the table comment via
	// ALTER TABLE ... SET TBLPROPERTIES('comment'='new'), so this is an
	// ALTER-able difference.
	comment string
}

// hiveColumn is one column definition: name + type + optional COMMENT.
type hiveColumn struct {
	name    string
	colType string
	comment string
}

// diffTableDDL compares two Hive CREATE TABLE statements (base = source,
// target = destination) and returns the variant SQL strategy:
//
//   - alterStmts: a sequence of ALTER TABLE statements that, applied in
//     order on the target side, will reconcile its structure with the base.
//     Empty when fallbackDropCreate is true.
//   - fallbackDropCreate: true when at least one structural difference falls
//     into the "B. DROP+CREATE" bucket per design.md §3.4 (partition key,
//     storage format, ROW FORMAT/SerDe, EXTERNAL toggle, incompatible
//     column type change, column deletion, or combinations thereof).
//   - err: non-nil only when the DDL strings could not be parsed.
//
// When base and target are semantically identical, returns (nil, false, nil).
func diffTableDDL(base, target string) (alterStmts []string, fallbackDropCreate bool, err error) {
	baseSchema, err := parseHiveTableDDL(base)
	if err != nil {
		return nil, false, fmt.Errorf("parse base DDL: %v", err)
	}
	targetSchema, err := parseHiveTableDDL(target)
	if err != nil {
		return nil, false, fmt.Errorf("parse target DDL: %v", err)
	}

	// Prefer the base table name when emitting ALTER statements; if base is
	// missing the name (unlikely but defensive) fall back to target.
	tableName := baseSchema.tableName
	if tableName == "" {
		tableName = targetSchema.tableName
	}

	// === DROP+CREATE triggers (design §3.4) ===

	// 1. Partition key changes (any add / remove / type-change / reorder).
	if !columnsEqual(baseSchema.partitionedBy, targetSchema.partitionedBy) {
		return nil, true, nil
	}
	// 2. STORED AS changes.
	if !equalFold(baseSchema.storedAs, targetSchema.storedAs) {
		return nil, true, nil
	}
	// 3. ROW FORMAT / SerDe changes.
	if !equalFold(baseSchema.rowFormat, targetSchema.rowFormat) {
		return nil, true, nil
	}
	// 4. EXTERNAL toggle. Per design §3.4 EXTERNAL toggling itself is
	//    "ALTER able" via TBLPROPERTIES('EXTERNAL'='TRUE'); we treat the
	//    keyword-level toggle in CREATE TABLE as an indicator and emit the
	//    corresponding ALTER SET TBLPROPERTIES statement rather than
	//    DROP+CREATE.
	// 5. Incompatible column type change OR column deletion.
	colAlters, colTypeIncompat, colDeleted := diffColumns(
		baseSchema.columns, targetSchema.columns, tableName)
	if colTypeIncompat || colDeleted {
		return nil, true, nil
	}

	// === ALTER path: collect statements ===

	// Columns first (ADD, CHANGE for rename/widen/comment).
	alterStmts = append(alterStmts, colAlters...)

	// EXTERNAL toggle via TBLPROPERTIES.
	if baseSchema.external != targetSchema.external {
		val := "FALSE"
		if baseSchema.external {
			val = "TRUE"
		}
		alterStmts = append(alterStmts, fmt.Sprintf(
			"ALTER TABLE %s SET TBLPROPERTIES ('EXTERNAL'='%s');", tableName, val))
	}

	// Table comment change.
	if baseSchema.comment != targetSchema.comment {
		alterStmts = append(alterStmts, fmt.Sprintf(
			"ALTER TABLE %s SET TBLPROPERTIES ('comment'='%s');",
			tableName, escapeProperty(baseSchema.comment)))
	}

	// LOCATION change (ALTER-able per design §3.4).
	if !equalFold(baseSchema.location, targetSchema.location) {
		if baseSchema.location != "" {
			alterStmts = append(alterStmts, fmt.Sprintf(
				"ALTER TABLE %s SET LOCATION '%s';", tableName, baseSchema.location))
		}
	}

	// TBLPROPERTIES (business-level keys; runtime keys already filtered).
	tblAlters := diffTBLProperties(baseSchema.tblProperties, targetSchema.tblProperties, tableName)
	alterStmts = append(alterStmts, tblAlters...)

	return alterStmts, false, nil
}

// diffColumns walks base/target column lists and returns ALTER statements to
// reconcile them. It also reports whether any difference is "incompatible"
// (forcing DROP+CREATE).
//
// Rules (design §3.4 row "TABLE / 改列类型"):
//   - same column count, same names, only type widening or comment change
//     -> ALTER CHANGE COLUMN per affected column (compatible).
//   - rename (same position, different name) -> ALTER CHANGE COLUMN.
//   - base has new columns at the end that target lacks -> ALTER ADD COLUMNS.
//   - target has columns that base lacks (target column deletion from
//     base perspective) -> column deletion, forces DROP+CREATE.
//   - column type change that is not a widening direction -> incompatible.
//
// "Widening" is conservatively limited to the pairs explicitly listed in
// design §3.4 ("int → bigint" etc.). Everything else is incompatible.
func diffColumns(base, target []hiveColumn, tableName string) (alters []string, incompatible, deleted bool) {
	// Detect deletion first: any target column whose name is not in base
	// (case-insensitive) signals a column removal from base → target.
	baseNames := make(map[string]int)
	for i, c := range base {
		baseNames[strings.ToLower(c.name)] = i
	}
	targetNames := make(map[string]int)
	for i, c := range target {
		targetNames[strings.ToLower(c.name)] = i
	}
	for _, c := range target {
		if _, ok := baseNames[strings.ToLower(c.name)]; !ok {
			deleted = true
			return
		}
	}

	// For each column in base, decide ADD vs CHANGE vs no-op.
	for _, bc := range base {
		idx, ok := targetNames[strings.ToLower(bc.name)]
		if !ok {
			// New column in base → ADD COLUMNS.
			col := fmt.Sprintf("%s %s", bc.name, bc.colType)
			if bc.comment != "" {
				col += fmt.Sprintf(" COMMENT '%s'", escapeProperty(bc.comment))
			}
			alters = append(alters, fmt.Sprintf(
				"ALTER TABLE %s ADD COLUMNS (%s);", tableName, col))
			continue
		}
		tc := target[idx]

		// Same name -> compare type + comment.
		if !typesEqual(bc.colType, tc.colType) {
			if !isCompatibleTypeChange(tc.colType, bc.colType) {
				incompatible = true
				return
			}
			// Compatible: emit CHANGE COLUMN.
			alters = append(alters, formatChangeColumn(tableName, bc))
			continue
		}
		if bc.comment != tc.comment {
			// Only comment differs.
			alters = append(alters, formatChangeColumn(tableName, bc))
		}
	}
	return alters, incompatible, deleted
}

// formatChangeColumn produces `ALTER TABLE t CHANGE COLUMN c c type COMMENT 'x';`
func formatChangeColumn(tableName string, c hiveColumn) string {
	s := fmt.Sprintf("ALTER TABLE %s CHANGE COLUMN %s %s %s",
		tableName, c.name, c.name, c.colType)
	if c.comment != "" {
		s += fmt.Sprintf(" COMMENT '%s'", escapeProperty(c.comment))
	}
	return s + ";"
}

// isCompatibleTypeChange reports whether changing a column from `from` to
// `to` is safe (no data loss). We support the widening pairs documented in
// design §3.4: int → bigint, smallint → int / bigint, tinyint → smallint /
// int / bigint, float → double, char/varchar → string (Hive treats string as
// the widest text). Everything else is treated as incompatible.
func isCompatibleTypeChange(from, to string) bool {
	f := normalizeType(from)
	t := normalizeType(to)
	if f == t {
		return true
	}
	widening := map[string]map[string]bool{
		"tinyint":  {"smallint": true, "int": true, "bigint": true},
		"smallint": {"int": true, "bigint": true},
		"int":      {"bigint": true},
		"float":    {"double": true},
		"char":     {"string": true, "varchar": true},
		"varchar":  {"string": true},
	}
	if dest, ok := widening[f]; ok {
		return dest[t]
	}
	return false
}

// typesEqual normalizes whitespace and case before comparing.
func typesEqual(a, b string) bool {
	return normalizeType(a) == normalizeType(b)
}

// normalizeType collapses whitespace, lowercases, and strips the size
// specifier from char/varchar so "varchar(20)" compares equal to "varchar".
// This matches Hive's permissive widening semantics: only the base type
// matters for the widening matrix.
func normalizeType(t string) string {
	t = strings.ToLower(strings.TrimSpace(t))
	if i := strings.IndexByte(t, '('); i >= 0 {
		t = t[:i]
	}
	return strings.TrimSpace(t)
}

// columnsEqual reports whether two column lists are identical in name,
// type, and order (comment is ignored; partition keys carry no comment).
func columnsEqual(a, b []hiveColumn) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !strings.EqualFold(a[i].name, b[i].name) {
			return false
		}
		if !typesEqual(a[i].colType, b[i].colType) {
			return false
		}
	}
	return true
}

// diffTBLProperties emits ALTER TABLE ... SET TBLPROPERTIES statements for
// keys whose values differ. Both inputs already have runtime keys filtered.
// The output is deterministic (keys sorted lexicographically).
func diffTBLProperties(base, target map[string]string, tableName string) []string {
	var changed []string
	keys := make(map[string]struct{})
	for k := range base {
		keys[k] = struct{}{}
	}
	for k := range target {
		keys[k] = struct{}{}
	}
	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var setParts []string
	for _, k := range sortedKeys {
		if base[k] != target[k] {
			// Even "delete from base" → emit empty-string set (Hive does not
			// support UNSET via SET; UNSET TBLPROPERTIES exists separately,
			// but for our diff we treat removal as setting empty).
			setParts = append(setParts, fmt.Sprintf("'%s'='%s'",
				escapeProperty(k), escapeProperty(base[k])))
		}
	}
	if len(setParts) > 0 {
		changed = append(changed, fmt.Sprintf(
			"ALTER TABLE %s SET TBLPROPERTIES (%s);",
			tableName, strings.Join(setParts, ", ")))
	}
	return changed
}

// escapeProperty replaces single quotes inside property values with the
// Hive-safe doubled form ''.
func escapeProperty(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// equalFold is shorthand for case-insensitive equality.
func equalFold(a, b string) bool {
	return strings.EqualFold(strings.TrimSpace(a), strings.TrimSpace(b))
}

// parseHiveTableDDL parses a Hive CREATE TABLE statement into a
// hiveTableSchema. The parser is intentionally simple and section-based:
// it walks the DDL token by token (keyword-driven) and collects each
// section's body into the corresponding hiveTableSchema field. It does NOT
// attempt to validate the SQL — invalid DDL still produces a partial schema.
//
// Why no regex (per dev.md "非必要少用正则表达式"):
//   - Hive DDL nesting (parentheses inside TBLPROPERTIES values, quoted
//     SerDe property strings) is difficult to express precisely with
//     regex without backtracking pitfalls.
//   - A section-driven loop is easier to extend (new section types added
//     in later Hive versions).
func parseHiveTableDDL(ddl string) (*hiveTableSchema, error) {
	if strings.TrimSpace(ddl) == "" {
		return nil, fmt.Errorf("empty DDL")
	}
	s := &hiveTableSchema{
		tblProperties: map[string]string{},
	}

	// Tokenize lightly: keep the original DDL but identify section
	// boundaries by uppercase keywords. We scan through the DDL once.
	upper := strings.ToUpper(ddl)

	// EXTERNAL flag and table name come from the CREATE TABLE prefix.
	cidx := strings.Index(upper, "CREATE")
	if cidx < 0 {
		return nil, fmt.Errorf("missing CREATE keyword")
	}
	// Check for EXTERNAL between CREATE and TABLE.
	tidx := strings.Index(upper[cidx:], "TABLE")
	if tidx < 0 {
		return nil, fmt.Errorf("missing TABLE keyword")
	}
	prefix := upper[cidx : cidx+tidx]
	if strings.Contains(prefix, "EXTERNAL") {
		s.external = true
	}
	// Table name: first token after TABLE (skipping IF NOT EXISTS).
	afterTable := strings.TrimSpace(ddl[cidx+tidx+len("TABLE"):])
	afterTableUpper := strings.ToUpper(afterTable)
	if strings.HasPrefix(afterTableUpper, "IF NOT EXISTS") {
		afterTable = strings.TrimSpace(afterTable[len("IF NOT EXISTS"):])
	}
	// Take everything up to the next whitespace or '(' as the name.
	nameEnd := strings.IndexAny(afterTable, " \t\n(")
	if nameEnd < 0 {
		nameEnd = len(afterTable)
	}
	s.tableName = stripQualifier(strings.Trim(afterTable[:nameEnd], "`"))

	// Body parsing: find the column-list parenthesized block first.
	openIdx := strings.IndexByte(afterTable, '(')
	if openIdx < 0 {
		// No columns; this is an unusual but valid DDL (e.g. CTAS).
		return s, nil
	}
	closeIdx := matchParen(afterTable, openIdx)
	if closeIdx < 0 {
		return nil, fmt.Errorf("unmatched column list parenthesis")
	}
	columnsBody := afterTable[openIdx+1 : closeIdx]
	s.columns = parseColumnList(columnsBody)

	// Remaining = everything after the column list.
	rest := strings.TrimSpace(afterTable[closeIdx+1:])
	parseSections(rest, s)
	return s, nil
}

// stripQualifier turns "db.t" into "t" (Hive supports schema-qualified
// names in SHOW CREATE output).
func stripQualifier(s string) string {
	if i := strings.LastIndex(s, "."); i >= 0 {
		return s[i+1:]
	}
	return s
}

// matchParen returns the index of the matching closing parenthesis for the
// '(' at openIdx, respecting single-quoted strings. Returns -1 on mismatch.
func matchParen(s string, openIdx int) int {
	if openIdx < 0 || openIdx >= len(s) || s[openIdx] != '(' {
		return -1
	}
	depth := 0
	inStr := false
	for i := openIdx; i < len(s); i++ {
		c := s[i]
		if c == '\'' && (i == 0 || s[i-1] != '\\') {
			inStr = !inStr
			continue
		}
		if inStr {
			continue
		}
		switch c {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// parseColumnList parses the comma-separated column list inside the
// parenthesized body. Each entry is "name type [COMMENT 'x']". Commas
// inside parentheses (decimal(10,2), struct<a:int, b:string>) are honored.
func parseColumnList(body string) []hiveColumn {
	var cols []hiveColumn
	for _, raw := range splitTopLevelCommas(body) {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		col := parseColumnEntry(entry)
		if col.name != "" {
			cols = append(cols, col)
		}
	}
	return cols
}

// splitTopLevelCommas splits a string by commas that are NOT inside
// parentheses or single-quoted strings.
func splitTopLevelCommas(s string) []string {
	var parts []string
	var buf strings.Builder
	depth := 0
	inStr := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' && (i == 0 || s[i-1] != '\\') {
			inStr = !inStr
		}
		if !inStr {
			switch c {
			case '(', '<':
				depth++
			case ')', '>':
				depth--
			case ',':
				if depth == 0 {
					parts = append(parts, buf.String())
					buf.Reset()
					continue
				}
			}
		}
		buf.WriteByte(c)
	}
	if buf.Len() > 0 {
		parts = append(parts, buf.String())
	}
	return parts
}

// parseColumnEntry parses a single column definition: name type [COMMENT 'x'].
func parseColumnEntry(entry string) hiveColumn {
	col := hiveColumn{}
	entry = strings.TrimSpace(strings.Trim(entry, "\n\r\t "))
	// Name: first whitespace-delimited token, may be backtick-quoted.
	nameEnd := indexAnyOutsideQuotes(entry, " \t\n")
	if nameEnd < 0 {
		col.name = strings.Trim(entry, "`")
		return col
	}
	col.name = strings.Trim(entry[:nameEnd], "`")
	rest := strings.TrimSpace(entry[nameEnd:])
	// COMMENT may appear after the type. Find " COMMENT " (case-insensitive)
	// outside of quotes.
	upperRest := strings.ToUpper(rest)
	cidx := strings.Index(upperRest, " COMMENT ")
	if cidx >= 0 {
		col.colType = strings.TrimSpace(rest[:cidx])
		commentPart := strings.TrimSpace(rest[cidx+len(" COMMENT "):])
		col.comment = stripQuotes(commentPart)
	} else {
		col.colType = strings.TrimSpace(rest)
	}
	return col
}

// indexAnyOutsideQuotes returns the index of the first byte in `chars`
// that occurs outside of a single-quoted string.
func indexAnyOutsideQuotes(s, chars string) int {
	inStr := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' {
			inStr = !inStr
			continue
		}
		if inStr {
			continue
		}
		if strings.IndexByte(chars, c) >= 0 {
			return i
		}
	}
	return -1
}

// stripQuotes removes surrounding single quotes (Hive comment literals).
func stripQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '\'' && s[len(s)-1] == '\'' {
		return s[1 : len(s)-1]
	}
	return s
}

// parseSections walks the post-column-list body and fills the schema with
// COMMENT, PARTITIONED BY, ROW FORMAT, STORED AS, LOCATION, TBLPROPERTIES.
// Sections are keyword-driven; their order in Hive DDL is fixed but the
// parser does not rely on the order.
func parseSections(rest string, s *hiveTableSchema) {
	cursor := 0
	for cursor < len(rest) {
		// Skip whitespace.
		for cursor < len(rest) && isWhitespace(rest[cursor]) {
			cursor++
		}
		if cursor >= len(rest) {
			break
		}

		remaining := rest[cursor:]
		upper := strings.ToUpper(remaining)

		switch {
		case strings.HasPrefix(upper, "COMMENT "):
			cursor += len("COMMENT ")
			val, consumed := consumeQuotedString(rest[cursor:])
			s.comment = val
			cursor += consumed
		case strings.HasPrefix(upper, "PARTITIONED BY"):
			cursor += len("PARTITIONED BY")
			cursor = skipWhitespace(rest, cursor)
			if cursor < len(rest) && rest[cursor] == '(' {
				closeIdx := matchParen(rest, cursor)
				if closeIdx > cursor {
					s.partitionedBy = parseColumnList(rest[cursor+1 : closeIdx])
					cursor = closeIdx + 1
				} else {
					cursor++
				}
			}
		case strings.HasPrefix(upper, "ROW FORMAT"):
			// Capture everything up to the next top-level section keyword.
			cursor += len("ROW FORMAT")
			next := findNextSectionKeyword(rest, cursor)
			s.rowFormat = strings.TrimSpace(rest[cursor:next])
			cursor = next
		case strings.HasPrefix(upper, "STORED AS"):
			cursor += len("STORED AS")
			next := findNextSectionKeyword(rest, cursor)
			s.storedAs = strings.TrimSpace(rest[cursor:next])
			cursor = next
		case strings.HasPrefix(upper, "STORED BY"):
			cursor += len("STORED BY")
			next := findNextSectionKeyword(rest, cursor)
			// Treat STORED BY like storedAs for diff purposes.
			s.storedAs = "STORED BY " + strings.TrimSpace(rest[cursor:next])
			cursor = next
		case strings.HasPrefix(upper, "LOCATION"):
			cursor += len("LOCATION")
			cursor = skipWhitespace(rest, cursor)
			val, consumed := consumeQuotedString(rest[cursor:])
			s.location = val
			cursor += consumed
		case strings.HasPrefix(upper, "TBLPROPERTIES"):
			cursor += len("TBLPROPERTIES")
			cursor = skipWhitespace(rest, cursor)
			if cursor < len(rest) && rest[cursor] == '(' {
				closeIdx := matchParen(rest, cursor)
				if closeIdx > cursor {
					parseTBLProperties(rest[cursor+1:closeIdx], s)
					cursor = closeIdx + 1
				} else {
					cursor++
				}
			}
		default:
			// Unknown keyword: skip one rune so we make progress.
			cursor++
		}
	}
}

// findNextSectionKeyword returns the index in `s` (>= start) of the next
// top-level Hive table-DDL section keyword (COMMENT, PARTITIONED BY, etc.),
// or len(s) if none is found. Keywords occurring inside parentheses are
// ignored. This allows ROW FORMAT bodies to span multiple lines safely.
func findNextSectionKeyword(s string, start int) int {
	keywords := []string{
		"COMMENT ", "PARTITIONED BY", "CLUSTERED BY", "SKEWED BY",
		"ROW FORMAT", "STORED AS", "STORED BY", "LOCATION", "TBLPROPERTIES",
	}
	depth := 0
	inStr := false
	for i := start; i < len(s); i++ {
		c := s[i]
		if c == '\'' && (i == 0 || s[i-1] != '\\') {
			inStr = !inStr
			continue
		}
		if inStr {
			continue
		}
		switch c {
		case '(':
			depth++
		case ')':
			depth--
		}
		if depth != 0 {
			continue
		}
		// Match keyword at the start of a word (preceded by whitespace).
		if i == start || isWhitespace(s[i-1]) {
			upperFrom := strings.ToUpper(s[i:])
			for _, kw := range keywords {
				if strings.HasPrefix(upperFrom, kw) {
					return i
				}
			}
		}
	}
	return len(s)
}

// parseTBLProperties parses the body of a TBLPROPERTIES (...) section and
// fills s.tblProperties, filtering out runtime keys (compat-RISK-6).
// Body format: 'key1'='value1', 'key2'='value2', ...
func parseTBLProperties(body string, s *hiveTableSchema) {
	for _, raw := range splitTopLevelCommas(body) {
		entry := strings.TrimSpace(raw)
		if entry == "" {
			continue
		}
		// Find '=' that is OUTSIDE quotes.
		eq := indexEqualsOutsideQuotes(entry)
		if eq < 0 {
			continue
		}
		k := stripQuotes(strings.TrimSpace(entry[:eq]))
		v := stripQuotes(strings.TrimSpace(entry[eq+1:]))
		if _, runtime := runtimeTBLPropertyKeys[k]; runtime {
			continue
		}
		// Hive also surfaces EXTERNAL via TBLPROPERTIES; treat that as the
		// external flag, not a property.
		if strings.EqualFold(k, "EXTERNAL") {
			if strings.EqualFold(v, "TRUE") {
				s.external = true
			}
			continue
		}
		// COMMENT in TBLPROPERTIES maps to the table comment field; prefer
		// the inline COMMENT clause if it was already set.
		if strings.EqualFold(k, "comment") {
			if s.comment == "" {
				s.comment = v
			}
			continue
		}
		s.tblProperties[k] = v
	}
}

// indexEqualsOutsideQuotes returns the index of the first '=' outside of
// single-quoted strings.
func indexEqualsOutsideQuotes(s string) int {
	inStr := false
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '\'' && (i == 0 || s[i-1] != '\\') {
			inStr = !inStr
			continue
		}
		if !inStr && c == '=' {
			return i
		}
	}
	return -1
}

// consumeQuotedString reads a single-quoted string starting at the current
// position (after optional leading whitespace) and returns the unquoted
// value plus the number of input bytes consumed.
func consumeQuotedString(s string) (string, int) {
	i := 0
	for i < len(s) && isWhitespace(s[i]) {
		i++
	}
	if i >= len(s) || s[i] != '\'' {
		return "", i
	}
	start := i + 1
	for j := start; j < len(s); j++ {
		if s[j] == '\'' && (j == start || s[j-1] != '\\') {
			return s[start:j], j + 1
		}
	}
	return s[start:], len(s)
}

// skipWhitespace returns the index of the next non-whitespace byte at or
// after `i`.
func skipWhitespace(s string, i int) int {
	for i < len(s) && isWhitespace(s[i]) {
		i++
	}
	return i
}

// isWhitespace reports whether c is one of the ASCII whitespace bytes.
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
