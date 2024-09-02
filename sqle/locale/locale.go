package locale

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/actiontech/sqle/sqle/log"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// todo: 迁移到公共地方，给dms, sqle 用

const (
	LocalizerCtxKey   = "localizer"
	AcceptLanguageKey = "Accept-Language"
)

//go:embed active.*.toml
var LocaleFS embed.FS

var bundle *i18n.Bundle

var newEntry = log.NewEntry()

var DefaultLang = language.Chinese // todo i18n make sure plugins support

func init() {
	bundle = i18n.NewBundle(DefaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, err := bundle.LoadMessageFileFS(LocaleFS, "active.zh.toml")
	if err != nil {
		panic(fmt.Sprintf("load i18n config failed, error: %v", err))
	}
	_, err = bundle.LoadMessageFileFS(LocaleFS, "active.en.toml")
	if err != nil {
		panic(fmt.Sprintf("load i18n config failed, error: %v", err))
	}
}

func EchoMiddlewareI18nByAcceptLanguage() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			acceptLang := c.Request().Header.Get(AcceptLanguageKey)
			localizer := i18n.NewLocalizer(bundle, acceptLang)

			langTag := DefaultLang
			for _, lang := range bundle.LanguageTags() {
				if strings.HasPrefix(c.Request().Header.Get(AcceptLanguageKey), lang.String()) {
					langTag = lang
				}
			}

			ctx := context.WithValue(c.Request().Context(), LocalizerCtxKey, localizer)
			ctx = context.WithValue(ctx, AcceptLanguageKey, langTag)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func ShouldLocalizeMsg(ctx context.Context, msg *i18n.Message) string {
	l, ok := ctx.Value(LocalizerCtxKey).(*i18n.Localizer)
	if !ok {
		l = i18n.NewLocalizer(bundle)
		newEntry.Warnf("No localizer in context when localize msg: %v, use default", msg.ID)
	}

	m, err := l.LocalizeMessage(msg)
	if err != nil {
		newEntry.Errorf("LocalizeMessage: %v failed: %v", msg.ID, err)
	}
	return m
}

func GetLangTagFromCtx(ctx context.Context) language.Tag {
	al, ok := ctx.Value(AcceptLanguageKey).(language.Tag)
	if ok {
		return al
	}
	return DefaultLang
}
