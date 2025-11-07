package sql_flash

import (
	sqlDriver "database/sql/driver"
	"encoding/json"
	"fmt"
)

// 创建SQL优化任务请求
type CreateOptimizeTaskReq struct {
	Type     string `json:"type"`     // SQL类型：SQL 或 MyBatis
	Content  string `json:"content"`  // SQL内容（文本形式）
	Metadata string `json:"metadata"` // 数据库元数据信息
	Explain  string `json:"explain"`  // 执行计划信息
}

// 创建SQL优化任务响应
type CreateOptimizeTaskResp struct {
	Code    json.Number `json:"code"`
	Message string      `json:"message"`
	Data    struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

// 查询优化结果响应
type GetOptimizeResultResp struct {
	Code    json.Number     `json:"code"`
	Message string          `json:"message"`
	Data    *OptimizeResult `json:"data"`
}

// 优化结果数据结构
type OptimizeResult struct {
	ID                 string          `json:"id"`
	OriginSQL          string          `json:"origin_sql"`
	Metadata           string          `json:"metadata"`
	TotalState         string          `json:"total_state"`
	OriginSQLQueryPlan *QueryPlan      `json:"origin_sql_query_plan"`
	Optimize           *OptimizeDetail `json:"optimize"`
	TotalAnalysis      *TotalAnalysis  `json:"total_analysis"`
	AdvisedIndex       *AdvisedIndex   `json:"advised_index"`
}

// 查询计划
type QueryPlan struct {
	State         string           `json:"state"`
	QueryPlanDesc []*QueryPlanNode `json:"query_plan_desc"`
}

// Value impl sql.driver.Valuer interface
func (q QueryPlan) Value() (sqlDriver.Value, error) {
	bytes, err := json.Marshal(q)
	return string(bytes), err
}

// Scan impl sql.Scanner interface
func (q *QueryPlan) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("sql.Scanner scan with nil")
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %s", value)
	}

	result := QueryPlan{}
	err := json.Unmarshal(data, &result)
	*q = result
	return err
}

// 查询计划节点
type QueryPlanNode struct {
	Summary  []string         `json:"summary"`
	Children []*QueryPlanNode `json:"children"`
	Operator string           `json:"operator"`
}

// 优化详情
type OptimizeDetail struct {
	State string          `json:"state"`
	Steps []*OptimizeStep `json:"steps"`
}

// Value impl sql.driver.Valuer interface
func (q OptimizeDetail) Value() (sqlDriver.Value, error) {
	bytes, err := json.Marshal(q)
	return string(bytes), err
}

// Scan impl sql.Scanner interface
func (q *OptimizeDetail) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("sql.Scanner scan with nil")
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %s", value)
	}

	result := OptimizeDetail{}
	err := json.Unmarshal(data, &result)
	*q = result
	return err
}

// 优化步骤
type OptimizeStep struct {
	RuleID       string     `json:"rule_id"`
	RuleName     string     `json:"rule_name"`
	RuleDesc     string     `json:"rule_desc"`
	OptimizedSQL string     `json:"optimized_sql"`
	ChatID       string     `json:"chat_id"`
	QueryPlan    *QueryPlan `json:"query_plan"`
	Analysis     *Analysis  `json:"analysis"`
}

// 分析结果
type Analysis struct {
	State           string            `json:"state"`
	ImprovementRate float64           `json:"improvement_rate"`
	ImprovementDesc string            `json:"improvement_desc"`
	Detail          []*AnalysisDetail `json:"detail"`
	TotalScore      float64           `json:"total_score"`
}

// 分析详情
type AnalysisDetail struct {
	Category       string  `json:"category"`
	OptimizedScore float64 `json:"optimized_score"`
	OriginalScore  float64 `json:"original_score"`
}

// 总体分析
type TotalAnalysis struct {
	State           string            `json:"state"`
	ImprovementRate float64           `json:"improvement_rate"`
	ImprovementDesc string            `json:"improvement_desc"`
	Detail          []*AnalysisDetail `json:"detail"`
	TotalScore      float64           `json:"total_score"`
}

// Value impl sql.driver.Valuer interface
func (q TotalAnalysis) Value() (sqlDriver.Value, error) {
	bytes, err := json.Marshal(q)
	return string(bytes), err
}

// Scan impl sql.Scanner interface
func (q *TotalAnalysis) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("sql.Scanner scan with nil")
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %s", value)
	}

	result := TotalAnalysis{}
	err := json.Unmarshal(data, &result)
	*q = result
	return err
}

// 索引建议
type AdvisedIndex struct {
	State       string       `json:"state"`
	HasAdvice   bool         `json:"has_advice"`
	OtherAdvice string       `json:"other_advice"`
	Indexes     []*IndexInfo `json:"indexes"`
}

// Value impl sql.driver.Valuer interface
func (q AdvisedIndex) Value() (sqlDriver.Value, error) {
	bytes, err := json.Marshal(q)
	return string(bytes), err
}

// Scan impl sql.Scanner interface
func (q *AdvisedIndex) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("sql.Scanner scan with nil")
	}
	var data []byte
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSON value: %s", value)
	}

	result := AdvisedIndex{}
	err := json.Unmarshal(data, &result)
	*q = result
	return err
}

// 索引信息
type IndexInfo struct {
	CreateIndexStatement string `json:"create_index_statement"`
	Reason               string `json:"reason"`
}

// ========== 重写任务相关结构 ==========

// 创建重写任务请求
type CreateRewriteTaskReq struct {
	Type     string `json:"type"`     // SQL类型：SQL 或 MyBatis
	Content  string `json:"content"`  // SQL内容（文本形式）
	Metadata string `json:"metadata"` // 数据库元数据信息
	Explain  string `json:"explain"`  // 执行计划信息
}

// 创建重写任务响应
type CreateRewriteTaskResp struct {
	Code    json.Number `json:"code"`
	Message string      `json:"message"`
	Data    struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

// 重写任务结果
type RewriteResult struct {
	ID            string          `json:"id"`
	TaskID        string          `json:"task_id"`
	OriginSQL     string          `json:"origin_sql"`
	Metadata      string          `json:"metadata"`
	OptimizeState string          `json:"optimize_state"` // 优化状态：running, rewrite_done, failed等
	OptimizeSteps []*OptimizeStep `json:"optimize_steps"` // 优化步骤
}

// 获取重写任务结果响应
type GetRewriteResultResp struct {
	Code    json.Number    `json:"code"`
	Message string         `json:"message"`
	Data    *RewriteResult `json:"data"`
}

// ========== 索引推荐任务相关结构 ==========

// 创建索引推荐任务请求
type CreateAdviseIndexTaskReq struct {
	Explain string `json:"explain"` // 执行计划信息
}

// 创建索引推荐任务响应
type CreateAdviseIndexTaskResp struct {
	Code    json.Number `json:"code"`
	Message string      `json:"message"`
	Data    struct {
		Message string `json:"message"`
	} `json:"data"`
}
