package driverV2

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	protoV2 "github.com/actiontech/sqle/sqle/driver/v2/proto"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"golang.org/x/text/language"
)

const (
	// grpc error code
	GrpcErrSQLIsNotSupported = 1000

	RuleVersionUnknown = 0
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
	DriverTypeHANA           = "HANA"
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
	OptionalModuleI18n
	OptionalBackup
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
	case OptionalModuleI18n:
		return "I18n"
	case OptionalBackup:
		return "Backup"
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
	RuleVersionIncluded      []uint32
	EnabledOptionalModule    []OptionalModule
}

func (d *DriverMetas) IsOptionalModuleEnabled(expectModule OptionalModule) bool {
	if d == nil {
		return false
	}
	for _, m := range d.EnabledOptionalModule {
		if m == expectModule {
			return true
		}
	}
	return false
}

func ConvertI18nAuditResultsFromProtoToDriver(pars []*protoV2.AuditResult, isI18n bool) ([]*AuditResult, error) {
	ars := make([]*AuditResult, len(pars))
	for k, par := range pars {
		ar, err := ConvertI18nAuditResultFromProtoToDriver(par, isI18n)
		if err != nil {
			return nil, err
		}
		ars[k] = ar
	}
	return ars, nil
}

func ConvertI18nAuditResultFromProtoToDriver(par *protoV2.AuditResult, isI18n bool) (*AuditResult, error) {
	ar := &AuditResult{
		RuleName:            par.RuleName,
		Level:               RuleLevel(par.Level),
		ExecutionFailed:     par.ExecutionFailed,
		I18nAuditResultInfo: make(map[language.Tag]AuditResultInfo, len(par.I18NAuditResultInfo)),
	}
	if !isI18n {
		// 对非多语言的插件支持
		ar.I18nAuditResultInfo = map[language.Tag]AuditResultInfo{
			i18nPkg.DefaultLang: {Message: par.Message, ErrorInfo: par.ErrorInfo},
		}
	} else {
		if _, exist := par.I18NAuditResultInfo[i18nPkg.DefaultLang.String()]; !exist {
			// 多语言的插件 需包含 i18nPkg.DefaultLang
			return nil, fmt.Errorf("client audit results must support language: %s", i18nPkg.DefaultLang.String())
		}
	}
	for langTag, ruleInfo := range par.I18NAuditResultInfo {
		tag, err := language.Parse(langTag)
		if err != nil {
			return nil, fmt.Errorf("fail to parse I18NAuditResultInfo tag [%s], error: %v", langTag, err)
		}
		ar.I18nAuditResultInfo[tag] = AuditResultInfo{
			Message:   ruleInfo.Message,
			ErrorInfo: ruleInfo.ErrorInfo,
		}
	}
	return ar, nil
}

func ConvertI18nAuditResultFromDriverToProto(ar *AuditResult) *protoV2.AuditResult {
	par := &protoV2.AuditResult{
		Message:             ar.I18nAuditResultInfo[i18nPkg.DefaultLang].Message,
		RuleName:            ar.RuleName,
		Level:               string(ar.Level),
		ExecutionFailed:     ar.ExecutionFailed,
		I18NAuditResultInfo: make(map[string]*protoV2.I18NAuditResultInfo, len(ar.I18nAuditResultInfo)),
	}
	for langTag, ruleInfo := range ar.I18nAuditResultInfo {
		par.I18NAuditResultInfo[langTag.String()] = &protoV2.I18NAuditResultInfo{
			Message:   ruleInfo.Message,
			ErrorInfo: ruleInfo.ErrorInfo,
		}
	}
	return par
}

func ConvertI18nRuleFromProtoToDriver(rule *protoV2.Rule, dbtype string, isI18n bool) (*Rule, error) {
	ps, err := ConvertProtoParamToParam(rule.Params)
	if err != nil {
		return nil, err
	}
	dRule := &Rule{
		Name:         rule.Name,
		Level:        RuleLevel(rule.Level),
		Params:       ps,
		I18nRuleInfo: make(I18nRuleInfo, len(rule.I18NRuleInfo)),
		Version:      rule.Version,
	}
	if dRule.Version == RuleVersionUnknown { // 正确标记旧插件规则的版本
		dRule.Version = GetDriverTypeDefaultRuleVersion(dbtype)
	}
	for langTag, ruleInfo := range rule.I18NRuleInfo {
		tag, err := language.Parse(langTag)
		if err != nil {
			return nil, fmt.Errorf("fail to parse I18NRuleInfo tag [%s], error: %v", langTag, err)
		}
		dRule.I18nRuleInfo[tag] = ConvertI18nRuleInfoFromProtoToDriver(ruleInfo)
	}
	if !isI18n {
		// 对非多语言的插件支持
		ruleInfo := &RuleInfo{
			Desc:       rule.Desc,
			Annotation: rule.Annotation,
			Category:   rule.Category,
		}
		if rule.Knowledge != nil {
			ruleInfo.Knowledge = RuleKnowledge{Content: rule.Knowledge.Content}
		}
		dRule.I18nRuleInfo = I18nRuleInfo{
			i18nPkg.DefaultLang: ruleInfo,
		}
	} else {
		if _, exist := rule.I18NRuleInfo[i18nPkg.DefaultLang.String()]; !exist {
			// 多语言的插件 需包含 i18nPkg.DefaultLang
			return nil, fmt.Errorf("client rule: %s does not support language: %s", rule.Name, i18nPkg.DefaultLang.String())
		}
	}
	return dRule, nil
}

