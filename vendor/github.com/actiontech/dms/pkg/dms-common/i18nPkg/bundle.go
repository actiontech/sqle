package i18nPkg

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

func (b *Bundle) localizeMsg(localizer *i18n.Localizer, msg *i18n.Message) string {
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

func (b *Bundle) LocalizeMsgByCtx(ctx context.Context, msg *i18n.Message) string {
	l, ok := ctx.Value(LocalizerCtxKey).(*i18n.Localizer)
	if !ok {
		l = b.DefaultLocalizer()
		b.logger.Errorf("i18nPkg No localizer in context when localize msg: %v, use default", msg.ID)
	}

	return b.localizeMsg(l, msg)
}

func (b *Bundle) LocalizeMsgByLang(lang language.Tag, msg *i18n.Message) string {
	l := b.GetLocalizer(lang)
	return b.localizeMsg(l, msg)
}

func (b *Bundle) LocalizeAll(msg *i18n.Message) I18nStr {
	result := make(I18nStr, len(b.localizers))
	for langTag, localizer := range b.localizers {
		result[langTag] = b.localizeMsg(localizer, msg)
	}
	return result
}

// LocalizeAllWithArgs if there is any i18n.Message or I18nStr in args, it will be Localized in the corresponding language too
func (b *Bundle) LocalizeAllWithArgs(fmtMsg *i18n.Message, args ...any) I18nStr {
	result := make(I18nStr, len(b.localizers))
	msgs := map[int]*i18n.Message{}
	i18nStrs := map[int]*I18nStr{}
	for k, v := range args {
		switch arg := v.(type) {
		case i18n.Message:
			msgs[k] = &arg
		case *i18n.Message:
			msgs[k] = arg
		case I18nStr:
			i18nStrs[k] = &arg
		case *I18nStr:
			i18nStrs[k] = arg
		default:
		}
	}
	for langTag, localizer := range b.localizers {
		for k := range msgs {
			args[k] = b.localizeMsg(localizer, msgs[k])
		}
		for k := range i18nStrs {
			args[k] = i18nStrs[k].GetStrInLang(langTag)
		}
		result[langTag] = fmt.Sprintf(b.localizeMsg(localizer, fmtMsg), args...)
	}
	return result
}

func GetLangByAcceptLanguage(c echo.Context) string {
	return c.Request().Header.Get(AcceptLanguageKey)
}

func (b *Bundle) EchoMiddlewareByCustomFunc(getLang ...func(c echo.Context) string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var acceptLang string
			for _, f := range getLang {
				if lang := f(c); lang != "" {
					acceptLang = lang
					break
				}
			}
			langTag := b.MatchLangTag(acceptLang)
			ctx := context.WithValue(c.Request().Context(), LocalizerCtxKey, b.localizers[langTag])
			ctx = context.WithValue(ctx, AcceptLanguageKey, langTag)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func (b *Bundle) MatchLangTag(s string) language.Tag {
	langTag := DefaultLang
	for _, lang := range b.Bundle.LanguageTags() {
		if strings.HasPrefix(s, lang.String()) {
			langTag = lang
			break
		}
	}
	return langTag
}

func (b *Bundle) JoinI18nStr(elems []I18nStr, sep string) I18nStr {
	var result = make(I18nStr, len(b.LanguageTags()))
	for _, langTag := range b.LanguageTags() {
		var langStr []string
		for _, v := range elems {
			langStr = append(langStr, v.GetStrInLang(langTag))
		}
		result.SetStrInLang(langTag, strings.Join(langStr, sep))
	}
	return result
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
