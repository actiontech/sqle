package driverV2

import (
	"context"
	sqlDriver "database/sql/driver"
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver/common"
	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"golang.org/x/text/language"

	goPlugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

const (
	ProtocolVersion = 2
	PluginSetName   = "driver"
)

var PluginSet = goPlugin.PluginSet{
	PluginSetName: &DriverPlugin{},
}

var HandshakeConfig = goPlugin.HandshakeConfig{
	ProtocolVersion:  ProtocolVersion,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

func ServePlugin(meta DriverMetas, fn func(cfg *Config) (Driver, error)) {
	goPlugin.Serve(&goPlugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: goPlugin.PluginSet{
			PluginSetName: &DriverPlugin{
				Meta:          meta,
				DriverFactory: fn,
			},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: common.NewGRPCServer,
	})
}

type DriverPlugin struct {
	goPlugin.NetRPCUnsupportedPlugin

	Meta          DriverMetas
	DriverFactory func(*Config) (Driver, error)
}

func (dp *DriverPlugin) GRPCServer(broker *goPlugin.GRPCBroker, s *grpc.Server) error {
	protoV2.RegisterDriverServer(s, &DriverGrpcServer{
		Meta:          dp.Meta,
		DriverFactory: dp.DriverFactory,
		Drivers:       map[string]Driver{},
	})
	return nil
}

func (dp *DriverPlugin) GRPCClient(ctx context.Context, broker *goPlugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return protoV2.NewDriverClient(c), nil
}

type Driver interface {
	Close(ctx context.Context)

	// Parse parse sqlText to Node array. sqlText may be single SQL or batch SQLs.
	Parse(ctx context.Context, sql string) ([]Node, error)

	// Audit sql with rules. sql is single SQL text or multi audit.
	Audit(ctx context.Context, sqls []string) ([]*AuditResults, error)

	// GenRollbackSQL generate sql's rollback SQL.
	GenRollbackSQL(ctx context.Context, sql string) (string, string, error)

	Ping(ctx context.Context) error
	Exec(ctx context.Context, sql string) (sqlDriver.Result, error)
	ExecBatch(ctx context.Context, sqls ...string) ([]sqlDriver.Result, error)

	// Tx execute sqls in transaction.
	// When just a sql execute failed, ErrSqlIndex of TxResponse.ExecErr should be set.
	Tx(ctx context.Context, sqls ...string) (*TxResponse, error)
	Query(ctx context.Context, sql string, conf *QueryConf) (*QueryResult, error)
	Explain(ctx context.Context, conf *ExplainConf) (*ExplainResult, error)

	GetDatabases(context.Context) ([]string, error)
	GetTableMeta(ctx context.Context, table *Table) (*TableMeta, error)
	ExtractTableFromSQL(ctx context.Context, sql string) ([]*Table, error)
	EstimateSQLAffectRows(ctx context.Context, sql string) (*EstimatedAffectRows, error)
	KillProcess(ctx context.Context) (*KillProcessInfo, error)
	GetDatabaseObjectDDL(ctx context.Context, objInfos []*DatabaseSchemaInfo) ([]*DatabaseSchemaObjectResult, error)
	GetDatabaseDiffModifySQL(ctx context.Context, calibratedDSN *DSN, objInfos []*DatabasCompareSchemaInfo) ([]*DatabaseDiffModifySQLResult, error)

	Backup(ctx context.Context, req *BackupReq) (*BackupRes, error)
	RecommendBackupStrategy(ctx context.Context, req *RecommendBackupStrategyReq) (*RecommendBackupStrategyRes, error)
	GetSelectivityOfSQLColumns(ctx context.Context, sql string) (map[string] /*table name*/ map[string] /*column name*/ float32, error)
}

const (
	BackupStrategyNone        string = "none"         // 不备份(不支持备份、无需备份、选择不备份)
	BackupStrategyReverseSql  string = "reverse_sql"  // 备份为反向SQL
	BackupStrategyOriginalRow string = "original_row" // 备份为原始行
	BackupStrategyManually    string = "manual"       // 标记为人工备份
)

type BackupReq struct {
	BackupStrategy string
	Sql            string
	BackupMaxRows  uint64
}

type BackupRes struct {
	BackupSql     []string
	ExecuteResult string
}

type RecommendBackupStrategyReq struct {
	Sql string
}

type RecommendBackupStrategyRes struct {
	BackupStrategy    string
	BackupStrategyTip string
	TablesRefer       []string
	SchemasRefer      []string
}

type Node struct {
	// Text is the raw SQL text of Node.
	Text string

	// Type is type of SQL, such as DML/DDL/DCL.
	Type string

	// Fingerprint is fingerprint of Node's raw SQL.
	Fingerprint string

	// StartLine is the starting row number of the Node's raw SQL.
	StartLine uint64

	// ExecBatchId represents the identifier for a group of SQL statements that should be executed within a single context using the ExecBatch method.
	ExecBatchId uint64
}

type RuleLevel string

const (
	RuleLevelNull   RuleLevel = "" // used to indicate no rank
	RuleLevelNormal RuleLevel = "normal"
	RuleLevelNotice RuleLevel = "notice"
	RuleLevelWarn   RuleLevel = "warn"
	RuleLevelError  RuleLevel = "error"
)

var ruleLevelMap = map[RuleLevel]int{
	RuleLevelNull:   -1,
	RuleLevelNormal: 0,
	RuleLevelNotice: 1,
	RuleLevelWarn:   2,
	RuleLevelError:  3,
}

func (r RuleLevel) LessOrEqual(l RuleLevel) bool {
	return ruleLevelMap[r] <= ruleLevelMap[l]
}

func (r RuleLevel) More(l RuleLevel) bool {
	return ruleLevelMap[r] > ruleLevelMap[l]
}

func (r RuleLevel) MoreOrEqual(l RuleLevel) bool {
	return ruleLevelMap[r] >= ruleLevelMap[l]
}

// RuleLevelLessOrEqual return level a <= level b
func RuleLevelLessOrEqual(a, b string) bool {
	return RuleLevel(a).LessOrEqual(RuleLevel(b))
}

type AuditResults struct {
	Results []*AuditResult
}

type AuditResult struct {
	Level               RuleLevel
	RuleName            string
	ExecutionFailed     bool
	I18nAuditResultInfo map[language.Tag]AuditResultInfo
}

type AuditResultInfo struct {
	Message   string
	ErrorInfo string
}

func NewAuditResults() *AuditResults {
	return &AuditResults{
		Results: []*AuditResult{},
	}
}

// Level find highest Level in result
func (rs *AuditResults) Level() RuleLevel {
	level := RuleLevelNull
	for _, curr := range rs.Results {
		if curr.ExecutionFailed {
			continue
		}
		if ruleLevelMap[curr.Level] > ruleLevelMap[level] {
			level = curr.Level
		}
	}
	return level
}

func (rs *AuditResults) Message() string {
	repeatCheck := map[string]struct{}{}
	messages := []string{}
	for _, result := range rs.Results {
		token := result.I18nAuditResultInfo[i18nPkg.DefaultLang].Message + string(result.Level)
		if _, ok := repeatCheck[token]; ok {
			continue
		}
		repeatCheck[token] = struct{}{}

		var message string
		match, _ := regexp.MatchString(fmt.Sprintf(`^\[%s|%s|%s|%s|%s\]`,
			RuleLevelError, RuleLevelWarn, RuleLevelNotice, RuleLevelNormal, "osc"),
			result.I18nAuditResultInfo[i18nPkg.DefaultLang].Message)
		if match {
			message = result.I18nAuditResultInfo[i18nPkg.DefaultLang].Message
		} else {
			message = fmt.Sprintf("[%s]%s", result.Level, result.I18nAuditResultInfo[i18nPkg.DefaultLang].Message)
		}
		messages = append(messages, message)
	}
	return strings.Join(messages, "\n")
}

func (rs *AuditResults) Add(level RuleLevel, ruleName string, i18nMsgPattern i18nPkg.I18nStr, args ...interface{}) {
	rs.AddResultWithError(level, ruleName, "", false, i18nMsgPattern, args...)
}

func (rs *AuditResults) AddResultWithError(level RuleLevel, ruleName, errorMsg string, executionFailed bool, i18nMsgPattern i18nPkg.I18nStr, args ...interface{}) {
	if level == "" || len(i18nMsgPattern) == 0 {
		return
	}

	defer rs.SortByLevel()

	if ruleName != "" {
		for _, v := range rs.Results {
			// 审核结果规则存在则更新
			if v.RuleName == ruleName {
				v.Level = level
				for langTag, msg := range i18nMsgPattern {
					if len(args) > 0 {
						msg = fmt.Sprintf(msg, args...)
					}
					v.I18nAuditResultInfo[langTag] = AuditResultInfo{
						Message:   msg,
						ErrorInfo: errorMsg,
					}
				}
				return
			}
		}
	}

	ar := &AuditResult{
		Level:               level,
		RuleName:            ruleName,
		ExecutionFailed:     executionFailed,
		I18nAuditResultInfo: make(map[language.Tag]AuditResultInfo, len(i18nMsgPattern)),
	}
	for langTag, msg := range i18nMsgPattern {
		if len(args) > 0 {
			msg = fmt.Sprintf(msg, args...)
		}
		ari := AuditResultInfo{
			Message:   msg,
			ErrorInfo: errorMsg,
		}
		ar.I18nAuditResultInfo[langTag] = ari
	}
	rs.Results = append(rs.Results, ar)
}

func (rs *AuditResults) SortByLevel() {
	sort.Slice(rs.Results, func(i, j int) bool {
		return rs.Results[i].Level.More(rs.Results[j].Level)
	})
}

func (rs *AuditResults) HasResult() bool {
	return len(rs.Results) != 0
}

type QueryConf struct {
	TimeOutSecond uint32
}

// The data location in Values should be consistent with that in Column
type QueryResult struct {
	Column params.Params
	Rows   []*QueryResultRow
}

type QueryResultRow struct {
	Values []*QueryResultValue
}

type QueryResultValue struct {
	Value string
}

// TabularData
// the field Columns represents the column name of a table
// the field Rows represents the data of the table
// their relationship is as follows
/*
	| Columns[0]  | Columns[1]  | Columns[2]  |
	| Rows[0][0] | Rows[0][1] | Rows[0][2] |
	| Rows[1][0] | Rows[1][1] | Rows[1][2] |
*/
type TabularDataHead struct {
	Name     string
	I18nDesc i18nPkg.I18nStr
}

type TabularData struct {
	Columns []TabularDataHead
	Rows    [][]string
}

type ExplainConf struct {
	// this SQL should be a single SQL
	Sql string
}

// ExplainJSONResult QueryBlock CostInfo ExplainTable CostInfoTable Explain JSON FORMAT的结果集
type ExplainJSONResult struct {
	QueryBlock QueryBlock `json:"query_block"`
}

type QueryBlock struct {
	SelectID int          `json:"select_id"`
	CostInfo CostInfo     `json:"cost_info"`
	Table    ExplainTable `json:"table"`
}

type CostInfo struct {
	QueryCost string `json:"query_cost"`
}

type ExplainTable struct {
	TableName           string   `json:"table_name"`
	AccessType          string   `json:"access_type"`
	PossibleKeys        []string `json:"possible_keys"`
	Key                 string   `json:"key"`
	UsedKeyParts        []string `json:"used_key_parts"`
	KeyLength           string   `json:"key_length"`
	Ref                 []string `json:"ref"`
	RowsExaminedPerScan int      `json:"rows_examined_per_scan"`
	RowsProducedPerJoin int      `json:"rows_produced_per_join"`
	Filtered            string   `json:"filtered"`
	CostInfo            CostInfo `json:"cost_info"`
	UsedColumns         []string `json:"used_columns"`
}

type CostInfoTable struct {
	ReadCost        string `json:"read_cost"`
	EvalCost        string `json:"eval_cost"`
	PrefixCost      string `json:"prefix_cost"`
	DataReadPerJoin string `json:"data_read_per_join"`
}

type ExplainClassicResult struct {
	TabularData
}

type ExplainResult struct {
	ClassicResult ExplainClassicResult
}

type ColumnsInfo struct {
	TabularData
}

type IndexesInfo struct {
	TabularData
}

type TableMeta struct {
	ColumnsInfo    ColumnsInfo
	IndexesInfo    IndexesInfo
	CreateTableSQL string
	Message        string
}

type Table struct {
	Name   string
	Schema string
}

type EstimatedAffectRows struct {
	Count      int64
	ErrMessage string
}

type KillProcessInfo struct {
	ErrMessage string
}

func NewKillProcessInfo(errorMessage string) *KillProcessInfo {
	return &KillProcessInfo{
		ErrMessage: errorMessage,
	}
}

type RuleKnowledge struct {
	Content string
}

type DatabaseSchemaInfo struct {
	SchemaName      string
	DatabaseObjects []*DatabaseObject
}

type DatabasCompareSchemaInfo struct {
	BaseSchemaName     string
	ComparedSchemaName string
	DatabaseObjects    []*DatabaseObject
}

const (
	ObjectType_TABLE     string = "TABLE"
	ObjectType_VIEW      string = "VIEW"
	ObjectType_PROCEDURE string = "PROCEDURE"
	ObjectType_TRIGGER   string = "TRIGGER"
	ObjectType_EVENT     string = "EVENT"
	ObjectType_FUNCTION  string = "FUNCTION"
)

type DatabaseObject struct {
	ObjectName string
	ObjectType string
}
type DatabaseSchemaObjectResult struct {
	SchemaName         string
	SchemaDDL          string
	DatabaseObjectDDLs []*DatabaseObjectDDL
}
type DatabaseObjectDDL struct {
	DatabaseObject *DatabaseObject
	ObjectDDL      string
}

type DatabaseDiffModifySQLResult struct {
	SchemaName string
	ModifySQLs []string
}

type TxResponse struct {
	// ExecResult indicates the result of successfully executed SQLs.
	ExecResult []sqlDriver.Result
	// ExecErr indicates the error when executing SQL.
	ExecErr *ExecErr
}

type ExecErr struct {
	// ErrSqlIndex indicates the index of the SQL that failed to execute.
	ErrSqlIndex uint32
	// SqlExecErrMsg indicates the error message when executing SQL.
	SqlExecErrMsg string
}
