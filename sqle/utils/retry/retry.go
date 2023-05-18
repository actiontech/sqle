package retry

import (
	"errors"
	"strings"
	"time"

	"github.com/actiontech/sqle/sqle/utils"
	"github.com/ungerik/go-dry"
)

type RetryableFunc func() error

func Do(retryableFunc RetryableFunc, doneChan chan struct{}, ops ...Option) error {

	cfg := NewDefaultRetryConfig()

	for i := range ops {
		op := ops[i]
		op(cfg)
	}

	{
		if cfg.attempts == 0 {
			if err := retryableFunc(); err != nil {
				return err
			}
			utils.TryClose(doneChan)
			return nil
		}
	}

	var idx uint = 0
	errList := errListType{}

	// cfg.attempts can not be 0.
	for idx < cfg.attempts {
		err := retryableFunc()
		if err == nil {
			utils.TryClose(doneChan)
			return nil
		}

		idx++
		errList = errList.AppendError(err.Error())
		time.Sleep(cfg.delay)
	}

	if len(errList) == 0 {
		utils.TryClose(doneChan)
		return nil
	}

	return errors.New(strings.Join(errList, "; "))
}

type errListType []string

func (e errListType) AppendError(errMsg string) errListType {
	if dry.StringInSlice(errMsg, e) {
		return e
	}
	return append(e, errMsg)
}
