//go:build enterprise
// +build enterprise

package optimization

import (
	"context"
	"fmt"
	"strings"

	dmsCommonHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
	"github.com/actiontech/sqle/sqle/config"
	"github.com/actiontech/sqle/sqle/model"
)

func getUserKey() string {
	return config.GetOptions().SqleOptions.OptimizationConfig.OptimizationKey
}
func getPawHost() string {
	return config.GetOptions().SqleOptions.OptimizationConfig.OptimizationURL
}

type BaseReply struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 创建空间
type CreateWorkspaceReq struct {
	UserKey    string `json:"userKey"`              // 激活码
	Mode       string `json:"mode"`                 // 线上: online/离线: offline
	DbType     string `json:"dbType"`               // 数据库类型: mysql/postgresql/opengauss/oracle
	Host       string `json:"host,omitempty"`       // 数据库主机地址
	Port       string `json:"port,omitempty"`       // 数据库端口号
	Database   string `json:"database,omitempty"`   // 数据库名称
	Schemas    string `json:"schemas,omitempty"`    // 模式列表(pg/opengauss/orcale需要填写)，若有多个按英文逗号分隔且区分大小写
	DbUser     string `json:"dbUser,omitempty"`     // 访问数据库的用户名
	DbPassword string `json:"dbPassword,omitempty"` // 用户名对应密码
	DdlText    string `json:"ddlText,omitempty"`    // 工作空间对应ddl文本(offline模式)
}

type CreateWorkspaceReply struct {
	BaseReply
	Data CreateWorkspace `json:"data"`
}
type CreateWorkspace struct {
	WorkspaceId string `json:"workspaceId"`
}

// 在线模式（online）
func (a *OptimizationPawSQLServer) createWorkspaceOnline(ctx context.Context, instance *model.Instance, schema string) (string, error) {
	req := CreateWorkspaceReq{
		UserKey:    getUserKey(),
		Mode:       "online",
		DbType:     strings.ToLower(instance.DbType),
		Host:       instance.Host,
		Port:       instance.Port,
		Database:   schema,
		Schemas:    "",
		DbUser:     instance.User,
		DbPassword: instance.Password,
	}
	if instance.DbType == "Oracle" {
		req.Schemas, req.Database = schema, ""
		if p := instance.AdditionalParams.GetParam("service_name"); p != nil {
			req.Database = p.Value
		}
	}

	reply := new(CreateWorkspaceReply)
	err := dmsCommonHttp.POST(ctx, getPawHost()+"/api/v1/createWorkspace", nil, req, reply)
	if err != nil {
		return "", err
	}
	// todo: 私有化部署响应无code字段 https://github.com/actiontech/sqle-ee/issues/1527
	if reply.Message != "操作成功" {
		return "", fmt.Errorf("code is %v,message is %v", reply.Code, reply.Message)
	}
	return reply.Data.WorkspaceId, nil
}

// 离线模式：定义DDL建表
func (a *OptimizationPawSQLServer) createWorkspaceOffline(ctx context.Context, dbType string, ddlText string) (string, error) {
	req := CreateWorkspaceReq{
		UserKey: getUserKey(),
		Mode:    "offline",
		DbType:  dbType,
		DdlText: ddlText,
	}
	reply := new(CreateWorkspaceReply)
	err := dmsCommonHttp.POST(ctx, getPawHost()+"/api/v1/workspaces", nil, req, reply)
	if err != nil {
		return "", err
	}
	// todo: 私有化部署响应无code字段 https://github.com/actiontech/sqle-ee/issues/1527
	if reply.Message != "操作成功" {
		return "", fmt.Errorf("code is %v,message is %v", reply.Code, reply.Message)
	}
	return reply.Data.WorkspaceId, nil
}

