package ast

import (
	"fmt"
	"strings"
	"unicode"
)

func replaceWhere (data string) string {
	lowerData := strings.ToLower(strings.TrimSpace(data))
	// find where keyword.
	hasWhereKeyword := false
	if strings.HasPrefix(lowerData,"where") {
		trimData := strings.TrimPrefix(lowerData, "where")
		if len(trimData) > 0 && unicode.IsSpace(rune(trimData[0])){
			hasWhereKeyword = true
		}
	}
	if hasWhereKeyword {
		return fmt.Sprintf("AND%s", strings.TrimSpace(data)[5:]/* string slice, skip 5 for `where` */)
	}
	return data
}
