//go:build enterprise
// +build enterprise

package differ

import (
	"fmt"
	"strings"
)

// Index represents a single index (primary key, unique secondary index, or non-
// unique secondard index) in a table.
type Index struct {
	Name           string      `json:"name"`
	Parts          []IndexPart `json:"parts"`
	PrimaryKey     bool        `json:"primaryKey,omitempty"`
	Unique         bool        `json:"unique,omitempty"`
	Invisible      bool        `json:"invisible,omitempty"` // MySQL 8+, also used for MariaDB 10.6's IGNORED indexes
	Comment        string      `json:"comment,omitempty"`
	Type           string      `json:"type"`
	FullTextParser string      `json:"parser,omitempty"`
}

// IndexPart represents an individual indexed column or expression. Each index
// has one or more IndexPart values.
type IndexPart struct {
	ColumnName   string `json:"columnName,omitempty"`   // name of column, or empty if expression
	Expression   string `json:"expression,omitempty"`   // expression value (MySQL 8+), or empty if column
	PrefixLength uint16 `json:"prefixLength,omitempty"` // nonzero if only a prefix of column is indexed
	Descending   bool   `json:"descending,omitempty"`   // if true, collation is descending (MySQL 8+)
}

// Definition returns this index's definition clause, for use as part of a DDL
// statement.
func (idx *Index) Definition(flavor Flavor) string {
	parts := make([]string, len(idx.Parts))
	for n := range idx.Parts {
		parts[n] = idx.Parts[n].Definition(flavor)
	}
	var typeAndName, comment, invis, parser string
	if idx.PrimaryKey {
		if !idx.Unique {
			// TODO 错误处理
			return "index is primary key, but isn't marked as unique"
		}
		typeAndName = "PRIMARY KEY"
	} else if idx.Unique {
		typeAndName = fmt.Sprintf("UNIQUE KEY %s", EscapeIdentifier(idx.Name))
	} else if idx.Type != "BTREE" && idx.Type != "" {
		typeAndName = fmt.Sprintf("%s KEY %s", idx.Type, EscapeIdentifier(idx.Name))
	} else {
		typeAndName = fmt.Sprintf("KEY %s", EscapeIdentifier(idx.Name))
	}
	if idx.Comment != "" {
		comment = fmt.Sprintf(" COMMENT '%s'", EscapeValueForCreateTable(idx.Comment))
	}
	if idx.Invisible {
		if flavor.IsMariaDB() {
			invis = " IGNORED"
		} else {
			invis = " /*!80000 INVISIBLE */"
		}
	}
	if idx.Type == "FULLTEXT" && idx.FullTextParser != "" {
		// Note the trailing space here is intentional -- it's always present in SHOW
		// CREATE TABLE for this particular clause
		parser = fmt.Sprintf(" /*!50100 WITH PARSER `%s` */ ", idx.FullTextParser)
	}
	return fmt.Sprintf("%s (%s)%s%s%s", typeAndName, strings.Join(parts, ","), comment, invis, parser)
}

// Equals returns true if two indexes are completely identical, false otherwise.
func (idx *Index) Equals(other *Index) bool {
	if idx == nil || other == nil {
		return idx == other // only equal if BOTH are nil
	}
	return idx.Name == other.Name && idx.Comment == other.Comment && idx.Invisible == other.Invisible && idx.Equivalent(other)
}

// sameParts returns true if two Indexes' Parts slices are identical.
func (idx *Index) sameParts(other *Index) bool {
	if len(idx.Parts) != len(other.Parts) {
		return false
	}
	for n := range idx.Parts {
		if idx.Parts[n] != other.Parts[n] {
			return false
		}
	}
	return true
}

// Equivalent returns true if two Indexes are functionally equivalent,
// regardless of whether or not they have the same names, comments, or
// visibility.
func (idx *Index) Equivalent(other *Index) bool {
	if idx == nil || other == nil {
		return idx == other // only equivalent if BOTH are nil
	}
	if idx.PrimaryKey != other.PrimaryKey || idx.Unique != other.Unique || idx.Type != other.Type || idx.FullTextParser != other.FullTextParser {
		return false
	}
	return idx.sameParts(other)
}

// RedundantTo returns true if idx is equivalent to, or a strict subset of,
// other. Both idx and other should be indexes of the same table.
// A non-unique index is considered redundant to any other same-type index
// having the same (or more) columns in the same order, unless its parts have a
// greater column prefix length. A unique index can only be redundant to the
// primary key or an exactly equivalent unique index; another unique index with
// more cols may coexist due to the desired constraint semantics. A primary key
// is never redundant to another index.
func (idx *Index) RedundantTo(other *Index) bool {
	if idx == nil || other == nil {
		return false
	}
	if idx.PrimaryKey || (idx.Unique && !other.Unique) || idx.Type != other.Type || idx.FullTextParser != other.FullTextParser {
		return false
	}
	if !idx.Invisible && other.Invisible {
		return false // a visible index is never redundant to an invisible one
	}
	if idx.Unique && other.Unique {
		// Since unique indexes are also unique *constraints*, two unique indexes are
		// non-redundant unless they have identical parts.
		return idx.sameParts(other)
	} else if idx.Type == "FULLTEXT" && len(idx.Parts) != len(other.Parts) {
		return false // FT composite indexes don't behave like BTREE in terms of left-right prefixing
	} else if len(idx.Parts) > len(other.Parts) {
		return false // can't be redundant to an index with fewer cols
	}
	for n, part := range idx.Parts {
		if part.ColumnName != other.Parts[n].ColumnName || part.Expression != other.Parts[n].Expression || part.Descending != other.Parts[n].Descending {
			return false
		}
		partPrefix, otherPrefix := part.PrefixLength, other.Parts[n].PrefixLength
		if otherPrefix > 0 && (partPrefix == 0 || partPrefix > otherPrefix) {
			return false
		}
	}
	return true
}

// Functional returns true if at least one IndexPart in idx is an expression
// rather than a column.
func (idx *Index) Functional() bool {
	for _, part := range idx.Parts {
		if part.Expression != "" {
			return true
		}
	}
	return false
}

// Definition returns this index part's definition clause.
func (part *IndexPart) Definition(_ Flavor) string {
	var base, prefix, collation string
	if part.ColumnName != "" {
		base = EscapeIdentifier(part.ColumnName)
	} else {
		base = fmt.Sprintf("(%s)", part.Expression)
	}
	if part.PrefixLength > 0 {
		prefix = fmt.Sprintf("(%d)", part.PrefixLength)
	}
	if part.Descending {
		collation = " DESC"
	}
	return fmt.Sprintf("%s%s%s", base, prefix, collation)
}
