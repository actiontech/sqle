package plocale

import (
	"embed"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/i18nPkg"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed active.*.toml
var localeFS embed.FS

var Bundle *i18nPkg.Bundle

func init() {
	b, err := i18nPkg.NewBundleFromTomlDir(localeFS, log.NewEntry())
	if err != nil {
		panic(err)
	}
	Bundle = b
}

// todo i18n 把调用下列函数的地方改成直接使用 Bundle对应方法

// ShouldLocalizeMsgByLang todo i18n MySQL插件用到这个函数的地方应该都没完成国际化（低优先级）
func ShouldLocalizeMsgByLang(lang language.Tag, msg *i18n.Message) string {
	return Bundle.ShouldLocalizeMsgByLang(lang, msg)
}

func ShouldLocalizeAll(msg *i18n.Message) i18nPkg.I18nStr {
	return Bundle.ShouldLocalizeAll(msg)
}

func ShouldLocalizeAllWithArgs(msg *i18n.Message, args ...any) i18nPkg.I18nStr {
	return Bundle.ShouldLocalizeAllWithArgs(msg, args...)
}

func ShouldLocalizeAllWithFmt(fmtMsg *i18n.Message, msg ...*i18n.Message) i18nPkg.I18nStr {
	return Bundle.ShouldLocalizeAllWithMsgArgs(fmtMsg, msg...)
}

func ConvertStr2I18n(s string) i18nPkg.I18nStr {
	return i18nPkg.ConvertStr2I18nAsDefaultLang(s)
}
