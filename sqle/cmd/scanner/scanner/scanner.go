package scanner

import (
	"actiontech.cloud/sqle/sqle/sqle/cmd/scanner/config"
)

func Run(cfg *config.Config) error {
	switch cfg.Typ {
	case "mybatis":
		return MybatisScanner(cfg)
	default:
		return nil
	}
}
