package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unsafe"

	"github.com/bwmarrin/snowflake"
)

// base64 encoding string to decode string

func DecodeString(base64Str string) (string, error) {
	sDec, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&sDec)), nil
}

func Md5String(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	md5Data := md5.Sum([]byte(nil))
	return hex.EncodeToString(md5Data)
}

func HasPrefix(s, prefix string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.HasPrefix(s, prefix)
	}
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
}

func HasSuffix(s, suffix string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.HasSuffix(s, suffix)
	}
	return strings.HasSuffix(strings.ToLower(s), strings.ToLower(suffix))
}

func GetDuplicate(c []string) []string {
	d := []string{}
	for i, v1 := range c {
		for j, v2 := range c {
			if i >= j {
				continue
			}
			if v1 == v2 {
				d = append(d, v1)
			}
		}
	}
	return RemoveDuplicate(d)
}

func RemoveDuplicate(c []string) []string {
	var tmpMap = map[string]struct{}{}
	var result = []string{}
	for _, v := range c {
		beforeLen := len(tmpMap)
		tmpMap[v] = struct{}{}
		AfterLen := len(tmpMap)
		if beforeLen != AfterLen {
			result = append(result, v)
		}
	}
	return result
}

func RemoveDuplicatePtrUint64(c []*uint64) []*uint64 {
	var tmpMap = map[uint64]struct{}{}
	var result = []*uint64{}
	for _, v := range c {
		beforeLen := len(tmpMap)
		tmpMap[*v] = struct{}{}
		AfterLen := len(tmpMap)
		if beforeLen != AfterLen {
			result = append(result, v)
		}
	}
	return result
}

func RemoveDuplicateUint(c []uint) []uint {
	var tmpMap = map[uint]struct{}{}
	var result = []uint{}
	for _, v := range c {
		beforeLen := len(tmpMap)
		tmpMap[v] = struct{}{}
		AfterLen := len(tmpMap)
		if beforeLen != AfterLen {
			result = append(result, v)
		}
	}
	return result
}

// Round rounds the argument f to dec decimal places.
func Round(f float64, dec int) float64 {
	shift := math.Pow10(dec)
	tmp := f * shift
	if math.IsInf(tmp, 0) {
		return f
	}

	result := math.RoundToEven(tmp) / shift
	if math.IsNaN(result) {
		return 0
	}
	return result
}

func AddDelTag(delTime *time.Time, target string) string {
	if delTime != nil {
		return target + "[x]"
	}
	return target
}

// sep example: ", "
func JoinUintSliceToString(s []uint, sep string) string {
	if len(s) == 0 {
		return ""
	}
	strSlice := make([]string, len(s))
	for i := range s {
		strSlice[i] = strconv.Itoa(int(s[i]))
	}

	return strings.Join(strSlice, sep)
}

// If there are no quotation marks (', ", `) at the beginning and end of the string, the string will be wrapped with "`"
// Need to be wary of the presence of "`" in the string
// do nothing if s is an empty string
func SupplementalQuotationMarks(s string) string {
	if s == "" {
		return ""
	}
	end := len(s) - 1
	if s[0] != s[end] {
		return fmt.Sprintf("`%s`", s)
	}
	if string(s[0]) != "'" && s[0] != '"' && s[0] != '`' {
		return fmt.Sprintf("`%s`", s)
	}
	return s
}

func NvlString(param *string) string {
	if param != nil {
		return *param
	}
	return ""
}

// IsUpperAndLowerLetterMixed
// return true if the string contains both uppercase and lowercase letters
func IsUpperAndLowerLetterMixed(s string) bool {
	if len(s) == 1 {
		return false
	}

	var isUpper bool
	var once sync.Once
	for _, v := range s {
		if !unicode.IsLetter(v) {
			continue
		}
		once.Do(func() {
			isUpper = unicode.IsUpper(v)
		})
		if unicode.IsUpper(v) != isUpper {
			return true
		}
	}

	return false
}

func StringsContains(array []string, ele string) bool {
	for _, a := range array {
		if ele == a {
			return true
		}
	}
	return false
}

var defaultNodeNo int64 = 1
var node *snowflake.Node

// InitSnowflake initiate Snowflake node singleton.
func InitSnowflake(nodeNo int64) error {
	// Create snowflake node
	n, err := snowflake.NewNode(nodeNo)
	if err != nil {
		return err
	}
	// Set node
	node = n
	return nil
}

// GenUid genUid为生成随机uid
func GenUid() (string, error) {
	if node == nil {
		if err := InitSnowflake(defaultNodeNo); err != nil {
			return "", err
		}
	}
	return node.Generate().String(), nil
}

type LowerCaseMap map[string] /*lower case string*/ struct{}

func (l LowerCaseMap) Add(key string) {
	if key == "" {
		return
	}
	l[strings.ToLower(key)] = struct{}{}
}

func (l LowerCaseMap) Exist(key string) bool {
	if key == "" {
		return false
	}
	_, ok := l[strings.ToLower(key)]
	return ok
}

func (l LowerCaseMap) Delete(key string) {
	if key == "" {
		return
	}
	delete(l, strings.ToLower(key))
}

func IsClosed(ch <-chan struct{}) bool {
	if ch == nil {
		return true
	}
	select {
	case _, ok := <-ch:
		if !ok {
			return true
		}
	default:
	}

	return false
}

func TryClose(ch chan struct{}) {
	if !IsClosed(ch) {
		close(ch)
	}
}

// 对比两个float64中更大的并返回
func MaxFloat64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// 计算float64变量的增量平均值
func IncrementalAverageFloat64(oldAverage, newValue float64, oldCount, newCount int) float64 {
	return (oldAverage*float64(oldCount) + newValue) / (float64(oldCount) + float64(newCount))
}

// 判断字符串是否是Git Http URL
func IsGitHttpURL(s string) bool {
	u, err := url.Parse(s)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if !strings.HasSuffix(u.Path, ".git") {
		return false
	}
	return true
}

func IsPrefixSubStrArray(arr []string, prefix []string) bool {
	if len(prefix) > len(arr) {
		return false
	}

	for i := 0; i < len(prefix); i++ {
		if arr[i] != prefix[i] {
			return false
		}
	}

	return true
}

// 全模糊匹配字符串，并且对大小写不敏感
func FullFuzzySearchRegexp(str string) *regexp.Regexp {
	return regexp.MustCompile(`^.*(?i)` + regexp.QuoteMeta(str) + `.*$`)
}
