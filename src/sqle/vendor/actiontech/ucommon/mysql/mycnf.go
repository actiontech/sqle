package mysql

import (
	"actiontech/ucommon/conf"
	"fmt"
	dry "github.com/ungerik/go-dry"
	"io/ioutil"
	"math"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

/*
	when update my.cnf, Mycnf should keep comment/format unchanged, which goconf doesn't satisfied
	Mycnf use text-replacement for this purpose
	precondition: option key is unique, even between different sections
*/

type Mycnf struct {
	content string
}

func LoadMycnf(mycnfPath string) (*Mycnf, error) {
	bs, err := ioutil.ReadFile(mycnfPath)
	if nil != err {
		return nil, fmt.Errorf("read my.cnf %v error: %v", mycnfPath, err)
	}
	return NewMycnf(string(bs))
}

func NewMycnf(mycnf string) (*Mycnf, error) {
	t := &Mycnf{content: mycnf}
	if _, err := t.reader(); nil != err {
		return nil, err
	}
	return t, nil
}

func (t *Mycnf) reader() (*conf.ConfigFile, error) {
	return conf.ReadConfigBytes([]byte(t.content))
}

func (t *Mycnf) GetMysqldInt(key string) int {
	r, _ := t.reader()
	val, _ := r.GetInt("mysqld", key)
	return val
}

func (t *Mycnf) GetUniverseInt(key string) int {
	r, _ := t.reader()
	val, _ := r.GetInt("universe", key)
	return val
}

func (t *Mycnf) GetMysqldOct(key string) uint32 {
	r, _ := t.reader()
	str, _ := r.GetString("mysqld", key)
	if "" == str {
		return 0
	}
	ret, _ := strconv.ParseUint(str, 8, 32)
	return uint32(ret)
}

func (t *Mycnf) GetUniverseOct(key string) uint32 {
	r, _ := t.reader()
	str, _ := r.GetString("universe", key)
	if "" == str {
		return 0
	}
	ret, _ := strconv.ParseUint(str, 8, 32)
	return uint32(ret)
}

func (t *Mycnf) GetMysqldString(key string) string {
	r, _ := t.reader()
	key = t.GetKeyInMycnf(key)
	if "" == key {
		return ""
	}
	val, _ := r.GetString("mysqld", key)
	return val
}

func (t *Mycnf) GetUniverseString(key string) string {
	r, _ := t.reader()
	val, _ := r.GetString("universe", key)
	return val
}

func (t *Mycnf) GetRequiredMysqldString(key string) (string, error) {
	val := t.GetMysqldString(key)
	if "" == val {
		return "", fmt.Errorf("option %v is required", key)
	}
	return val, nil
}

func (t *Mycnf) GetKeyInMycnf(option string) string {
	r, _ := t.reader()
	keys, _ := r.GetOptions("mysqld")
	option = strings.Replace(option, "-", "_", -1)
	for _, key := range keys {
		if strings.Replace(key, "-", "_", -1) == option {
			return key
		}
	}
	return ""
}

func (t *Mycnf) GetAbsPathOrRelPathOnBaseDir(option string, baseOption string) (string, error) {
	val := t.GetMysqldString(option)
	if strings.HasPrefix(val, "/") {
		return val, nil
	}
	base, err := t.GetRequiredMysqldString(baseOption)
	if nil != err {
		return "", fmt.Errorf("base option %v is required", baseOption)
	}
	return filepath.Join(base, val), nil
}

var (
	TRUE_VALUES  = []string{"TRUE", "true", "True", "1", "ON", "on", "On"}
	FALSE_VALUES = []string{"FALSE", "false", "False", "0", "OFF", "off", "Off"}
)

func (t *Mycnf) IsValueDefaultOrEq(option string, expect string) bool {
	val := t.GetMysqldString(option)
	if "" == val {
		return true
	}
	val = strings.Trim(val, `"'`)
	expect = strings.Trim(expect, `"'`)

	if dry.StringInSlice(val, TRUE_VALUES) && dry.StringInSlice(expect, TRUE_VALUES) {
		return true
	}
	if dry.StringInSlice(val, FALSE_VALUES) && dry.StringInSlice(expect, FALSE_VALUES) {
		return true
	}

	//trim ending "/" in case of file path
	val = filepath.Clean(strings.ToLower(val))
	expect = filepath.Clean(strings.ToLower(expect))
	if val == expect {
		return true
	}
	if t.isSizeEq(val, expect) {
		return true
	}
	if t.isFloatEq(val, expect) {
		return true
	}
	
	// https://dev.mysql.com/doc/refman/8.0/en/replication-options-slave.html#option_mysqld_slave-skip-errors
	// slave-skip-errors's value ddl_exist_errors equals to 1007,1008,1050,1051,1054,1060,1061,1068,1091,1146
	if option == "slave_skip_errors" && val == "ddl_exist_errors" {
		val = "1007,1008,1050,1051,1054,1060,1061,1068,1091,1146"
	}
	if option == "slave_skip_errors" && expect == "ddl_exist_errors" {
		expect = "1007,1008,1050,1051,1054,1060,1061,1068,1091,1146"
	}
	arrayTypeOptions := []string{
		"sql_mode", "disabled_storage_engines", "log_output", "slave_skip_errors",
	}
	for _, op := range arrayTypeOptions {
		if option == op && t.isArrayEq(val, expect) {
			return true
		}
	}
	return false
}

func (t *Mycnf) isSizeEq(a, b string) bool {
	numberRegexp := regexp.MustCompile("^(\\d+)([g|G|m|M|k|K]?)$")
	if !numberRegexp.MatchString(a) || !numberRegexp.MatchString(b) {
		return false
	}
	aMatches := numberRegexp.FindStringSubmatch(a)
	bMatches := numberRegexp.FindStringSubmatch(b)
	aNum := t.sizeToNumber(aMatches[1], aMatches[2])
	bNum := t.sizeToNumber(bMatches[1], bMatches[2])
	return aNum == bNum
}

func (t *Mycnf) isFloatEq(a, b string) bool {
	aNum, err := strconv.ParseFloat(a, 32)
	if nil != err {
		return false
	}
	bNum, err := strconv.ParseFloat(b, 32)
	if nil != err {
		return false
	}
	return math.Abs(aNum-bNum) <= 1e-5
}

func (t *Mycnf) isArrayEq(a, b string) bool {
	return t.isSubArrayOf(a, b) && t.isSubArrayOf(b, a)
}

func (t *Mycnf) isSubArrayOf(a, b string) bool {
	aa := strings.Split(a, ",")
	bb := strings.Split(b, ",")

	for _, a := range aa {
		found := false
		for _, b := range bb {
			if a == b {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func (t *Mycnf) sizeToNumber(num string, unit string) uint64 {
	n, _ := strconv.ParseUint(num, 10, 64)
	switch unit {
	case "k":
		n *= 1024
	case "K":
		n *= 1024
	case "m":
		n *= 1024 * 1024
	case "M":
		n *= 1024 * 1024
	case "g":
		n *= 1024 * 1024 * 1024
	case "G":
		n *= 1024 * 1024 * 1024
	}
	return n
}

func (t *Mycnf) SetMysqldOption(key, newVal string) {
	keyInMycnf := t.GetKeyInMycnf(key)
	if "" == keyInMycnf {
		//new option
		t.content = strings.Replace(t.content, "[mysqld]",
			fmt.Sprintf("[mysqld]\n%v = %v", key, newVal), 1)
	} else {
		//existing option
		t.content = regexp.MustCompile("(?m)^\\s*"+strings.Replace(keyInMycnf, "-", "\\-", -1)+"\\s*=.*$").
			ReplaceAllString(t.content, fmt.Sprintf("%v = %v", keyInMycnf, newVal))
	}
}

func (t *Mycnf) SetUniverseOption(key, newVal string) {
	r, _ := t.reader()
	if !r.HasSection("universe") {
		t.content = "[universe]\n\n" + t.content
	}

	if "" == newVal {
		//delete option
		t.content = regexp.MustCompile("(?m)^\\s*"+strings.Replace(key, "-", "\\-", -1)+"\\s*=.*\n").
			ReplaceAllString(t.content, "")
	} else if !r.HasOption("universe", key) {
		//new option
		t.content = strings.Replace(t.content, "[universe]",
			fmt.Sprintf("[universe]\n%v = %v", key, newVal), 1)
	} else {
		//existing option
		t.content = regexp.MustCompile("(?m)^\\s*"+strings.Replace(key, "-", "\\-", -1)+"\\s*=.*$").
			ReplaceAllString(t.content, fmt.Sprintf("%v = %v", key, newVal))
	}
}

func (t *Mycnf) SetOption(section, key, newVal string) {
	r, _ := t.reader()
	if !r.HasSection(section) {
		t.content = fmt.Sprintf("[%v]\n\n", section) + t.content
	}

	if "" == newVal {
		//delete option
		t.content = regexp.MustCompile("(?m)^\\s*"+strings.Replace(key, "-", "\\-", -1)+"\\s*=.*\n").
			ReplaceAllString(t.content, "")
	} else if !r.HasOption(section, key) {
		//new option
		t.content = strings.Replace(t.content, fmt.Sprintf("[%v]", section),
			fmt.Sprintf("[%v]\n%v = %v", section, key, newVal), 1)
	} else {
		//existing option
		t.content = regexp.MustCompile("(?m)^\\s*"+strings.Replace(key, "-", "\\-", -1)+"\\s*=.*$").
			ReplaceAllString(t.content, fmt.Sprintf("%v = %v", key, newVal))
	}
}

func (t *Mycnf) GetUniverseOptions() []string {
	r, _ := t.reader()
	options, _ := r.GetOptionsWithoutDefaultSection("universe")
	return options
}

func (t *Mycnf) GetMysqldOptions() []string {
	r, _ := t.reader()
	opts, _ := r.GetOptions("mysqld")
	return opts
}

func (t *Mycnf) GetString(section string, option string) string {
	r, _ := t.reader()
	val, _ := r.GetString(section, option)
	return val
}

func (t *Mycnf) String() string {
	return t.content
}

// this may be need lock
func (t *Mycnf) DelOption(section, option string) {
	r, _ := t.reader()
	if r.HasOption(section, option) {
		r.RemoveOption(section, option)
		t.content = string(r.WriteConfigBytes(""))
	}
}

func (t *Mycnf) DelMysqlOption(option string) {
	t.DelOption("mysqld", option)
}
