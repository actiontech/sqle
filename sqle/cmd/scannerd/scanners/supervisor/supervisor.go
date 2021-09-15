package supervisor

import (
	"context"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/config"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/mybatis"
	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners/slowquery"
	pkg "github.com/actiontech/sqle/sqle/pkg/scanner"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Start(cfg *config.Config) error {
	l := logrus.WithField("scanner_type", cfg.Typ)

	var err error
	var scanner scanners.Scanner
	client := pkg.NewSQLEClient(time.Second, cfg).WithToken(cfg.Token)

	switch cfg.Typ {
	case config.ScannerTypeMyBatis:
		p := &mybatis.Params{
			XMLDir: cfg.Dir,
			APName: cfg.AuditPlanName,
		}
		scanner, err = mybatis.New(p, l, client)

	case config.ScannerTypeSlowQuery:
		p := &slowquery.Params{
			APName:      cfg.AuditPlanName,
			LogFilePath: cfg.LogFilePath,
		}
		scanner, err = slowquery.New(p, l, client)

	default:
		err = errors.Errorf("unsupported scanner type %s.", cfg.Typ)
	}

	if err != nil {
		return err
	}

	return start(context.TODO(), scanner, 30, 1024)
}

func start(ctx context.Context, scanner scanners.Scanner, leastPushSecond, pushBufferSize int) error {
	go scanner.Run(ctx)
	logrus.StandardLogger().Infoln("scanner started...")

	t := time.NewTicker(time.Second * time.Duration(leastPushSecond))
	defer t.Stop()

	sqlCh := scanner.SQLs()
	batch := make([]scanners.SQL, 0, pushBufferSize)
	for {
		select {
		case <-ctx.Done():
			logrus.StandardLogger().Infoln("context done, exited")
			return nil

		case sql, ok := <-sqlCh:
			if !ok {
				logrus.StandardLogger().Infoln("SQL channel closed")
				if len(batch) != 0 {
					err := scanner.Upload(context.TODO(), batch)
					if err != nil {
						return errors.Wrap(err, "failed to upload sql")
					}
				}
				return nil
			}
			batch = append(batch, sql)
			if len(batch) != pushBufferSize {
				continue
			}

		case <-t.C:
			if len(batch) == 0 {
				continue
			}
		}

		err := scanner.Upload(context.TODO(), batch)
		if err != nil {
			return errors.Wrap(err, "failed to upload sql")
		}
		batch = make([]scanners.SQL, 0, pushBufferSize)
	}
}
