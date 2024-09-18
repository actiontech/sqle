package rule

import (
	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/driver/mysql/plocale"
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

type SourceRule struct {
	Name       string
	Desc       *i18n.Message
	Annotation *i18n.Message
	Category   *i18n.Message
	Level      driverV2.RuleLevel
	Params     []*SourceParam
	Knowledge  driverV2.RuleKnowledge
}

type SourceHandler struct {
	Rule                 SourceRule
	Message              *i18n.Message
	Func                 RuleHandlerFunc
	AllowOffline         bool
	NotAllowOfflineStmts []ast.Node
	// 开始事后审核时将会跳过这个值为ture的规则
	OnlyAuditNotExecutedSQL bool
	// 事后审核时将会跳过下方列表中的类型
	NotSupportExecutedSQLAuditStmts []ast.Node
}

// 通过 source* 生成多语言版本的 RuleHandler
func generateI18nRuleHandlersFromSource(shs []*SourceHandler) []RuleHandler {
	rhs := make([]RuleHandler, len(shs))
	for k, v := range shs {
		rhs[k] = RuleHandler{
			Rule:                            *ConvertSourceRule(&v.Rule),
			Message:                         v.Message,
			Func:                            v.Func,
			AllowOffline:                    v.AllowOffline,
			NotAllowOfflineStmts:            v.NotAllowOfflineStmts,
			OnlyAuditNotExecutedSQL:         v.OnlyAuditNotExecutedSQL,
			NotSupportExecutedSQLAuditStmts: v.NotSupportExecutedSQLAuditStmts,
		}
	}
	return rhs
}

func ConvertSourceRule(sr *SourceRule) *driverV2.Rule {
	r := &driverV2.Rule{
		Name:         sr.Name,
		Level:        sr.Level,
		Params:       make(params.Params, 0, len(sr.Params)),
		I18nRuleInfo: genAllI18nRuleInfo(sr),
	}
	for _, v := range sr.Params {
		r.Params = append(r.Params, &params.Param{
			Key:      v.Key,
			Value:    v.Value,
			Desc:     plocale.Bundle.LocalizeMsgByLang(i18nPkg.DefaultLang, v.Desc),
			I18nDesc: plocale.Bundle.LocalizeAll(v.Desc),
			Type:     v.Type,
			Enums:    nil, // all nil now
		})
	}

	return r
}

func genAllI18nRuleInfo(sr *SourceRule) map[language.Tag]*driverV2.RuleInfo {
	result := make(map[language.Tag]*driverV2.RuleInfo, len(plocale.Bundle.LanguageTags()))
	for _, langTag := range plocale.Bundle.LanguageTags() {
		newInfo := &driverV2.RuleInfo{
			Desc:       plocale.Bundle.LocalizeMsgByLang(langTag, sr.Desc),
			Annotation: plocale.Bundle.LocalizeMsgByLang(langTag, sr.Annotation),
			Category:   plocale.Bundle.LocalizeMsgByLang(langTag, sr.Category),
			Knowledge:  driverV2.RuleKnowledge{Content: sr.Knowledge.Content}, //todo i18n Knowledge
		}

		result[langTag] = newInfo
	}
	return result
}