// 创建SQL优化任务
// Post http://${server-host}:${server-port}/api/v1/createAnalysis
func (a *OptimizationPawSQLServer) createOptimization(ctx context.Context, workspaceId string, instance *model.Instance, workload string, queryMode string) (string, error) {
	rules, err := getOptimizationReqRules(instance)
	if err != nil {
		return "", err
	}
	req := CreateOptimizationReq{
		UserKey:      getUserKey(),
		Workspace:    workspaceId,
		DBType:       instance.DbType,
		Workload:     workload,
		QueryMode:    queryMode,
		ValidateFlag: true,
		Rules:        rules,
	}

	if len(rules) == 0 {
		req.CloseRewrite = true
	}

	reply := new(CreateOptimizationReply)
	err = dmsCommonHttp.POST(context.TODO(), getPawHost()+"/api/v1/createAnalysis", nil, req, reply)
	if err != nil {
		return "", err
	}
	// todo: 私有化部署响应无code字段 https://github.com/actiontech/sqle-ee/issues/1527
	if reply.Message != "操作成功" {
		return "", fmt.Errorf("code is %v,message is %v", reply.Code, reply.Message)
	}
	if reply.Data.Status != "success1" {
		return "", fmt.Errorf("create optimization failed")
	}
	return reply.Data.OptimizationId, nil
}

type CreateOptimizationReq struct {
	UserKey      string                     `json:"userKey"`
	Workspace    string                     `json:"workspace"`
	DBType       string                     `json:"dbType"`
	Workload     string                     `json:"workload"`
	QueryMode    string                     `json:"queryMode"`
	ValidateFlag bool                       `json:"validateFlag"`
	CloseRewrite bool                       `json:"closeRewrite"`
	Rules        []*CreateOptimizationRules `json:"rules"`
}

type CreateOptimizationRules struct {
	RuleCode  string `json:"ruleCode"`
	Rewrite   bool   `json:"rewrite"`
	Threshold string `json:"threshold"`
}

type CreateOptimizationReply struct {
	Data CreateOptimizationData `json:"data"`
	BaseReply
}
type CreateOptimizationData struct {
	OptimizationId string `json:"analysisId"`
	Status         string `json:"status"`
}

// 查询优化SQL概览
//  Post http://${server-host}:${server-port}/api/v1/getAnalysisSummary

func (a *OptimizationPawSQLServer) getOptimizationSummary(ctx context.Context, optimizationId string) (ret OptimizationSummaryBody, err error) {
	req := OptimizationListReq{
		UserKey:        getUserKey(),
		OptimizationId: optimizationId,
	}

	reply := new(OptimizationSummaryReply)
	err = dmsCommonHttp.POST(ctx, getPawHost()+"/api/v1/getAnalysisSummary", nil, req, reply)
	if err != nil {
		return ret, err
	}
	// todo: 私有化部署响应无code字段 https://github.com/actiontech/sqle-ee/issues/1527
	if reply.Message != "操作成功" {
		return ret, fmt.Errorf("code is %v,message is %v", reply.Code, reply.Message)
	}
	return reply.Data, nil
}

type OptimizationListReq struct {
	UserKey        string `json:"userKey"`
	OptimizationId string `json:"analysisId"`
}

type OptimizationSummaryReply struct {
	BaseReply
	Data OptimizationSummaryBody `json:"data"`
}

type OptimizationSummaryBody struct {
	Status                string                 `json:"status"`
	BasicSummary          OptimizationSummary    `json:"basicSummary"`
	OptimizationRuleInfo  []OptimizationRuleInfo `json:"analysisRuleInfo"`
	OptimizationIndexInfo []string               `json:"analysisIndexInfo"`
	SummaryStatementInfo  []SummaryStatementInfo `json:"summaryStatementInfo"`
}

type OptimizationSummary struct {
	OptimizationSummaryId  string  `json:"analysisSummaryId"`
	OptimizationId         string  `json:"analysisId"`
	NumberOfQuery          int     `json:"numberOfQuery"`
	NumberOfSyntaxError    int     `json:"numberOfSyntaxError"`
	NumberOfRewrite        int     `json:"numberOfRewrite"`
	NumberOfRewrittenQuery int     `json:"numberOfRewrittenQuery"`
	NumberOfViolations     int     `json:"numberOfViolations"`
	NumberOfViolatedQuery  int     `json:"numberOfViolatedQuery"`
	NumberOfIndex          int     `json:"numberOfIndex"`
	NumberOfQueryIndex     int     `json:"numberOfQueryIndex"`
	PerformanceImprove     float64 `json:"performanceImprove"`
	SummaryMarkdown        string  `json:"summaryMarkdown"`
	SummaryMarkdownZh      string  `json:"summaryMarkdownZh"`
	CommentCount           string  `json:"commentCount"`
	NeedReply              string  `json:"needReply"`
}

