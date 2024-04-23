//go:build enterprise
// +build enterprise

package model

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/actiontech/sqle/sqle/errors"
)

func init() {
	autoMigrateList = append(autoMigrateList, &SQLOptimizationRecord{})
	autoMigrateList = append(autoMigrateList, &OptimizationSQL{})
}

type SQLOptimizationRecord struct {
	Model
	OptimizationId     string  `json:"optimization_id"`
	OptimizationName   string  `json:"optimization_name"`
	DBType             string  `json:"db_type"`
	ProjectId          string  `json:"project_id"`
	InstanceName       string  `json:"instance_name"`
	SchemaName         string  `json:"schema_name"`
	Creator            string  `json:"creator"`
	Status             string  `json:"status"`
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
	OriginalSQL              string                  `json:"original_sql" gorm:"type:text;not null"`
	OptimizedSQL             string                  `json:"optimized_sql" gorm:"type:text;not null"`
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

var optimizationQueryTpl = `
SELECT sql_optimization_records.optimization_id,sql_optimization_records.optimization_name,  sql_optimization_records.db_type,
sql_optimization_records.instance_name, sql_optimization_records.schema_name,
sql_optimization_records.creator, sql_optimization_records.performance_improve,sql_optimization_records.created_at

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var optimizationCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var optimizationBodyTpl = `
{{ define "body" }}
FROM sql_optimization_records 

WHERE sql_optimization_records.deleted_at IS NULL

{{- if not .current_user_is_admin }}
AND sql_optimization_records.creator = :current_user
{{- end }}


{{- if .fuzzy_search }}
AND (
sql_optimization_records.optimization_id LIKE '%{{ .fuzzy_search }}%'
OR
sql_optimization_records.creator LIKE '%{{ .fuzzy_search }}%'
)
{{- end }}

{{- if .filter_instance_name }}
AND sql_optimization_records.instance_name = :filter_instance_name
{{- end }}

{{- if .filter_project_id }}
AND sql_optimization_records.project_id = :filter_project_id
{{- end }}

{{- if .filter_create_time_from }}
AND sql_optimization_records.create_at >= :filter_create_time_from
{{- end }}

{{- if .filter_create_time_to }}
AND sql_optimization_records.create_at <= :filter_create_time_to
{{- end }}

{{ end }}
`

func (s *Storage) GetOptimizationRecordsByReq(data map[string]interface{}) (
	list []*SQLOptimizationRecord, count uint64, err error) {
	err = s.getListResult(optimizationBodyTpl, optimizationQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(optimizationBodyTpl, optimizationCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

func (s *Storage) GetOptimizationRecordId(optimizationId string) (*SQLOptimizationRecord, error) {
	optimization_record := &SQLOptimizationRecord{}
	err := s.db.
		Where("optimization_id = ?", optimizationId).Find(&optimization_record).Error
	return optimization_record, errors.New(errors.ConnectStorageError, err)
}

var optimizationSQLQueryTpl = `
SELECT id,optimization_id,original_sql,optimized_sql,number_of_rewrite,number_of_syntax_error,number_of_index,number_of_hit_index,performance,contributing_indices,triggered_rules,index_recommendations,explain_validation_details

{{- template "body" . -}} 

{{- if .limit }}
LIMIT :limit OFFSET :offset
{{- end -}}
`

var optimizationSQLCountTpl = `
SELECT COUNT(*)

{{- template "body" . -}}
`

var optimizationSQLBodyTpl = `
{{ define "body" }}
FROM optimization_sqls 

WHERE optimization_sqls.deleted_at IS NULL
AND optimization_sqls.optimization_id = :optimization_id

{{ end }}
`

func (s *Storage) GetOptimizationSQLsByReq(data map[string]interface{}) (
	list []*OptimizationSQL, count uint64, err error) {
	err = s.getListResult(optimizationSQLBodyTpl, optimizationSQLQueryTpl, data, &list)
	if err != nil {
		return nil, 0, err
	}
	count, err = s.getCountResult(optimizationSQLBodyTpl, optimizationSQLCountTpl, data)
	if err != nil {
		return nil, 0, err
	}
	return
}

func (s *Storage) GetOptimizationSQLById(optimizationId string, number int) (*OptimizationSQL, error) {
	optimization_sql := &OptimizationSQL{}
	err := s.db.
		Where("optimization_id = ? AND id = ?", optimizationId, number).Find(&optimization_sql).Error
	return optimization_sql, errors.New(errors.ConnectStorageError, err)
}
