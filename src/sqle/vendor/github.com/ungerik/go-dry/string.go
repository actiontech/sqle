package dry

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// StringMarshalJSON marshals data to an indented string.
func StringMarshalJSON(data interface{}, indent string) string {
	buffer, err := json.MarshalIndent(data, "", indent)
	if err != nil {
		return ""
	}
	return string(buffer)
}

func StringListContains(l []string, s string) bool {
	for i := range l {
		if l[i] == s {
			return true
		}
	}
	return false
}

func StringListContainsCaseInsensitive(l []string, s string) bool {
	s = strings.ToLower(s)
	for i := range l {
		if strings.ToLower(l[i]) == s {
			return true
		}
	}
	return false
}

func StringPrettifyJSON(compactJSON string) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, []byte(compactJSON), "", "\t")
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

func StringEscapeJSON(jsonString string) string {
	jsonString = strings.Replace(jsonString, `\`, `\\`, -1)
	jsonString = strings.Replace(jsonString, `"`, `\"`, -1)
	return jsonString
}

// StringStripHTMLTags strips HTML/XML tags from text.
func StringStripHTMLTags(text string) (plainText string) {
	var buf *bytes.Buffer
	tagClose := -1
	tagStart := -1
	for i, char := range text {
		if char == '<' {
			if buf == nil {
				buf = bytes.NewBufferString(text)
				buf.Reset()
			}
			buf.WriteString(text[tagClose+1 : i])
			tagStart = i
		} else if char == '>' && tagStart != -1 {
			tagClose = i
			tagStart = -1
		}
	}
	if buf == nil {
		return text
	}
	buf.WriteString(text[tagClose+1:])
	return buf.String()
}

// StringReplaceHTMLTags replaces HTML/XML tags from text with replacement.
func StringReplaceHTMLTags(text, replacement string) (plainText string) {
	var buf *bytes.Buffer
	tagClose := -1
	tagStart := -1
	for i, char := range text {
		if char == '<' {
			if buf == nil {
				buf = bytes.NewBufferString(text)
				buf.Reset()
			}
			buf.WriteString(text[tagClose+1 : i])
			tagStart = i
		} else if char == '>' && tagStart != -1 {
			buf.WriteString(replacement)
			tagClose = i
			tagStart = -1
		}
	}
	if buf == nil {
		return text
	}
	buf.WriteString(text[tagClose+1:])
	return buf.String()
}

// StringMD5Hex returns the hex encoded MD5 hash of data
func StringMD5Hex(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// StringSHA1Base64 returns the base64 encoded SHA1 hash of data
func StringSHA1Base64(data string) string {
	hash := sha1.Sum([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func StringAddURLParam(url, name, value string) string {
	var separator string
	if strings.IndexRune(url, '?') == -1 {
		separator = "?"
	} else {
		separator = "&"
	}
	return url + separator + name + "=" + value
}

func StringConvertTime(timeString, formatIn, formatOut string) (resultTime string, err error) {
	if timeString == "" {
		return "", nil
	}
	t, err := time.Parse(formatIn, timeString)
	if err != nil {
		return "", err
	}
	return t.Format(formatOut), nil
}

func StringCSV(records [][]string) string {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.WriteAll(records)
	if err != nil {
		return ""
	}
	return buf.String()
}

func StringToInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func StringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func StringToBool(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

func StringInSlice(s string, slice []string) bool {
	for i := range slice {
		if slice[i] == s {
			return true
		}
	}
	return false
}

// StringJoinFormat formats every value in values with format
// and joins the result with sep as separator.
// values must be a slice of a formatable type
func StringJoinFormat(format string, values interface{}, sep string) string {
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice {
		panic("values is not a slice")
	}
	var buffer bytes.Buffer
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(fmt.Sprintf(format, v.Index(i).Interface()))
	}
	return buffer.String()
}

// StringJoin formats every value in values according to its default formatting
// and joins the result with sep as separator.
// values must be a slice of a formatable type
func StringJoin(values interface{}, sep string) string {
	v := reflect.ValueOf(values)
	if v.Kind() != reflect.Slice {
		panic("values is not a slice")
	}
	var buffer bytes.Buffer
	for i := 0; i < v.Len(); i++ {
		if i > 0 {
			buffer.WriteString(sep)
		}
		buffer.WriteString(fmt.Sprint(v.Index(i).Interface()))
	}
	return buffer.String()
}

func StringFormatBigInt(mem uint64) string {
	switch {
	case mem >= 10e12:
		return fmt.Sprintf("%dT", mem/1e12)
	case mem >= 1e12:
		return strings.TrimSuffix(fmt.Sprintf("%.1fT", float64(mem)/1e12), ".0")

	case mem >= 10e9:
		return fmt.Sprintf("%dG", mem/1e9)
	case mem >= 1e9:
		return strings.TrimSuffix(fmt.Sprintf("%.1fG", float64(mem)/1e9), ".0")

	case mem >= 10e6:
		return fmt.Sprintf("%dM", mem/1e6)
	case mem >= 1e6:
		return strings.TrimSuffix(fmt.Sprintf("%.1fM", float64(mem)/1e6), ".0")

	case mem >= 10e3:
		return fmt.Sprintf("%dk", mem/1e3)
	case mem >= 1e3:
		return strings.TrimSuffix(fmt.Sprintf("%.1fk", float64(mem)/1e3), ".0")
	}
	return fmt.Sprintf("%d", mem)
}

func StringFormatMemory(mem uint64) string {
	return StringFormatBigInt(mem) + "B"
}

func StringReplaceMulti(str string, fromTo ...string) string {
	if len(fromTo)%2 != 0 {
		panic("Need even number of fromTo arguments")
	}
	for i := 0; i < len(fromTo); i += 2 {
		str = strings.Replace(str, fromTo[i], fromTo[i+1], -1)
	}
	return str
}

func StringToUpperCamelCase(str string) string {
	var buf bytes.Buffer
	var last byte = '_'
	for _, c := range []byte(str) {
		if c != '_' {
			if last == '_' {
				c = byte(unicode.ToUpper(rune(c)))
			} else {
				c = byte(unicode.ToLower(rune(c)))
			}
			buf.WriteByte(c)
		}
		last = c
	}
	return buf.String()
}

func StringToLowerCamelCase(str string) string {
	var buf bytes.Buffer
	var last byte
	for _, c := range []byte(str) {
		if c != '_' {
			if last == '_' {
				c = byte(unicode.ToUpper(rune(c)))
			} else {
				c = byte(unicode.ToLower(rune(c)))
			}
			buf.WriteByte(c)
		}
		last = c
	}
	return buf.String()
}

func StringMapSortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func StringMapGroupedNumberPostfixSortedKeys(m map[string]string) []string {
	keys := make(StringGroupedNumberPostfixSorter, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Sort(keys)
	return keys
}

func StringMapGroupedNumberPostfixSortedValues(m map[string]string) []string {
	values := make(StringGroupedNumberPostfixSorter, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	sort.Sort(values)
	return values
}

func StringEndsWithNumber(s string) bool {
	if s == "" {
		return false
	}
	c := s[len(s)-1]
	return c >= '0' && c <= '9'
}

func StringSplitNumberPostfix(s string) (base, number string) {
	if s == "" {
		return "", ""
	}
	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]
		if c < '0' || c > '9' {
			if i == len(s)-1 {
				return s, ""
			}
			return s[:i+1], s[i+1:]
		}
	}
	return "", s
}

func StringSplitOnce(s, sep string) (pre, post string) {
	parts := strings.SplitN(s, sep, 1)
	if len(parts) == 2 {
		return parts[0], parts[1]
	} else {
		return parts[0], ""
	}
}

func StringSplitOnceChar(s string, sep byte) (pre, post string) {
	i := strings.IndexByte(s, sep)
	if i == -1 {
		return s, ""
	}
	return s[:i], s[i+1:]
}

func StringSplitOnceRune(s string, sep rune) (pre, post string) {
	sepIndex := -1
	postSepIndex := -1
	for i, c := range s {
		if sepIndex != -1 {
			postSepIndex = i
			break // we got the index after the sep rune
		}
		if c == sep {
			sepIndex = i
			// continue to get index after the current UTF8 rune
		}
	}
	if sepIndex == -1 {
		return s, ""
	}
	return s[:sepIndex], s[postSepIndex:]
}

type StringGroupedNumberPostfixSorter []string

// Len is the number of elements in the collection.
func (s StringGroupedNumberPostfixSorter) Len() int {
	return len(s)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (s StringGroupedNumberPostfixSorter) Less(i, j int) bool {
	bi, ni := StringSplitNumberPostfix(s[i])
	bj, nj := StringSplitNumberPostfix(s[j])

	if bi == bj {
		if len(ni) == len(nj) {
			inti, _ := strconv.Atoi(ni)
			intj, _ := strconv.Atoi(nj)
			return inti < intj
		} else {
			return len(ni) < len(nj)
		}
	}

	return bi < bj
}

// Swap swaps the elements with indexes i and j.
func (s StringGroupedNumberPostfixSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Map a function on each element of a slice of strings.
func StringMap(f func(string) string, data []string) []string {
	size := len(data)
	result := make([]string, size, size)
	for i := 0; i < size; i++ {
		result[i] = f(data[i])
	}
	return result
}

// Filter out all strings where the function does not return true.
func StringFilter(f func(string) bool, data []string) []string {
	result := make([]string, 0, 0)
	for _, element := range data {
		if f(element) {
			result = append(result, element)
		}
	}
	return result
}

// StringFindBetween returns the string between the first occurrences of the tokens start and stop.
// The remainder of the string after the stop token will be returned if found.
// If the tokens couldn't be found, then the whole string will be returned as remainder.
func StringFindBetween(s, start, stop string) (between, remainder string, found bool) {
	begin := strings.Index(s, start)
	if begin == -1 {
		return "", s, false
	}
	between = s[begin+len(start):]
	end := strings.Index(between, stop)
	if end == -1 {
		return "", s, false
	}
	return between[:end], s[begin+len(start)+end+len(stop):], true
}

// StringFind returns in found if token has been found in s,
// and returns the remaining string afte token in remainder.
// The whole string s will be returned if found is false.
func StringFind(s, token string) (remainder string, found bool) {
	i := strings.Index(s, token)
	if i == -1 {
		return s, false
	}
	return s[i+len(token):], true
}
