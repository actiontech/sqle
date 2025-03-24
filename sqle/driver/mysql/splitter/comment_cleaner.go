package splitter

import (
	"strings"
	"unicode"
)

/*
该方法用于清除sql语句中的注释信息
*/
func removeSQLComments(sql string) string {
	var result []rune
	inSingleQuote, inDoubleQuote, inBackQuote := false, false, false
	inLineComment, inBlockComment := false, false
	runes := []rune(sql)
	n := len(runes)

	for i := 0; i < n; i++ {
		c := runes[i]

		// 结束行注释
		if inLineComment {
			if c == '\n' || c == '\r' {
				inLineComment = false
				result = append(result, c)
			}
			continue
		}

		// 结束块注释
		if inBlockComment {
			if c == '*' && i+1 < n && runes[i+1] == '/' {
				inBlockComment = false
				i++ // 跳过 '/'
				// 在块注释结束后，检查下一个字符
				if i+1 < n && !unicode.IsSpace(runes[i+1]) {
					// 只有当 result 最后字符不是空格时，才插入一个空格
					if len(result) == 0 || !unicode.IsSpace(result[len(result)-1]) {
						result = append(result, ' ')
					}
				}
			}
			continue
		}

		// 进入引号状态（进入前保留引号）
		if !inSingleQuote && !inDoubleQuote && !inBackQuote {
			// 行注释：-- 开始
			if c == '-' && i+1 < n && runes[i+1] == '-' {
				inLineComment = true
				i++ // 跳过第二个 -
				continue
			}
			// 行注释：# 开始
			if c == '#' {
				inLineComment = true
				continue
			}
			if c == '/' && i+1 < n && runes[i+1] == '*' {
				// 判断是否为 Hint：检查是否有 "+" 符号
				if i+2 < n && runes[i+2] == '+' {
					// Hint 不删除，原样输出整个 Hint 注释
					startIndex := i
					// 找到 Hint 结束位置
					endPos := findCommentEnd(runes, i)
					result = append(result, runes[startIndex:endPos]...)
					i = endPos - 1
					continue
				} else {
					// 非 Hint 块注释，进入删除状态
					inBlockComment = true
					i++ // 跳过 '*'
					continue
				}
			}
		}

		// 状态切换：引号内不检查注释
		if c == '\'' && !inDoubleQuote && !inBackQuote {
			inSingleQuote = !inSingleQuote
		} else if c == '"' && !inSingleQuote && !inBackQuote {
			inDoubleQuote = !inDoubleQuote
		} else if c == '`' && !inSingleQuote && !inDoubleQuote {
			inBackQuote = !inBackQuote
		}

		result = append(result, c)
	}
	// 处理dump文件可能存在多行注释结尾有分号问题
	if strings.TrimSpace(string(result)) == ";" {
		return ""
	}
	return string(result)
}

// findCommentEnd 返回从 pos 开始的块注释结束位置（包括 "*/"），如果没找到则返回 n
func findCommentEnd(runes []rune, pos int) int {
	n := len(runes)
	for pos < n-1 {
		if runes[pos] == '*' && runes[pos+1] == '/' {
			return pos + 2
		}
		pos++
	}
	return n
}
