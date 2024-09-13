package locale

import (
	"embed"

	"github.com/actiontech/dms/pkg/dms-common/i18nPkg"
	"github.com/actiontech/sqle/sqle/log"
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
