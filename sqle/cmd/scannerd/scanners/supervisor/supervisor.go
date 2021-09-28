package supervisor

import (
	"context"
	"time"

	"github.com/actiontech/sqle/sqle/cmd/scannerd/scanners"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func Start(ctx context.Context, scanner scanners.Scanner, leastPushSecond, pushBufferSize int) error {
	runErrCh := make(chan error)
	go func() {
		err := scanner.Run(ctx)
		runErrCh <- err
	}()

	logrus.StandardLogger().Infoln("scanner started...")

	t := time.NewTicker(time.Second * time.Duration(leastPushSecond))
	defer t.Stop()

	sqlCh := scanner.SQLs()
	batch := make([]scanners.SQL, 0, pushBufferSize)
	for {
		select {
		case err := <-runErrCh:
			return err

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
