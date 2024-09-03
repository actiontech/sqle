package driverV2

import (
	"fmt"
	"math/rand"
	"time"

	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/pkg/params"
)

const (
	// grpc error code
	GrpcErrSQLIsNotSupported = 1000
)

const (
	SQLTypeDML = "dml"
	SQLTypeDDL = "ddl"
	SQLTypeDQL = "dql"
)

const (
	DriverTypeMySQL          = "MySQL"
	DriverTypePostgreSQL     = "PostgreSQL"
	DriverTypeTiDB           = "TiDB"
	DriverTypeSQLServer      = "SQL Server"
	DriverTypeOracle         = "Oracle"
	DriverTypeDB2            = "DB2"
	DriverTypeOceanBase      = "OceanBase For MySQL"
	DriverTypeTDSQLForInnoDB = "TDSQL For InnoDB"
	DriverTypeTBase          = "TBase"
)

type DriverNotSupportedError struct {
	DriverTyp string
}

func (e *DriverNotSupportedError) Error() string {
	return fmt.Sprintf("driver type %v is not supported", e.DriverTyp)
}

type OptionalModule uint32

const (
	OptionalModuleGenRollbackSQL = iota
	OptionalModuleQuery
	OptionalModuleExplain
	OptionalModuleGetTableMeta
	OptionalModuleExtractTableFromSQL
	OptionalModuleEstimateSQLAffectRows
	OptionalModuleKillProcess
	OptionalExecBatch
)

func (m OptionalModule) String() string {
	switch m {
	case OptionalModuleGenRollbackSQL:
		return "GenRollbackSQL"
	case OptionalModuleQuery:
		return "Query"
	case OptionalModuleExplain:
		return "Explain"
	case OptionalModuleGetTableMeta:
		return "GetTableMeta"
	case OptionalModuleExtractTableFromSQL:
		return "ExtractTableFromSQL"
	case OptionalModuleEstimateSQLAffectRows:
		return "EstimateSQLAffectRows"
	case OptionalModuleKillProcess:
		return "KillProcess"
	case OptionalExecBatch:
		return "ExecBatch"
	default:
		return "Unknown"
	}
}

type DriverMetas struct {
	PluginName               string
	DatabaseDefaultPort      int64
	Logo                     []byte
	DatabaseAdditionalParams params.Params
	Rules                    []*Rule
	EnabledOptionalModule    []OptionalModule
}

func ConvertI18nAuditResultsFromProtoToDriver(pars []*protoV2.AuditResult) ([]*AuditResult, error) {
	ars := make([]*AuditResult, len(pars))
	for k, par := range pars {
		ar, err := ConvertI18nAuditResultFromProtoToDriver(par)
		if err != nil {
			return nil, err
		}
		ars[k] = ar
	}
	return ars, nil
}

func ConvertI18nAuditResultFromProtoToDriver(par *protoV2.AuditResult) (*AuditResult, error) {
	ar := &AuditResult{
		RuleName:            par.RuleName,
		Level:               RuleLevel(par.Level),
		I18nAuditResultInfo: make(map[string]AuditResultInfo, len(par.I18NAuditResultInfo)),
	}
	if len(par.I18NAuditResultInfo) == 0 {
		// 对非多语言的插件支持
		ar.I18nAuditResultInfo = map[string]AuditResultInfo{
			locale.DefaultLang.String(): {Message: par.Message},
		}
	} else {
		if _, exist := par.I18NAuditResultInfo[locale.DefaultLang.String()]; !exist {
			// 多语言的插件 审核结果需包含 locale.DefaultLang
			return nil, fmt.Errorf("client audit results must support language: %s", locale.DefaultLang.String())
		}
	}
	for langTag, ruleInfo := range par.I18NAuditResultInfo {
		ar.I18nAuditResultInfo[langTag] = AuditResultInfo{
			Message: ruleInfo.Message,
		}
	}
	return ar, nil
}

