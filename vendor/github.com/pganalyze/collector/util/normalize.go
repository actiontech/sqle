package util

import (
	pg_query "github.com/pganalyze/pg_query_go/v4"
)

// TruncatedQueryMarker - Added to queries that were truncated and caused a
// parsing error, but we were able to "fix" to allow grouping in the UI
const TruncatedQueryMarker string = "/* truncated-query */ "

// NormalizeQuery - Normalizes the query text with an improved variant of the
// pg_stat_statements normalization logic.
func NormalizeQuery(query string, filterQueryText string, trackActivityQuerySize int) string {
	normalizedQuery, err := pg_query.Normalize(query)
	if err == nil {
		return normalizedQuery
	}

	fixedQuery := fixTruncatedQuery(query)
	normalizedQuery, err = pg_query.Normalize(fixedQuery)
	if err == nil {
		return TruncatedQueryMarker + normalizedQuery
	}

	if filterQueryText == "none" {
		return query
	} else if len(query) == trackActivityQuerySize-1 {
		return QueryTextTruncated
	} else {
		return QueryTextUnparsable
	}
}
