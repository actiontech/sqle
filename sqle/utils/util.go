package utils

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
	"unsafe"

	sqleErrors "github.com/actiontech/sqle/sqle/errors"
	goGit "github.com/go-git/go-git/v5"
	goGitTransport "github.com/go-git/go-git/v5/plumbing/transport/http"
	"golang.org/x/crypto/ssh"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/bwmarrin/snowflake"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

func MergeAndDeduplicateSort(arr1, arr2 []string) []string {
	// 合并两个数组
	merged := append(arr1, arr2...)

	// 如果合并后的数组是空的，直接返回空切片
	if len(merged) == 0 {
		return []string{}
	}

	// 使用map去重
	seen := make(map[string]struct{})
	// 预分配足够的空间，避免多次内存分配
	result := make([]string, 0, len(merged))

	for _, str := range merged {
		if _, exists := seen[str]; !exists {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}

	// 排序
	sort.Strings(result)

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

// 全模糊匹配字符串，对大小写不敏感，匹配多行，且防止正则注入
func FullFuzzySearchRegexp(str string) *regexp.Regexp {
	/*
		1. (?is)是一个正则表达式修饰符,其中：
			i表示忽略大小写(case-insensitive)
			s表示开启单行模式，开启后.可以匹配换行符，让整个字符串作为一行
		2. ^.*匹配字符串的开头,其中：
			^表示起始位置,
			.表示匹配任何字符(除了换行符)
			*表示匹配前面的模式零次或多次
		3. .*$匹配字符串的结尾,其中：
			$表示结束位置
	*/
	return regexp.MustCompile(`(?is)^.*` + regexp.QuoteMeta(str) + `.*$`)
}

var ErrUnknownEncoding = errors.New("unknown encoding")

var encodings = []transform.Transformer{
	simplifiedchinese.GBK.NewDecoder(),
}

func ConvertToUtf8(in []byte) ([]byte, error) {
	if utf8.Valid(in) {
		return in, nil
	}

	for _, enc := range encodings {
		reader := transform.NewReader(bytes.NewReader(in), enc)
		out, err := io.ReadAll(reader)
		if err == nil {
			return out, nil
		}
		log.NewEntry().Errorf("ConvertToUtf8 failed: %v", err)
	}

	return nil, ErrUnknownEncoding
}

// 生成随机字符串，生成长度是halfLength的两倍
func GenerateRandomString(halfLength int) string {
	bytes := make([]byte, halfLength)
	//nolint:errcheck
	rand.Read(bytes)
	return fmt.Sprintf("%x", bytes)
}

// TruncateStringByRunes 按字符数截取字符串
func TruncateStringByRunes(s string, maxRunes uint) string {
	// 字节数不大于 maxRunes ，那字符数肯定不大于 maxRunes
	if uint(len(s)) <= maxRunes {
		return s
	}

	// UTF-8一个字符的字节数是不确定的，如：s="a一b二c"，汉字为多字节字符，len(s)=9
	//    s的hexdump结果：
	//    00000000  61 e4 b8 80 62 e4 ba 8c 63                       |a...b...c|
	//
	//    当想截取头两个字符：“a一”，即 maxRunes 为2时，
	//    直接返回s[:maxRunes]的话得到是：“61 e4”这两个字节组成的字符串，并非“a一”，“a一”是“61 e4 b8 80”这四个字节，此时应取s[:4]
	//
	//    为得到s[:4]中4这个索引，可以“range s”：逐个rune遍历s，i为每个rune起始的字节索引
	//    i依次为　0　　1　　4　　5　　8
	//          　^ａ　^一　^ｂ　^二　^ｃ
	//    遍历 maxRunes (2)次后，i为下一个字符(b)的起始索引，即4，此时s[:i]就是要截取的头两个字符“a一”
	var runesCount uint
	for i := range s {
		if runesCount == maxRunes {
			// 达到截取的字符数了，将字符截取至此时rune的字节索引
			return s[:i]
		}
		// 未达到要截取的字符数，继续获取下一个rune
		runesCount++
	}
	// 字符串字符数不足 maxRunes
	return s
}

const excelCellMaxRunes = 32766

// TruncateAndMarkForExcelCell 对超长字符串进行截取，以符合Excel类工具对单元格字符数上限的限制
func TruncateAndMarkForExcelCell(s string) string {
	truncated := TruncateStringByRunes(s, excelCellMaxRunes-4)
	if truncated != s {
		// 截取了的话，做标记
		return truncated + " ..."
	}
	return s
}

func IntersectionStringSlice(slice1, slice2 []string) []string {
	// 用 map 来存储第一个切片的元素
	elemMap := make(map[string]bool)
	for _, v := range slice1 {
		elemMap[v] = true
	}

	// 遍历第二个切片，找到交集
	var intersection []string
	for _, v := range slice2 {
		if elemMap[v] {
			intersection = append(intersection, v)
			// 删除元素以防重复添加
			delete(elemMap, v)
		}
	}
	return intersection
}

func CloneGitRepository(ctx context.Context, url, username, password string) (repository *goGit.Repository, directory string, cleanup func() error, err error) {

	// http 协议
	// git 协议
	// ssh 协议

	if !IsGitHttpURL(url) {
		return nil, "", nil, sqleErrors.New(sqleErrors.DataInvalid, fmt.Errorf("url is not a git url"))
	}
	// 创建一个临时目录用于存放克隆的仓库
	directory, err = os.MkdirTemp("./", "git-repo-")
	if err != nil {
		return nil, "", nil, err
	}
	// 定义清理函数，用于删除临时目录
	cleanup = func() error {
		return os.RemoveAll(directory)
	}
	cloneOpts := &goGit.CloneOptions{
		URL: url,
	}
	// http协议下：
	//   1. 账号密码登录
	//       username/password
	//       不需要密码的方式
	//   2. token 方式
	//       gitlab：
	//       github:

	// ssh 协议
	//     	前置条件：
	// 			1. 生成密钥【用户手动执行】，用什么用户生成？
	// 		      -	mkdir -p ../keys/ (700权限)
	// 			  - sudo -u actiontech-universe ssh-keygen -t rsa -b 4096 -f ../keys/id_rsa -N ""
	// 			2. 查看公钥【用户手动执行】
	// 			   查看公钥匙内容
	// 			2. 仓库配置密钥
	//             TODO 目前不支持该步骤，只能用户手动执行
	//
	// git协议
	//     不需要校验权限
	if username != "" {
		cloneOpts.Auth = &goGitTransport.BasicAuth{
			Username: username,
			Password: password,
		}
	}
	repository, err = goGit.PlainCloneContext(ctx, directory, false, cloneOpts)
	if err != nil {
		err = cleanup()
		return nil, directory, nil, err
	}
	return repository, directory, cleanup, nil
}

func GeneratePublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) (string, error) {
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}
	return string(ssh.MarshalAuthorizedKey(publicKey)), nil
}

func GenerateSSHKeyPair() (privateKeyStr, publicKeyStr string, err error) {
	// 1. 生成 4096-bit RSA 私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return "", "", err
	}

	// 2. 编码私钥为 PEM 格式，与ssh-keygen -N "" 生成的格式保持一致（无密码保护）
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privatePEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	// 3. 生成 SSH 公钥，格式：ssh-rsa AAAA...
	publicKeyStr, err = GeneratePublicKeyFromPrivateKey(privateKey)
	if err != nil {
		return "", "", err
	}

	return string(privatePEM), publicKeyStr, nil
}