func ConvertI18nAuditResultFromDriverToProto(ar *AuditResult) *protoV2.AuditResult {
	par := &protoV2.AuditResult{
		Message:             ar.I18nAuditResultInfo[locale.DefaultLang.String()].Message,
		RuleName:            ar.RuleName,
		Level:               string(ar.Level),
		I18NAuditResultInfo: make(map[string]*protoV2.I18NAuditResultInfo, len(ar.I18nAuditResultInfo)),
	}
	for langTag, ruleInfo := range ar.I18nAuditResultInfo {
		par.I18NAuditResultInfo[langTag] = &protoV2.I18NAuditResultInfo{
			Message: ruleInfo.Message,
		}
	}
	return par
}

func ConvertI18nRuleFromProtoToDriver(rule *protoV2.Rule) *Rule {
	dRule := &Rule{
		Name:         rule.Name,
		Level:        RuleLevel(rule.Level),
		Params:       ConvertProtoParamToParam(rule.Params),
		I18nRuleInfo: make(I18nRuleInfo, len(rule.I18NRuleInfo)),
	}
	for langTag, ruleInfo := range rule.I18NRuleInfo {
		dRule.I18nRuleInfo[langTag] = ConvertI18nRuleInfoFromProtoToDriver(ruleInfo)
	}
	if len(rule.I18NRuleInfo) == 0 {
		ruleInfo := &RuleInfo{
			Desc:       rule.Desc,
			Annotation: rule.Annotation,
			Category:   rule.Category,
			Params:     ConvertProtoParamToParam(rule.Params),
		}
		if rule.Knowledge != nil {
			ruleInfo.Knowledge = RuleKnowledge{Content: rule.Knowledge.Content}
		}
		dRule.I18nRuleInfo = I18nRuleInfo{
			locale.DefaultLang.String(): ruleInfo,
		}
	}
	return dRule
}

func ConvertI18nRuleInfoFromProtoToDriver(ruleInfo *protoV2.I18NRuleInfo) *RuleInfo {
	return &RuleInfo{
		Desc:       ruleInfo.Desc,
		Category:   ruleInfo.Category,
		Params:     ConvertProtoParamToParam(ruleInfo.Params),
		Annotation: ruleInfo.Annotation,
		Knowledge:  RuleKnowledge{Content: ruleInfo.Knowledge.Content},
	}
}

func ConvertI18nRulesFromDriverToProto(rules []*Rule) []*protoV2.Rule {
	rs := make([]*protoV2.Rule, len(rules))
	for i, rule := range rules {
		rs[i] = ConvertI18nRuleFromDriverToProto(rule)
	}
	return rs
}

func ConvertI18nRuleFromDriverToProto(rule *Rule) *protoV2.Rule {
	// 填充默认语言以支持非多语言插件
	pRule := &protoV2.Rule{
		Name:       rule.Name,
		Desc:       rule.I18nRuleInfo[locale.DefaultLang.String()].Desc,
		Level:      string(rule.Level),
		Category:   rule.I18nRuleInfo[locale.DefaultLang.String()].Category,
		Params:     ConvertParamToProtoParam(rule.Params),
		Annotation: rule.I18nRuleInfo[locale.DefaultLang.String()].Annotation,
		Knowledge: &protoV2.Knowledge{
			Content: rule.I18nRuleInfo[locale.DefaultLang.String()].Knowledge.Content,
		},
		I18NRuleInfo: make(map[string]*protoV2.I18NRuleInfo, len(rule.I18nRuleInfo)),
	}
	for langTag, ruleInfo := range rule.I18nRuleInfo {
		pRule.I18NRuleInfo[langTag] = ConvertI18nRuleInfoFromDriverToProto(ruleInfo)
	}
	return pRule
}

func ConvertI18nRuleInfoFromDriverToProto(ruleInfo *RuleInfo) *protoV2.I18NRuleInfo {
	return &protoV2.I18NRuleInfo{
		Desc:       ruleInfo.Desc,
		Category:   ruleInfo.Category,
		Params:     ConvertParamToProtoParam(ruleInfo.Params),
		Annotation: ruleInfo.Annotation,
		Knowledge: &protoV2.Knowledge{
			Content: ruleInfo.Knowledge.Content,
		},
	}
}

