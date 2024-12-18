//go:build enterprise
// +build enterprise

package differ

import (
	"fmt"
	"strings"
)

// PartitionListMode values control edge-cases for how the list of partitions
// is represented in SHOW CREATE TABLE.
type PartitionListMode string

// Constants enumerating valid PartitionListMode values.
const (
	PartitionListDefault  PartitionListMode = ""          // Default behavior based on partitioning method
	PartitionListExplicit PartitionListMode = "full"      // List each partition individually
	PartitionListCount    PartitionListMode = "countOnly" // Just use a count of partitions
	PartitionListNone     PartitionListMode = "omit"      // Omit partition list and count, implying just 1 partition
)

// TablePartitioning stores partitioning configuration for a partitioned table.
// Note that despite subpartitioning fields being present and possibly
// populated, the rest of this package does not fully support subpartitioning
// yet.
type TablePartitioning struct {
	Method             string            `json:"method"`              // one of "RANGE", "RANGE COLUMNS", "LIST", "LIST COLUMNS", "HASH", "LINEAR HASH", "KEY", or "LINEAR KEY"
	SubMethod          string            `json:"subMethod,omitempty"` // one of "" (no sub-partitioning), "HASH", "LINEAR HASH", "KEY", or "LINEAR KEY"; not fully supported yet
	Expression         string            `json:"expression"`
	SubExpression      string            `json:"subExpression,omitempty"` // empty string if no sub-partitioning; not fully supported yet
	Partitions         []*Partition      `json:"partitions"`
	ForcePartitionList PartitionListMode `json:"forcePartitionList,omitempty"`
	AlgoClause         string            `json:"algoClause,omitempty"` // full text of optional ALGORITHM clause for KEY or LINEAR KEY
}

// Definition returns the overall partitioning definition for a table.
func (tp *TablePartitioning) Definition(flavor Flavor) string {
	if tp == nil {
		return ""
	}

	plMode := tp.ForcePartitionList
	if plMode == PartitionListDefault {
		plMode = PartitionListCount
		for n, p := range tp.Partitions {
			if p.Values != "" || p.Comment != "" || p.DataDir != "" || p.Name != fmt.Sprintf("p%d", n) {
				plMode = PartitionListExplicit
				break
			}
		}
	}
	var partitionsClause string
	if plMode == PartitionListExplicit {
		pdefs := make([]string, len(tp.Partitions))
		for n, p := range tp.Partitions {
			pdefs[n] = p.Definition(flavor, tp.Method)
		}
		partitionsClause = fmt.Sprintf("\n(%s)", strings.Join(pdefs, ",\n "))
	} else if plMode == PartitionListCount {
		partitionsClause = fmt.Sprintf("\nPARTITIONS %d", len(tp.Partitions))
	}

	opener, closer := "/*!50100", " */"
	if flavor.MinMariaDB(10, 2) {
		// MariaDB stopped wrapping partitioning clauses in version-gated comments
		// in 10.2.
		opener, closer = "", ""
	} else if strings.HasSuffix(tp.Method, "COLUMNS") {
		// RANGE COLUMNS and LIST COLUMNS were introduced in 5.5
		opener = "/*!50500"
	}

	return fmt.Sprintf("\n%s PARTITION BY %s%s%s", opener, tp.partitionBy(flavor), partitionsClause, closer)
}

// partitionBy returns the partitioning method and expression, formatted to
// match SHOW CREATE TABLE's extremely arbitrary, completely inconsistent way.
func (tp *TablePartitioning) partitionBy(flavor Flavor) string {
	method, expr := fmt.Sprintf("%s ", tp.Method), tp.Expression

	if tp.Method == "RANGE COLUMNS" {
		method = "RANGE  COLUMNS"
	} else if tp.Method == "LIST COLUMNS" {
		method = "LIST  COLUMNS"
	}

	// MySQL (any version) and MariaDB 10.1 (but not later) normally omit the
	// backticks around column names in the partitioning expression, if the method
	// is RANGE COLUMNS, LIST COLUMNS, KEY, or LINEAR KEY.
	// TODO handle edge cases where the backticks are still present: column name is
	// a keyword (even if not a *reserved* word) or contains special characters.
	// See https://github.com/skeema/skeema/issues/199
	if (strings.HasSuffix(tp.Method, "COLUMNS") || strings.HasSuffix(tp.Method, "KEY")) && !flavor.MinMariaDB(10, 2) {
		expr = strings.Replace(expr, "`", "", -1)
	}

	return fmt.Sprintf("%s%s(%s)", method, tp.AlgoClause, expr)
}

