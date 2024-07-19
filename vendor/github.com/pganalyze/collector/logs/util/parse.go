package util

import "strings"

func WasTruncated(line string) bool {
	// RDS limits log lines to 1MB, and inserts a message if the line overflows. We may
	// stitch the line together with the truncation message, so we need to look for it
	// anywhere in the line.
	return strings.Contains(line, "[Your log message was truncated]")
}
