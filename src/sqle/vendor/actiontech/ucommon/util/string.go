package util

import (
	"math/rand"
	"strings"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GetSection(str string, idx int) string {
	segs := strings.Split(str, "/")
	if idx >= 0 {
		return segs[idx]
	} else {
		return segs[len(segs)+idx]
	}
}

func StringNvl(strs ...string) string {
	for _, str := range strs {
		if "" != str {
			return str
		}
	}
	return ""
}

// SectionMatch supports:
// 1. * pattern, only one in one section, like a/b/*/d
// 2. ** pattern, only one in the end section, listed in matches, like a/b/c/**
// 3. ~~ pattern, only one in the end section, not listed in matches, like a/b/c/??
// 4. | pattern, multi in one section, like a/b|c|d/e
// 5. \* \| \\ \~ pattern, escape chars for * | \, like "a/b/\\*/c"
// 6. TODO: character level match support, but not section level
// return {true, "{matching string of *}"} or {false, ""}
func SectionMatch(str string, pattern string) (bool, []string) {
	if strings.Contains(pattern, `\`) {
		str = strings.Replace(str, "*", "\r", -1)
		str = strings.Replace(str, "~", "\r", -1)
		str = strings.Replace(str, "|", "\n", -1)
		pattern = strings.Replace(pattern, `\*`, "\r", -1)
		pattern = strings.Replace(pattern, `\~`, "\r", -1)
		pattern = strings.Replace(pattern, `\|`, "\n", -1)
		pattern = strings.Replace(pattern, `\\`, "\\", -1)
	}
	segs := strings.Split(str, "/")
	patternSegs := strings.Split(pattern, "/")

	if len(segs) < len(patternSegs) {
		return false, []string{}
	}

	matchs := []string{}
PATTERN_SEGS_LOOP:
	for idx, patternSeg := range patternSegs {
		switch {
		case "*" == patternSegs[idx]:
			matchs = append(matchs, segs[idx])
			continue
		case "**" == patternSegs[idx] && idx == len(patternSegs)-1:
			matchs = append(matchs, strings.Join(segs[idx:], "/"))
			return true, matchs
		case "~~" == patternSegs[idx] && idx == len(patternSegs)-1:
			return true, matchs
		case strings.Contains(patternSeg, "|"):
			pSegs := strings.Split(patternSeg, "|")
			for i := range pSegs {
				if segs[idx] == pSegs[i] {
					matchs = append(matchs, segs[idx])
					continue PATTERN_SEGS_LOOP
				}
			}
			return false, []string{}
		case segs[idx] != patternSegs[idx]:
			return false, []string{}
		default:
		}
	}

	if len(patternSegs) == len(segs) {
		return true, matchs
	}

	return false, []string{}
}

func ArrayMatchSections(strs []string, pattern string) map[string][]string {
	m := map[string][]string{}
	for _, str := range strs {
		if ok, matchs := SectionMatch(str, pattern); ok {
			m[str] = matchs
		}
	}
	return m
}

func MapKeysMatchSections(strs map[string]string, pattern string) map[string][]string {
	m := map[string][]string{}
	for str := range strs {
		if ok, matchs := SectionMatch(str, pattern); ok {
			m[str] = matchs
		}
	}
	return m
}

func MapKeysMatchSectionsUniq(strs map[string]string, pattern string) [][]string {
	ret := [][]string{}

OUTER:
	for str := range strs {
		if ok, matches := SectionMatch(str, pattern); ok {
			for _, res := range ret {
				if StringSliceEqual(res, matches) {
					continue OUTER
				}
			}
			ret = append(ret, matches)
		}
	}
	return ret
}

func StringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for idx, _ := range a {
		if a[idx] != b[idx] {
			return false
		}
	}
	return true
}

// MapKeyMatchFilter filter for map from ucore kvTree or other map that key can split by '/'
func MapKeyMatchFilter(strs map[string]string, pattern string) map[string]string {
	m := map[string]string{}
	for k, v := range strs {
		if ok, _ := SectionMatch(k, pattern); ok {
			m[k] = v
		}
	}
	return m
}

func StringMapEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, _ := range a {
		if a[k] != b[k] {
			return false
		}
	}
	return true
}

func RandStr(prefix, srcLetter string, n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = srcLetter[rand.Intn(len(srcLetter))]
	}
	return prefix + string(b)
}

func MysqlVersionLess(v1, v2 string) bool {
	v1Slice := strings.Split(v1, ".")
	v2Slice := strings.Split(v2, ".")
	for index, _ := range v1Slice {
		if len(v2Slice) < index+1 {
			return false
		}
		val1, _ := strconv.Atoi(v1Slice[index])
		val2, _ := strconv.Atoi(v2Slice[index])
		if val1 < val2 {
			return true
		}
		if val1 > val2 {
			return false
		}
	}
	return true
}