type OptimizationRuleInfo struct {
	RuleName    string `json:"ruleName"`
	StmtNameStr string `json:"stmtNameStr"`
}

type SummaryStatementInfo struct {
	OptimizationStmtId  string  `json:"analysisStmtId"`
	StmtId              string  `json:"stmtId"`
	StmtName            string  `json:"stmtName"`
	StmtType            string  `json:"stmtType"`
	StmtText            string  `json:"stmtText"`
	CostBefore          float64 `json:"costBefore"`
	CostAfter           float64 `json:"costAfter"`
	NumberOfRewrite     int     `json:"numberOfRewrite"`
	NumberOfViolations  int     `json:"numberOfViolations"`
	NumberOfSyntaxError int     `json:"numberOfSyntaxError"`
	NumberOfIndex       int     `json:"numberOfIndex"`
	NumberOfHitIndex    int     `json:"numberOfHitIndex"`
	Performance         float64 `json:"performance"`
	ContributingIndices string  `json:"contributingIndices"`
	CommentCount        string  `json:"commentCount"`
	NeedReply           string  `json:"needReply"`
}

// 查询优化详情
//
//	Post http://${server-host}:${server-port}/api/v1/getStatementDetails
func (a *OptimizationPawSQLServer) getOptimizationDetail(ctx context.Context, optimizationStmtId string) (ret OptimizationDetail, err error) {
	req := OptimizationDetailReq{
		UserKey:            getUserKey(),
		OptimizationStmtId: optimizationStmtId,
	}
	reply := new(OptimizationDetailReply)
	err = dmsCommonHttp.POST(ctx, getPawHost()+"/api/v1/getStatementDetails", nil, req, reply)
	if err != nil {
		return ret, err
	}
	// todo: 私有化部署响应无code字段 https://github.com/actiontech/sqle-ee/issues/1527
	if reply.Message != "操作成功" {
		return ret, fmt.Errorf("code is %v,message is %v", reply.Code, reply.Message)
	}
	return reply.Data, nil
}

type OptimizationDetailReq struct {
	UserKey            string `json:"userKey"`
	OptimizationStmtId string `json:"analysisStmtId"`
}

type OptimizationDetailReply struct {
	BaseReply
	Data OptimizationDetail `json:"data"`
}
type OptimizationDetail struct {
	OptimizationId       string   `json:"analysisId"`
	OptimizationName     string   `json:"analysisName"`
	StmtId               string   `json:"stmtId"`
	StatementName        string   `json:"statementName"`
	StmtText             string   `json:"stmtText"`
	DetailMarkdown       string   `json:"detailMarkdown"`
	DetailMarkdownZh     string   `json:"detailMarkdownZh"`
	OpenaiOptimizeTextEn string   `json:"openaiOptimizeTextEn"`
	OpenaiOptimizeTextZh string   `json:"openaiOptimizeTextZh"`
	IndexRecommended     []string `json:"indexRecommended"`
	RewrittenQuery       []struct {
		RuleCode            string `json:"ruleCode"`
		RuleNameZh          string `json:"ruleNameZh"`
		RuleNameEn          string `json:"ruleNameEn"`
		RewrittenQueriesStr string `json:"rewrittenQueriesStr"`
		ViolatedQueriesStr  string `json:"violatedQueriesStr"`
	} `json:"rewrittenQuery"`
	ViolationRule []struct {
		RuleName     string `json:"ruleName"`
		FragmentsStr string `json:"fragmentsStr"`
	} `json:"violationRule"`
	ValidationDetails struct {
		BeforeCost        float64 `json:"beforeCost"`
		AfterCost         float64 `json:"afterCost"`
		BeforePlan        string  `json:"beforePlan"`
		AfterPlan         string  `json:"afterPlan"`
		PerformImprovePer float64 `json:"performImprovePer"`
	} `json:"validationDetails"`
}
