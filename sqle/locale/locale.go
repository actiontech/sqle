package locale

import (
	"context"
	"embed"

	"github.com/actiontech/sqle/sqle/log"
	"github.com/actiontech/sqle/sqle/pkg/i18nPkg"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed active.*.toml
var localeFS embed.FS

var Bundle *i18nPkg.Bundle

var DefaultLang = i18nPkg.DefaultLang // todo i18n make sure plugins support

func init() {
	b, err := i18nPkg.NewBundleFromTomlDir(localeFS, log.NewEntry())
	if err != nil {
		panic(err)
	}
	Bundle = b
}

// todo i18n 把调用下列函数的地方改成直接使用 Bundle对应方法

func EchoMiddlewareI18nByAcceptLanguage() echo.MiddlewareFunc {
	return Bundle.EchoMiddlewareByAcceptLanguage()
}

func ShouldLocalizeMsg(ctx context.Context, msg *i18n.Message) string {
	return Bundle.ShouldLocalizeMsg(ctx, msg)
}

func GetLangTagFromCtx(ctx context.Context) language.Tag {
	return Bundle.GetLangTagFromCtx(ctx)
}

func ShouldLocalizeAll(msg *i18n.Message) i18nPkg.I18nStr {
	return Bundle.ShouldLocalizeAll(msg)
}

func ShouldLocalizeAllWithArgs(fmtMsg *i18n.Message, args ...any) i18nPkg.I18nStr {
	return Bundle.ShouldLocalizeAllWithArgs(fmtMsg, args)
}
