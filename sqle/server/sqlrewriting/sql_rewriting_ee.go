//go:build enterprise
// +build enterprise

package sqlrewriting

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"time"

	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/driver"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/model"
	"github.com/ungerik/go-dry"
	"golang.org/x/text/language"
)

func getUrl() string {
	return config.GetOptions().SqleOptions.SQLRewritingConfig.RewritingURL
}

// Suggestion 定义单个重写建议的结构
type Suggestion struct {
	// 重写建议的规则ID
	RuleID string `json:"rule_id"`
	// 重写建议的类型
	Type string `json:"type"`
	// 重写描述
	Description string `json:"desc,omitempty"`
	// 重写SQL
	RewrittenSql string `json:"rewritten_sql,omitempty"`
	// 需要执行的DDL DCL描述
	DDLDCLDesc string `json:"ddl_dcl_desc,omitempty"`
	// 需要执行的DDL DCL
	DDLDCL string `json:"ddl_dcl,omitempty"`
}

// 获取外部进程完成重写的响应结构
type GetRewriteResponse struct {
	// 重写前SQL的业务描述
	BusinessDesc string `json:"business_desc"`
	// 重写前SQL的执行逻辑描述
	LogicDesc string `json:"logic_desc"`
	// 重写建议列表
	Suggestions []Suggestion `json:"suggestions"`
	// 重写后的SQL的业务差异描述
	BusinessNonEquivalent string `json:"business_non_equivalent,omitempty"`
	// 重写后SQL的业务描述
	BusinessDescAfterOptimize string `json:"business_desc_after_optimize"`
	// 重写后的SQL执行逻辑描述
	LogicDescAfterOptimize string `json:"logic_desc_after_optimize"`
}

type Rule struct {
	// 审核规则Id
	Id string `json:"id"`
	// 审核规则参数
	Params map[string]string `json:"params"`
	// 规则审核结果信息
	Msg string `json:"msg"`
}

// 调用外部进程完成重写的请求结构
type CallRewriteSQLRequest struct {
	DBType              string `json:"db_type"`
	Rules               []Rule `json:"rules"`
	SQL                 string `json:"sql"`
	TableStructures     string `json:"table_structures"`
	Explain             string `json:"explain"`
	EnableStructureType bool   `json:"enable_structure_type"`
}

// 现有规则名 到 重写功能使用的 新规则 ID 的映射
// 只映射了部分规则，用于demo版本
// TODO: 当SQLE使用 新规则ID后，可以去除这个映射
var ruleIDConvert = map[string]string{
	"dml_check_fuzzy_search":                      "00086",
	"dml_check_limit_offset_num":                  "00045",
	"dml_check_select_with_order_by":              "00101",
	"dml_check_sort_direction":                    "00104",
	"dml_check_using_index":                       "00110",
	"dml_not_recommend_func_in_where":             "00009",
	"dml_not_recommend_order_by_rand":             "00131",
	"dml_check_where_exist_implicit_conversion":   "00179",
	"dml_check_where_exist_not":                   "00113",
	"dml_check_where_exist_null":                  "00001",
	"dml_disable_select_all_column":               "00053",
	"dml_check_math_computation_or_func_on_index": "00009",
}

// TODO: 更多预检查规则补充
var rulePreCheck = []string{"dml_enable_explain_pre_check"}

func ConvertRuleIDToRuleName(ruleID string) string {
	for ruleName, id := range ruleIDConvert {
		if id == ruleID {
			return ruleName
		}
	}
	return ""
}

type SQLRewritingParams struct {
	DBType              string                  // 数据库类型
	SQL                 *model.ExecuteSQL       // 需要重写的SQL
	TableStructures     []*driver.TableMeta     // 表结构
	Explain             *driverV2.ExplainResult // SQL 执行计划
	EnableStructureType bool                    // 是否启用涉及数据库结构化的重写
}

func SQLRewriting(ctx context.Context, params *SQLRewritingParams) (*GetRewriteResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params is nil")
	}
	if params.SQL == nil {
		return nil, fmt.Errorf("sql is nil")
	}
	if params.TableStructures == nil {
		return nil, fmt.Errorf("table structures is nil")
	}
	// TODO: 重写功能需要Explain信息，暂时不支持
	// if params.Explain == nil {
	//  return nil, fmt.Errorf("explain is nil")
	// }

	s := model.GetStorage()

	var rules []Rule
	for _, ar := range params.SQL.AuditResults {
		if ar.RuleName == "" {
			continue
		}
		// 不对预检查规则进行重写
		if dry.StringInSlice(ar.RuleName, rulePreCheck) {
			continue
		}

		ruleID, ok := ruleIDConvert[ar.RuleName]
		if !ok {
			return nil, fmt.Errorf("can't find rule id for rule name: %s", ar.RuleName)
		}
		r, exist, err := s.GetRule(ar.RuleName, params.DBType)
		if err != nil {
			return nil, fmt.Errorf("get rule failed: %v", err)
		}
		if !exist {
			return nil, fmt.Errorf("rule not found: %s", ar.RuleName)
		}
		ruleParams := map[string]string{}
		for _, p := range r.Params {
			ruleParams[p.Key] = p.Value
		}

		rules = append(rules, Rule{
			Id:     ruleID,
			Params: ruleParams,
			Msg:    ar.GetAuditMsgByLangTag(language.Chinese),
		})

	}

	// 定义要发送的参数
	req := &CallRewriteSQLRequest{
		DBType:              params.DBType,
		Rules:               rules,
		SQL:                 params.SQL.Content,
		EnableStructureType: params.EnableStructureType,
	}

	if s, err := json.Marshal(params.Explain); err != nil {
		return nil, fmt.Errorf("marshal explain failed: %v", err)
	} else {
		req.Explain = string(s)
	}

	for _, table := range params.TableStructures {
		s, err := json.Marshal(table)
		if err != nil {
			return nil, fmt.Errorf("marshal table structure failed: %v", err)
		}
		req.TableStructures += string(s) + "\n\n"
	}

	reply := &GetRewriteResponse{}

	// 定义 API 端点
	apiURL := getUrl()

	// 先测试API端点的连通性
	connectivityTimeout := 5 * time.Second // 设置连通性测试的超时时间
	if err := checkAPIConnectivity(apiURL, connectivityTimeout); err != nil {
		return nil, fmt.Errorf("API connectivity check failed: %v", err)
	}

	// 发送 HTTP POST 请求
	var callRewriteSQLTimeout int64 = 300 // 设置重写请求的超时时间
	if err := pkgHttp.POST(pkgHttp.SetTimeoutValueContext(ctx, callRewriteSQLTimeout), apiURL, nil, req, reply); err != nil {
		return nil, fmt.Errorf("failed to call %v: %v", apiURL, err)
	}

	return reply, nil

}

func checkAPIConnectivity(apiURL string, timeout time.Duration) error {
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return fmt.Errorf("parse rewriting_url failed: %v, please check config.yml", err)
	}

	// 仅支持http和https协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("unsupported URL scheme: %s", parsedURL.Scheme)
	}

	host := parsedURL.Host
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return fmt.Errorf("cannot connect to API endpoint %s: %v", apiURL, err)
	}
	defer conn.Close()
	return nil
}
