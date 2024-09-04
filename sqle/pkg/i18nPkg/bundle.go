package i18nPkg

// todo i18n 将该包换到可以同时给dms sqle provision使用的地方

import (
	"context"
	"fmt"
	"io/fs"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	LocalizerCtxKey   = "localizer"
	AcceptLanguageKey = "Accept-Language"
)

var DefaultLang = language.Chinese

func NewBundleFromTomlDir(fsys fs.FS, log Log) (*Bundle, error) {
	if log == nil {
		log = &StdLogger{}
	}
	b := &Bundle{
		Bundle:     i18n.NewBundle(DefaultLang),
		localizers: nil,
		logger:     log,
	}
	b.Bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), "toml") {
			return nil
		}
		_, err = b.Bundle.LoadMessageFileFS(fsys, path)
		return err
	})
	if err != nil {
		return nil, err
	}

	b.localizers = make(map[language.Tag]*i18n.Localizer, len(b.Bundle.LanguageTags()))
	for _, tag := range b.Bundle.LanguageTags() {
		b.localizers[tag] = i18n.NewLocalizer(b.Bundle, tag.String())
	}

	return b, nil
}

type Bundle struct {
	*i18n.Bundle
	localizers map[language.Tag]*i18n.Localizer
	logger     Log
}

func (b *Bundle) DefaultLocalizer() *i18n.Localizer {
	return b.localizers[DefaultLang]
}

func (b *Bundle) GetLocalizer(tag language.Tag) *i18n.Localizer {
	if localizer, ok := b.localizers[tag]; ok {
		return localizer
	}
	return b.DefaultLocalizer()
}

func (b *Bundle) shouldLocalizeMsg(localizer *i18n.Localizer, msg *i18n.Message) string {
	if msg == nil {
		b.logger.Errorf("i18nPkg localize nil msg")
		return ""
	}
	m, err := localizer.LocalizeMessage(msg)
	if err != nil {
		b.logger.Errorf("i18nPkg LocalizeMessage %v failed: %v", msg.ID, err)
	}
	return m
}

func (b *Bundle) ShouldLocalizeMsg(ctx context.Context, msg *i18n.Message) string {
	l, ok := ctx.Value(LocalizerCtxKey).(*i18n.Localizer)
	if !ok {
		l = b.DefaultLocalizer()
		b.logger.Errorf("i18nPkg No localizer in context when localize msg: %v, use default", msg.ID)
	}

	return b.shouldLocalizeMsg(l, msg)
}

func (b *Bundle) ShouldLocalizeMsgByLang(lang language.Tag, msg *i18n.Message) string {
	l := b.GetLocalizer(lang)
	return b.shouldLocalizeMsg(l, msg)
}

func (b *Bundle) ShouldLocalizeAll(msg *i18n.Message) I18nStr {
	result := make(I18nStr, len(b.localizers))
	for langTag, localizer := range b.localizers {
		result[langTag] = b.shouldLocalizeMsg(localizer, msg)
	}
	return result
}

func (b *Bundle) ShouldLocalizeAllWithArgs(fmtMsg *i18n.Message, args ...any) I18nStr {
	result := make(I18nStr, len(b.localizers))
	for langTag, localizer := range b.localizers {
		result[langTag] = fmt.Sprintf(b.shouldLocalizeMsg(localizer, fmtMsg), args...)
	}
	return result
}

func (b *Bundle) ShouldLocalizeAllWithMsgArgs(fmtMsg *i18n.Message, msgArgs ...*i18n.Message) I18nStr {
	result := make(I18nStr, len(b.localizers))
	for langTag, localizer := range b.localizers {
		strs := make([]any, len(msgArgs))
		for k, m := range msgArgs {
			strs[k] = b.shouldLocalizeMsg(localizer, m)
		}
		result[langTag] = fmt.Sprintf(b.shouldLocalizeMsg(localizer, fmtMsg), strs...)
	}
	return result
}

func (b *Bundle) EchoMiddlewareByAcceptLanguage() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			acceptLang := c.Request().Header.Get(AcceptLanguageKey)
			langTag := DefaultLang
			for _, lang := range b.Bundle.LanguageTags() {
				if strings.HasPrefix(acceptLang, lang.String()) {
					langTag = lang
				}
			}

			ctx := context.WithValue(c.Request().Context(), LocalizerCtxKey, b.localizers[langTag])
			ctx = context.WithValue(ctx, AcceptLanguageKey, langTag)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func (b *Bundle) GetLangTagFromCtx(ctx context.Context) language.Tag {
	al, ok := ctx.Value(AcceptLanguageKey).(language.Tag)
	if ok {
		if _, ok = b.localizers[al]; ok {
			return al
		}
	}
	return DefaultLang
}
