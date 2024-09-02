package plocale

import (
	"embed"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/actiontech/sqle/sqle/locale"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

//go:embed active.*.toml
var LocaleFS embed.FS

var bundle *i18n.Bundle

var defaultLang = locale.DefaultLang

var DefaultLocalizer *i18n.Localizer

var AllLocalizers map[string]*i18n.Localizer

var newEntry = log.NewEntry()

func init() {
	bundle = i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFileFS(LocaleFS, "active.zh.toml")
	if err != nil {
		panic(fmt.Sprintf("MySQL plugin load i18n config failed, error: %v", err))
	}
	_, err = bundle.LoadMessageFileFS(LocaleFS, "active.en.toml")
	if err != nil {
		panic(fmt.Sprintf("MySQL plugin load i18n config failed, error: %v", err))
	}

	AllLocalizers = make(map[string]*i18n.Localizer, len(bundle.LanguageTags()))
	for _, langTag := range bundle.LanguageTags() {
		AllLocalizers[langTag.String()] = GetLocalizer(langTag.String())
	}
	DefaultLocalizer = AllLocalizers[defaultLang.String()]
}

func GetLocalizer(langs ...string) *i18n.Localizer {
	l := i18n.NewLocalizer(bundle, langs...)
	return l
}

func ShouldLocalizeMessage(localizer *i18n.Localizer, msg *i18n.Message) string {
	m, err := localizer.LocalizeMessage(msg)
	if err != nil {
		newEntry.Errorf("MySQL plugin LocalizeMessage %v failed: %v", msg.ID, err)
	}
	return m
}

func ShouldLocalizeAll(msg *i18n.Message) map[string]string {
	result := make(map[string]string, len(AllLocalizers))
	for langTag, localizer := range AllLocalizers {
		result[langTag] = ShouldLocalizeMessage(localizer, msg)
	}
	return result
}

func ShouldLocalizeAllWithArgs(msg *i18n.Message, args ...any) map[string]string {
	result := make(map[string]string, len(AllLocalizers))
	for langTag, localizer := range AllLocalizers {
		result[langTag] = fmt.Sprintf(ShouldLocalizeMessage(localizer, msg), args...)
	}
	return result
}

func ShouldLocalizeAllWithFmt(fmtMsg *i18n.Message, msg ...*i18n.Message) map[string]string {
	result := make(map[string]string, len(AllLocalizers))
	for langTag, localizer := range AllLocalizers {
		strs := make([]any, len(msg))
		for k, m := range msg {
			strs[k] = ShouldLocalizeMessage(localizer, m)
		}
		result[langTag] = fmt.Sprintf(ShouldLocalizeMessage(localizer, fmtMsg), strs...)
	}
	return result
}

func ConvertStr2I18n(s string) map[string]string {
	if s == "" {
		return nil
	}
	result := make(map[string]string, len(AllLocalizers))
	for langTag := range AllLocalizers {
		result[langTag] = s
	}
	return result
}
