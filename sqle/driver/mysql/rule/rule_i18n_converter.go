package rule

import (
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	driverV2 "github.com/actiontech/sqle/sqle/driver/v2"
	"github.com/actiontech/sqle/sqle/pkg/params"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pingcap/parser/ast"
	"golang.org/x/text/language"
)

type SourceEnum struct {
	Value string        `json:"value"`
	Desc  *i18n.Message `json:"desc"`
}

type SourceParam struct {
	Key   string           `json:"key"`
	Value string           `json:"value"`
	Desc  *i18n.Message    `json:"desc"`
	Type  params.ParamType `json:"type"`
	Enums []SourceEnum     `json:"enums"`
}

// SourceRule 用于初始化时定义国际化的规则
type SourceRule struct {
	Name         string
	Desc         *i18n.Message
	Annotation   *i18n.Message
	Category     *i18n.Message // Deprecated, use CategoryTags instead
	CategoryTags map[string][]string
	Level        driverV2.RuleLevel
	Params       []*SourceParam
	Knowledge    driverV2.RuleKnowledge
	AllowOffline bool
	Version      uint32
}

type SourceHandler struct {
	Rule                 SourceRule
	Message              *i18n.Message
	Func                 RuleHandlerFunc
	NotAllowOfflineStmts []ast.Node
	// 开始事后审核时将会跳过这个值为ture的规则
	OnlyAuditNotExecutedSQL bool
	// 事后审核时将会跳过下方列表中的类型
	NotSupportExecutedSQLAuditStmts []ast.Node
}

// GenerateI18nRuleHandlers 根据规则初始化时定义的 SourceHandler 生成支持多语言的 RuleHandler
func GenerateI18nRuleHandlers(bundle *i18nPkg.Bundle, shs []*SourceHandler) []RuleHandler {
	rhs := make([]RuleHandler, len(shs))
	for k, v := range shs {
		rhs[k] = RuleHandler{
			Rule:                            *ConvertSourceRule(bundle, &v.Rule),
			Message:                         v.Message,
			Func:                            v.Func,
			NotAllowOfflineStmts:            v.NotAllowOfflineStmts,
			OnlyAuditNotExecutedSQL:         v.OnlyAuditNotExecutedSQL,
			NotSupportExecutedSQLAuditStmts: v.NotSupportExecutedSQLAuditStmts,
		}
	}
	return rhs
}

// ConvertSourceRule 将规则初始化时定义的 SourceRule 转换成 driverV2.Rule
func ConvertSourceRule(bundle *i18nPkg.Bundle, sr *SourceRule) *driverV2.Rule {
	r := &driverV2.Rule{
		Name:         sr.Name,
		Level:        sr.Level,
		CategoryTags: sr.CategoryTags,
		Params:       make(params.Params, 0, len(sr.Params)),
		I18nRuleInfo: genAllI18nRuleInfo(bundle, sr),
		AllowOffline: sr.AllowOffline,
		Version:      sr.Version,
	}
	for _, v := range sr.Params {
		r.Params = append(r.Params, &params.Param{
			Key:      v.Key,
			Value:    v.Value,
			Desc:     bundle.LocalizeMsgByLang(i18nPkg.DefaultLang, v.Desc),
			I18nDesc: bundle.LocalizeAll(v.Desc),
			Type:     v.Type,
			Enums:    nil, // all nil now
		})
	}

	return r
}

func genAllI18nRuleInfo(bundle *i18nPkg.Bundle, sr *SourceRule) map[language.Tag]*driverV2.RuleInfo {
	result := make(map[language.Tag]*driverV2.RuleInfo, len(bundle.LanguageTags()))
	for _, langTag := range bundle.LanguageTags() {
		newInfo := &driverV2.RuleInfo{
			Desc:       bundle.LocalizeMsgByLang(langTag, sr.Desc),
			Annotation: bundle.LocalizeMsgByLang(langTag, sr.Annotation),
			Category:   bundle.LocalizeMsgByLang(langTag, sr.Category),
			Knowledge:  driverV2.RuleKnowledge{Content: sr.Knowledge.Content}, //todo i18n Knowledge
		}

		result[langTag] = newInfo
	}
	return result
}