func ConvertI18nRuleInfoFromProtoToDriver(ruleInfo *protoV2.I18NRuleInfo) *RuleInfo {
	return &RuleInfo{
		Desc:       ruleInfo.Desc,
		Category:   ruleInfo.Category,
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
		Desc:       rule.I18nRuleInfo[i18nPkg.DefaultLang].Desc,
		Level:      string(rule.Level),
		Category:   rule.I18nRuleInfo[i18nPkg.DefaultLang].Category,
		Params:     ConvertParamToProtoParam(rule.Params),
		Annotation: rule.I18nRuleInfo[i18nPkg.DefaultLang].Annotation,
		Knowledge: &protoV2.Knowledge{
			Content: rule.I18nRuleInfo[i18nPkg.DefaultLang].Knowledge.Content,
		},
		I18NRuleInfo: make(map[string]*protoV2.I18NRuleInfo, len(rule.I18nRuleInfo)),
		Version:      rule.Version,
	}
	for langTag, ruleInfo := range rule.I18nRuleInfo {
		pRule.I18NRuleInfo[langTag.String()] = ConvertI18nRuleInfoFromDriverToProto(ruleInfo)
	}
	return pRule
}

func ConvertI18nRuleInfoFromDriverToProto(ruleInfo *RuleInfo) *protoV2.I18NRuleInfo {
	return &protoV2.I18NRuleInfo{
		Desc:       ruleInfo.Desc,
		Category:   ruleInfo.Category,
		Annotation: ruleInfo.Annotation,
		Knowledge: &protoV2.Knowledge{
			Content: ruleInfo.Knowledge.Content,
		},
	}
}