// Diff returns a set of differences between this TablePartitioning and another
// TablePartitioning. If supported==true, the returned clauses (if executed)
// would transform tp into other.
func (tp *TablePartitioning) Diff(other *TablePartitioning) (clauses []TableAlterClause, supported bool) {
	// Handle cases where one or both sides are nil, meaning one or both tables are
	// unpartitioned
	if tp == nil && other == nil {
		return nil, true
	} else if tp == nil {
		return []TableAlterClause{PartitionBy{Partitioning: other}}, true
	} else if other == nil {
		return []TableAlterClause{RemovePartitioning{}}, true
	}

	// Modifications to partitioning method or expression: re-partition
	if tp.Method != other.Method || tp.SubMethod != other.SubMethod ||
		tp.Expression != other.Expression || tp.SubExpression != other.SubExpression ||
		tp.AlgoClause != other.AlgoClause {
		clause := PartitionBy{
			Partitioning: other,
			RePartition:  true,
		}
		return []TableAlterClause{clause}, true
	}

	// Modifications to partition list: ignored for RANGE, RANGE COLUMNS, LIST,
	// LIST COLUMNS via generation of a no-op placeholder clause. This is done
	// to side-step the safety mechanism at the end of Table.Diff() which treats 0
	// clauses as indicative of an unsupported diff.
	// For other partitioning methods, changing the partition list is currently
	// unsupported.
	var foundPartitionsDiff bool
	if len(tp.Partitions) != len(other.Partitions) {
		foundPartitionsDiff = true
	} else {
		for n := range tp.Partitions {
			// all Partition fields are scalars, so simple comparison is fine
			if *tp.Partitions[n] != *other.Partitions[n] {
				foundPartitionsDiff = true
				break
			}
		}
	}
	if foundPartitionsDiff && (strings.HasPrefix(tp.Method, "RANGE") || strings.HasPrefix(tp.Method, "LIST")) {
		clause := PartitionBy{
			Partitioning: other,
			RePartition:  true,
		}
		return []TableAlterClause{clause}, true
	}
	return nil, !foundPartitionsDiff
}

// Partition stores information on a single partition.
type Partition struct {
	Name    string `json:"name"`
	SubName string `json:"subName,omitempty"` // empty string if no sub-partitioning; not fully supported yet
	Values  string `json:"values,omitempty"`  // only populated for RANGE or LIST
	Comment string `json:"comment,omitempty"`
	Engine  string `json:"engine"`
	DataDir string `json:"dataDir,omitempty"`
}

// Definition returns this partition's definition clause, for use as part of a
// DDL statement.
func (p *Partition) Definition(flavor Flavor, method string) string {
	// MariaDB 10.2+ wraps partition names in backticks.
	// TODO MySQL (any version) and MariaDB 10.1 will also wrap a partition name in
	// backticks if the name is a keyword (even if not a *reserved* word) or has
	// special characters. See https://github.com/skeema/skeema/issues/175
	name := p.Name
	if flavor.MinMariaDB(10, 2) {
		name = EscapeIdentifier(name)
	}

	var values string
	if method == "RANGE" && p.Values == "MAXVALUE" {
		values = "VALUES LESS THAN MAXVALUE "
	} else if strings.Contains(method, "RANGE") {
		values = fmt.Sprintf("VALUES LESS THAN (%s) ", p.Values)
	} else if strings.Contains(method, "LIST") {
		values = fmt.Sprintf("VALUES IN (%s) ", p.Values)
	}

	var dataDir string
	if p.DataDir != "" {
		dataDir = fmt.Sprintf("DATA DIRECTORY = '%s' ", p.DataDir) // any necessary escaping is already present in p.DataDir
	}

	var comment string
	if p.Comment != "" {
		comment = fmt.Sprintf("COMMENT = '%s' ", EscapeValueForCreateTable(p.Comment))
	}

	return fmt.Sprintf("PARTITION %s %s%s%sENGINE = %s", name, values, dataDir, comment, p.Engine)
}
