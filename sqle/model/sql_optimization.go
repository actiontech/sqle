package model

import (
	"database/sql/driver"
	"encoding/json"
)

type SQLOptimizationRecord struct {
	Model
	OptimizationId     string  `json:"optimization_id"`
	OptimizationName   string  `json:"optimization_name"`
	DBType             string  `json:"db_type"`
	ProjectId          string  `json:"project_id"`
	InstanceName       string  `json:"instance_name"`
	SchemaName         string  `json:"schema_name"`
	Creator            string  `json:"creator"`
	PerformanceImprove float64 `json:"performance_improve"`
	// summary
	NumberOfQuery          int       `json:"number_of_query"`
	NumberOfSyntaxError    int       `json:"number_of_syntax_error"`
	NumberOfRewrite        int       `json:"number_of_rewrite"`
	NumberOfRewrittenQuery int       `json:"number_of_rewritten_query"`
	NumberOfIndex          int       `json:"number_of_index"`
	NumberOfQueryIndex     int       `json:"number_of_query_index"`
	IndexRecommendations   DBStrings `json:"index_recommendations" gorm:"type:json"`

	OptimizationSQLs []*OptimizationSQL `json:"-" gorm:"foreignkey:OptimizationId;association_foreignkey:OptimizationId"`
}

func (sm SQLOptimizationRecord) TableName() string {
	return "sql_optimization_records"
}

type OptimizationSQL struct {
	Model
	OptimizationId           string                  `json:"optimization_id"`
	OriginalSQL              string                  `json:"original_sql"`
	OptimizedSQL             string                  `json:"optimized_sql"`
	NumberOfRewrite          int                     `json:"number_of_rewrite"`
	NumberOfSyntaxError      int                     `json:"number_of_syntax_error"`
	NumberOfIndex            int                     `json:"number_of_index"`
	NumberOfHitIndex         int                     `json:"number_of_hit_index"`
	Performance              float64                 `json:"performance"`
	ContributingIndices      string                  `json:"contributing_indices"`
	TriggeredRules           RewriteRules            `json:"triggered_rules" gorm:"type:json"`       // 触发的规则
	IndexRecommendations     DBStrings               `json:"index_recommendations" gorm:"type:json"` // 索引建议
	ExplainValidationDetails ExplainValidationDetail `json:"explain_validation_details" gorm:"type:json"`
}

func (sm OptimizationSQL) TableName() string {
	return "optimization_sqls"
}

type DBStrings []string

func (r DBStrings) Value() (driver.Value, error) {
	v, err := json.Marshal(r)
	return string(v), err
}

func (r *DBStrings) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), r)
}

type RewriteRule struct {
	RuleName            string `json:"rule_name"`
	Message             string `json:"message"`
	RewrittenQueriesStr string `json:"rewritten_queries_str"`
	ViolatedQueriesStr  string `json:"violated_queries_str"`
}

type RewriteRules []RewriteRule

func (a RewriteRules) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *RewriteRules) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), a)
}

type ExplainValidationDetail struct {
	BeforeCost        float64 `json:"before_cost"`
	AfterCost         float64 `json:"after_cost"`
	BeforePlan        string  `json:"before_plan"`
	AfterPlan         string  `json:"after_plan"`
	PerformImprovePer float64 `json:"perform_improve_per"`
}

func (e ExplainValidationDetail) Value() (driver.Value, error) {
	v, err := json.Marshal(e)
	return string(v), err
}

func (e *ExplainValidationDetail) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), e)
}
