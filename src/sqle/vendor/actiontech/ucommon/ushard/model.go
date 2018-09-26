package ushard

import (
	"regexp"
	"fmt"
	"strings"
	"time"
	"math"
	"strconv"
	"encoding/xml"
)

/*
	Config properties' json name follows MyCat standard, which violates with json common sense
 */

type Config struct {
	DataHosts      []*DataHost `json:"dataHosts"`
	DataHostsError string      `json:"dataHosts_error"`
	DataNodes      []*DataNode `json:"dataNodes"`
	DataNodesError string      `json:"dataNodes_error"`
	Schemas        []*Schema   `json:"schemas"`
	SchemasError   string      `json:"schemas_error"`
	Users          []*User     `json:"users"`
	UsersError     string      `json:"users_error"`
}

func (c *Config) Validate(adminUser string) bool {
	pass := true
	if nil != c.DataHosts {
		for _, host := range c.DataHosts {
			if !host.validate(c) {
				pass = false
			}
		}
	}
	if nil != c.DataNodes {
		for _, node := range c.DataNodes {
			if !node.validate(c) {
				pass = false
			}
		}
	}
	if nil != c.Schemas {
		for _, schema := range c.Schemas {
			if !schema.validate(c) {
				pass = false
			}
		}
	}
	if nil != c.Users {
		for _, user := range c.Users {
			if !user.validate(c, adminUser) {
				pass = false
			}
		}
	}
	return pass
}

var namePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]*$`)
var nameListPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]*(,[a-zA-Z][a-zA-Z0-9_\-]*)*$`)
var urlPattern = regexp.MustCompile(`^[a-zA-Z_0-9.\-]+:\d+`)
//var mycatSpecPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]*(\$\d+-\d+)?$`)
var mycatSpecListPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_\-]*(\$\d+-\d+)?(,[a-zA-Z][a-zA-Z0-9_\-]*(\$\d+-\d+)?)*$`)

func validateMycatSpecList(spec string) string {
	if "" == spec {
		return "should not be empty"
	}
	if !mycatSpecListPattern.MatchString(spec) {
		return "name should match `^[a-zA-Z][a-zA-Z0-9_\\-]*(\\$\\d+-\\d+)?(,[a-zA-Z][a-zA-Z0-9_\\-]*(\\$\\d+-\\d+)?)*$`"
	}
	return ""
}

func validateName(name string) string {
	if "" == name {
		return "should not be empty"
	}
	if !namePattern.MatchString(name) {
		return "name should match `^[a-zA-Z][a-zA-Z0-9_\\-]*$`"
	}
	return ""
}

func validateNameList(nameList string) string {
	if "" == nameList {
		return "should not be empty"
	}
	if !nameListPattern.MatchString(nameList) {
		return "name should match `^[a-zA-Z][a-zA-Z0-9_\\-]*$`"
	}
	return ""
}

func validateUrl(url string) string {
	if "" == url {
		return "should not be empty"
	}
	if !urlPattern.MatchString(url) {
		return "url should match [a-zA-Z_]+:\\d+"
	}
	return ""
}

func validateNotEmpty(a string) string {
	if "" == a {
		return "should not be empty"
	}
	return ""
}

type DataHost struct {
	Name                  string       `json:"name"`
	NameError             string       `json:"name_error"`
	Balance               int          `json:"balance"`
	BalanceError          string       `json:"balance_error"`
	SwitchType            int          `json:"switchType"`
	SwitchTypeError       string       `json:"switchType_error"`
	SlaveThreshold        int          `json:"slaveThreshold"`
	SlaveThresholdError   string       `json:"slaveThreshold_error"`
	TempReadHostAvailable bool         `json:"tempReadHostAvailable"`
	MaxCon                int          `json:"maxCon"`
	MaxConError           string       `json:"maxCon_error"`
	MinCon                int          `json:"minCon"`
	MinConError           string       `json:"minCon_error"`
	WriteHosts            []*WriteHost `json:"writeHosts"`
	WriteHostsError       string       `json:"writeHosts_error"`
}

func (d *DataHost) validate(c *Config) bool {
	pass := true
	d.NameError = validateName(d.Name)
	if "" != d.NameError {
		pass = false
	} else {
		//duplicate check
		for _, host := range c.DataHosts {
			if host != d && host.Name == d.Name {
				d.NameError = "duplicate name"
				pass = false
			}
		}
	}

	if !(d.Balance >= 0 && d.Balance <= 3) {
		d.BalanceError = "should >= 0 and <= 3"
		pass = false
	}
	if !(d.SwitchType == -1 || d.SwitchType == 1 || d.SwitchType == 2) {
		d.SwitchTypeError = "should be -1 or 1 or 2"
		pass = false
	}
	if !(-1 == d.SlaveThreshold || d.SlaveThreshold > 0) {
		d.SlaveThresholdError = "should be -1 or >0"
		pass = false
	}
	if -1 != d.SlaveThreshold && 2 != d.SwitchType {
		d.SlaveThresholdError = "only effect when SwitchType == 2"
		pass = false
	}
	if !(d.MaxCon > 0) {
		d.MaxConError = "should >0"
		pass = false
	}
	if !(d.MinCon > 0) {
		d.MinConError = "should >0"
		pass = false
	}
	if !(d.MaxCon >= d.MinCon) {
		d.MaxConError = "MaxCon should >= MinCon"
		pass = false
	}
	if nil == d.WriteHosts || 0 == len(d.WriteHosts) {
		d.WriteHostsError = "should have WriteHosts"
		pass = false
	} else if len(d.WriteHosts) > 1 {
		d.WriteHostsError = "should only 1 WriteHost"
		pass = false
	} else {
		for _, host := range d.WriteHosts {
			if !host.validate(c) {
				pass = false
			}
		}
	}
	return pass
}

type WriteHost struct {
	Name                            string      `json:"name"`
	Host                            string      `json:"host"`
	HostError                       string      `json:"host_error"`
	Url                             string      `json:"url"`
	UrlError                        string      `json:"url_error"`
	User                            string      `json:"user"`
	UserError                       string      `json:"user_error"`
	Password                        string      `json:"password"`
	PasswordError                   string      `json:"password_error"`
	ReadHosts                       []*ReadHost `json:"readHosts"`
	TakeoverMysqlGroup              string      `json:"takeover_mysql_group"`
	TakeoverMysqlInstance           string      `json:"takeover_mysql_instance"`
	NoWriteHostInTakeoverMysqlGroup bool        `json:"no_write_host_in_takeover_mysql_group"`
}

func (d *WriteHost) ToReadHost() *ReadHost {
	return &ReadHost{
		Name:                            d.Name,
		Host:                            d.Host,
		Url:                             d.Url,
		User:                            d.User,
		Password:                        d.Password,
		TakeoverMysqlGroup:              d.TakeoverMysqlGroup,
		TakeoverMysqlInstance:           d.TakeoverMysqlInstance,
		NoWriteHostInTakeoverMysqlGroup: d.NoWriteHostInTakeoverMysqlGroup,
	}
}

func (d *WriteHost) IsTakeover() bool {
	return "" != d.TakeoverMysqlGroup
}

func (d *WriteHost) validate(c *Config) bool {
	pass := true
	d.HostError = validateName(d.Host)
	d.UrlError = validateUrl(d.Url)
	d.UserError = validateNotEmpty(d.User)
	d.PasswordError = validateNotEmpty(d.Password)
	if "" != d.HostError+d.UrlError+d.UserError+d.PasswordError {
		pass = false
	}
	if nil != d.ReadHosts {
		for _, host := range d.ReadHosts {
			if !host.validate(c) {
				pass = false
			}
		}
	}
	return pass
}

type ReadHost struct {
	Name                            string `json:"name"`
	Host                            string `json:"host"`
	HostError                       string `json:"host_error"`
	Url                             string `json:"url"`
	UrlError                        string `json:"url_error"`
	User                            string `json:"user"`
	UserError                       string `json:"user_error"`
	Password                        string `json:"password"`
	PasswordError                   string `json:"password_error"`
	TakeoverMysqlGroup              string `json:"takeover_mysql_group"`
	TakeoverMysqlInstance           string `json:"takeover_mysql_instance"`
	NoWriteHostInTakeoverMysqlGroup bool   `json:"no_write_host_in_takeover_mysql_group"`
}

func (d *ReadHost) validate(c *Config) bool {
	d.HostError = validateName(d.Host)
	d.UrlError = validateUrl(d.Url)
	d.UserError = validateNotEmpty(d.User)
	d.PasswordError = validateNotEmpty(d.Password)
	return "" == (d.HostError + d.UrlError + d.UserError + d.PasswordError)
}

func (d *ReadHost) ToWriteHost() *WriteHost {
	return &WriteHost{
		Name:                            d.Name,
		Host:                            d.Host,
		Url:                             d.Url,
		User:                            d.User,
		Password:                        d.Password,
		TakeoverMysqlGroup:              d.TakeoverMysqlGroup,
		TakeoverMysqlInstance:           d.TakeoverMysqlInstance,
		NoWriteHostInTakeoverMysqlGroup: d.NoWriteHostInTakeoverMysqlGroup,
	}
}

type DataNode struct {
	Name          string `json:"name"`
	NameError     string `json:"name_error"`
	DataHost      string `json:"dataHost"`
	DataHostError string `json:"dataHost_error"`
	Database      string `json:"database"`
	DatabaseError string `json:"database_error"`
}

func (d *DataNode) validate(c *Config) bool {
	d.NameError = validateName(d.Name)
	d.DataHostError = validateNotEmpty(d.DataHost)
	if "" == d.DataHostError {
		found := false
		for _, host := range c.DataHosts {
			if host.Name == d.DataHost {
				found = true
			}
		}
		if !found {
			d.DataHostError = fmt.Sprintf("no DataHost named \"%v\"", d.DataHost)
		}
	}
	d.DatabaseError = validateNotEmpty(d.Database)
	if "" == d.DatabaseError {
		d.DatabaseError = validateMycatSpecList(d.Database)
	}

	return "" == d.NameError+d.DataHostError+d.DatabaseError
}

type Schema struct {
	Name             string   `json:"name"`
	NameError        string   `json:"name_error"`
	DataNode         string   `json:"dataNode"`
	DataNodeError    string   `json:"dataNode_error"`
	CheckSqlSchema   bool     `json:"checkSQLschema"`
	SqlMaxLimit      int      `json:"sqlMaxLimit"`
	SqlMaxLimitError string   `json:"sqlMaxLimit_error"`
	Tables           []*Table `json:"tables"`
}

func (s *Schema) validate(c *Config) bool {
	s.NameError = validateName(s.Name)
	if "" != s.DataNode {
		s.DataNodeError = validateMycatSpecList(s.DataNode)
	}
	if !(s.SqlMaxLimit >= 0) {
		s.SqlMaxLimitError = "should be >=0"
	}
	pass := "" == s.NameError+s.DataNodeError+s.SqlMaxLimitError
	if nil != s.Tables {
		for _, table := range s.Tables {
			if !table.validate(c) {
				pass = false
			}
		}
	}
	return pass
}

type Table struct {
	Name             string        `json:"name"`
	NameError        string        `json:"name_error"`
	Type             string        `json:"type"`
	TypeError        string        `json:"type_error"`
	PrimaryKey       string        `json:"primaryKey"`
	PrimaryKeyError  string        `json:"primaryKey_error"`
	AutoIncrement    bool          `json:"autoIncrement"`
	NeedAddLimit     bool          `json:"needAddLimit"`
	DataNode         string        `json:"dataNode"`
	DataNodeError    string        `json:"dataNode_error"`
	Rule             *Rule         `json:"rule"`
	RuleError        string        `json:"rule_error"`
	ChildTables      []*ChildTable `json:"childTable"`
	ChildTablesError string        `json:"childTable_error"`
}

func (t *Table) partitionCount() int {
	ret := 0
	for _, seg := range strings.Split(t.DataNode, ",") {
		matches := regexp.MustCompile(`.*\$(\d+)-(\d+)`).FindStringSubmatch(seg)
		if nil == matches {
			ret++
			continue
		}

		b, _ := strconv.Atoi(matches[1])
		e, _ := strconv.Atoi(matches[2])
		ret += e - b + 1
	}
	return ret
}

func (t *Table) validate(c *Config) bool {
	pass := true
	t.NameError = validateName(t.Name)
	if "" != t.NameError {
		pass = false
	}

	if !("default" == t.Type || "global" == t.Type) {
		t.TypeError = "should be 'default' or 'global'"
		pass = false
	}

	t.PrimaryKeyError = validateNameList(t.PrimaryKey)
	if "" != t.PrimaryKeyError {
		pass = false
	}

	if "" != t.DataNode {
		t.DataNodeError = validateMycatSpecList(t.DataNode)
		if "" != t.DataNodeError {
			pass = false
		}
	}

	if "global" == t.Type {
		if nil != t.Rule {
			t.RuleError = "global table should not have rule"
			pass = false
		}
		if nil != t.ChildTables {
			t.ChildTablesError = "global table should not have child tables"
			pass = false
		}
	} else {
		if nil == t.Rule {
			t.RuleError = "should not be null"
			pass = false
		} else {
			if !t.Rule.validate(c) {
				pass = false
			}
		}

		if pass {
			expectCount := t.Rule.partitionCount()
			if -1 != expectCount {
				actualCount := t.partitionCount()
				if actualCount < expectCount {
					t.RuleError = fmt.Sprintf("config has %v data nodes, rule has %v partitions, config < rule", expectCount, actualCount)
					pass = false
				}
			} else {
				//no need to check violation
			}
		}

		if nil != t.ChildTables {
			for _, table := range t.ChildTables {
				if !table.validate(c) {
					pass = false
				}
			}
		}
	}

	return pass
}

type ChildTable struct {
	Name            string        `json:"name"`
	NameError       string        `json:"name_error"`
	PrimaryKey      string        `json:"primaryKey"`
	PrimaryKeyError string        `json:"primaryKey_error"`
	AutoIncrement   bool          `json:"autoIncrement"`
	JoinKey         string        `json:"joinKey"`
	JoinKeyError    string        `json:"joinKey_error"`
	ParentKey       string        `json:"parentKey"`
	ParentKeyError  string        `json:"parentKey_error"`
	ChildTables     []*ChildTable `json:"childTable"`
}

func (t *ChildTable) validate(c *Config) bool {
	pass := true
	t.NameError = validateName(t.Name)
	if "" != t.NameError {
		pass = false
	}

	t.PrimaryKeyError = validateNameList(t.PrimaryKey)
	if "" != t.PrimaryKeyError {
		pass = false
	}

	t.JoinKeyError = validateNameList(t.JoinKey)
	if "" != t.JoinKeyError {
		pass = false
	}

	t.ParentKeyError = validateNameList(t.ParentKey)
	if "" != t.ParentKeyError {
		pass = false
	}

	if nil != t.ChildTables {
		for _, table := range t.ChildTables {
			if !table.validate(c) {
				pass = false
			}
		}
	}

	return pass
}

type Rule struct {
	Columns                      []string                      `json:"columns"`
	ColumnsError                 string                        `json:"columns_error"`
	Algorithm                    string                        `json:"algorithm"`
	AlgorithmError               string                        `json:"algorithm_error"`
	AlgorithmPartitionByMap      *AlgorithmPartitionByMap      `json:"algorithmPartitionByMap"`
	AlgorithmAutoPartitionByLong *AlgorithmAutoPartitionByLong `json:"algorithmAutoPartitionByLong"`
	AlgorithmPartitionByLong     *AlgorithmPartitionByLong     `json:"algorithmPartitionByLong"`
	AlgorithmPartitionByString   *AlgorithmPartitionByString   `json:"algorithmPartitionByString"`
	AlgorithmPartitionByDate     *AlgorithmPartitionByDate     `json:"algorithmPartitionByDate"`
}

func (r *Rule) IsSame(a *Rule) bool {
	if len(r.Columns) != len(a.Columns) {
		return false
	}

	if r.Algorithm != a.Algorithm {
		return false
	}

	{
		if (nil == r.AlgorithmPartitionByMap && nil != a.AlgorithmPartitionByMap ||
			nil != r.AlgorithmPartitionByMap && nil == a.AlgorithmPartitionByMap) {
			return false
		}
		if nil != r.AlgorithmPartitionByMap && !r.AlgorithmPartitionByMap.IsSame(a.AlgorithmPartitionByMap) {
			return false
		}
	}

	{
		if (nil == r.AlgorithmAutoPartitionByLong && nil != a.AlgorithmAutoPartitionByLong ||
			nil != r.AlgorithmAutoPartitionByLong && nil == a.AlgorithmAutoPartitionByLong) {
			return false
		}
		if nil != r.AlgorithmAutoPartitionByLong && !r.AlgorithmAutoPartitionByLong.IsSame(a.AlgorithmAutoPartitionByLong) {
			return false
		}
	}

	{
		if (nil == r.AlgorithmPartitionByLong && nil != a.AlgorithmPartitionByLong ||
			nil != r.AlgorithmPartitionByLong && nil == a.AlgorithmPartitionByLong) {
			return false
		}
		if nil != r.AlgorithmPartitionByLong && !r.AlgorithmPartitionByLong.IsSame(a.AlgorithmPartitionByLong) {
			return false
		}
	}

	{
		if (nil == r.AlgorithmPartitionByString && nil != a.AlgorithmPartitionByString ||
			nil != r.AlgorithmPartitionByString && nil == a.AlgorithmPartitionByString) {
			return false
		}
		if nil != r.AlgorithmPartitionByString && !r.AlgorithmPartitionByString.IsSame(a.AlgorithmPartitionByString) {
			return false
		}
	}

	{
		if (nil == r.AlgorithmPartitionByDate && nil != a.AlgorithmPartitionByDate ||
			nil != r.AlgorithmPartitionByDate && nil == a.AlgorithmPartitionByDate) {
			return false
		}
		if nil != r.AlgorithmPartitionByDate && !r.AlgorithmPartitionByDate.IsSame(a.AlgorithmPartitionByDate) {
			return false
		}
	}

	return true
}

func (r *Rule) partitionCount() int {
	switch r.Algorithm {
	case "algorithmPartitionByMap":
		return r.AlgorithmPartitionByMap.partitionCount()
	case "algorithmAutoPartitionByLong":
		return r.AlgorithmAutoPartitionByLong.partitionCount()
	case "algorithmPartitionByLong":
		return r.AlgorithmPartitionByLong.partitionCount()
	case "algorithmPartitionByString":
		return r.AlgorithmPartitionByString.partitionCount()
	case "algorithmPartitionByDate":
		return r.AlgorithmPartitionByDate.partitionCount()
	default:
		return -1
	}
}

func (r *Rule) validate(c *Config) bool {
	pass := true
	if nil == r.Columns || 0 == len(r.Columns) {
		r.ColumnsError = "should not be empty"
		pass = false
	}

	switch r.Algorithm {
	case "algorithmPartitionByMap":
		if nil == r.AlgorithmPartitionByMap {
			r.AlgorithmError = "no algorithm detail"
			pass = false
		} else if !r.AlgorithmPartitionByMap.validate(c) {
			pass = false
		}
	case "algorithmAutoPartitionByLong":
		if nil == r.AlgorithmAutoPartitionByLong {
			r.AlgorithmError = "no algorithm detail"
			pass = false
		} else if !r.AlgorithmAutoPartitionByLong.validate(c) {
			pass = false
		}
	case "algorithmPartitionByLong":
		if nil == r.AlgorithmPartitionByLong {
			r.AlgorithmError = "no algorithm detail"
			pass = false
		} else if !r.AlgorithmPartitionByLong.validate(c) {
			pass = false
		}
	case "algorithmPartitionByString":
		if nil == r.AlgorithmPartitionByString {
			r.AlgorithmError = "no algorithm detail"
			pass = false
		} else if !r.AlgorithmPartitionByString.validate(c) {
			pass = false
		}
	case "algorithmPartitionByDate":
		if nil == r.AlgorithmPartitionByDate {
			r.AlgorithmError = "no algorithm detail"
			pass = false
		} else if !r.AlgorithmPartitionByDate.validate(c) {
			pass = false
		}
	default:
		r.AlgorithmError = "invaild algorithm"
		pass = false
	}

	return pass
}

type AlgorithmPartitionByMap struct {
	Type        int            `json:"type"` // 0-number 1-string
	TypeError   string         `json:"type_error"`
	Map         map[string]int `json:"map"`
	MapError    string         `json:"map_error"`
	DefaultNode int            `json:"defaultNode"`
}

func (a *AlgorithmPartitionByMap) IsSame(t *AlgorithmPartitionByMap) bool {
	if a.Type != t.Type {
		return false
	}
	if len(a.Map) != len(t.Map) {
		return false
	}
	for k, v := range a.Map {
		if v != t.Map[k] {
			return false
		}
	}
	if a.DefaultNode != t.DefaultNode {
		return false
	}
	return true
}

func (a *AlgorithmPartitionByMap) partitionCount() int {
	ret := 0
	for _, val := range a.Map {
		if val > ret {
			ret = val
		}
	}
	return ret
}

func (a *AlgorithmPartitionByMap) validate(c *Config) bool {
	pass := true
	if !(0 == a.Type || 1 == a.Type) {
		a.TypeError = "should be 0 or 1"
		pass = false
	}
	if nil == a.Map || 0 == len(a.Map) {
		a.MapError = "should not be empty"
		pass = false
	}
	return pass
}

type AlgorithmAutoPartitionByLong struct {
	Ranges      []*AlgorithmAutoPartitionByLongRange `json:"ranges"`
	RangesError string                               `json:"ranges_error"`
	DefaultNode int                                  `json:"defaultNode"`
}

func (a *AlgorithmAutoPartitionByLong) IsSame(t *AlgorithmAutoPartitionByLong) bool {
	if len(a.Ranges) != len(t.Ranges) {
		return false
	}
	for idx := range a.Ranges {
		if !a.Ranges[idx].IsSame(t.Ranges[idx]) {
			return false
		}
	}
	if a.DefaultNode != t.DefaultNode {
		return false
	}
	return true
}

func (a *AlgorithmAutoPartitionByLong) partitionCount() int {
	ret := 0
	for _, r := range a.Ranges {
		if r.NodeIndex > ret {
			ret = r.NodeIndex
		}
	}
	return ret
}

func (a *AlgorithmAutoPartitionByLong) validate(c *Config) bool {
	pass := true
	if nil == a.Ranges || 0 == len(a.Ranges) {
		a.RangesError = "should not be empty"
		pass = false
	} else {
		for _, r := range a.Ranges {
			if r.Start > r.End {
				a.RangesError = "start should <= end"
				pass = false
			}
			for _, r2 := range a.Ranges {
				if r == r2 {
					continue
				}
				if (r.Start <= r2.Start && r.End >= r2.Start) || (r.Start <= r2.End && r.End >= r2.End) {
					a.RangesError = "should not have overlap"
					pass = false
				}
			}
		}
	}
	return pass
}

type AlgorithmAutoPartitionByLongRange struct {
	Start     int `json:"start"`
	End       int `json:"end"`
	NodeIndex int `json:"nodeIndex"`
}

func (a *AlgorithmAutoPartitionByLongRange) IsSame(t *AlgorithmAutoPartitionByLongRange) bool {
	if a.Start != t.Start {
		return false
	}
	if a.End != t.End {
		return false
	}
	if a.NodeIndex != t.NodeIndex {
		return false
	}
	return true
}

type AlgorithmPartitionByLong struct {
	Partitions      []*AlgorithmPartitionByLongPartition `json:"partitions"`
	PartitionsError string                               `json:"partitions_error"`
}

func (a *AlgorithmPartitionByLong) IsSame(t *AlgorithmPartitionByLong) bool {
	if len(a.Partitions) != len(t.Partitions) {
		return false
	}
	for idx := range a.Partitions {
		if !a.Partitions[idx].IsSame(t.Partitions[idx]) {
			return false
		}
	}
	return true
}

type AlgorithmPartitionByLongPartition struct {
	Count  int `json:"count"`
	Length int `json:"length"`
}

func (a *AlgorithmPartitionByLongPartition) IsSame(t *AlgorithmPartitionByLongPartition) bool {
	if a.Count != t.Count {
		return false
	}
	if a.Length != t.Length {
		return false
	}
	return true
}

func (a *AlgorithmPartitionByLong) partitionCount() int {
	ret := 0
	for _, p := range a.Partitions {
		ret += p.Count
	}
	return ret
}

func (a *AlgorithmPartitionByLong) validate(c *Config) bool {
	pass := true
	if nil == a.Partitions || 0 == len(a.Partitions) {
		a.PartitionsError = "should not be empty"
		pass = false
	} else {
		sum := 0
		for _, p := range a.Partitions {
			sum += p.Count * p.Length
		}
		if sum > 2880 {
			a.PartitionsError = "total should <= 2880"
			pass = false
		}
	}
	return pass
}

var hashSlicePattern = regexp.MustCompile(`^\d*:(-?\d*)$`)

type AlgorithmPartitionByString struct {
	Partitions      []*AlgorithmPartitionByStringPartition `json:"partitions"`
	PartitionsError string                                 `json:"partitions_error"`
	HashSlice       string                                 `json:"hashSlice"`
	HashSliceError  string                                 `json:"hashSlice_error"`
}

func (a *AlgorithmPartitionByString) IsSame(t *AlgorithmPartitionByString) bool {
	if len(a.Partitions) != len(t.Partitions) {
		return false
	}
	for idx := range a.Partitions {
		if a.Partitions[idx] != t.Partitions[idx] {
			return false
		}
	}
	if a.HashSlice != t.HashSlice {
		return false
	}
	return true
}

type AlgorithmPartitionByStringPartition struct {
	Count  int `json:"count"`
	Length int `json:"length"`
}

func (a *AlgorithmPartitionByStringPartition) IsSame(t *AlgorithmPartitionByStringPartition) bool {
	if a.Count != t.Count {
		return false
	}
	if a.Length != t.Length {
		return false
	}
	return true
}

func (a *AlgorithmPartitionByString) partitionCount() int {
	ret := 0
	for _, p := range a.Partitions {
		ret += p.Count
	}
	return ret
}

func (a *AlgorithmPartitionByString) validate(c *Config) bool {
	pass := true
	if nil == a.Partitions || 0 == len(a.Partitions) {
		a.PartitionsError = "should not be empty"
		pass = false
	} else {
		sum := 0
		for _, p := range a.Partitions {
			sum += p.Count * p.Length
		}
		if 1024 != sum {
			a.PartitionsError = "total is not 1024 partitions"
			pass = false
		}
	}

	if !hashSlicePattern.MatchString(a.HashSlice) {
		a.HashSliceError = "should match `^\\d*:\\d*$`"
		pass = false
	}

	return pass
}

type AlgorithmPartitionByDate struct {
	DateFormat        string `json:"dateFormat"`
	BeginDate         string `json:"sBeginDate"`
	BeginDateError    string `json:"sBeginDate_error"`
	EndDate           string `json:"sEndDate"`
	EndDateError      string `json:"sEndDate_error"`
	PartitionDay      int    `json:"sPartionDay"`
	PartitionDayError string `json:"sPartionDay_error"`
}

func (a *AlgorithmPartitionByDate) IsSame(t *AlgorithmPartitionByDate) bool {
	if a.DateFormat != t.DateFormat {
		return false
	}
	if a.BeginDate != t.BeginDate {
		return false
	}
	if a.EndDate != t.EndDate {
		return false
	}
	if a.PartitionDay != t.PartitionDay {
		return false
	}
	return true
}

var dateReplacer = strings.NewReplacer("yyyy", "2006", "MM", "01", "dd", "02", "HH", "15", "mm", "04", "ss", "05")

func (a *AlgorithmPartitionByDate) partitionCount() int {
	if "" == a.EndDate {
		//no need to check
		return -1
	}
	dateFormat := dateReplacer.Replace(a.DateFormat)
	beginDate, _ := time.Parse(dateFormat, a.BeginDate)
	endDate, _ := time.Parse(dateFormat, a.EndDate)
	partitions := endDate.Sub(beginDate).Hours() / float64(a.PartitionDay) / 24
	return int(math.Ceil(partitions))
}

func (a *AlgorithmPartitionByDate) validate(c *Config) bool {
	pass := true
	dateFormat := dateReplacer.Replace(a.DateFormat)

	if "" == a.BeginDate {
		a.BeginDateError = "should not be empty"
		pass = false
	} else if _, err := time.Parse(dateFormat, a.BeginDate); nil != err {
		a.BeginDateError = "validate format failed"
		pass = false
	}

	if "" != a.EndDate {
		if _, err := time.Parse(dateFormat, a.EndDate); nil != err {
			a.EndDateError = "validate format failed"
			pass = false
		}
	}

	if !(a.PartitionDay > 0) {
		a.PartitionDayError = "should >0"
		pass = false
	}
	return pass
}

type User struct {
	Name           string `json:"name"`
	NameError      string `json:"name_error"`
	Password       string `json:"password"`
	ReadOnly       bool   `json:"readOnly"`
	Schemas        string `json:"schemas"`
	SchemasError   string `json:"schemas_error"`
	Benchmark      int    `json:"benchmark"`
	BenchmarkError string `json:"benchmark_error"`
}

func (s *User) validate(c *Config, adminUser string) bool {
	s.NameError = validateName(s.Name)
	s.SchemasError = validateMycatSpecList(s.Schemas)
	if !(s.Benchmark > 0) {
		s.BenchmarkError = "should >0"
	}
	if "" == s.NameError && s.Name == adminUser {
		s.NameError = "duplicate with admin user"
	}
	return "" == s.NameError+s.SchemasError+s.BenchmarkError
}

/* config in xml file */
type CoreRuleXml struct {
	XMLName    xml.Name         `xml:"mycat:rule"`
	Xmlns      string           `xml:"xmlns:mycat,attr"`
	TableRules []*CoreTableRule `xml:"tableRule"`
	Functions  []*CoreFunction  `xml:"function"`
}

type CoreTableRule struct {
	XMLName xml.Name  `xml:"tableRule"`
	Name    string    `xml:"name,attr"`
	Rule    *CoreRule `xml:"rule"`
}

type CoreRule struct {
	XMLName   xml.Name `xml:"rule"`
	Columns   string   `xml:"columns"`
	Algorithm string   `xml:"algorithm"`
}

type CoreFunction struct {
	XMLName    xml.Name        `xml:"function"`
	Name       string          `xml:"name,attr"`
	Class      string          `xml:"class,attr"`
	Properties []*CoreProperty `xml:"property"`
}

type CoreServerXml struct {
	XMLName xml.Name    `xml:"mycat:server"`
	Xmlns   string      `xml:"xmlns:mycat,attr"`
	System  *CoreSystem `xml:"system"`
	Users   []*CoreUser `xml:"user"`
	Alarm   *DbleAlarm  `xml:"alarm"`
}

type DbleAlarm struct {
	XMLName       xml.Name `xml:"alarm"`
	Url           string   `xml:"url"`
	Port          string   `xml:"port"`
	Level         string   `xml:"level"`
	ServerId      string   `xml:"serverId"`
	ComponentId   string   `xml:"componentId"`
	ComponentType string   `xml:"componentType"`
}

type CoreSystem struct {
	XMLName    xml.Name        `xml:"system"`
	Properties []*CoreProperty `xml:"property"`
}

type CoreProperty struct {
	XMLName xml.Name `xml:"property"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",innerxml"`
}

type CoreUser struct {
	XMLName    xml.Name        `xml:"user"`
	Name       string          `xml:"name,attr"`
	Properties []*CoreProperty `xml:"dataHost"`
}

type CoreSchemaXml struct {
	XMLName   xml.Name        `xml:"mycat:schema"`
	Xmlns     string          `xml:"xmlns:mycat,attr"`
	Schemas   []*CoreSchema   `xml:"schema"`
	DataNodes []*CoreDataNode `xml:"dataNode"`
	DataHosts []*CoreDataHost `xml:"dataHost"`
}

type DbleSchemaXmlForUnmarshal struct {
	//struct for unmarshal: https://github.com/golang/go/issues/9519
	XMLName xml.Name `xml:"dble schema"`
	//Xmlns     string   `xml:"xmlns:dble,attr"` //remove Xmlns because of xml unmarshal bug
	Schemas   []*CoreSchema   `xml:"schema"`
	DataNodes []*CoreDataNode `xml:"dataNode"`
	DataHosts []*CoreDataHost `xml:"dataHost"`
}

type CoreSchema struct {
	XMLName        xml.Name     `xml:"schema"`
	Name           string       `xml:"name,attr"`
	CheckSqlSchema bool         `xml:"checkSQLschema,attr"`
	SqlMaxLimit    int          `xml:"sqlMaxLimit,attr,omitempty"`
	DataNode       string       `xml:"dataNode,attr,omitempty"`
	Tables         []*CoreTable `xml:"table"`
}

type CoreTable struct {
	XMLName       xml.Name          `xml:"table"`
	Name          string            `xml:"name,attr,omitempty"`
	Type          string            `xml:"type,attr,omitempty"`
	PrimaryKey    string            `xml:"primaryKey,attr,omitempty"`
	AutoIncrement bool              `xml:"autoIncrement,attr,omitempty"`
	NeedAddLimit  bool              `xml:"needAddLimit,attr,omitempty"`
	Rule          string            `xml:"rule,attr,omitempty"`
	RuleRequired  bool              `xml:"ruleRequired,attr,omitempty"`
	DataNode      string            `xml:"dataNode,attr,omitempty"`
	ChildTables   []*CoreChildTable `xml:"childTable"`
}

type CoreChildTable struct {
	XMLName       xml.Name          `xml:"childTable"`
	Name          string            `xml:"name,attr,omitempty"`
	JoinKey       string            `xml:"joinKey,attr,omitempty"`
	ParentKey     string            `xml:"parentKey,attr,omitempty"`
	PrimaryKey    string            `xml:"primaryKey,attr,omitempty"`
	AutoIncrement bool              `xml:"autoIncrement,attr,omitempty"`
	ChildTables   []*CoreChildTable `xml:"childTable"`
}

type CoreDataNode struct {
	XMLName  xml.Name `xml:"dataNode"`
	Name     string   `xml:"name,attr"`
	DataHost string   `xml:"dataHost,attr"`
	Database string   `xml:"database,attr"`
}

type CoreDataHost struct {
	XMLName               xml.Name         `xml:"dataHost"`
	Name                  string           `xml:"name,attr"`
	DbType                string           `xml:"dbType,attr"`
	DbDriver              string           `xml:"dbDriver,attr"`
	MaxCon                int              `xml:"maxCon,attr"`
	MinCon                int              `xml:"minCon,attr"`
	Balance               int              `xml:"balance,attr"`
	SwitchType            int              `xml:"switchType,attr"`
	SlaveThreshold        int              `xml:"slaveThreshold,attr"`
	TempReadHostAvailable int              `xml:"tempReadHostAvailable,attr"`
	Heartbeat             string           `xml:"heartbeat"`
	WriteHosts            []*CoreWriteHost `xml:"writeHost"`
	WriteType             string           `xml:"writeType,attr"`
}

type CoreWriteHost struct {
	XMLName      xml.Name        `xml:"writeHost"`
	Host         string          `xml:"host,attr"`
	Url          string          `xml:"url,attr"`
	User         string          `xml:"user,attr"`
	Password     string          `xml:"password,attr"`
	ReadHosts    []*CoreReadHost `xml:"readHost"`
	UsingDecrypt string          `xml:"usingDecrypt,attr"`
}

func (c *CoreWriteHost) toReadHost() *CoreReadHost {
	return &CoreReadHost{
		Host:     c.Host,
		Url:      c.Url,
		User:     c.User,
		Password: c.Password,
	}
}

type CoreReadHost struct {
	XMLName      xml.Name `xml:"readHost"`
	Host         string   `xml:"host,attr"`
	Url          string   `xml:"url,attr"`
	User         string   `xml:"user,attr"`
	Password     string   `xml:"password,attr"`
	Weight       string   `xml:"weight,attr"`
	UsingDecrypt string   `xml:"usingDecrypt,attr"`
}

func (c *CoreReadHost) toWriteHost() *CoreWriteHost {
	return &CoreWriteHost{
		Host:     c.Host,
		Url:      c.Url,
		User:     c.User,
		Password: c.Password,
	}
}

func ParseSchemaXml(content string) (*DbleSchemaXmlForUnmarshal, error) {
	//remove xmlns:dble because of xml unmarshal bug
	content = strings.Replace(content, `xmlns:dble="http://dble.cloud/"`, ``, -1)

	schemaXml := &DbleSchemaXmlForUnmarshal{
		Schemas: []*CoreSchema{},
	}
	if err := xml.Unmarshal([]byte(content), schemaXml); nil != err {
		return nil, fmt.Errorf("invalid xml: %v", err.Error())
	}
	return schemaXml, nil
}