func ConvertRuleFromProtoToDriver(rule *protoV2.Rule) *Rule {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		ps = append(ps, &params.Param{
			Key:   p.Key,
			Value: p.Value,
			Desc:  p.Desc,
			Type:  params.ParamType(p.Type),
		})
	}
	dr := &Rule{
		Name:         rule.Name,
		Level:        RuleLevel(rule.Level),
		Params:       ps,
		I18nRuleInfo: make(I18nRuleInfo, len(rule.I18NRuleInfo)),
	}
	for langTag, v := range rule.I18NRuleInfo {
		ri := &RuleInfo{
			Desc:       v.Desc,
			Annotation: v.Annotation,
			Category:   v.Category,
			Params:     ps,
		}
		if v.Knowledge != nil {
			ri.Knowledge = RuleKnowledge{Content: v.Knowledge.Content}
		}
		dr.I18nRuleInfo[langTag] = ri
	}
	return dr
}

func ConvertRuleFromDriverToProto(rule *Rule) *protoV2.Rule {
	pr := &protoV2.Rule{
		Name:         rule.Name,
		Desc:         "",
		Level:        string(rule.Level),
		Category:     "",
		Params:       ConvertParamToProtoParam(rule.Params),
		Annotation:   "",
		Knowledge:    nil,
		I18NRuleInfo: make(map[string]*protoV2.I18NRuleInfo, len(rule.I18nRuleInfo)),
	}
	for k, v := range rule.I18nRuleInfo {
		pr.I18NRuleInfo[k] = &protoV2.I18NRuleInfo{
			Desc:       v.Desc,
			Category:   v.Category,
			Params:     ConvertParamToProtoParam(v.Params),
			Annotation: v.Annotation,
			Knowledge:  &protoV2.Knowledge{Content: v.Knowledge.Content},
		}
	}
	return pr
}

func ConvertParamToProtoParam(p params.Params) []*protoV2.Param {
	pp := make([]*protoV2.Param, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		pp[i] = &protoV2.Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  string(v.Type),
		}
	}
	return pp
}

func ConvertProtoParamToParam(p []*protoV2.Param) params.Params {
	pp := make(params.Params, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		pp[i] = &params.Param{
			Key:   v.Key,
			Value: v.Value,
			Desc:  v.Desc,
			Type:  params.ParamType(v.Type),
		}
	}
	return pp
}

func ConvertTabularDataToProto(td TabularData) *protoV2.TabularData {
	columns := make([]*protoV2.TabularDataHead, 0, len(td.Columns))
	for _, c := range td.Columns {
		columns = append(columns, &protoV2.TabularDataHead{
			Name: c.Name,
			Desc: c.Desc,
		})
	}

	rows := make([]*protoV2.TabularDataRows, 0, len(td.Rows))
	for _, r := range td.Rows {
		rows = append(rows, &protoV2.TabularDataRows{
			Items: r,
		})
	}

	return &protoV2.TabularData{
		Columns: columns,
		Rows:    rows,
	}
}

func ConvertProtoTabularDataToDriver(pTd *protoV2.TabularData) TabularData {
	columns := make([]TabularDataHead, 0, len(pTd.Columns))
	for _, c := range pTd.Columns {
		columns = append(columns, TabularDataHead{
			Name: c.Name,
			Desc: c.Desc,
		})
	}

	rows := make([][]string, 0, len(pTd.Rows))
	for _, r := range pTd.Rows {
		rows = append(rows, r.Items)
	}

	return TabularData{
		Columns: columns,
		Rows:    rows,
	}
}

func ConvertTableMetaToProto(meta *TableMeta) *protoV2.TableMeta {
	return &protoV2.TableMeta{
		ColumnsInfo:    &protoV2.ColumnsInfo{Data: ConvertTabularDataToProto(meta.ColumnsInfo.TabularData)},
		IndexesInfo:    &protoV2.IndexesInfo{Data: ConvertTabularDataToProto(meta.IndexesInfo.TabularData)},
		CreateTableSQL: meta.CreateTableSQL,
	}
}

func ConvertProtoTableMetaToDriver(meta *protoV2.TableMeta) *TableMeta {
	return &TableMeta{
		ColumnsInfo:    ColumnsInfo{TabularData: ConvertProtoTabularDataToDriver(meta.ColumnsInfo.Data)},
		IndexesInfo:    IndexesInfo{TabularData: ConvertProtoTabularDataToDriver(meta.IndexesInfo.Data)},
		CreateTableSQL: meta.CreateTableSQL,
	}
}

func RandStr(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