func ConvertRuleFromProtoToDriver(rule *protoV2.Rule) (*Rule, error) {
	var ps = make(params.Params, 0, len(rule.Params))
	for _, p := range rule.Params {
		i18nDesc, err := i18nPkg.ConvertStrMap2I18nStr(p.I18NDesc)
		if err != nil {
			return nil, err
		}
		ps = append(ps, &params.Param{
			Key:      p.Key,
			Value:    p.Value,
			Desc:     p.Desc,
			I18nDesc: i18nDesc,
			Type:     params.ParamType(p.Type),
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
		}
		if v.Knowledge != nil {
			ri.Knowledge = RuleKnowledge{Content: v.Knowledge.Content}
		}
		tag, err := language.Parse(langTag)
		if err != nil {
			return nil, fmt.Errorf("fail to parse I18NRuleInfo tag [%s], error: %v", langTag, err)
		}
		dr.I18nRuleInfo[tag] = ri
	}
	return dr, nil
}

func ConvertRuleFromDriverToProto(rule *Rule) *protoV2.Rule {
	pr := &protoV2.Rule{
		Name:         rule.Name,
		Desc:         rule.I18nRuleInfo[i18nPkg.DefaultLang].Desc,
		Level:        string(rule.Level),
		Category:     rule.I18nRuleInfo[i18nPkg.DefaultLang].Category,
		Params:       ConvertParamToProtoParam(rule.Params),
		Annotation:   rule.I18nRuleInfo[i18nPkg.DefaultLang].Annotation,
		Knowledge:    nil,
		I18NRuleInfo: make(map[string]*protoV2.I18NRuleInfo, len(rule.I18nRuleInfo)),
	}
	for langTag, v := range rule.I18nRuleInfo {
		pr.I18NRuleInfo[langTag.String()] = &protoV2.I18NRuleInfo{
			Desc:       v.Desc,
			Category:   v.Category,
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
			Key:      v.Key,
			Value:    v.Value,
			Desc:     v.GetDesc(i18nPkg.DefaultLang),
			I18NDesc: v.I18nDesc.StrMap(),
			Type:     string(v.Type),
		}
	}
	return pp
}

func ConvertProtoParamToParam(p []*protoV2.Param) (params.Params, error) {
	pp := make(params.Params, len(p))
	for i, v := range p {
		if v == nil {
			continue
		}
		i18nDesc, err := i18nPkg.ConvertStrMap2I18nStr(v.I18NDesc)
		if err != nil {
			return nil, fmt.Errorf("fail to convert I18NDesc: %v", err)
		}
		pp[i] = &params.Param{
			Key:      v.Key,
			Value:    v.Value,
			Desc:     v.Desc,
			I18nDesc: i18nDesc,
			Type:     params.ParamType(v.Type),
		}
	}
	return pp, nil
}

func ConvertTabularDataToProto(td TabularData) *protoV2.TabularData {
	columns := make([]*protoV2.TabularDataHead, 0, len(td.Columns))
	for _, c := range td.Columns {
		columns = append(columns, &protoV2.TabularDataHead{
			Name:     c.Name,
			I18NDesc: c.I18nDesc.StrMap(),
			Desc:     c.I18nDesc.GetStrInLang(i18nPkg.DefaultLang),
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

func ConvertProtoTabularDataToDriver(pTd *protoV2.TabularData, isI18n bool) (TabularData, error) {
	columns := make([]TabularDataHead, 0, len(pTd.Columns))
	for _, c := range pTd.Columns {
		h := TabularDataHead{
			Name:     c.Name,
			I18nDesc: nil,
		}
		if isI18n {
			if len(c.I18NDesc) > 0 { // 列描述可以为空，为空时跳过
				if _, exist := c.I18NDesc[i18nPkg.DefaultLang.String()]; !exist {
					// 多语言的插件 需包含 i18nPkg.DefaultLang
					return TabularData{}, fmt.Errorf("client TabularDataHead: %s does not support language: %s", c.Name, i18nPkg.DefaultLang.String())
				}
				i18nDesc, err := i18nPkg.ConvertStrMap2I18nStr(c.I18NDesc)
				if err != nil {
					return TabularData{}, fmt.Errorf("TabularData: %w", err)
				}
				h.I18nDesc = i18nDesc
			}
		} else {
			// 对非多语言的插件支持
			h.I18nDesc.SetStrInLang(i18nPkg.DefaultLang, c.Desc)
		}
		columns = append(columns, h)
	}

	rows := make([][]string, 0, len(pTd.Rows))
	for _, r := range pTd.Rows {
		rows = append(rows, r.Items)
	}

	return TabularData{
		Columns: columns,
		Rows:    rows,
	}, nil
}

func ConvertTableMetaToProto(meta *TableMeta) *protoV2.TableMeta {
	return &protoV2.TableMeta{
		ColumnsInfo:    &protoV2.ColumnsInfo{Data: ConvertTabularDataToProto(meta.ColumnsInfo.TabularData)},
		IndexesInfo:    &protoV2.IndexesInfo{Data: ConvertTabularDataToProto(meta.IndexesInfo.TabularData)},
		CreateTableSQL: meta.CreateTableSQL,
	}
}

func ConvertProtoTableMetaToDriver(meta *protoV2.TableMeta, isI18n bool) (*TableMeta, error) {
	columnsInfo, err := ConvertProtoTabularDataToDriver(meta.ColumnsInfo.Data, isI18n)
	if err != nil {
		return nil, fmt.Errorf("ColumnsInfo: %w", err)
	}
	indexesInfo, err := ConvertProtoTabularDataToDriver(meta.IndexesInfo.Data, isI18n)
	if err != nil {
		return nil, fmt.Errorf("IndexesInfo: %w", err)
	}
	return &TableMeta{
		ColumnsInfo:    ColumnsInfo{TabularData: columnsInfo},
		IndexesInfo:    IndexesInfo{TabularData: indexesInfo},
		CreateTableSQL: meta.CreateTableSQL,
	}, nil
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

func ConvertDatabasSchemaInfoToProto(infos []*DatabasCompareSchemaInfo) []*protoV2.DatabasDiffSchemaInfo {
	dbInfoReq := make([]*protoV2.DatabasDiffSchemaInfo, len(infos))
	for i, dbSchema := range infos {
		dbObjs := make([]*protoV2.DatabaseObject, len(dbSchema.DatabaseObjects))
		for j, dbObj := range dbSchema.DatabaseObjects {
			dbObjs[j] = &protoV2.DatabaseObject{
				ObjectName: dbObj.ObjectName,
				ObjectType: dbObj.ObjectType,
			}
		}
		dbInfoReq[i] = &protoV2.DatabasDiffSchemaInfo{
			BaseSchemaName:     dbSchema.BaseSchemaName,
			ComparedSchemaName: dbSchema.ComparedSchemaName,
			DatabaseObject:     dbObjs,
		}
	}
	return dbInfoReq
}

func ConvertProtoDatabaseDiffReqToDriver(infos []*protoV2.DatabasDiffSchemaInfo) []*DatabasCompareSchemaInfo {
	dbInfoReq := make([]*DatabasCompareSchemaInfo, len(infos))
	for i, dbSchema := range infos {
		dbObjs := make([]*DatabaseObject, len(dbSchema.DatabaseObject))
		for j, dbObj := range dbSchema.DatabaseObject {
			dbObjs[j] = &DatabaseObject{
				ObjectName: dbObj.ObjectName,
				ObjectType: dbObj.ObjectType,
			}
		}
		dbInfoReq[i] = &DatabasCompareSchemaInfo{
			BaseSchemaName:     dbSchema.BaseSchemaName,
			ComparedSchemaName: dbSchema.ComparedSchemaName,
			DatabaseObjects:    dbObjs,
		}
	}
	return dbInfoReq
}

func GetDriverTypeDefaultRuleVersion(dbType string) uint32 {
	switch dbType {
	case DriverTypeTBase, DriverTypeHANA:
		// 这两个插件规则是直接基于知识库开发的，基于知识库的规则版本统一为 2
		return 2
	default:
		return 1
	}
}
